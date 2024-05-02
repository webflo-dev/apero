package hyprland

import (
	"log"
	"reflect"
	"strconv"
)

func callHandler(handler interface{}, target string, rawValues []string) {
	value := reflect.ValueOf(handler)
	if !value.IsValid() {
		log.Println("handler not found")
		return
	}
	log.Printf("value: %+v", value)

	method := reflect.ValueOf(handler).MethodByName(target)
	if !method.IsValid() {
		log.Println("method not found")
		return
	}
	log.Printf("method(%s) %+v", target, method)

	methodType := method.Type()
	if !method.IsValid() {
		log.Println("method type invalid")
		return
	}
	log.Printf("type: %+v", methodType)

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

	method.Call(in)
}
