package op

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"zhongxin/internal/client"
	"zhongxin/internal/conf"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/util"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type WxLoginResponseMeta struct {
	OpenID  string `json:"openid"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func UserGetWXIDByWxToken(wxToken string) (string, error, int) {
	//TobeImplemented
	if conf.Conf.WXConfig.AppID == "" || conf.Conf.WXConfig.AppSecret == "" {
		return "", errors.New("WxAPP Config Not implemented"), _error.OpConfigNotImplemented
	}

	u, err := url.Parse(conf.WxLoginURL)
	if err != nil {
		return "", errors.WithStack(err), _error.OpURLParseError
	}
	q := u.Query()
	q.Set("appid", conf.Conf.WXConfig.AppID)
	q.Set("secret", conf.Conf.WXConfig.AppSecret)
	q.Set("js_code", wxToken)
	q.Set("grant_type", "authorization_code")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		return "", errors.WithStack(err), _error.OpWxClientError
	}
	defer resp.Body.Close()

	var result WxLoginResponseMeta
	data, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &result); err != nil {
		return "", errors.WithStack(err), _error.OpJSONError
	}
	if resp.StatusCode != http.StatusOK || result.ErrCode != 0 {
		if result.ErrCode == 40029 {
			return "", errors.New("invalid code"), _error.OpWxTokenExpired
		}
		return "", errors.New(result.ErrMsg), _error.OpWxResponseError
	}

	return result.OpenID, nil, 0
}

func GetFullUserInfoByWXID(wxid string) (model.User, []string, error, int) {
	userDB, _, err := db.GetUserByWXID(wxid)
	if err != nil {
		return model.User{}, nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	workPeriods := util.DBToStringList(userDB.WorkPeriods)
	return userDB, workPeriods, nil, 0
}

func GetFullUserInfoByID(id string) (model.User, []string, error, int) {
	userDB, _, err := db.GetUserByID(id)
	if err != nil {
		return model.User{}, nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	workPeriods := util.DBToStringList(userDB.WorkPeriods)
	return userDB, workPeriods, nil, 0
}

func GetAllFullWorkerInfo() ([]model.User, [][]string, error, int) {
	userDB, _, err := db.GetUserByFilter([]clause.Eq{{
		Column: "type", Value: "worker",
	}})
	if err != nil {
		return nil, nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	var workPeriods [][]string
	for _, user := range userDB {
		workPeriods = append(workPeriods, util.DBToStringList(user.WorkPeriods))
	}
	return userDB, workPeriods, nil, 0
}

func NewUserLoginByID(id string) (string, int64, error, int) {
	token := util.NewToken(conf.TokenLen)
	expiredOn := util.TimeNow().AddDate(0, 0, 7).Unix()
	userDB, _, err := db.GetUserByID(id)
	if err != nil {
		return "", 0, errors.WithStack(err), _error.ConvertGormError(err)
	}
	if err := db.NewToken(id, token, userDB.Type, expiredOn); err != nil {
		return "", 0, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return token, expiredOn, nil, 0
}

func UserBindWXIDByNameAndPhone(name, phone, wxid string) (error, int) {
	userDB, ok, err := db.GetUserByFilter([]clause.Eq{
		clause.Eq{Column: "name", Value: name},
		clause.Eq{Column: "phone", Value: phone},
	})
	if err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	if !ok || len(userDB) == 0 {
		return errors.New("user not found"), _error.DBRecordNotFound
	}
	if userDB[0].WXID != "" &&
		userDB[0].Name != "wxAdmin" &&
		userDB[0].Name != "wxWorker" {
		return errors.New("user already bound"), _error.OpUserAlreadyBinded
	}
	userDB[0].WXID = wxid
	if err := db.UpdateUser(userDB[0]); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func GetUserByID(userID string) (model.User, error, int) {
	userDB, _, err := db.GetUserByID(userID)
	if err != nil {
		return model.User{}, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return userDB, nil, 0
}

func GetUserByName(name string) (model.User, error, int) {
	userDB, _, err := db.GetUserByName(name)
	if err != nil {
		return model.User{}, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return userDB, nil, 0
}

func UpdateUser(user model.User) (error, int) {
	if err := db.UpdateUser(user); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func CreateWorkerByNameAndPhone(name, phone string) (model.User, error, int) {
	user := model.User{
		ID:        util.NewUUID(),
		Name:      name,
		Phone:     phone,
		Type:      "worker",
		IsWorking: false,
	}
	if err := db.UpdateUser(user); err != nil {
		return model.User{}, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return user, nil, 0
}

func DeleteUserByID(id string) (error, int) {
	if err := db.DeleteUserByID(id); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func GetUserSalaryByIDAndTime(id string, startTime, endTime int64) (float64, model.User, []model.WorkPeriod, error, int) {
	user, _, err := db.GetUserByID(id)
	if err != nil {
		return 0, model.User{}, nil, errors.WithStack(err), _error.ConvertGormError(err)
	}

	workPeriods, _, err := db.GetWorkPeriodsByClauses([]clause.Expression{
		clause.Eq{Column: "workerID", Value: id},
		clause.Gte{Column: "startTime", Value: startTime},
		clause.Lte{Column: "startTime", Value: endTime},
		clause.Eq{Column: "isMachineDataMerged", Value: true},
	})

	if err != nil {
		return 0, model.User{}, nil, errors.WithStack(err), _error.ConvertGormError(err)
	}

	salary := float64(0)
	for _, workPeriod := range workPeriods {
		salary += float64(workPeriod.ValidTimeSeconds) / float64(3600) * workPeriod.UnitPriceYuan
	}
	return salary, user, workPeriods, nil, 0
}

func GetAllAdmins() ([]model.User, error, int) {
	admins, _, err := db.GetUserByFilter([]clause.Eq{
		clause.Eq{Column: "type", Value: "admin"},
	})
	if err != nil {
		return nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return admins, nil, 0
}
