package conf

import (
	"path/filepath"
	"zhongxin/cmd/flags"
)

type Time struct {
	Hour   int `json:"hour" env:"HOUR"`
	Minute int `json:"minute" env:"MINUTE"`
}

type LocalDatabase struct {
	TablePrefix string `json:"tablePrefix" env:"TABLE_PREFIX"`
	DBFile      string `json:"db_file" env:"DB_FILE"`
}

type RemoteDatabase struct {
	Host          string `json:"host" env:"DB_HOST"`
	Port          int    `json:"port" env:"DB_PORT"`
	User          string `json:"user" env:"DB_USER"`
	Password      string `json:"password" env:"DB_PASSWORD"`
	DBName        string `json:"db_name" env:"DB_NAME"`
	SyncTimeDay   Time   `json:"sync_time_day" envPrefix:"SYNC_TIME_DAY_"`
	SyncTimeNight Time   `json:"sync_time_night" envPrefix:"SYNC_TIME_NIGHT_"`
}

type WXConfig struct {
	AppID     string `json:"app_id" env:"APP_ID"`
	AppSecret string `json:"app_secret" env:"APP_SECRET"`
}

type Schema struct {
	Port     int    `json:"port" env:"PORT"`
	Listen   string `json:"listen" env:"LISTEN"`
	URL      string `json:"url" env:"URL"`
	AppState string `json:"app_state" env:"APP_STATE"`
}

type Config struct {
	LocalDatabase  LocalDatabase  `json:"local_database" envPrefix:"LocalDB_"`
	RemoteDatabase RemoteDatabase `json:"remote_database" envPrefix:"RemoteDB_"`
	WXConfig       WXConfig       `json:"wx_config" envPrefix:"WX_"`
	Schema         Schema         `json:"schema" envPrefix:"SCHEMA"`
}

func DefaultConfig() *Config {
	dbPath := filepath.Join(flags.DataDir, "local.database")
	return &Config{
		LocalDatabase: LocalDatabase{
			TablePrefix: "x_",
			DBFile:      dbPath,
		},
		Schema: Schema{
			Port:     2233,
			Listen:   "127.0.0.1",
			URL:      "https://mydomain.com",
			AppState: "developer",
		},
		RemoteDatabase: RemoteDatabase{
			SyncTimeDay: Time{
				Hour:   8,
				Minute: 10,
			},
			SyncTimeNight: Time{
				Hour:   20,
				Minute: 10,
			},
		},
	}
}
