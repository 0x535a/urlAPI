package cmd

import (
	"fmt"
	"log"
	"urlAPI/internal/bootstrap"
	"urlAPI/internal/database"
)

func admin(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing admin command")
	}
	if err := bootstrap.Init(); err != nil {
		return err
	}
	defer bootstrap.Release()
	switch args[0] {
	case "repwd":
		resetPassword()
		log.Println("Password has been reset to 123456, please change it ASAP.")
	case "clear":
		database.ClearTask()
		log.Println("Task Cleared")
	case "logout":
		database.ClearSession()
		log.Println("Session Restored")
	case "clear_ip_restriction":
		clearIPRestrict()
		log.Println("Cleared IP restriction")
	default:
		return fmt.Errorf("unknown admin command %q", args[0])
	}
	return nil
}

func resetPassword() {
	settings := database.SettingsStore.Get()
	settings.Security.DashboardPasswordHash = "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"
	if err := database.SaveAppSettings(settings); err != nil {
		log.Fatal(err)
	}
	database.ClearSession()
}

func clearIPRestrict() {
	settings := database.SettingsStore.Get()
	settings.Security.DashboardAllowedIPs = []string{"*"}
	if err := database.SaveAppSettings(settings); err != nil {
		log.Fatal(err)
	}
}
