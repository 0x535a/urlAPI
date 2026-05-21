package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"zhongxin/internal/bootstrap"
	"zhongxin/internal/database"
	"zhongxin/internal/op"
	"zhongxin/util"
)

var pid = -1
var pidFile string

func Init() {
	bootstrap.InitConfig()
	bootstrap.InitOp()
	if err := bootstrap.InitDB(); err != nil {
		util.Log.Error(err)
		Release()
		os.Exit(1)
	}

	util.Log.Info("initialized all module")
}

func Release() {
	database.Close()
	op.Close()

	util.Log.Infoln("Server gracefully stopped")
}

func initDaemon() {
	ex, err := os.Executable()
	if err != nil {
		util.Log.Error(err)
		return
	}
	exPath := filepath.Dir(ex)
	_ = os.MkdirAll(filepath.Join(exPath, "daemon"), 0700)
	pidFile = filepath.Join(exPath, "daemon/pid")
	if util.FileExist(pidFile) {
		bytes, err := os.ReadFile(pidFile)
		if err != nil {
			util.Log.Fatal("failed to read pid file", err)
		}
		id, err := strconv.Atoi(string(bytes))
		if err != nil {
			util.Log.Fatal("failed to parse pid data", err)
		}
		pid = id
	}
}
