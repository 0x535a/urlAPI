package handles

import (
	"fmt"
	"reflect"
	"time"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/internal/op"
	"zhongxin/util"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type WorkerRequestMeta struct {
	WorkerID  string `json:"workerID"`
	WorkerIDs string `form:"workerIDs"`
	StartTime int64  `form:"startTime"`
	EndTime   int64  `form:"endTime"`
	TimeYear  int    `form:"timeYear"`
	TimeMonth int    `form:"timeMonth"`
	model.User
}

type WorkerSingleResponseMeta struct {
	Code int `json:"code"`
	Data struct {
		model.User
		WorkPeriod []string `json:"workPeriod"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type WorkerMultiResponseMeta struct {
	Code int `json:"code"`
	Data []struct {
		model.User
		WorkingPeriod []string `json:"workingPeriod"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type WorkerSalaryDetailResponseMeta struct {
	Code int `json:"code"`
	Data []struct {
		SalaryYuan float64            `json:"salaryYuan"`
		Details    []model.WorkPeriod `json:"details"`
		User       struct {
			model.User
			WorkingPeriod []string `json:"workingPeriod"`
		} `json:"user"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type WorkerSalaryOverviewItem struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	SalaryYuan   float64 `json:"salaryYuan"`
	WorkingTime  float64 `json:"workingTime"`
	StoppingTime float64 `json:"stoppingTime"`
	TotalTime    float64 `json:"totalTime"`
	Efficiency   float64 `json:"efficiency"`
}

type WorkerSalaryOverviewResponseMeta struct {
	Code      int                        `json:"code"`
	Data      []WorkerSalaryOverviewItem `json:"data"`
	TimeYear  int                        `json:"timeYear"`
	TimeMonth int                        `json:"timeMonth"`
	Msg       string                     `json:"msg"`
}

func WorkerSalaryDetailHandler(c *gin.Context) {
	var req WorkerRequestMeta
	var res WorkerSalaryDetailResponseMeta
	if err := c.ShouldBind(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userT, _ := c.Get("userType")
	userI, _ := c.Get("user")
	userType := userT.(string)
	userID := userI.(string)

	var userList []string
	if userType == "worker" {
		userList = []string{userID}
	} else {
		userList = util.DBToStringList(req.WorkerIDs)
	}

	for index, user := range userList {
		salary, userInfo, workPeriods, err, code := op.GetUserSalaryByIDAndTime(user, req.StartTime, req.EndTime)
		if err != nil {
			_error.Print(err)
			ErrorResponse(c, code, err.Error(), &res)
			return
		}

		res.Data = append(res.Data, struct {
			SalaryYuan float64            `json:"salaryYuan"`
			Details    []model.WorkPeriod `json:"details"`
			User       struct {
				model.User
				WorkingPeriod []string `json:"workingPeriod"`
			} `json:"user"`
		}{
			SalaryYuan: salary,
			User: struct {
				model.User
				WorkingPeriod []string `json:"workingPeriod"`
			}{
				User:          userInfo,
				WorkingPeriod: util.DBToStringList(userInfo.WorkPeriods),
			},
		})
		res.Data[index].Details = workPeriods
	}

	res.Code = 0
	res.Msg = "Get Worker Salary Detail Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func WorkerGetHandler(c *gin.Context) {
	var req WorkerRequestMeta
	var res WorkerMultiResponseMeta
	if err := c.ShouldBind(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userT, _ := c.Get("userType")
	userI, _ := c.Get("user")
	userType := userT.(string)
	userID := userI.(string)

	var userList []string
	if userType == "worker" {
		userList = []string{userID}
	} else {
		userList = util.DBToStringList(req.WorkerIDs)
	}

	for _, user := range userList {
		userInfo, workPeriods, err, code := op.GetFullUserInfoByID(user)
		if err != nil {
			_error.Print(err)
			ErrorResponse(c, code, err.Error(), &res)
			return
		}
		res.Data = append(res.Data, struct {
			model.User
			WorkingPeriod []string `json:"workingPeriod"`
		}{
			User:          userInfo,
			WorkingPeriod: workPeriods,
		})
	}

	res.Code = 0
	res.Msg = "Get Worker Info Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func WorkerListHandler(c *gin.Context) {
	var res WorkerMultiResponseMeta

	userInfos, workPeriods, err, code := op.GetAllFullWorkerInfo()
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, code, err.Error(), &res)
		return
	}

	for i, userInfo := range userInfos {
		res.Data = append(res.Data, struct {
			model.User
			WorkingPeriod []string `json:"workingPeriod"`
		}{
			User:          userInfo,
			WorkingPeriod: workPeriods[i],
		})
	}
	res.Code = 0
	res.Msg = "Get All Worker Info Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func WorkerDeleteHandler(c *gin.Context) {
	var req WorkerRequestMeta
	var res WorkerSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	if err, errCode := op.DeleteUserByID(req.WorkerID); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Delete Worker Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func WorkerUpdateHandler(c *gin.Context) {
	var req WorkerRequestMeta
	var res WorkerSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userInfo, err, errCode := op.GetUserByID(req.WorkerID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	v := reflect.ValueOf(req.User)
	t := reflect.TypeOf(req.User)
	uv := reflect.ValueOf(&userInfo).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if util.IsEmpty(value) {
			continue
		}

		uf := uv.FieldByName(field.Name)
		if !uf.CanSet() || !uf.IsValid() {
			continue
		}
		uf.Set(v.Field(i))
	}
	if err, errCode := op.UpdateUser(userInfo); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Update Worker Info Successfully"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func WorkerAddHandler(c *gin.Context) {
	var req WorkerRequestMeta
	var res WorkerSingleResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userInfo, err, errCode := op.CreateWorkerByNameAndPhone(req.Name, req.Phone)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Create Worker Successfully"
	res.Data.User = userInfo
	res.Data.WorkPeriod = []string{}

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func getWorkerSalaryOverviewData(c *gin.Context) ([]WorkerSalaryOverviewItem, int, int, error) {
	var req WorkerRequestMeta
	if err := c.ShouldBind(&req); err != nil {
		return nil, 0, 0, err
	}

	startTime := req.StartTime
	endTime := req.EndTime

	if req.TimeYear != 0 && req.TimeMonth != 0 {
		startTime = time.Date(req.TimeYear, time.Month(req.TimeMonth), 1, 8, 0, 0, 0, time.Local).Unix()
		endTime = time.Date(req.TimeYear, time.Month(req.TimeMonth)+1, 1, 7, 59, 59, 0, time.Local).Unix()
	}

	userList := util.DBToStringList(req.WorkerIDs)

	var data []WorkerSalaryOverviewItem
	var totWorkingTime, totTotalTime, totSalary float64
	for _, user := range userList {
		salary, userInfo, workPeriods, err, _ := op.GetUserSalaryByIDAndTime(user, startTime, endTime)
		if err != nil {
			return nil, 0, 0, err
		}

		var workingTime, totalTime float64
		for _, wp := range workPeriods {
			workingTime += float64(wp.ValidTimeSeconds)
			totalTime += float64(wp.EndTime - wp.StartTime)
		}

		stoppingTime := totalTime - workingTime
		efficiency := 0.0
		if totalTime > 0 {
			efficiency = workingTime / totalTime
		}

		if salary != salary {
			util.Log.Errorf("NaN detected in salary for user %s (%s)", userInfo.Name, userInfo.ID)
		}
		if efficiency != efficiency {
			util.Log.Errorf("NaN detected in efficiency for user %s (%s). WorkingTime: %f, TotalTime: %f", userInfo.Name, userInfo.ID, workingTime, totalTime)
		}

		data = append(data, WorkerSalaryOverviewItem{
			ID:           userInfo.ID,
			Name:         userInfo.Name,
			SalaryYuan:   salary,
			WorkingTime:  workingTime / 3600.0,
			StoppingTime: stoppingTime / 3600.0,
			TotalTime:    totalTime / 3600.0,
			Efficiency:   efficiency,
		})
		totWorkingTime += workingTime
		totTotalTime += totalTime
		totSalary += salary
	}

	if totTotalTime > 0 {
		totEfficiency := totWorkingTime / totTotalTime
		if totEfficiency != totEfficiency {
			util.Log.Errorf("NaN detected in total efficiency. totWorkingTime: %f, totTotalTime: %f", totWorkingTime, totTotalTime)
		}
		data = append(data, WorkerSalaryOverviewItem{
			ID:           "0",
			Name:         "总计",
			SalaryYuan:   totSalary,
			WorkingTime:  totWorkingTime / 3600.0,
			StoppingTime: (totTotalTime - totWorkingTime) / 3600.0,
			TotalTime:    totTotalTime / 3600.0,
			Efficiency:   totEfficiency,
		})
	}

	util.ReplaceNilWithZeroValue(&data)

	return data, req.TimeYear, req.TimeMonth, nil
}

func WorkerSalaryOverviewHandler(c *gin.Context) {
	data, year, month, err := getWorkerSalaryOverviewData(c)
	var res WorkerSalaryOverviewResponseMeta
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerOtherError, err.Error(), &res)
		return
	}

	res.Data = data
	res.TimeYear = year
	res.TimeMonth = month
	res.Code = 0
	res.Msg = "Get Worker Salary Overview Successfully"
	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func WorkerSalaryExcelHandler(c *gin.Context) {
	data, year, month, err := getWorkerSalaryOverviewData(c)
	if err != nil {
		var res WorkerSalaryOverviewResponseMeta
		_error.Print(err)
		ErrorResponse(c, _error.ServerOtherError, err.Error(), &res)
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			_error.Print(err)
		}
	}()

	// Set value of a cell.
	title := fmt.Sprintf("%d年%d月分红汇总表（导出于%s）", year, month, time.Now().Format("2006年1月2日15:04"))
	f.MergeCell("Sheet1", "A1", "F1")
	f.SetCellValue("Sheet1", "A1", title)
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err == nil {
		f.SetCellStyle("Sheet1", "A1", "F1", style)
	}

	f.SetCellValue("Sheet1", "A2", "员工姓名")
	f.SetCellValue("Sheet1", "B2", "合计分红")
	f.SetCellValue("Sheet1", "C2", "开机时间")
	f.SetCellValue("Sheet1", "D2", "关机时间")
	f.SetCellValue("Sheet1", "E2", "总时间")
	f.SetCellValue("Sheet1", "F2", "效率")

	for i, item := range data {
		row := i + 3
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), item.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), item.SalaryYuan)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), item.WorkingTime)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), item.StoppingTime)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), item.TotalTime)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), item.Efficiency)
	}

	c.Header("Content-Disposition", "attachment; filename=salary.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err := f.Write(c.Writer); err != nil {
		_error.Print(err)
	}
}
