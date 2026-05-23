package op

import (
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/util"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func NewWorkPeriod(user model.User, machine model.Machine, order model.Order) (model.WorkPeriod, error, int) {
	workPeriod := model.WorkPeriod{
		ID:                  util.NewUUID(),
		WorkerID:            user.ID,
		MachineID:           machine.ID,
		OrderID:             order.ID,
		StartTime:           util.TimeNow().Unix(),
		IsMachineDataMerged: false,
		MachineName:         machine.Name,
		OrderName:           order.Name,
		UnitPriceYuan:       machine.UnitPriceYuan,
	}
	if err := db.CreateWorkPeriod(workPeriod); err != nil {
		return model.WorkPeriod{}, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return workPeriod, nil, 0
}

func GetWorkPeriodByID(id string) (model.WorkPeriod, error, int) {
	workPeriodDB, _, err := db.GetWorkPeriodByID(id)
	if err != nil {
		return model.WorkPeriod{}, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return workPeriodDB, nil, 0
}

func UpdateWorkPeriod(wp model.WorkPeriod) (error, int) {
	if err := db.UpdateWorkPeriod(wp); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func DeleteWorkPeriodByID(id string) (error, int) {
	err := db.DeleteWorkPeriodByID(id)
	if err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func GetWorkPeriodsByClauses(cls []clause.Expression) ([]model.WorkPeriod, error, int) {
	workPeriods, _, err := db.GetWorkPeriodsByClauses(cls)
	if err != nil {
		return nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return workPeriods, nil, 0
}
