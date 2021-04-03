package AlgoeDB

import "reflect"

func IsNumber(value interface{}) bool {
	switch value.(type) {
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, complex64, complex128:
		return true
	default:
		return false
	}
}

func IsString(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.String
}

func IsBoolean(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Bool
}

func IsNil(value interface{}) bool {
	return value == nil
}

type QueryFunc func(value int) bool

func IsFunction(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Func
}
