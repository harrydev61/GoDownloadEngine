package gormc

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tranTriDev61/GoDownloadEngine/component/gormc/dialets"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDBType int

const (
	GormDBTypeMySQL GormDBType = iota + 1
	GormDBTypePostgres
	GormDBTypeSQLite
	GormDBTypeMSSQL
	GormDBTypeNotSupported
)

type GormOpt struct {
	dsn                   string
	dbType                string
	maxOpenConnections    int
	maxIdleConnections    int
	maxConnectionIdleTime int
}

type gormDB struct {
	id     string
	prefix string
	logger core.Logger
	db     *gorm.DB
	*GormOpt
}

func NewGormDB(id, prefix string) *gormDB {
	return &gormDB{
		GormOpt: new(GormOpt),
		id:      id,
		prefix:  strings.TrimSpace(prefix),
	}
}

func (gdb *gormDB) ID() string {
	return gdb.id
}

func (gdb *gormDB) InitFlags() {
	prefix := gdb.prefix
	if gdb.prefix != "" {
		prefix += "_"
	}

	flag.StringVar(
		&gdb.dsn,
		fmt.Sprintf("%sdb_dsn", prefix),
		"",
		"Database dsn",
	)

	flag.StringVar(
		&gdb.dbType,
		fmt.Sprintf("%sdb_driver", prefix),
		"mysql",
		"Database driver (mysql, postgres, sqlite, mssql) - Default mysql",
	)

	flag.IntVar(
		&gdb.maxOpenConnections,
		fmt.Sprintf("%sdb_max_conn", prefix),
		30,
		"maximum number of open connections to the database - Default 30",
	)

	flag.IntVar(
		&gdb.maxIdleConnections,
		fmt.Sprintf("%sdb_max_ide_conn", prefix),
		10,
		"maximum number of database connections in the idle - Default 10",
	)

	flag.IntVar(
		&gdb.maxConnectionIdleTime,
		fmt.Sprintf("%sdb_max_conn_ide_time", prefix),
		3600,
		"maximum amount of time a connection may be idle in seconds - Default 3600",
	)
}

func (gdb *gormDB) isDisabled() bool {
	return gdb.dsn == ""
}

func (gdb *gormDB) Activate(_ core.ServiceContext) error {
	gdb.logger = core.GlobalLogger().GetLogger(gdb.id)

	dbType := getDBType(gdb.dbType)
	if dbType == GormDBTypeNotSupported {
		return errors.WithStack(errors.New("Database type not supported."))
	}
	gdb.logger.Info("Connecting to database:", gdb.dsn)
	var err error
	gdb.db, err = gdb.getDBConn(dbType)

	if err != nil {
		gdb.logger.Error("Cannot connect to database", err.Error())
		return err
	}

	return nil
}

func (gdb *gormDB) Stop() error {
	return nil
}

func (gdb *gormDB) GetDB() *gorm.DB {
	if gdb.logger.GetLevel() == "debug" || gdb.logger.GetLevel() == "trace" {
		return gdb.db.Session(&gorm.Session{NewDB: true}).Debug()
	}

	newSessionDB := gdb.db.Session(&gorm.Session{NewDB: true, Logger: gdb.db.Logger.LogMode(logger.Silent)})

	if db, err := newSessionDB.DB(); err == nil {
		db.SetMaxOpenConns(gdb.maxOpenConnections)
		db.SetMaxIdleConns(gdb.maxIdleConnections)
		db.SetConnMaxIdleTime(time.Second * time.Duration(gdb.maxConnectionIdleTime))
	}

	return newSessionDB
}

func getDBType(dbType string) GormDBType {
	switch strings.ToLower(dbType) {
	case "mysql":
		return GormDBTypeMySQL
	case "postgres":
		return GormDBTypePostgres
	case "sqlite":
		return GormDBTypeSQLite
	case "mssql":
		return GormDBTypeMSSQL
	}

	return GormDBTypeNotSupported
}

func (gdb *gormDB) getDBConn(t GormDBType) (dbConn *gorm.DB, err error) {
	switch t {
	case GormDBTypeMySQL:
		return dialets.MySqlDB(gdb.dsn)
	case GormDBTypePostgres:
		return dialets.PostgresDB(gdb.dsn)
	case GormDBTypeSQLite:
		return dialets.SQLiteDB(gdb.dsn)
	case GormDBTypeMSSQL:
		return dialets.MSSqlDB(gdb.dsn)
	}

	return nil, nil
}
