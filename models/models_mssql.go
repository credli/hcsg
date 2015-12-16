package models

import (
	//ODBC driver
	_ "github.com/denisenkom/go-mssqldb"
)

func init() {
	EnableMSSQL = true
}
