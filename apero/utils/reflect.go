package utils

// These functions are not used. They are here for reference only about how to use reflect module

import (
	"log"
	"reflect"
	"strconv"
)

func implementsMethod(obj any, methodName string, checkEmbed bool) bool {

	objType := reflect.TypeOf(obj)

	// Check if MyStruct implements the method with the given name
	for i := 0; i < objType.NumMethod(); i++ {
		method := objType.Method(i)
		if method.Name == methodName {
			return true
		}
	}

	if checkEmbed {
		// Check if MyStruct embeds a type that implements the method with the given name
		if objType.Kind() == reflect.Pointer {
			objType = objType.Elem()
			for i := 0; i < objType.NumField(); i++ {
				field := objType.Field(i)
				if field.Anonymous {
					for j := 0; j < field.Type.NumMethod(); j++ {
						method := field.Type.Method(j)
						if method.Name == methodName {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func callMethod(handler interface{}, target interface{}, rawValues []string) {
	method := reflect.ValueOf(target)

	// reflect.TypeOf(target).MethodByName("")
	// method.MethodByName("WorkspaceV2")

	methodType := method.Type()
	if !method.IsValid() {
		log.Println("method type invalid")
		return
	}
	// log.Printf("type: %+v", methodType)

	in := make([]reflect.Value, methodType.NumIn())

	for i := 0; i < method.Type().NumIn(); i++ {
		value := rawValues[i]
		t := methodType.In(i)
		switch t.Kind() {
		case reflect.Bool:
			in[i] = reflect.ValueOf(value == "1")
			break
		case reflect.Int:
			intValue, _ := strconv.Atoi(value)
			in[i] = reflect.ValueOf(intValue)
		default:
			in[i] = reflect.ValueOf(value)
			break
		}
	}

	// reflect.ValueOf(handler).Call(in)
	method.Call(in)
}
