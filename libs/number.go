package libs

import "reflect"

func IsNumber(v interface{}) bool {
		var (
				getValue = reflect.ValueOf(v)
				kindName = getValue.Kind()
		)
		switch kindName {
		case reflect.Int:
				fallthrough
		case reflect.Int8:
				fallthrough
		case reflect.Int16:
				fallthrough
		case reflect.Int32:
				fallthrough
		case reflect.Int64:
				fallthrough
		case reflect.Uint:
				fallthrough
		case reflect.Uint8:
				fallthrough
		case reflect.Uint16:
				fallthrough
		case reflect.Uint32:
				fallthrough
		case reflect.Uint64:
				fallthrough
		case reflect.Uintptr:
				fallthrough
		case reflect.Float32:
				fallthrough
		case reflect.Float64:
				return true
		}
		return false
}
