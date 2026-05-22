package database

import "gorm.io/gorm"

type SQLiteAdapter struct {
	db *gorm.DB
}

var localDB = &SQLiteAdapter{}

func SetLocalDB(db *gorm.DB) *SQLiteAdapter {
	localDB.db = db
	return localDB
}

func GetLocalDB() *SQLiteAdapter {
	return localDB
}

func (adapter *SQLiteAdapter) DB() *gorm.DB {
	return adapter.db
}
