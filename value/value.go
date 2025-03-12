package value

import "fmt"

type ValueType int

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
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

func NUMBER_VAL(n int) Value {
	return Value{type_: VAL_NUMBER, val: n}
}

func AS_BOOL(v Value) bool {
	b, ok := v.val.(bool)
	if !ok {
		panic("Error: AS_BOOL() expects a boolean value")
	}
	return b
}

func AS_NUMBER(v Value) int {
	n, ok := v.val.(int)
	if !ok {
		panic("Error: AS_BOOL() expects an int value")
	}
	return n
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
	fmt.Printf("%g", value)
}
