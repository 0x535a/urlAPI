package op

import (
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func GetMachineByID(machineID string) (model.Machine, error, int) {
	machineID = util.RemoveLeadingZeros(machineID)
	machine, _, err := db.GetMachineByID(machineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Machine{}, errors.WithStack(errors.New("machine not found")), _error.DBRecordNotFound
		}
		return model.Machine{}, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return machine, nil, 0
}

func UpdateMachine(machine model.Machine) (error, int) {
	if err := db.UpdateMachine(machine); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func UpdateMachineWithOrder(machineIDs []string, orderID string, userID string) (error, int) {
	for _, machineID := range machineIDs {
		machineID = util.RemoveLeadingZeros(machineID)
		machine, _, err := db.GetMachineByID(machineID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.WithStack(errors.New("machine not found")), _error.DBRecordNotFound
			}
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
		if machine.OrderID != "" {
			orderInfo, err, errCode := GetOrderByID(machine.OrderID)
			if err != nil {
				return errors.WithStack(err), errCode
			}
			orderInfo.AssignedMachinesDB = util.RemoveFromDBList(orderInfo.AssignedMachinesDB, machine.ID)
			if err, errCode := UpdateOrder(orderInfo, userID); err != nil {
				return errors.WithStack(err), errCode
			}
		}
		machine.OrderID = orderID
		if err := db.UpdateMachine(machine); err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
	}
	return nil, 0
}

func DeleteMachineByID(machineID string) (error, int) {
	if err := db.DeleteMachineByID(machineID); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func NewMachine(machine model.Machine) (error, int) {
	machine.IsAvailable = true
	if err := db.NewMachine(machine); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func GetAllMachines() ([]model.Machine, error, int) {
	machines, err := db.GetAllMachines()
	if err != nil {
		return nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return machines, nil, 0
}
