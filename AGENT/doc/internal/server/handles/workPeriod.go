package handles

import (
	"reflect"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/internal/op"
	"zhongxin/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type WorkPeriodRequestMeta struct {
	model.WorkPeriod
	WorkPeriodID  string `json:"workPeriodID"`
	WorkPeriodIDs string `form:"workPeriodIDs"`
}

type WorkPeriodSingleResponseMeta struct {
	Code int              `json:"code"`
	Data model.WorkPeriod `json:"data"`
	Msg  string           `json:"msg"`
}

type WorkPeriodMultiResponseMeta struct {
	Code int                `json:"code"`
	Data []model.WorkPeriod `json:"data"`
	Msg  string             `json:"msg"`
}

func WorkPeriodStartHandler(c *gin.Context) {
	// Attention: Hasn't checked whether the target order is finished
	var req WorkPeriodRequestMeta
	var res WorkPeriodSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userI, _ := c.Get("user")
	userID := userI.(string)
	userInfo, err, errCode := op.GetUserByID(userID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	machineInfo, err, errCode := op.GetMachineByID(req.MachineID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}
	if !machineInfo.IsAvailable {
		ErrorResponse(c, _error.OpMachineNotAvailable, "Machine is not available", &res)
		return
	}
	if machineInfo.WorkingPeriodID != "" {
		ErrorResponse(c, _error.OpMachineAlreadyInUse, "Machine is already in use", &res)
		return
	}
	if machineInfo.OrderID == "" {
		ErrorResponse(c, _error.OpMachineNotAssignedToOrder, "Machine is not assigned to any order", &res)
	} else {
		req.OrderID = machineInfo.OrderID
	}

	orderInfo, err, errCode := op.GetOrderByID(req.OrderID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}
	if !util.ExistInDBList(orderInfo.AssignedToes, userID) {
		ErrorResponse(c, _error.OpOrderNotAssignedToWorker, "User is not assigned to this order", &res)
		return
	}

	workPeriod, err, errCode := op.NewWorkPeriod(userInfo, machineInfo, orderInfo)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	userInfo.IsWorking = true
	userInfo.WorkPeriods = util.StringListToDB(append(util.DBToStringList(userInfo.WorkPeriods), workPeriod.ID))
	machineInfo.WorkingPeriodID = workPeriod.ID
	if orderInfo.Status == 0 {
		orderInfo.Status = 1
	}

	if err, errCode = op.UpdateUser(userInfo); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}
	if err, errCode = op.UpdateMachine(machineInfo); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}
	if err, errCode = op.UpdateOrder(orderInfo, userID); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Work Period Started Successfully"
	res.Data = workPeriod

	c.JSON(200, res)
}

func WorkPeriodStopHandler(c *gin.Context) {
	// Attention: Order markedFinishedTime hasn't been set
	var req WorkPeriodRequestMeta
	var ret WorkPeriodSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &ret)
		return
	}

	userI, _ := c.Get("user")
	userID := userI.(string)
	userInfo, err, errCode := op.GetUserByID(userID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &ret)
		return
	}

	workPeriod, err, errCode := op.GetWorkPeriodByID(req.WorkPeriodID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &ret)
		return
	}

	machine, err, errCode := op.GetMachineByID(workPeriod.MachineID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &ret)
		return
	}

	userInfo.WorkPeriods = util.RemoveFromDBList(userInfo.WorkPeriods, req.WorkPeriodID)
	if len(userInfo.WorkPeriods) == 0 {
		userInfo.IsWorking = false
	}
	workPeriod.EndTime = util.TimeNow().Unix()
	machine.WorkingPeriodID = ""

	if err, errCode = op.UpdateUser(userInfo); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &ret)
		return
	}

	if err, errCode = op.UpdateWorkPeriod(workPeriod); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &ret)
		return
	}
	if err, errCode = op.UpdateMachine(machine); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &ret)
		return
	}

	ret.Code = 0
	ret.Msg = "Work Period Stopped Successfully"
	ret.Data = workPeriod
	c.JSON(200, ret)
}

func WorkPeriodGetHandler(c *gin.Context) {
	var req WorkPeriodRequestMeta
	var res WorkPeriodMultiResponseMeta
	if err := c.ShouldBind(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userT, _ := c.Get("userType")
	userI, _ := c.Get("user")
	userType := userT.(string)
	userID := userI.(string)

	workPeriods := make([]model.WorkPeriod, 0)
	workPeriodList := util.DBToStringList(req.WorkPeriodIDs)
	for _, workPeriod := range workPeriodList {
		workPeriodDB, err, errCode := op.GetWorkPeriodByID(workPeriod)
		if err != nil {
			_error.Print(err)
			ErrorResponse(c, errCode, err.Error(), &res)
			return
		}
		if userType == "admin" || userID == workPeriodDB.WorkerID {
			workPeriods = append(workPeriods, workPeriodDB)
		}
	}

	res.Code = 0
	res.Msg = "Get Work Period(s) Successfully"
	res.Data = workPeriods

	c.JSON(200, res)
}

func WorkPeriodDeleteHandler(c *gin.Context) {
	var req WorkPeriodRequestMeta
	var res WorkPeriodSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userI, _ := c.Get("user")
	userID := userI.(string)
	userInfo, err, errCode := op.GetUserByID(userID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	userInfo.WorkPeriods = util.RemoveFromDBList(userInfo.WorkPeriods, req.WorkPeriodID)
	if err, errCode = op.UpdateUser(userInfo); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	if err, errCode = op.DeleteWorkPeriodByID(req.WorkPeriodID); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Work Period Deleted Successfully"

	c.JSON(200, res)
}

func WorkPeriodUpdateHandler(c *gin.Context) {
	var req WorkPeriodRequestMeta
	var res WorkPeriodSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	workPeriodDB, err, errCode := op.GetWorkPeriodByID(req.WorkPeriodID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	t := reflect.TypeOf(req.WorkPeriod)
	v := reflect.ValueOf(req.WorkPeriod)
	wv := reflect.ValueOf(&workPeriodDB).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if util.IsEmpty(value) {
			continue
		}

		wf := wv.FieldByName(field.Name)
		if !wf.IsValid() || !wf.CanSet() {
			continue
		}
		wf.Set(v.Field(i))
	}
	if err, errCode = op.UpdateWorkPeriod(workPeriodDB); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Work Period Updated Successfully"
	res.Data = workPeriodDB

	c.JSON(200, res)
}

func WorkPeriodListHandler(c *gin.Context) {
	var req WorkPeriodRequestMeta
	var res WorkPeriodMultiResponseMeta
	if err := c.ShouldBind(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userT, _ := c.Get("userType")
	userI, _ := c.Get("user")
	userType := userT.(string)
	userID := userI.(string)
	if userType != "admin" {
		req.WorkerID = userID
	}

	cls := make([]clause.Expression, 0)
	t := reflect.TypeOf(req)
	v := reflect.ValueOf(req)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if util.IsEmpty(value) {
			continue
		}
		switch field.Name {
		case "startTime":
			cls = append(cls, clause.Gte{Column: field.Tag.Get("json"), Value: value})
		case "endTime":
			cls = append(cls, clause.Lte{Column: field.Tag.Get("json"), Value: value})
		default:
			cls = append(cls, clause.Eq{Column: field.Tag.Get("json"), Value: value})
		}
	}
	t = reflect.TypeOf(req.WorkPeriod)
	v = reflect.ValueOf(req.WorkPeriod)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if util.IsEmpty(value) {
			continue
		}
		cls = append(cls, clause.Eq{Column: field.Tag.Get("json"), Value: value})
	}

	workPeriods, err, errCode := op.GetWorkPeriodsByClauses(cls)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Get Work Period(s) Successfully"
	res.Data = workPeriods

	c.JSON(200, res)
}
