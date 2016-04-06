package goconf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	LoadDefaults(configObject)

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

func LoadDefaults(configObject interface{}) {
	v := reflect.Indirect(reflect.ValueOf(configObject))

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		defaultField := typeField.Tag.Get("default")
		if len(defaultField) == 0 {
			continue
		}

		switch v.Field(i).Kind() {
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(defaultField)
			if err == nil {
				v.Field(i).SetBool(boolVal)
			}
		case reflect.Int:
			intVal, err := strconv.ParseInt(defaultField, 10, 0)
			if err == nil {
				v.Field(i).SetInt(intVal)
			}
		case reflect.String:
			v.Field(i).SetString(defaultField)
		}

	}
}

func RuntimeTest(configObject interface{}) {
	v := reflect.Indirect(reflect.ValueOf(configObject))

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		if testable, err := strconv.ParseBool(typeField.Tag.Get("testable")); err != nil || testable == false {
			continue
		}

		kindField := v.Field(i).Kind()
		if kindField != reflect.Struct && kindField != reflect.Ptr {
			log.Printf("Warning: Config property `%s` cannot be marked as testable, only `struct` and `ptr` allowed.\n", typeField.Name)
			continue
		}

		ptrField := v.Field(i)
		if ptrField.Kind() == reflect.Struct {
			ptrField = ptrField.Addr()
		}

		method := ptrField.MethodByName(TestMethodName)
		if !method.IsValid() {
			log.Printf("Warning: Config property `%s` marked as testable, but no method `%s` was found.\n", typeField.Name, TestMethodName)
			continue
		}

		methodType, _ := ptrField.Type().MethodByName(TestMethodName)
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
