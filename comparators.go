package AlgoeDB

import (
	"errors"
	"regexp"
)

func MoreThan(value float64) QueryFunc {
	return func(target interface{}) bool {
		number, err := getNumber(target)
		if err != nil {
			return false
		}

		return number > value
	}
}

func MoreThanOrEqual(value float64) QueryFunc {
	return func(target interface{}) bool {
		number, err := getNumber(target)
		if err != nil {
			return false
		}

		return number >= value
	}
}

func LessThan(value float64) QueryFunc {
	return func(target interface{}) bool {
		number, err := getNumber(target)
		if err != nil {
			return false
		}

		return number < value
	}
}

func LessThanOrEqual(value float64) QueryFunc {
	return func(target interface{}) bool {
		number, err := getNumber(target)
		if err != nil {
			return false
		}

		return number <= value
	}
}

func Between(lowValue float64, highValue float64) QueryFunc {
	return func(target interface{}) bool {
		number, err := getNumber(target)
		if err != nil {
			return false
		}

		return number < highValue && number > lowValue
	}
}

func BetweenOrEqual(lowValue float64, highValue float64) QueryFunc {
	return func(target interface{}) bool {
		number, err := getNumber(target)
		if err != nil {
			return false
		}

		return number <= highValue && number >= lowValue
	}
}

func Exists() QueryFunc {
	return func(target interface{}) bool {
		return target != nil
	}
}

func Matches(pattern string) QueryFunc {
	return func(target interface{}) bool {
		switch x := target.(type) {
		case string:
			ok, _ := regexp.MatchString(pattern, x)
			if ok {
				return true
			}
			return false
		default:
			return false
		}
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

func getNumber(value interface{}) (float64, error) {
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
		return 0, errors.New("could not convert number")
	}
}
