package AlgoeDB

import (
	"reflect"
)

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

func IsFunction(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Func
}

func MoreThan(value float64) QueryFunc {
	return func(target interface{}) bool {
		if IsNumber(target) {
			target := reflect.ValueOf(target).Interface().(float64)
			return target > value
		}

		return false
	}
}

func MoreThanOrEqual(value float64) QueryFunc {
	return func(target interface{}) bool {
		if IsNumber(target) {
			target := reflect.ValueOf(target).Interface().(float64)
			return target >= value
		}

		return false
	}
}

func LessThan(value float64) QueryFunc {
	return func(target interface{}) bool {
		if IsNumber(target) {
			target := reflect.ValueOf(target).Interface().(float64)
			return target < value
		}

		return false
	}
}

func LessThanOrEqual(value float64) QueryFunc {
	return func(target interface{}) bool {
		if IsNumber(target) {
			target := reflect.ValueOf(target).Interface().(float64)
			return target < value
		}

		return false
	}
}

func Between(lowValue float64, highValue float64) QueryFunc {
	return func(target interface{}) bool {
		if IsNumber(target) {
			target := reflect.ValueOf(target).Interface().(float64)
			return target < highValue && target > lowValue
		}

		return false
	}
}

func BetweenOrEqual(lowValue float64, highValue float64) QueryFunc {
	return func(target interface{}) bool {
		if IsNumber(target) {
			target := reflect.ValueOf(target).Interface().(float64)
			return target <= highValue && target >= lowValue
		}

		return false
	}
}

func Exists() QueryFunc {
	return func(target interface{}) bool {
		return target != nil
	}
}

func And(queryValues ...QueryFunc) QueryFunc {
	return func(target interface{}) bool {
		for _, queryValue := range queryValues {
			if !matchValues(queryValue, target) {
				return false
			}
		}
		return true
	}
}

func Or(queryValues ...QueryFunc) QueryFunc {
	return func(target interface{}) bool {
		for _, queryValue := range queryValues {
			if matchValues(queryValue, target) {
				return true
			}
		}
		return false
	}
}

func Not(queryValue QueryFunc) QueryFunc {
	return func(target interface{}) bool {
		return !matchValues(queryValue, target)
	}
}
