package handles

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"reflect"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/internal/op"
	"zhongxin/util"
)

type VersionRequestMeta struct {
	model.Version
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

type VersionResponseMeta struct {
	Code int             `json:"code"`
	Data []model.Version `json:"data"`
	Msg  string          `json:"msg"`
}

func VersionQueryHandler(c *gin.Context) {
	var req VersionRequestMeta
	var res VersionResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	cls := make([]clause.Expression, 0)
	v := reflect.ValueOf(req.Version)
	t := reflect.TypeOf(req.Version)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		if util.IsEmpty(value) {
			continue
		}

		switch field.Name {
		case "StartTime":
			cls = append(cls, clause.Gte{Column: "time", Value: value})
		case "EndTime":
			cls = append(cls, clause.Lte{Column: "time", Value: value})
		default:
			cls = append(cls, clause.Eq{Column: field.Name, Value: value})
		}
	}

	versionList, err, errCode := op.GetVersionByFilter(cls)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Data = versionList
	res.Msg = "Query Version Success"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}
