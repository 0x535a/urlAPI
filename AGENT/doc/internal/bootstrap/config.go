package bootstrap

import (
	"encoding/json"
	"os"
	"path/filepath"
	"zhongxin/cmd/flags"
	"zhongxin/internal/conf"
	"zhongxin/util"
)

func InitConfig() {
	configPath := filepath.Join(flags.DataDir, "config.json")
	conf.Conf = conf.DefaultConfig()
	if !util.FileExist(configPath) {
		util.Log.Infof("config file not exists, setting default config file")
		_, err := util.CreateFile(configPath)
		if err != nil {
			util.Log.Errorf("init config file failed: %s", err)
			return
		}
	} else {
		configBytes, err := os.ReadFile(configPath)
		if err != nil {
			util.Log.Errorf("load config file failed: %s", err)
			return
		}
		if err := json.Unmarshal(configBytes, conf.Conf); err != nil {
			util.Log.Errorf("init config file failed: %s", err)
			return
		}
	}
	if err := util.JsonToFile(configPath, conf.Conf); err != nil {
		util.Log.Errorf("write config file failed: %s", err)
	}

	util.Log.Info("initialized config module")
}
