package database

import "gorm.io/gorm"

var remoteDB MSSQLAdapter

func SetRemoteDB(d *gorm.DB) {
	remoteDB = MSSQLAdapter{
		db: d,
	}
}

func GetRemoteDB() *MSSQLAdapter {
	return &remoteDB
}

func (m *MSSQLAdapter) Close() error {
	if m.db == nil {
		return nil
	}
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (m *MSSQLAdapter) Init() error {
	return nil
}
