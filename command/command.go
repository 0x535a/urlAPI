package command

import (
	"log"
	"urlAPI/database"
)

func Arg(args []string) {
	log.Println("The default password is 123456")
	for index, arg := range args {
		if index == 0 {
			continue
		}
		switch arg {
		case "port":
			Port = args[index+1]
		case "repwd":
			repwd()
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
		}
	}
}

func repwd() {
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
