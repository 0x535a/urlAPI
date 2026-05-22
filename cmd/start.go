package cmd

import (
	"urlAPI/internal/bootstrap"
	"urlAPI/internal/server"
)

func start() error {
	if err := bootstrap.Init(); err != nil {
		return err
	}
	defer bootstrap.Release()
	return server.Run(Port)
}
