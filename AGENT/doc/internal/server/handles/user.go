package handles

import (
	"net/http"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/internal/op"

	"github.com/gin-gonic/gin"
)

type UserRequestMeta struct {
	WxToken string `json:"wxToken"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
}

type UserResponseMeta struct {
	Code int `json:"code"`
	Data struct {
		ExpiresOn   int64  `json:"expiresOn"`
		BearerToken string `json:"bearerToken"`
		User        struct {
			model.User
			WorkingPeriod []string `json:"workingPeriod"`
		} `json:"user"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func UserLoginHandler(c *gin.Context) {
	var req UserRequestMeta
	var res UserResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	wxid, err, errCode := op.UserGetWXIDByWxToken(req.WxToken)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	userInfo, workPeriods, err, errCode := op.GetFullUserInfoByWXID(wxid)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	token, expiredOn, err, errCode := op.NewUserLoginByID(userInfo.ID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Login Success"
	res.Data.BearerToken = token
	res.Data.ExpiresOn = expiredOn
	res.Data.User.WorkingPeriod = workPeriods
	res.Data.User.User = userInfo
	c.JSON(http.StatusOK, res)
}

func UserBindHandler(c *gin.Context) {
	var req UserRequestMeta
	var res UserResponseMeta
	if err := c.ShouldBindJSON(&req); err != nil {
		_error.Print(err)
		ErrorResponse(c, _error.ServerInvalidParams, err.Error(), &res)
		return
	}

	wxid, err, errCode := op.UserGetWXIDByWxToken(req.WxToken)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	if err, errCode := op.UserBindWXIDByNameAndPhone(req.Name, req.Phone, wxid); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	userInfo, workPeriods, err, errCode := op.GetFullUserInfoByWXID(wxid)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	token, expiredOn, err, errCode := op.NewUserLoginByID(userInfo.ID)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &res)
		return
	}

	res.Code = 0
	res.Msg = "Bind Success"
	res.Data.BearerToken = token
	res.Data.ExpiresOn = expiredOn
	res.Data.User.WorkingPeriod = workPeriods
	res.Data.User.User = userInfo
	c.JSON(http.StatusOK, res)
}

func UserUnbindHandler(c *gin.Context) {
	name := c.Query("name")
	user, err, errCode := op.GetUserByName(name)
	if err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &UserResponseMeta{})
		return
	}
	user.WXID = ""
	if err, errCode := op.UpdateUser(user); err != nil {
		_error.Print(err)
		ErrorResponse(c, errCode, err.Error(), &UserResponseMeta{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Unbind Success"})
}
