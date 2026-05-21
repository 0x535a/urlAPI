package database

import (
	"zhongxin/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *SQLiteAdapter) CreateWorkPeriod(wp model.WorkPeriod) error {
	return s.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&wp).Error
}

func (s *SQLiteAdapter) UpdateWorkPeriod(wp model.WorkPeriod) error {
	if err := s.db.Where(clause.Eq{
		Column: "workPeriodID",
		Value:  wp.ID,
	}).Delete(&model.MachinePeriod{}).Error; err != nil {
		return err
	}
	return s.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&wp).Error
}

func (s *SQLiteAdapter) GetWorkPeriodByID(id string) (model.WorkPeriod, bool, error) {
	var wp model.WorkPeriod
	result := s.db.Where(clause.Eq{Column: "id", Value: id}).
		Preload(clause.Associations).
		First(&wp)
	if result.Error != nil {
		return wp, false, result.Error
	}
	if result.RowsAffected == 0 {
		return wp, false, nil
	}
	return wp, true, nil
}

// Attention
func (s *SQLiteAdapter) GetWorkPeriodsByClauses(cls []clause.Expression) ([]model.WorkPeriod, bool, error) {
	var workPeriods []model.WorkPeriod
	result := s.db
	for _, c := range cls {
		result = result.Where(c)
	}
	result = result.Preload(clause.Associations).Find(&workPeriods)
	if result.Error != nil {
		return workPeriods, false, result.Error
	}
	if result.RowsAffected == 0 {
		return workPeriods, false, nil
	}
	return workPeriods, true, nil
}

func (s *SQLiteAdapter) DeleteWorkPeriodByID(id string) error {
	return s.db.Select(clause.Associations).Delete(&model.WorkPeriod{ID: id}).Error
}
