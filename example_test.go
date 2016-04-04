package goconf

import "testing"

type testConfig struct {
	IsDebug bool `json:"is_debug"`
	Mysql   *MysqlConnection `json:"mysql" testable:"true"`
	MysqlNoTest *MysqlConnection `json:"mysql_no_test" testable:"false"`
}

// TODO must somehow use moc mysql connection
var exampleConfig = []byte(`
{
  "IsDebug"           : true,
  "Mysql"             : {
    "Address"   : "localhost:3306",
    "User"      : "root",
    "Pass"      : "usbw",
    "Database"  : ""
  }
}
`)

func TestConfig(t *testing.T) {
	configObj := &testConfig{}
	LoadConfig(exampleConfig, configObj)
}