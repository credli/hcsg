package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/credli/hcsg/settings"
)

var (
	db    *sql.DB
	DbCfg struct {
		Type string
		// mssql
		Host   string
		Name   string
		User   string
		Passwd string

		// odbc
		Driver     string
		Server     string
		UID        string
		Pwd        string
		Database   string
		TdsVersion string
		Port       string

		// sqlite
		Path string
	}
	Connected     bool
	EnableSQLite3 bool
	EnableODBC    bool
	EnableMSSQL   bool
)

func GetDb() (*sql.DB, error) {
	if db != nil && Connected {
		return db, nil
	}
	var err error
	db, err = openDbConnection()
	if err != nil {
		log.Panicf("ERROR: Could not initialize database connection: %v", err)
		return nil, err
	}
	return db, nil
}

func LoadConfigs() {
	sec := settings.Cfg.Section("database")
	DbCfg.Type = sec.Key("DB_TYPE").MustString("sqlite3")
	switch DbCfg.Type {
	case "sqlite3":
		settings.UseSQLite3 = true

		DbCfg.Path = sec.Key("PATH").MustString("data/hcsg.db")
	case "mssql":
		settings.UseMSSQL = true

		DbCfg.Host = sec.Key("HOST").String()
		DbCfg.Name = sec.Key("NAME").String()
		DbCfg.User = sec.Key("USER").String()
		if len(DbCfg.Passwd) == 0 {
			DbCfg.Passwd = sec.Key("PASSWD").String()
		}
	case "odbc":
		settings.UseODBC = true

		DbCfg.Driver = sec.Key("DRIVER").String()
		DbCfg.Server = sec.Key("SERVER").String()
		DbCfg.Port = sec.Key("PORT").String()
		DbCfg.UID = sec.Key("UID").String()
		DbCfg.Pwd = sec.Key("PWD").String()
		DbCfg.Database = sec.Key("DATABASE").String()
		DbCfg.TdsVersion = sec.Key("TDS_VERSION").String()
	}

	openDbConnection() //initialize database on startup
}

func openDbConnection() (*sql.DB, error) {
	connStr := ""
	switch DbCfg.Type {
	case "sqlite3":
		if !EnableSQLite3 {
			return nil, fmt.Errorf("SQLite3 is not enabled: %s", DbCfg.Type)
		}
		if err := os.MkdirAll(path.Dir(DbCfg.Path), os.ModePerm); err != nil {
			return nil, fmt.Errorf("Fail to create directories: %v", err)
		}
		connStr = "file:" + DbCfg.Path + "?cache=shared&mode=rwc"
	case "mssql":
		if !EnableMSSQL {
			return nil, fmt.Errorf("MSSQL is not enabled: %s", DbCfg.Type)
		}
		connStr = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
			DbCfg.Host, DbCfg.User, DbCfg.Passwd, DbCfg.Name)
	case "odbc":
		if !EnableODBC {
			return nil, fmt.Errorf("ODBC is not enabled: %s", DbCfg.Type)
		}
		connStr = fmt.Sprintf("DRIVER=%s;SERVER=%s;UID=%s;PWD=%s;DATABASE=%s;TDS_Version=%s;Port=%s",
			DbCfg.Driver, DbCfg.Server, DbCfg.UID, DbCfg.Passwd, DbCfg.Database, DbCfg.TdsVersion, DbCfg.Port)
	default:
		return nil, fmt.Errorf("Unsupported database type: %s", DbCfg.Type)
	}

	if !settings.ProdMode {
		log.Println("Connection String: " + connStr)
	}

	var err error
	db, err = sql.Open(DbCfg.Type, connStr)
	if err != nil {
		return nil, err
	}
	Connected = true
	return db, nil
}

func Ping() error {
	return db.Ping()
}
