package handles

import (
	"net/http"
	"reflect"
	"zhongxin/util"

	"github.com/gin-gonic/gin"
)

func APINoRoute(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"error": "No Such API",
	})
}

func ErrorResponse(c *gin.Context, errorCode int, msg string, resp interface{}) {
	v := reflect.ValueOf(resp)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		util.Log.Errorln("需要传入结构体指针")
		return
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		util.Log.Errorln("指针需要指向结构体")
		return
	}

	codeField := v.FieldByName("Code")
	if codeField.IsValid() && codeField.CanSet() && codeField.Kind() == reflect.Int {
		codeField.SetInt(int64(errorCode))
	}

	msgField := v.FieldByName("Msg")
	if msgField.IsValid() && msgField.CanSet() && msgField.Kind() == reflect.String {
		msgField.SetString(msg)
	}

	c.JSON(http.StatusBadRequest, resp)
}
