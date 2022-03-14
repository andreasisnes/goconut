package goconut

import (
	"reflect"
	"time"

	"github.com/spf13/cast"
)

func CastAndTryAssignValue(from interface{}, to interface{}) interface{} {
	casted := castValue(to, from)

	tValue := reflect.ValueOf(to)
	fValue := reflect.ValueOf(casted)

	if tValue.Kind() == reflect.Ptr {
		tValue = tValue.Elem()
		if tValue.Kind() == fValue.Kind() {
			if tValue.Elem().CanSet() {
				tValue.Set(fValue)
			}
		}
	}

	return casted
}

func castValue(to interface{}, from interface{}) interface{} {
	switch to.(type) {
	case string:
		return cast.ToString(from)
	case int:
		return cast.ToInt(from)
	case int64:
		return cast.ToInt64(from)
	case int32:
		return cast.ToInt32(from)
	case int16:
		return cast.ToInt16(from)
	case int8:
		return cast.ToInt8(from)
	case uint:
		return cast.ToUint(from)
	case uint64:
		return cast.ToUint16(from)
	case uint32:
		return cast.ToUint32(from)
	case uint16:
		return cast.ToUint16(from)
	case uint8:
		return cast.ToUint8(from)
	case float64:
		return cast.ToFloat64(from)
	case float32:
		return cast.ToFloat32(from)
	case bool:
		return cast.ToBool(from)
	case time.Time:
		return cast.ToTime(from)
	case time.Duration:
		return cast.ToDuration(from)
	default:
		return from
	}
}
