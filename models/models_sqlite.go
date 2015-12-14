package models

import (
	//SQLite3 Driver
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	EnableSQLite3 = true
}
