package request

import (
	"urlAPI/database"
)

type DB struct {
	Task    database.Task
	Repo    database.Repo
	Session database.Session
}
