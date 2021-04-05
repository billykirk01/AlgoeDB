package AlgoeDB

import (
	"errors"
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

func GetNumber(value interface{}) (float64, error) {
	switch x := value.(type) {
	case uint8:
		return float64(x), nil
	case int8:
		return float64(x), nil
	case uint16:
		return float64(x), nil
	case int16:
		return float64(x), nil
	case uint32:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case uint64:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case int:
		return float64(x), nil
	case float32:
		return float64(x), nil
	case float64:
		return float64(x), nil
	default:
		return 0, errors.New("Could not convert number")
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
		number, err := GetNumber(target)
		if err != nil {
			return false
		}

		return number > value
	}
}

func MoreThanOrEqual(value float64) QueryFunc {
	return func(target interface{}) bool {
		number, err := GetNumber(target)
		if err != nil {
			return false
		}

		return number >= value
	}
}

func LessThan(value float64) QueryFunc {
	return func(target interface{}) bool {
		number, err := GetNumber(target)
		if err != nil {
			return false
		}

		return number < value
	}
}

func LessThanOrEqual(value float64) QueryFunc {
	return func(target interface{}) bool {
		number, err := GetNumber(target)
		if err != nil {
			return false
		}

		return number <= value
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
