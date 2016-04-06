package goconf

import "testing"

type testConfig struct {
	IsDebug     bool             `json:"is_debug"`
	MysqlNoTest *MysqlConnection `json:"mysql_no_test" testable:"false"`
	Mysql         *MysqlConnection `json:"mysql" testable:"true"`
	MysqlStruct   MysqlConnection  `json:"mysql_struct" tastable:"true"`
}

// TODO must somehow use mock mysql connection
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
