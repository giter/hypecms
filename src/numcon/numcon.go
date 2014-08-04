// numcon helps you to convert arbitrarily typed numeric data to the desired type.
package numcon

import (
	"fmt"
)

const errMsg = "numcon: can't convert to %v"

func bigInt(n interface{}) (int64, error) {
	switch val := n.(type) {
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return int64(val), nil
	default:
		return 0, fmt.Errorf(errMsg, "int64")
	}
	return 0, nil
}

func bigFloat(n interface{}) (float64, error) {
	switch val := n.(type) {
		case float32:
		return float64(val), nil
	case float64:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	default:
		return 0, fmt.Errorf(errMsg, "float64")
	}
	return 0, nil
}

func Int8(n interface{}) (int8, error) {
	val, err := bigInt(n)
	if err != nil {
		return 0, err
	}
	return int8(val), nil
}

func Int8P(n interface{}) int8 {
	val, err := Int8(n)
	if err != nil {
		panic(err)
	}
	return val
}

func Int16(n interface{}) (int16, error) {
	val, err := bigInt(n)
	if err != nil {
		return 0, err
	}
	return int16(val), nil
}

func Int16P(n interface{}) int16 {
	val, err := Int16(n)
	if err != nil {
		panic(err)
	}
	return val
}

func Int(n interface{}) (int, error) {
	val, err := bigInt(n)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func IntP(n interface{}) int {
	val, err := Int(n)
	if err != nil {
		panic(err)
	}
	return val
}

func Int32(n interface{}) (int32, error) {
	val, err := bigInt(n)
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}

func Int32P(n interface{}) int32 {
	val, err := Int32(n)
	if err != nil {
		panic(err)
	}
	return val
}

func Int64(n interface{}) (int64, error) {
	val, err := bigInt(n)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

func Int64P(n interface{}) int64 {
	val, err := Int64(n)
	if err != nil {
		panic(err)
	}
	return val
}

func Float32(n interface{}) (float32, error) {
	val, err := bigFloat(n)
	if err != nil {
		return 0, err
	}
	return float32(val), nil
}

func Float32P(n interface{}) float32 {
	val, err := Float32(n)
	if err != nil {
		panic(err)
	}
	return val
}

func Float64(n interface{}) (float64, error) {
	val, err := bigFloat(n)
	if err != nil {
		return 0, err
	}
	return float64(val), nil
}

func Float64P(n interface{}) float64 {
	val, err := Float64(n)
	if err != nil {
		panic(err)
	}
	return val
}
