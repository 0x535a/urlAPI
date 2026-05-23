package database

import (
	"zhongxin/internal/model"
	"zhongxin/util"

	"gorm.io/gorm/clause"
)

func (s *SQLiteAdapter) GetNotificationsByClauses(cls []clause.Expression) ([]model.Notification, bool, error) {
	var notifications []model.Notification
	result := s.db
	for _, c := range cls {
		result = result.Where(c)
	}
	result = result.Find(&notifications)
	if result.Error != nil {
		return notifications, false, result.Error
	}
	if result.RowsAffected == 0 {
		return notifications, false, nil
	}
	return notifications, true, nil

}

func (s *SQLiteAdapter) GetNotificationByID(id string) (model.Notification, bool, error) {
	var notification model.Notification
	result := s.db.Where(clause.Eq{Column: "id", Value: id}).First(&notification)
	if result.Error != nil {
		return notification, false, result.Error
	}
	if result.RowsAffected == 0 {
		return notification, false, nil
	}
	return notification, true, nil
}

func (s *SQLiteAdapter) GetNotificationByUserID(userID string) ([]model.Notification, bool, error) {
	var notifications []model.Notification
	result := s.db.Where(clause.Eq{Column: "userID", Value: userID}).Find(&notifications)
	if result.Error != nil {
		return notifications, false, result.Error
	}
	if result.RowsAffected == 0 {
		return notifications, false, nil
	}
	return notifications, true, nil
}

func (s *SQLiteAdapter) DeleteNotificationByID(id string) error {
	return s.db.Where(clause.Eq{Column: "id", Value: id}).Delete(&model.Notification{}).Error
}

func (s *SQLiteAdapter) DeleteNotificationByUserID(userID string) error {
	return s.db.Where(clause.Eq{Column: "userID", Value: userID}).Delete(&model.Notification{}).Error
}

func (s *SQLiteAdapter) NewNotification(n model.Notification) error {
	n.ID = util.NewUUID()
	return s.db.Create(&n).Error
}

func (s *SQLiteAdapter) UpdateNotification(n model.Notification) error {
	return s.db.Save(&n).Error
}
