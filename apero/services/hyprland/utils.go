package hyprland

import (
	"reflect"
	"strconv"
)

func callHandler(handler interface{}, target string, rawValues []string) {
	method := reflect.ValueOf(handler).MethodByName(target)
	methodType := method.Type()

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
