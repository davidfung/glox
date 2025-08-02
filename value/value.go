package value

import (
	"fmt"

	"github.com/davidfung/glox/object"
)

type ValueType int

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
	VAL_OBJ
)

type Value struct {
	type_ ValueType // promote to export field?
	val   any       // promote to export field?
}

func BOOL_VAL(b bool) Value {
	return Value{type_: VAL_BOOL, val: b}
}

func NIL_VAL() Value {
	return Value{type_: VAL_NIL, val: nil}
}

func NUMBER_VAL(n float64) Value {
	return Value{type_: VAL_NUMBER, val: n}
}

func OBJ_VAL(obj object.Obj) Value {
	return Value{type_: VAL_OBJ, val: obj}
}

func IS_BOOL(v Value) bool {
	return v.type_ == VAL_BOOL
}

func IS_NIL(v Value) bool {
	return v.type_ == VAL_NIL
}

func IS_NUMBER(v Value) bool {
	return v.type_ == VAL_NUMBER
}

func IS_OBJ(v Value) bool {
	return v.type_ == VAL_OBJ
}

func IS_STRING(v Value) bool {
	return IsObjType(v, object.OBJ_STRING)
}

func IsObjType(val Value, type_ object.ObjType) bool {
	return IS_OBJ(val) && AS_OBJ(val).Type_ == type_
}

func AS_OBJ(v Value) object.Obj {
	obj, ok := v.val.(object.Obj)
	if !ok {
		panic("Error: AS_OBJ() expects an object value")
	}
	return obj
}

func AS_STRING(v Value) string {
	s, ok := v.val.(string)
	if !ok {
		panic("Error: AS_STRING() expects a string object")
	}
	return s
}

func AS_BOOL(v Value) bool {
	b, ok := v.val.(bool)
	if !ok {
		panic("Error: AS_BOOL() expects a boolean value")
	}
	return b
}

func AS_NUMBER(v Value) float64 {
	n, ok := v.val.(float64)
	if !ok {
		panic("Error: AS_NUMBER() expects an int value")
	}
	return n
}

func OBJ_TYPE(val Value) object.ObjType {
	return AS_OBJ(val).Type_
}

type ValueArray struct {
	Values []Value
}

func InitValueArray(array *ValueArray) {
	array.Values = nil
}

func WriteValueArray(array *ValueArray, value Value) {
	array.Values = append(array.Values, value)
}

func FreeValueArrary(array *ValueArray) {
	InitValueArray(array)
}

func PrintValue(value Value) {
	switch value.type_ {
	case VAL_BOOL:
		if AS_BOOL(value) {
			fmt.Printf("true")
		} else {
			fmt.Printf("false")
		}
	case VAL_NIL:
		fmt.Printf("nil")
	case VAL_NUMBER:
		fmt.Printf("%g", AS_NUMBER(value))
	}
}

func ValuesEqual(a Value, b Value) bool {
	if a.type_ != b.type_ {
		return false
	}
	switch a.type_ {
	case VAL_BOOL:
		return AS_BOOL(a) == AS_BOOL(b)
	case VAL_NIL:
		return true
	case VAL_NUMBER:
		return AS_NUMBER(a) == AS_NUMBER(b)
	default:
		return false
	}
}
