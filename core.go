package goconf

import (
	"os"
	"io/ioutil"
	"log"
	"encoding/json"
	"reflect"
	"strconv"
)

const (
	TestMethodName = "RuntimeTest"
)

func LoadConfigFromFile(fileName string, configObject interface{}) {
	configDataBytes, err := loadFile(fileName)
	if err != nil {
		log.Fatalf("Configuration file read error: %s", err)
	}

	LoadConfig(configDataBytes, configObject)
}

func LoadConfig(configBytes []byte, configObject interface{}) {
	err := json.Unmarshal(configBytes, &configObject)
	if err != nil {
		log.Fatalf("Configuration file marshal error: %s", err)
	}
	configBytes = nil

	RuntimeTest(configObject)
}

func loadFile(fileName string) (fileBytes []byte, err error) {
	configFilePath := getValidFullPath(fileName, getExecutePath())

	_, err = os.Stat(configFilePath)
	if err != nil {
		return
	}

	fileBytes, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		return
	}

	return
}


func RuntimeTest(configObject interface{}) {
	v := reflect.Indirect(reflect.ValueOf(configObject))

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		if testable, err := strconv.ParseBool(typeField.Tag.Get("testable")); err != nil || testable == false {
			continue
		}

		method := v.Field(i).MethodByName(TestMethodName)
		if !method.IsValid() {
			log.Printf("Warning: Config property `%s` marked as testable, but no method `%s` was found.\n", typeField.Name, TestMethodName)
			continue
		}

		methodType, _ := v.Field(i).Type().MethodByName(TestMethodName)
		methodTypeType := methodType.Type
		if methodTypeType.NumOut() != 2 || methodTypeType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() || methodTypeType.Out(1) != reflect.TypeOf((*bool)(nil)).Elem() {
			log.Printf("Warning: Method `%s` in property `%s` has invalid parameters, should be (bool, error).", TestMethodName, typeField.Name)
			continue
		}

		testOutput := method.Call([]reflect.Value{})

		var err error = nil
		if ie := testOutput[0].Interface(); ie != nil {
			err = ie.(error)
		}

		fatal := testOutput[1].Interface().(bool)

		if err != nil {
			if fatal {
				log.Fatalln(err)
			} else {
				log.Println(err)
			}
		}

	}
}