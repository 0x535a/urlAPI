package database

import (
	"zhongxin/util"
)

func Close() {
	if err := localDB.Close(); err != nil {
		util.Log.Errorf("failed closing connection to database: %s", err.Error())
		return
	}

	if err := remoteDB.Close(); err != nil {
		util.Log.Errorf("failed closing connection to database: %s", err.Error())
		return
	}

	util.Log.Info("closed connection to database")
}
