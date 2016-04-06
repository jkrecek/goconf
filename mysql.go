package goconf

import (
	"fmt"
	"strings"

	"database/sql"

	_ "github.com/jkrecek/mysql"
	"sync"
)

type MysqlConnection struct {
	Address  string
	User     string
	Pass     string
	Database string

	instance *sql.DB
	lock sync.Mutex
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

func (mc *MysqlConnection) CreateConnection() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", mc.getConnectionString())
	if err != nil {
		return
	}

	err = db.Ping()
	return
}

func (mc *MysqlConnection) GetInstance() (*sql.DB, error) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	var err error
	if mc.instance == nil {
		mc.instance, err = mc.CreateConnection()
	}

	return mc.instance, err
}


func (mc *MysqlConnection) RuntimeTest() (err error, fatal bool) {
	conn, err := mc.CreateConnection()
	if err != nil {
		fatal = true
		return
	}

	conn.Close()
	return
}
