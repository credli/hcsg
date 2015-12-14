package models

import (
	//ODBC driver
	_ "github.com/alexbrainman/odbc"
)

func init() {
	EnableODBC = true
}
