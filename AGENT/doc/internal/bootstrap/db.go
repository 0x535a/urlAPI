package bootstrap

import (
	"fmt"
	stdlog "log"
	"net/url"
	"strings"
	"time"
	"zhongxin/cmd/flags"
	"zhongxin/internal/conf"
	"zhongxin/internal/database"
	"zhongxin/util"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() error {
	logLevel := logger.Silent
	newLogger := logger.New(
		stdlog.New(log.StandardLogger().Out, "\r\n", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	if err := initLocalDB(newLogger); err != nil {
		return errors.Errorf("init local database failed: %v", err)
	}
	if !flags.NoRemote {
		if err := initRemoteDB(newLogger); err != nil {
			return errors.Errorf("init remote database failed: %v", err)
		}
	}

	util.Log.Info("initialized db module")
	return nil
}

func initLocalDB(logger logger.Interface) error {
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.Conf.LocalDatabase.TablePrefix,
		},
		Logger: logger,
	}

	dbConf := conf.Conf.LocalDatabase
	if !(strings.HasSuffix(dbConf.DBFile, ".database") && len(dbConf.DBFile) > 3) {
		return errors.Errorf("invalid database file name: %s", dbConf.DBFile)
	}
	// _journal=WAL设置为写前日志，提升并发性能，特别适合读多写少
	// _vacuum=incremental增量清理，使空间回收更加平缓，避免性能波动
	dB, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental&?_foreign_keys=on",
		dbConf.DBFile)), gormConfig)

	if err != nil {
		return errors.Errorf("failed to connect database:%s", err.Error())
	}

	if err := database.SetLocalDB(dB).Init(); err != nil {
		return errors.Errorf("failed to initialize database:%s", err.Error())
	}

	return nil
}

func initRemoteDB(logger logger.Interface) error {
	dbConf := conf.Conf.RemoteDatabase
	if dbConf.Host == "" || dbConf.Port == 0 || dbConf.User == "" || dbConf.DBName == "" || dbConf.Password == "" {
		return errors.New("incomplete remote database configuration")
	}
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
		url.QueryEscape(dbConf.User),
		url.QueryEscape(dbConf.Password),
		dbConf.Host,
		dbConf.Port,
		dbConf.DBName,
	)
	count := 0
	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: conf.Conf.LocalDatabase.TablePrefix,
			},
			Logger: logger,
		})
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
	database.SetRemoteDB(db)
	return nil
}
