package AlgoeDB

func isNumber(value interface{}) bool {
	switch value.(type) {
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, complex64, complex128:
		return true
	default:
		return false
	}
}

func isString(value interface{}) bool {
	switch value.(type) {
	case string:
		return true
	default:
		return false
	}
}
