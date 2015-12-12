package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/mattn/go-sqlite3"

	"github.com/credli/hcsg/settings"
)

var (
	db    *sql.DB
	DbCfg struct {
		Type, Host, Name, User, Passwd, Path string
	}
	Connected     bool
	EnableSQLite3 bool
	EnableODBC    bool
)

func GetDb() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	var err error
	db, err = getDB()
	if err != nil {
		log.Panicf("ERROR: Could not initialize database connection: %v", err)
		return nil, err
	}
	return db, nil
}

func LoadConfigs() {
	sec := settings.Cfg.Section("database")
	DbCfg.Type = sec.Key("DB_TYPE").String()
	switch DbCfg.Type {
	case "sqlite3":
		settings.UseSQLite3 = true
	case "odbc":
		settings.UseODBC = true
	}
	DbCfg.Host = sec.Key("HOST").String()
	DbCfg.Name = sec.Key("NAME").String()
	DbCfg.User = sec.Key("USER").String()
	if len(DbCfg.Passwd) == 0 {
		DbCfg.Passwd = sec.Key("PASSWD").String()
	}
	DbCfg.Path = sec.Key("PATH").MustString("data/hcsg.db")
}

func getDB() (*sql.DB, error) {
	connStr := ""
	switch DbCfg.Type {
	case "sqlite3":
		if !EnableSQLite3 {
			return nil, fmt.Errorf("Unknown database type: %s", DbCfg.Type)
		}
		if err := os.MkdirAll(path.Dir(DbCfg.Path), os.ModePerm); err != nil {
			return nil, fmt.Errorf("Fail to create directories: %v", err)
		}
		connStr = "file:" + DbCfg.Path + "?cache=shared&mode=rwc"
	case "odbc":
		connStr = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", DbCfg.Host, DbCfg.User, DbCfg.Passwd, DbCfg.Name)
	default:
		return nil, fmt.Errorf("Unsupported database type: %s", DbCfg.Type)
	}
	return sql.Open(DbCfg.Type, connStr)
}

func Ping() error {
	return db.Ping()
}
