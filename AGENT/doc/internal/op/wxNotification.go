package op

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"zhongxin/internal/client"
	"zhongxin/internal/conf"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/util"

	"github.com/pkg/errors"
)

type OrderConfirmNotification struct {
	OrderID    StringValue `json:"character_string8"`
	TimeString StringValue `json:"time4"`
	Status     StringValue `json:"phrase6"`
	Content    StringValue `json:"thing3"`
}

type OrderCompleteNotification struct {
	OrderID    StringValue `json:"character_string5"`
	TimeString StringValue `json:"time4"`
	Content    StringValue `json:"thing1"`
	Remark     StringValue `json:"thing11"`
}

type OrderNewNotification struct {
	OrderID            StringValue `json:"character_string13"`
	AssignTimeString   StringValue `json:"date4"`
	PlaceTimeString    StringValue `json:"date3"`
	CompleteTimeString StringValue `json:"time32"`
	Content            StringValue `json:"thing6"`
}

type OrderUpdateNotification struct {
	OrderID      StringValue `json:"character_string12"`
	OrderName    StringValue `json:"thing3"`
	TimeString   StringValue `json:"date4"`
	UpdateTypeCN StringValue `json:"thing5"`
}

type OrderDeleteNotification struct {
	OrderID    StringValue `json:"character_string1"`
	Content    StringValue `json:"thing7"`
	TimeString StringValue `json:"time13"`
}

type StringValue struct {
	Value string `json:"value"`
}

type WxTokenResponseMeta struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type WxNotificationRequestMeta struct {
	ToUser           string      `json:"touser"`
	TemplateID       string      `json:"template_id"`
	Page             string      `json:"page"`
	MiniProgramState string      `json:"miniprogram_state"`
	Lang             string      `json:"lang"`
	Data             interface{} `json:"data"`
}

func NewNotificationOrderNew(order model.Order, userIDs []string) (error, int) {
	// Attention: haven't comfirmed the content & extraInfo
	for _, toUser := range userIDs {
		noti := OrderNewNotification{
			OrderID: StringValue{
				Value: order.FactoryID[:min(31, len(order.FactoryID)-1)],
			},
			AssignTimeString: StringValue{
				Value: util.UnixTimeToString(order.AssignedTime),
			},
			PlaceTimeString: StringValue{
				Value: util.UnixTimeToString(order.PlannedStartTime),
			},
			CompleteTimeString: StringValue{
				Value: util.UnixTimeToString(order.PlannedEndTime),
			},
			Content: StringValue{
				Value: order.Name,
			},
		}
		notiDB := model.Notification{
			UserID:    toUser,
			TypeCode:  model.OrderNew,
			ContentID: order.ID,
			ExtraInfo: order.Name,
		}
		userInfo, _, err := db.GetUserByID(toUser)
		if err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
		go SendNotificationTask(conf.OrderNewNotificationTemplateID, userInfo.WXID, fmt.Sprintf(conf.WxAppOrderDetailPageTemplate, order.ID), noti, notiDB)
	}
	return nil, 0
}

func NewNotificationOrderConfirm(order model.Order, receiverIDs []string) (error, int) {
	// Attention: haven't confirmed the content & extraInfo
	for _, toUser := range receiverIDs {
		noti := OrderConfirmNotification{
			OrderID: StringValue{
				Value: util.Substr(order.FactoryID, 30),
			},
			TimeString: StringValue{
				Value: util.UnixTimeToString(util.TimeNow().Unix()),
			},
			Status: StringValue{
				Value: "待审核",
			},
			Content: StringValue{
				Value: util.Substr(order.Name, 20),
			},
		}
		notiDB := model.Notification{
			UserID:    toUser,
			TypeCode:  model.OrderConfirm,
			ContentID: order.ID,
			ExtraInfo: order.Name,
		}
		userInfo, _, err := db.GetUserByID(toUser)
		if err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
		go SendNotificationTask(conf.OrderConfirmNotificationTemplateID, userInfo.WXID, conf.WxAppIndex, noti, notiDB)
	}
	return nil, 0
}

func NewNotificationOrderComplete(order model.Order) (error, int) {
	// Attention: haven't confirmed the content & extraInfo
	receivers, _, err := db.GetUsersByType("admin")
	if err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	for _, worker := range util.DBToStringList(order.AssignedToes) {
		userInfo, _, err := db.GetUserByID(worker)
		if err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
		receivers = append(receivers, userInfo)
	}

	//for _, receiver := range receivers {
	//noti := OrderCompleteNotification{
	//	OrderID: StringValue{
	//		Value: order.FactoryID[:min(31, len(order.FactoryID)-1)],
	//	},
	//	TimeString: StringValue{
	//		Value: util.UnixTimeToString(util.TimeNow().Unix()),
	//	},
	//	Content: StringValue{
	//		Value: order.Name,
	//	},
	//	Remark: StringValue{
	//		Value: "等待审核",
	//	},
	//}
	//notiDB := model.Notification{
	//	UserID:    receiver.ID,
	//	TypeCode:  model.OrderComplete,
	//	ContentID: order.ID,
	//	ExtraInfo: order.Name,
	//}
	//go SendNotificationTask(conf.OrderCompleteNotificationTemplateID, receiver.WXID, fmt.Sprintf(conf.WxAppOrderDetailPageTemplate, order.ID), noti, notiDB)
	//}
	return nil, 0
}

func NewNotificationOrderUpdate(order model.Order, updateInfo string, receiverIDs []string) (error, int) {
	// Attention: haven't confirmed the content & extraInfo
	for _, toUser := range receiverIDs {
		noti := OrderUpdateNotification{
			OrderID: StringValue{
				Value: util.Substr(order.FactoryID, 30),
			},
			OrderName: StringValue{
				Value: util.Substr(order.Name, 20),
			},
			UpdateTypeCN: StringValue{
				Value: updateInfo,
			},
			TimeString: StringValue{
				Value: util.UnixTimeToString(util.TimeNow().Unix()),
			},
		}
		notiDB := model.Notification{
			UserID:    toUser,
			TypeCode:  model.OrderUpdate,
			ContentID: order.ID,
			ExtraInfo: updateInfo,
		}
		userInfo, _, err := db.GetUserByID(toUser)
		if err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
		go SendNotificationTask(conf.OrderUpdateNotificationTemplateID, userInfo.WXID, fmt.Sprintf(conf.WxAppOrderDetailPageTemplate, order.ID), noti, notiDB)
	}
	return nil, 0
}

func NewNotificationOrderDelete(order model.Order, userIDs []string) (error, int) {
	// Attention: haven't confirmed the content & extraInfo
	//for _, toUser := range userIDs {
	//	noti := OrderDeleteNotification{
	//		OrderID: StringValue{
	//			Value: order.FactoryID[:min(31, len(order.FactoryID)-1)],
	//		},
	//		Content: StringValue{
	//			Value: order.Name,
	//		},
	//		TimeString: StringValue{
	//			Value: util.UnixTimeToString(util.TimeNow().Unix()),
	//		},
	//	}
	//	notiDB := model.Notification{
	//		UserID:    toUser,
	//		TypeCode:  model.OrderDelete,
	//		ContentID: order.ID,
	//		ExtraInfo: order.Name,
	//	}
	//	userInfo, _, err := db.GetUserByID(toUser)
	//if err != nil {
	//	return errors.WithStack(err), _error.ConvertGormError(err)
	//}
	//go SendNotificationTask(conf.OrderDeleteNotificationTemplateID, userInfo.WXID, fmt.Sprintf(conf.WxAppOrderDetailPageTemplate, order.ID), noti, notiDB)
	//}
	return nil, 0
}

func UpdateWxToken() (error, int) {
	if conf.Conf.WXConfig.AppID == "" || conf.Conf.WXConfig.AppSecret == "" {
		return errors.New("WxAPP Config Not implemented"), _error.OpConfigNotImplemented
	}

	if util.TimeNow().Unix() < conf.WxTokenExpireTime-600 && conf.WxToken != "" {
		return nil, 0
	}

	u, err := url.Parse(conf.WxTokenURL)
	if err != nil {
		return errors.WithStack(err), _error.OpURLParseError
	}
	q := u.Query()
	q.Set("appid", conf.Conf.WXConfig.AppID)
	q.Set("secret", conf.Conf.WXConfig.AppSecret)
	q.Set("grant_type", "client_credential")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		return errors.WithStack(err), _error.OpWxClientError
	}
	defer resp.Body.Close()

	var result WxTokenResponseMeta
	data, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &result); err != nil {
		return errors.WithStack(err), _error.OpJSONError
	}
	if resp.StatusCode != http.StatusOK || result.ErrCode != 0 {
		return errors.New(result.ErrMsg), _error.OpWxResponseError
	}

	conf.WxToken = result.AccessToken
	conf.WxTokenExpireTime = util.TimeNow().Unix() + int64(result.ExpiresIn)
	return nil, 0
}

func SendNotificationTask(templateID, wxid, page string, data interface{}, notiDB model.Notification) {
	for i := 1; i <= 5; i++ {
		if err, _ := SendNotification(templateID, wxid, page, data); err != nil {
			_error.Print(err)
			time.Sleep(1 * time.Second)
		} else {
			notiDB.Time = util.TimeNow().Unix()
			if err := db.NewNotification(notiDB); err != nil {
				_error.Print(err)
			}
			break
		}
	}
}

func SendNotification(templateID, userWXID, page string, data interface{}) (error, int) {
	var err error
	var errCode int
	flags := false
	for i := 1; i <= 5; i++ {
		if err, errCode = UpdateWxToken(); err != nil {
			_error.Print(err)
			time.Sleep(time.Second)
		} else {
			flags = true
			break
		}
	}
	if !flags {
		return errors.New("update wx token failed"), errCode
	}
	if userWXID == "" {
		return errors.New("user wxid is empty"), _error.OpUserWXIDNotFound
	}

	request := WxNotificationRequestMeta{
		ToUser:           userWXID,
		TemplateID:       templateID,
		Page:             page,
		MiniProgramState: conf.Conf.Schema.AppState,
		Lang:             "zh_CN",
		Data:             data,
	}
	jsonPayload, _ := json.Marshal(request)
	u, err := url.Parse(conf.WxNotificationURL)
	if err != nil {
		return errors.WithStack(err), _error.OpURLParseError
	}
	q := u.Query()
	q.Set("access_token", conf.WxToken)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		return errors.WithStack(err), _error.OpWxClientError
	}
	defer resp.Body.Close()

	var result WxTokenResponseMeta
	dataResp, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(dataResp, &result); err != nil {
		return errors.WithStack(err), _error.OpJSONError
	}
	if resp.StatusCode != http.StatusOK || result.ErrCode != 0 {
		return errors.New(result.ErrMsg), _error.OpWxResponseError
	}

	return nil, 0
}
