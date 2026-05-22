package bootstrap

import (
	"urlAPI/internal/database"
	"urlAPI/internal/op"
)

func Init() error {
	if err := database.Init(); err != nil {
		return err
	}
	return op.Init()
}

func Release() {
	database.Disconnect()
}
