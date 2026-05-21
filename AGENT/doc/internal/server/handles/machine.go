package handles

import (
	"reflect"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/internal/op"
	"zhongxin/util"

	"github.com/gin-gonic/gin"
)

type MachineRequestMeta struct {
	model.Machine
	MachineIDs     string `form:"machineIDs"`
	MachineID      string `json:"machineID"`
	IsAvailablePtr *bool  `json:"isAvailable"`
}

type MachineSingleResponseMeta struct {
	Code int           `json:"code"`
	Data model.Machine `json:"data"`
	Msg  string        `json:"msg"`
}

type MachineMultiResponseMeta struct {
	Code int             `json:"code"`
	Data []model.Machine `json:"data"`
	Msg  string          `json:"msg"`
}

func MachineGetHandler(c *gin.Context) {
	var req MachineRequestMeta
	var res MachineMultiResponseMeta
	if err := c.ShouldBind(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), nil)
		return
	}

	machineList := util.DBToStringList(req.MachineIDs)
	for _, machine := range machineList {
		machineInfo, err, errCode := op.GetMachineByID(machine)
		if err != nil {
			_error.Print(err)
			ErrorResponse(c, errCode, err.Error(), &res)
			return
		}
		res.Data = append(res.Data, machineInfo)
	}
	res.Code = 0
	res.Msg = "Get Machine Info Successfully"
	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func MachineUpdateHandler(c *gin.Context) {
	var req MachineRequestMeta
	var res MachineSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	machineInfo, err, errCode := op.GetMachineByID(req.MachineID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	if req.IsAvailablePtr != nil {
		machineInfo.IsAvailable = *req.IsAvailablePtr
	}
	v := reflect.ValueOf(req.Machine)
	t := reflect.TypeOf(req.Machine)
	mv := reflect.ValueOf(&machineInfo).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if util.IsEmpty(value) {
			continue
		}

		mf := mv.FieldByName(field.Name)
		if !mf.CanSet() || !mf.IsValid() {
			continue
		}
		mf.Set(v.Field(i))
	}
	if err, errCode := op.UpdateMachine(machineInfo); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Update Machine Successfully"
	res.Data = machineInfo

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func MachineDeleteHandler(c *gin.Context) {
	var req MachineRequestMeta
	var res MachineSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	if err, errCode := op.DeleteMachineByID(req.MachineID); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Delete Machine Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func MachineAddHandler(c *gin.Context) {
	var req MachineRequestMeta
	var res MachineSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	if err, errCode := op.NewMachine(req.Machine); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Machine Created Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func MachineListHandler(c *gin.Context) {
	var res MachineMultiResponseMeta
	machines, err, errCode := op.GetAllMachines()
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}
	res.Code = 0
	res.Msg = "Get Machine List Successfully"
	res.Data = machines

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}
