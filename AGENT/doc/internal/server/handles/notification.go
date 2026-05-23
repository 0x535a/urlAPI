package handles

import (
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/internal/op"
	"zhongxin/util"

	"github.com/gin-gonic/gin"
)

type NotificationRequestMeta struct {
	UserID          string   `form:"userID"`
	NotificationIDs []string `json:"notificationIDs"`
}

type NotificationResponseMeta struct {
	Code int                  `json:"code"`
	Data []model.Notification `json:"data"`
	Msg  string               `json:"msg"`
}

func NotificationListHandler(c *gin.Context) {
	var req NotificationRequestMeta
	var res NotificationResponseMeta
	if err := c.ShouldBind(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	userI, _ := c.Get("user")
	userID := userI.(string)
	notifications, err, errCode := op.GetNotificationByUserID(userID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Success"
	res.Data = notifications

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func NotificationReadHandler(c *gin.Context) {
	var req NotificationRequestMeta
	var res NotificationResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	for _, notificationID := range req.NotificationIDs {
		if err, errCode := op.DeleteNotificationByID(notificationID); err != nil {
			_error.Print(err)
			ErrorResponse(c, errCode, err.Error(), &res)
			return
		}
	}

	res.Code = 0
	res.Msg = "Notifications marked as read"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}

func NotificationReadAllHandler(c *gin.Context) {
	var res NotificationResponseMeta

	userI, _ := c.Get("user")
	userID := userI.(string)

	if err, errCode := op.DeleteNotificationsByUserID(userID); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "All notifications marked as read"

	util.ReplaceNilWithZeroValue(&res)
	c.JSON(200, res)
}
