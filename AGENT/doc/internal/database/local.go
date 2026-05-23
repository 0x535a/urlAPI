package database

import (
	"zhongxin/internal/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var localDB SQLiteAdapter

func SetLocalDB(d *gorm.DB) *SQLiteAdapter {
	localDB = SQLiteAdapter{
		db: d,
	}
	return &localDB
}

func GetLocalDB() *SQLiteAdapter {
	return &localDB
}

func (s *SQLiteAdapter) Init() error {
	return errors.WithStack(s.db.AutoMigrate(
		new(model.Machine),
		new(model.Notification),
		new(model.Order),
		new(model.ColorAndQuantity),
		new(model.Token),
		new(model.User),
		new(model.Version),
		new(model.WorkPeriod),
		new(model.MachinePeriod),
		new(model.KV),
	))
}

func (s *SQLiteAdapter) Close() error {
	if s.db == nil {
		return nil
	}
	sqlDB, err := s.db.DB()
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(sqlDB.Close())
}
