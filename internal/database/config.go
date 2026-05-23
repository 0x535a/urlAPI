package database

import (
	"urlAPI/internal/model"
)

var (
	dbPath    = "assets/database.db"
	PromptMap = map[string]int{
		"laugh":    0,
		"poem":     1,
		"sentence": 2,
	}
	RepoMap    = make(map[string][]string)
	SessionMap = make(map[string]model.Session)
)

type Repo = model.Repo
type Session = model.Session
type Task = model.Task
type AppSetting = model.AppSetting
type Provider = model.Provider
type ServiceConfig = model.ServiceConfig
type Prompt = model.Prompt
type ConfigListItem = model.ConfigListItem
type DBList = model.DBList
