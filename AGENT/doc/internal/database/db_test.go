package database

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
	"time"
	"zhongxin/util"
)

func testInit() error {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		"sa",
		"SqlServer2022",
		"192.168.6.185", 1433, "lisDb")
	count := 0
	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open(sqlserver.Open(dsn))
		if err != nil {
			count++
			if count >= 5 {
				return errors.Errorf("failed to connect remote database:%s", err.Error())
			}
			util.Log.Warnf("failed to connect remote database, retrying... (%d/5)", count)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	SetRemoteDB(db)
	return nil
}

func TestMSSQLAdapter_GetMachineONLogByFilter(t *testing.T) {
	if err := testInit(); err != nil {
		t.Fatal(err)
	}
	db := GetRemoteDB()
	logs, _, err := db.GetMachineONLogByFilter([]clause.Expression{
		clause.Eq{
			Column: "MacID",
			Value:  "2",
		},
		clause.Gte{
			Column: "BeginTime",
			Value:  time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, log := range logs {
		fmt.Printf("%+v\n", log)
	}
}
