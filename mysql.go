package goconf

import (
	"fmt"
	"strings"

	"database/sql"

	_ "github.com/jkrecek/mysql"
)

type MysqlConnection struct {
	Address  string
	User     string
	Pass     string
	Database string
}

func (mc *MysqlConnection) getConnectionString() string {
	var conType string
	if strings.HasPrefix(mc.Address, "/") {
		conType = "unix"
	} else {
		conType = "tcp"
	}

	return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true", mc.User, mc.Pass, conType, mc.Address, mc.Database)
}

func (mc *MysqlConnection) GetDatabaseConnection() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", mc.getConnectionString())
	if err != nil {
		return
	}

	err = db.Ping()
	return
}

func (mc *MysqlConnection) RuntimeTest() (err error, fatal bool) {
	conn, err := mc.GetDatabaseConnection()
	if err != nil {
		fatal = true
		return
	}

	conn.Close()
	return
}