package database

import (
	"zhongxin/internal/model"

	"gorm.io/gorm/clause"
)

func (s *SQLiteAdapter) NewMachine(m model.Machine) error {
	return s.db.Create(&m).Error
}

func (s *SQLiteAdapter) GetMachineByID(id string) (model.Machine, bool, error) {
	var m model.Machine
	result := s.db.Where(clause.Eq{Column: "id", Value: id}).First(&m)
	if result.Error != nil {
		return m, false, result.Error
	}
	if result.RowsAffected == 0 {
		return m, false, nil
	}
	return m, true, nil
}

func (s *SQLiteAdapter) UpdateMachine(user model.Machine) error {
	return s.db.Save(&user).Error
}

func (s *SQLiteAdapter) DeleteMachineByID(id string) error {
	return s.db.Where(clause.Eq{Column: "id", Value: id}).Delete(&model.Machine{}).Error
}

func (s *SQLiteAdapter) GetAllMachinePasID() ([]string, error) {
	var ids []string
	result := s.db.Model(&model.Machine{}).Pluck("id", &ids)
	if result.Error != nil {
		return ids, result.Error
	}
	return ids, nil
}

func (s *SQLiteAdapter) GetAllMachines() ([]model.Machine, error) {
	var machines []model.Machine
	result := s.db.Find(&machines)
	if result.Error != nil {
		return machines, result.Error
	}
	return machines, nil
}
