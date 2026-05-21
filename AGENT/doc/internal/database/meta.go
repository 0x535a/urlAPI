package database

import (
	"zhongxin/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SQLiteAdapter struct {
	db *gorm.DB
}

type MSSQLAdapter struct {
	db *gorm.DB
}

type Localer interface {
	// DB itself
	Init() error
	Close() error

	// Token
	GetToken(token string) (model.Token, bool, error)
	NewToken(id, token, typ string, expireSeconds int64) error
	CleanExpiredTokens() error

	// User
	NewUser(user model.User) error
	DeleteUserByID(id string) error
	DeleteUserByName(name string) error
	UpdateUser(user model.User) error
	GetUserByID(id string) (model.User, bool, error)
	GetUserByWXID(wxid string) (model.User, bool, error)
	GetUserByName(name string) (model.User, bool, error)
	GetUsersByType(userType string) ([]model.User, bool, error)
	GetUserByFilter(cls []clause.Eq) ([]model.User, bool, error)

	// Machine
	NewMachine(machine model.Machine) error
	GetMachineByID(id string) (model.Machine, bool, error)
	GetAllMachinePasID() ([]string, error)
	UpdateMachine(user model.Machine) error
	DeleteMachineByID(id string) error

	// Order
	CreateOrder(order model.Order) error
	DeleteOrderByID(id string) error
	UpdateOrder(order model.Order) (string, string, error)
	UpdateOrderWithoutVersion(order model.Order) error
	GetOrderByID(id string) (model.Order, bool, error)
	GetOrdersByFilter(cls []clause.Expression) ([]model.Order, bool, error)

	//WorkPeriod
	CreateWorkPeriod(wp model.WorkPeriod) error
	UpdateWorkPeriod(wp model.WorkPeriod) error
	GetWorkPeriodByID(id string) (model.WorkPeriod, bool, error)
	GetWorkPeriodsByClauses(cls []clause.Expression) ([]model.WorkPeriod, bool, error)
	DeleteWorkPeriodByID(id string) error

	//Notification
	GetNotificationsByClauses(cls []clause.Expression) ([]model.Notification, bool, error)
	GetNotificationByID(id string) (model.Notification, bool, error)
	GetNotificationByUserID(userID string) ([]model.Notification, bool, error)
	DeleteNotificationByID(id string) error
	DeleteNotificationByUserID(userID string) error
	NewNotification(n model.Notification) error
	UpdateNotification(n model.Notification) error

	//Version
	NewVersion(userID string, contentID string, from string, to string, vType string) error
	GetVersionByFilter(cls []clause.Expression) ([]model.Version, error)
	DeleteVersionByContentID(id string) error

	// KV
	SetKV(key string, val string) error
	DeleteKVByKey(key string) error
	GetAllKVs() ([]model.KV, error)
	GetKVByKey(key string) (model.KV, bool, error)
	GetKVsByFilter(cls []clause.Expression) ([]model.KV, error)
}

type Remoter interface {
	Init() error
	Close() error

	GetMachineONLogByFilter(cls []clause.Expression) ([]model.MachineLog, bool, error)
	GetAllMachinePasID() ([]model.MachinePasID, bool, error)
}
