package logger

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type fieldType uint8

const (
	skipType fieldType = iota
	boolType
	floatType
	intType
	uintType
	stringType
	objectMarshalerType
)

// Field 结构化日志的基本单位
type Field struct {
	key     string
	ftype   fieldType
	integer int64
	str     string
	inf     interface{}
}

func (f Field) String() string {
	val := "null"
	switch f.ftype {
	case intType:
		val = strconv.FormatInt(f.integer, 10)
	case uintType:
		val = strconv.FormatUint(uint64(f.integer), 10)
	case boolType:
		val = strconv.FormatBool(f.integer == 1)
	case floatType:
		val = strconv.FormatFloat(math.Float64frombits(uint64(f.integer)), 'f', -1, 32)
	case stringType:
		b, _ := json.Marshal(f.str)
		val = string(b)
	case objectMarshalerType:
		b, _ := json.Marshal(f.inf)
		val = string(b)
	default:
		return ""
	}
	return fmt.Sprintf("\"%s\":%s", f.key, val)
}

func Skip() Field {
	return Field{ftype: skipType}
}

func ErrorField(err error) Field {
	return NamedError("error", err)
}

func NamedError(key string, err error) Field {
	if err == nil {
		return Skip()
	}
	return Field{key: key, ftype: stringType, str: err.Error()}
}

func Int64(key string, i int64) Field {
	return Field{key: key, ftype: intType, integer: i}
}

func Uint64(key string, i uint64) Field {
	return Field{key: key, ftype: uintType, integer: int64(i)}
}

func Int32(key string, i int32) Field {
	return Int64(key, int64(i))
}

func Uint32(key string, i uint32) Field {
	return Uint64(key, uint64(i))
}

func Int(key string, i int) Field {
	return Int64(key, int64(i))
}

func Uint(key string, i uint) Field {
	return Uint64(key, uint64(i))
}

func Float64(key string, f float64) Field {
	return Field{key: key, ftype: floatType, integer: int64(math.Float64bits(f))}
}

func Float32(key string, f float32) Field {
	return Float64(key, float64(f))
}

func String(key string, s string) Field {
	return Field{key: key, ftype: stringType, str: s}
}

func Object(key string, val interface{}) Field {
	return Field{key: key, ftype: objectMarshalerType, inf: val}
}
