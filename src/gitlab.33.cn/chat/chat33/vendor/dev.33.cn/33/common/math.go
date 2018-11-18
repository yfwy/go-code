package common

import (
	"math"
	"reflect"
	"strconv"
)

func Max(first int64, args ...int64) int64 {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func Min(first int64, args ...int64) int64 {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func Round(f float64, n int) float64 {
	pow10n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10n)*pow10n) / pow10n
}

func RoundMin(f float64, n int) float64 {
	pow10n := math.Pow10(n)
	return math.Trunc(f*pow10n) / pow10n
}

func ToInt(v interface{}) int {
	return int(ToInt32(v))
}

func ToInt32(o interface{}) int32 {
	switch t := o.(type) {
	case int:
		return int32(t)
	case int32:
		return t
	case int64:
		return int32(t)
	case float64:
		return int32(t)
	case string:
		return StringToInt32(t)
	default:
		panic(reflect.TypeOf(t).String())
	}
}

func ToInt64(v interface{}) int64 {
	switch t := v.(type) {
	case int:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return t
	case float64:
		return int64(t)
	case string:
		return StringToInt64(t)
	default:
		panic(reflect.TypeOf(t).String())
	}
}

func ToUint64(v interface{}) uint64 {
	switch t := v.(type) {
	case int:
		return uint64(t)
	case int32:
		return uint64(t)
	case int64:
		return uint64(t)
	case uint64:
		return t
	case float32:
		return uint64(t)
	case float64:
		return uint64(t)
	case string:
		return StringToUint64(t)
	default:
		panic(reflect.TypeOf(t).String())
	}
}

func ToFloat32(v interface{}) float32 {
	switch t := v.(type) {
	case int:
		return float32(t)
	case int32:
		return float32(t)
	case int64:
		return float32(t)
	case float32:
		return t
	case float64:
		return float32(t)
	case string:
		return StringToFloat32(t)
	default:
		panic(reflect.TypeOf(t).String())
	}
}

func ToFloat64(v interface{}) float64 {
	switch t := v.(type) {
	case int:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case float32:
		return float64(t)
	case float64:
		return t
	case string:
		return StringToFloat64(t)
	default:
		panic(reflect.TypeOf(t).String())
	}
}

func StringToInt32(str string) int32 {
	if str == "" {
		return 0
	}
	temp, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(temp)
}

func StringToUint32(str string) uint32 {
	if str == "" {
		return 0
	}
	temp, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(temp)
}

func StringToInt64(str string) int64 {
	if str == "" {
		return 0
	}
	temp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return temp
}

func StringToUint64(str string) uint64 {
	if str == "" {
		return 0
	}
	temp, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return temp
}

func StringToFloat32(str string) float32 {
	if str == "" {
		return 0
	}
	temp, err := strconv.ParseFloat(str, 32)
	if err != nil {
		panic(err)
	}
	return float32(temp)
}

func StringToFloat64(str string) float64 {
	if str == "" {
		return 0
	}
	temp, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return temp
}

func StringToBool(str string) bool {
	if str == "" {
		return false
	}
	b, err := strconv.ParseBool(str)
	if nil != err {
		panic(err)
	}
	return b
}
