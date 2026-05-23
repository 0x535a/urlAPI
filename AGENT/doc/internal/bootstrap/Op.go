package bootstrap

import (
	"zhongxin/internal/client"
	"zhongxin/internal/op"
	"zhongxin/util"
)

func InitOp() {
	op.Init()
	client.Init()

	util.Log.Info("initialized op module")
}
