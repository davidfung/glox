package objval

// This package is to workaround the import cycle problem
// between object and value packages.  All code that depends
// on both should be placed in this package.

import (
	"fmt"

	"github.com/davidfung/glox/object"
	"github.com/davidfung/glox/value"
)

func BOOL_VAL(b bool) value.Value {
	return value.Value{Type_: value.VAL_BOOL, Val: b}
}

func NIL_VAL() value.Value {
	return value.Value{Type_: value.VAL_NIL, Val: nil}
}

func NUMBER_VAL(n float64) value.Value {
	return value.Value{Type_: value.VAL_NUMBER, Val: n}
}

func OBJ_VAL(obj object.Obj) value.Value {
	return value.Value{Type_: value.VAL_OBJ, Val: obj}
}

func IS_BOOL(v value.Value) bool {
	return v.Type_ == value.VAL_BOOL
}

func IS_NIL(v value.Value) bool {
	return v.Type_ == value.VAL_NIL
}

func IS_NUMBER(v value.Value) bool {
	return v.Type_ == value.VAL_NUMBER
}

func IS_OBJ(v value.Value) bool {
	return v.Type_ == value.VAL_OBJ
}

func IS_FUNCTION(v value.Value) bool {
	return IsObjType(v, object.OBJ_FUNCTION)
}

func IS_NATIVE(v value.Value) bool {
	return IsObjType(v, object.OBJ_NATIVE)
}

func IS_STRING(v value.Value) bool {
	return IsObjType(v, object.OBJ_STRING)
}

func IsObjType(val value.Value, type_ object.ObjType) bool {
	return IS_OBJ(val) && AS_OBJ(val).Type_ == type_
}

func AS_OBJ(v value.Value) object.Obj {
	obj, ok := v.Val.(object.Obj)
	if !ok {
		panic("Error: AS_OBJ() expects an object value.Value")
	}
	return obj
}

func AS_FUNCTION(v value.Value) object.ObjFunction {
	obj, ok := v.Val.(object.Obj)
	if !ok {
		panic("Error: AS_FUNCTION() expects an object in a value.Value")
	}
	objFunction, ok := obj.Val.(object.ObjFunction)
	if !ok {
		panic("Error: AS_FUNCTION() expects a function object")
	}
	return objFunction
}

func AS_NATIVE(v value.Value) object.NativeFn {
	obj, ok := v.Val.(object.Obj)
	if !ok {
		panic("Error: AS_NATIVE() expects an object in a value.Value")
	}
	if obj.Type_ != object.OBJ_NATIVE {
		panic("Error: AS_NATIVE() expects a native function object")
	}
	native := obj.Val.(object.NativeFn)
	return native
}

func AS_STRING(v value.Value) object.ObjString {
	obj, ok := v.Val.(object.Obj)
	if !ok {
		panic("Error: AS_STRING() expects an object in a value.Value")
	}
	strobj, ok := obj.Val.(object.ObjString)
	if !ok {
		panic("Error: AS_STRING() expects a string object")
	}
	return strobj
}

func AS_BOOL(v value.Value) bool {
	b, ok := v.Val.(bool)
	if !ok {
		panic("Error: AS_BOOL() expects a boolean value.Value")
	}
	return b
}

func AS_NUMBER(v value.Value) float64 {
	n, ok := v.Val.(float64)
	if !ok {
		panic("Error: AS_NUMBER() expects an int value.Value")
	}
	return n
}

func OBJ_TYPE(val value.Value) object.ObjType {
	return AS_OBJ(val).Type_
}

func PrintValue(val value.Value) {
	switch val.Type_ {
	case value.VAL_BOOL:
		if AS_BOOL(val) {
			fmt.Printf("true")
		} else {
			fmt.Printf("false")
		}
	case value.VAL_NIL:
		fmt.Printf("nil")
	case value.VAL_NUMBER:
		fmt.Printf("%g", AS_NUMBER(val))
	case value.VAL_OBJ:
		printObject(val)
	}
}

func ValuesEqual(a value.Value, b value.Value) bool {
	if a.Type_ != b.Type_ {
		return false
	}
	switch a.Type_ {
	case value.VAL_BOOL:
		return AS_BOOL(a) == AS_BOOL(b)
	case value.VAL_NIL:
		return true
	case value.VAL_NUMBER:
		return AS_NUMBER(a) == AS_NUMBER(b)
	case value.VAL_OBJ:
		s1 := AS_STRING(a)
		s2 := AS_STRING(b)
		return s1 == s2
	default:
		return false
	}
}

func printObject(val value.Value) {
	switch OBJ_TYPE(val) {
	case object.OBJ_FUNCTION:
		object.PrintFunction(AS_FUNCTION(val))
	case object.OBJ_NATIVE:
		fmt.Printf("<native fn>") // can we also print the native function name?
	case object.OBJ_STRING:
		fmt.Printf("%s", AS_STRING(val))
	}
}
