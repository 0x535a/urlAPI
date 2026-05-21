package database

import (
	"zhongxin/internal/model"

	"gorm.io/gorm/clause"
)

func (m *MSSQLAdapter) GetMachineONLogByFilter(cls []clause.Expression) ([]model.MachineLog, bool, error) {
	var logs []model.MachineLog
	result := m.db
	for _, c := range cls {
		result = result.Where(c)
	}
	result = result.Where(clause.Eq{
		Column: "LogMachineryType",
		Value:  1,
	}).Find(&logs)
	if result.Error != nil {
		return logs, false, result.Error
	}
	if result.RowsAffected == 0 {
		return logs, false, nil
	}
	return logs, true, nil
}

func (m *MSSQLAdapter) GetAllMachinePasID() ([]model.MachinePasID, bool, error) {
	var ids []model.MachinePasID
	result := m.db.Model(&model.MachinePasID{}).Find(&ids)
	if result.Error != nil {
		return ids, false, result.Error
	}
	if result.RowsAffected == 0 {
		return ids, false, nil
	}
	return ids, true, nil
}
