package database

import (
	"zhongxin/internal/model"
	"zhongxin/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *SQLiteAdapter) CreateOrder(order model.Order) error {
	order.VersionID = util.NewUUID()
	order.IsLatest = true
	return s.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&order).Error
}

func (s *SQLiteAdapter) DeleteOrderByID(id string) error {
	var orders []model.Order
	result := s.db.Where(clause.Eq{Column: "id", Value: id}).Find(&orders)
	if result.Error != nil {
		return result.Error
	}
	for _, order := range orders {
		if err := s.db.Select(clause.Associations).Delete(&order).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLiteAdapter) UpdateOrder(order model.Order) (string, string, error) {
	var oldVersion model.Order
	result := s.db.Where(clause.Eq{Column: "id", Value: order.ID}).
		Where(clause.Eq{Column: "isLatest", Value: true}).
		Preload(clause.Associations).
		First(&oldVersion)
	if result.Error != nil {
		return "", "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", "", nil
	}

	oldVersion.IsLatest = false
	oldVersion.ColorAndQuantities = nil
	order.IsLatest = true
	order.VersionID = util.NewUUID()
	if err := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&oldVersion).Error; err != nil {
		return "", "", err
	}
	if err := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&order).Error; err != nil {
		return "", "", err
	}

	return oldVersion.VersionID, order.VersionID, nil
}

func (s *SQLiteAdapter) UpdateOrderWithoutVersion(order model.Order) error {
	return s.db.Save(&order).Error
}

func (s *SQLiteAdapter) GetOrderByID(id string) (model.Order, bool, error) {
	var o model.Order
	result := s.db.Where(clause.Eq{Column: "id", Value: id}).
		Where(clause.Eq{Column: "isLatest", Value: true}).
		Preload(clause.Associations).
		First(&o)
	if result.Error != nil {
		return o, false, result.Error
	}
	if result.RowsAffected == 0 {
		return o, false, nil
	}
	return o, true, nil
}

func (s *SQLiteAdapter) GetOrdersByFilter(cls []clause.Expression) ([]model.Order, bool, error) {
	var orders []model.Order
	result := s.db
	for _, c := range cls {
		result = result.Where(c)
	}
	result = result.Preload(clause.Associations).
		Find(&orders)
	if result.Error != nil {
		return orders, false, result.Error
	}
	if result.RowsAffected == 0 {
		return orders, false, nil
	}
	return orders, true, nil
}
