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
		if res, ok := from.(string); ok {
			return res
		}
		return cast.ToString(from)
	case int:
		if res, ok := from.(int); ok {
			return res
		}
		return cast.ToInt(from)
	case int64:
		if res, ok := from.(int64); ok {
			return res
		}
		return cast.ToInt64(from)
	case int32:
		if res, ok := from.(int32); ok {
			return res
		}
		return cast.ToInt32(from)
	case int16:
		if res, ok := from.(int16); ok {
			return res
		}
		return cast.ToInt16(from)
	case int8:
		if res, ok := from.(int8); ok {
			return res
		}
		return cast.ToInt8(from)
	case uint:
		if res, ok := from.(uint); ok {
			return res
		}
		return cast.ToUint(from)
	case uint64:
		if res, ok := from.(uint16); ok {
			return res
		}
		return cast.ToUint16(from)
	case uint32:
		if res, ok := from.(uint32); ok {
			return res
		}
		return cast.ToUint32(from)
	case uint16:
		if res, ok := from.(uint16); ok {
			return res
		}
		return cast.ToUint16(from)
	case uint8:
		if res, ok := from.(uint8); ok {
			return res
		}
		return cast.ToUint8(from)
	case float64:
		if res, ok := from.(float64); ok {
			return res
		}
		return cast.ToFloat64(from)
	case float32:
		if res, ok := from.(float32); ok {
			return res
		}
		return cast.ToFloat32(from)
	case bool:
		if res, ok := from.(bool); ok {
			return res
		}
		return cast.ToBool(from)
	case time.Time:
		if res, ok := from.(time.Time); ok {
			return res
		}
		return cast.ToTime(from)
	case time.Duration:
		if res, ok := from.(time.Duration); ok {
			return res
		}
		return cast.ToDuration(from)
	default:
		return from
	}
}
