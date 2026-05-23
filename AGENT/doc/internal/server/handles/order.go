package handles

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"zhongxin/internal/conf"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/internal/op"
	"zhongxin/util"

	"github.com/gin-gonic/gin"
)

type OrderRequestMeta struct {
	model.Order
	OrderID             string    `json:"orderID" form:"orderID"`
	OrderIDs            string    `form:"orderIDs"`
	OrderTypes          []string  `json:"orderTypes"`
	StatusPtr           *int      `json:"status"`
	NotesPtr            *string   `json:"notes"`
	AssignedTo          []string  `json:"assignedTo"`
	WorkPeriods         []string  `json:"workPeriods"`
	AssignedMachinesPtr *[]string `json:"assignedMachines"`
	Page                int       `json:"page"`
	Limit               int       `json:"limit"`
	Conditions          []struct {
		Key string      `json:"key"`
		Rel string      `json:"rel"`
		Val interface{} `json:"val"`
	} `json:"conditions"`
}

type OrderMultiResponseMeta struct {
	Code int `json:"code"`
	Data []struct {
		model.Order
		AssignedTo           []string `json:"assignedTo"`
		WorkPeriods          []string `json:"workPeriods"`
		AssignedMachines     []string `json:"assignedMachines"`
		AssignedMachineNames []string `json:"assignedMachineNames"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type OrderSingleResponseMeta struct {
	Code int `json:"code"`
	Data struct {
		OrderTypes []string `json:"orderTypes"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func OrderAddHandler(c *gin.Context) {
	var req OrderRequestMeta
	var res OrderMultiResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	assignedMachines := []string{}
	notes := ""
	req.ID = util.NewUUID()
	if req.NotesPtr != nil {
		notes = *req.NotesPtr
	}
	if req.AssignedMachinesPtr != nil {
		assignedMachines = *req.AssignedMachinesPtr
	}
	req.Notes = notes
	req.AssignedMachinesDB = util.StringListToDB(assignedMachines)
	req.AssignedToes = util.StringListToDB(req.AssignedTo)
	if err, errCode := op.NewOrder(req.Order); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	userI, _ := c.Get("user")
	if err, errCode := op.UpdateMachineWithOrder(assignedMachines, req.ID, userI.(string)); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Data = nil
	res.Msg = "Order Created Successfully"
	//if err, _ := op.NewNotificationOrderNew(req.Order, util.DBToStringList(req.Order.AssignedToes)); err != nil {
	//	res.Code = _error.ServerNotificationFailed
	//	res.Msg = fmt.Sprintf("Order Created Successfully, but notification failed: %s", err.Error())
	//}
	if err, _ := op.NewNotificationOrderUpdate(req.Order, conf.WxNotificationOrderNewTemplate, req.AssignedTo); err != nil {
		res.Code = _error.ServerNotificationFailed
		res.Msg = fmt.Sprintf("Order Created Successfully, but notification failed: %s", err.Error())
	}

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func OrderUpdateHandler(c *gin.Context) {
	var req OrderRequestMeta
	var res OrderMultiResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}
	res.Code = 0
	res.Msg = "Order Updated Successfully"

	orderInfo, err, errCode := op.GetOrderByID(req.OrderID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	userT, _ := c.Get("userType")
	if userT.(string) == "worker" {
		req = OrderRequestMeta{}
		if orderInfo.Status == 1 {
			req.StatusPtr = util.GetPtr(2)
		}
		adminUser, err, _ := op.GetAllAdmins()
		if err != nil {
			res.Code = _error.ServerNotificationFailed
			res.Msg = fmt.Sprintf("Order Deleted Successfully, but notification failed: %s", err.Error())
		} else {
			var adminIDs []string
			for _, admin := range adminUser {
				adminIDs = append(adminIDs, admin.ID)
			}
			err, _ = op.NewNotificationOrderConfirm(orderInfo, adminIDs)
			if err != nil {
				res.Code = _error.ServerNotificationFailed
				res.Msg = fmt.Sprintf("Order Deleted Successfully, but notification failed: %s", err.Error())
			}
		}
	}

	userI, _ := c.Get("user")
	var newWorkers, delWorkers []string
	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if util.IsEmpty(value) {
			continue
		}

		// Attention: number & int64
		switch field.Name {
		case "AssignedTo":
			newWorkers, delWorkers = util.DiffSlice(req.AssignedTo, util.DBToStringList(orderInfo.AssignedToes))
			orderInfo.AssignedToes = util.StringListToDB(req.AssignedTo)
		case "WorkPeriods":
			orderInfo.WorkPeriodsDB = util.StringListToDB(req.WorkPeriods)
		case "StatusPtr":
			orderInfo.Status = *req.StatusPtr
			switch *req.StatusPtr {
			case 2:
				orderInfo.MarkedFinishedTime = util.TimeNow().Unix()
			case 3:
				orderInfo.ConfirmedFinishedTime = util.TimeNow().Unix()
				if err, _ = op.NewNotificationOrderUpdate(orderInfo, conf.WxNotificationOrderUpdateTemplate, util.DBToStringList(orderInfo.AssignedToes)); err != nil {
					res.Code = _error.ServerNotificationFailed
					res.Msg = fmt.Sprintf("Order Comfirmed Successfully, but notification failed: %s", err.Error())
				}
			}
		case "NotesPtr":
			orderInfo.Notes = *req.NotesPtr
		case "AssignedMachinesPtr":
			if err, errCode := op.UpdateMachineWithOrder(util.DBToStringList(orderInfo.AssignedMachinesDB), "", ""); err != nil {
				_error.Print(err)
				ErrorResponse(c, errCode, err.Error(), &res)
				return
			}
			orderInfo.AssignedMachinesDB = util.StringListToDB(*req.AssignedMachinesPtr)
			if err, errCode := op.UpdateMachineWithOrder(*req.AssignedMachinesPtr, req.OrderID, userI.(string)); err != nil {
				_error.Print(err)
				ErrorResponse(c, errCode, err.Error(), &res)
				return
			}
		}
	}

	v = reflect.ValueOf(req.Order)
	t = reflect.TypeOf(req.Order)
	ov := reflect.ValueOf(&orderInfo).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if util.IsEmpty(value) {
			continue
		}

		of := ov.FieldByName(field.Name)
		if !of.IsValid() || !of.CanSet() {
			continue
		}

		of.Set(v.Field(i))
		//if cnMsgs, ok := conf.OrderStatusMap[field.Name]; ok {
		//	updateField = append(updateField, fmt.Sprintf("%s: %s;", field.Tag.Get("cn"), cnMsgs[of.Int()]))
		//} else {
		//	updateField = append(updateField, fmt.Sprintf("%s: %s;", field.Tag.Get("cn"), value))
		//}
	}

	userI, _ = c.Get("user")
	if err, errCode := op.UpdateOrder(orderInfo, userI.(string)); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	//for _, notiType := range notiList {
	//	switch notiType {
	//	case conf.OrderNewNotificationTemplateID:
	//		if err, _ = op.NewNotificationOrderNew(orderInfo, newWorker); err != nil {
	//			_error.Print(err)
	//		}
	//	case conf.OrderCompleteNotificationTemplateID:
	//		if err, _ = op.NewNotificationOrderComplete(orderInfo); err != nil {
	//			_error.Print(err)
	//		}
	//	case conf.OrderConfirmNotificationTemplateID:
	//		if err, _ = op.NewNotificationOrderConfirm(orderInfo); err != nil {
	//			_error.Print(err)
	//		}
	//	case conf.OrderDeleteNotificationTemplateID:
	//		if err, _ = op.NewNotificationOrderDelete(orderInfo, delWorker); err != nil {
	//			_error.Print(err)
	//		}
	//	}
	//	if err != nil {
	//		res.Code = _error.ServerNotificationFailed
	//		res.Msg = fmt.Sprintf("Order Created Successfully, but notification failed: %s", err.Error())
	//	}
	//	if len(updateField) > 0 {
	//		if err, _ = op.NewNotificationOrderUpdate(orderInfo, util.StringListToDB(updateField)); err != nil {
	//			_error.Print(err)
	//		}
	//	}
	if len(newWorkers) > 0 {
		err, _ = op.NewNotificationOrderUpdate(orderInfo, conf.WxNotificationOrderNewTemplate, newWorkers)
	}
	if len(delWorkers) > 0 {
		err, _ = op.NewNotificationOrderUpdate(orderInfo, conf.WxNotificationOrderDeleteTemplate, delWorkers)
	}
	//err, _ = op.NewNotificationOrderUpdate(orderInfo, conf.WxNotificationOrderUpdateTemplate, util.DBToStringList(orderInfo.AssignedToes))
	if err != nil {
		res.Code = _error.ServerNotificationFailed
		res.Msg = fmt.Sprintf("Order Created Successfully, but notification failed: %s", err.Error())
	}

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func OrderDeleteHandler(c *gin.Context) {
	var req OrderRequestMeta
	var res OrderMultiResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Order Deleted Successfully"

	orderInfo, err, _ := op.GetOrderByID(req.OrderID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}
	if err, _ := op.NewNotificationOrderDelete(orderInfo, util.DBToStringList(orderInfo.AssignedToes)); err != nil {
		_error.Print(err)
		res.Code = _error.ServerNotificationFailed
		res.Msg = fmt.Sprintf("Order Deleted Successfully, but notification failed: %s", err.Error())
	}

	userI, _ := c.Get("user")
	if err, errCode := op.UpdateMachineWithOrder(util.DBToStringList(orderInfo.AssignedMachinesDB), "", userI.(string)); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}
	if err, errCode := op.DeleteOrderByID(req.OrderID); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func OrderQueryHandler(c *gin.Context) {
	var req OrderRequestMeta
	var res OrderMultiResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	keys, rels, vals := make([]string, len(req.Conditions)), make([]string, len(req.Conditions)), make([]interface{}, len(req.Conditions))
	for index, condition := range req.Conditions {
		keys[index] = condition.Key
		rels[index] = condition.Rel
		vals[index] = condition.Val
	}
	userT, _ := c.Get("userType")
	if userT.(string) == "worker" {
		keys = append(keys, "assignedToes")
		rels = append(rels, "has")
		userI, _ := c.Get("user")
		vals = append(vals, userI.(string))
	}

	orders, err, errCode := op.GetOrdersByFilter(keys, rels, vals)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	adminStatusConvert := [4]int{1, 2, 0, 3}
	switch userT.(string) {
	case "worker":
		sort.Slice(orders, func(i, j int) bool {
			switch {
			case orders[i].Status != orders[j].Status:
				return orders[i].Status < orders[j].Status
			case orders[i].Priority != orders[j].Priority:
				return orders[i].Priority > orders[j].Priority
			case orders[i].PlannedStartTime != orders[j].PlannedStartTime:
				return orders[i].PlannedStartTime > orders[j].PlannedStartTime
			default:
				return strings.Compare(orders[i].ID, orders[j].ID) > 0
			}
		})
	case "admin":
		sort.Slice(orders, func(i, j int) bool {
			switch {
			case adminStatusConvert[orders[i].Status] != adminStatusConvert[orders[j].Status]:
				return adminStatusConvert[orders[i].Status] < adminStatusConvert[orders[j].Status]
			case orders[i].Priority != orders[j].Priority:
				return orders[i].Priority > orders[j].Priority
			case orders[i].MarkedFinishedTime != orders[j].MarkedFinishedTime:
				return orders[i].MarkedFinishedTime > orders[j].MarkedFinishedTime
			default:
				return strings.Compare(orders[i].ID, orders[j].ID) > 0
			}
		})
	default:
		sort.Slice(orders, func(i, j int) bool {
			return strings.Compare(orders[i].ID, orders[j].ID) > 0
		})
	}

	if util.IsEmpty(req.Page) {
		req.Page = 1
	}
	if util.IsEmpty(req.Limit) {
		req.Limit = 50
	}
	start := (req.Page - 1) * req.Limit
	end := min(start+req.Limit, len(orders))
	start = min(start, end)
	for _, order := range orders[start:end] {
		assignedMachineIDs := util.DBToStringList(order.AssignedMachinesDB)
		assignedMachineNames := make([]string, len(assignedMachineIDs))
		for i, machineID := range assignedMachineIDs {
			machineInfo, err, _ := op.GetMachineByID(machineID)
			if err != nil {
				_error.Print(err)
				continue
			}
			assignedMachineNames[i] = machineInfo.Name
		}
		res.Data = append(res.Data, struct {
			model.Order
			AssignedTo           []string `json:"assignedTo"`
			WorkPeriods          []string `json:"workPeriods"`
			AssignedMachines     []string `json:"assignedMachines"`
			AssignedMachineNames []string `json:"assignedMachineNames"`
		}{
			Order:                order,
			AssignedTo:           util.DBToStringList(order.AssignedToes),
			WorkPeriods:          util.DBToStringList(order.WorkPeriodsDB),
			AssignedMachines:     assignedMachineIDs,
			AssignedMachineNames: assignedMachineNames,
		})
	}

	res.Code = 0
	res.Msg = "Query Orders Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func OrderGetHandler(c *gin.Context) {
	var req OrderRequestMeta
	var res OrderMultiResponseMeta
	if err := c.ShouldBind(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userT, _ := c.Get("userType")
	userI, _ := c.Get("user")
	userType := userT.(string)
	userID := userI.(string)
	orderList := util.DBToStringList(req.OrderIDs)
	for _, order := range orderList {
		orderInfo, err, errCode := op.GetOrderByID(order)
		if err != nil {
			_error.Print(err)
			ErrorResponse(c, errCode, err.Error(), &res)
			return
		}

		if userType == "worker" && !util.ExistInDBList(orderInfo.AssignedToes, userID) {
			continue
		}
		assignedMachineIDs := util.DBToStringList(orderInfo.AssignedMachinesDB)
		assignedMachineNames := make([]string, len(assignedMachineIDs))
		for i, machineID := range assignedMachineIDs {
			machineInfo, err, _ := op.GetMachineByID(machineID)
			if err != nil {
				_error.Print(err)
				continue
			}
			assignedMachineNames[i] = machineInfo.Name
		}
		res.Data = append(res.Data, struct {
			model.Order
			AssignedTo           []string `json:"assignedTo"`
			WorkPeriods          []string `json:"workPeriods"`
			AssignedMachines     []string `json:"assignedMachines"`
			AssignedMachineNames []string `json:"assignedMachineNames"`
		}{
			Order:                orderInfo,
			AssignedTo:           util.DBToStringList(orderInfo.AssignedToes),
			WorkPeriods:          util.DBToStringList(orderInfo.WorkPeriodsDB),
			AssignedMachines:     assignedMachineIDs,
			AssignedMachineNames: assignedMachineNames,
		})
	}

	res.Code = 0
	res.Msg = "Get Orders Info Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func OrderGetTypeHandler(c *gin.Context) {
	var res OrderSingleResponseMeta

	if orderTypes, err, errCode := op.GetOrderTypes(); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	} else {
		res.Data.OrderTypes = orderTypes
	}

	res.Code = 0
	res.Msg = "Get Order Types Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func OrderUpdateTypeHandler(c *gin.Context) {
	var req OrderRequestMeta
	var res OrderSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	if err, errCode := op.UpdateOrderTypes(req.OrderTypes); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Update Order Types Successfully"
	res.Data.OrderTypes = req.OrderTypes

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}
