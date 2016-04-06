package goconf

import (
	"testing"
)

type testConfig struct {
	ExampleBool   bool             `json:"example_bool" default:"true"`
	ExampleString string           `json:"example_str" default:"empty"`
	ExampleInt    int              `json:"example_int" default:"50"`
	Mysql         *MysqlConnection `json:"mysql" testable:"true"`
	MysqlStruct   MysqlConnection  `json:"mysql_struct" tastable:"true"`
	MysqlNoTest   MysqlConnection  `json:"mysql_no_test" testable:"false"`
}

// TODO must somehow use mock mysql connection
var exampleConfig = []byte(`
{
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

	if configObj.ExampleBool == false {
		t.Fatal("ExampleBool has unexpected value")
	}

	if configObj.ExampleString != "empty" {
		t.Fatal("ExampleString has unexpected value")
	}

	if configObj.ExampleInt != 50 {
		t.Fatal("ExampleInt has unexpected value")
	}

}
