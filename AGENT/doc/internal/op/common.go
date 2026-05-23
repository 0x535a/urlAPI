package op

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"
	"zhongxin/internal/database"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/util"
)

var db *database.SQLiteAdapter
var remoteDB *database.MSSQLAdapter
var scheduler gocron.Scheduler
var job gocron.Job

func Init() {
	db = database.GetLocalDB()
	remoteDB = database.GetRemoteDB()
	InitTask()
}

func Close() {
	if err := scheduler.Shutdown(); err != nil {
		util.Log.Errorf("%s", err.Error())
		return
	}

	util.Log.Info("released op module")
}

func GetToken(token string) (model.Token, error, int) {
	userDB, _, err := db.GetToken(token)
	if err != nil {
		return model.Token{}, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return userDB, nil, 0
}
