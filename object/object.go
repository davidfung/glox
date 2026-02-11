package object

import (
	"fmt"

	"github.com/davidfung/glox/chunk"
	"github.com/davidfung/glox/value"
)

type ObjType int

const (
	_ ObjType = iota
	OBJ_FUNCTION
	OBJ_NATIVE
	OBJ_STRING
)

type Obj struct {
	Type_ ObjType
	Val   any
}

type ObjFunction struct {
	Type_ ObjType
	Arity int
	Chun  chunk.Chunk
	Name  ObjString
}

type NativeFn func(argCount int, args []value.Value) value.Value

type ObjNative struct {
	Type_    ObjType
	Function NativeFn
}

type ObjString string

// Given a segment of a string, return a object whose value is a string.
// It can be used to convert a scanner token into a string object.
func CopyString(s *string, start int, length int) Obj {
	d := ObjString((*s)[start : start+length])
	o := Obj{Type_: OBJ_STRING, Val: d}
	return o
}

func PrintFunction(function ObjFunction) {
	fmt.Printf("<fn %s>", function.Name)
}

// FNV-1a hash algorithm
func hashString(s string, length int) uint32 {
	var hash uint32 = 2166136261
	for i := 0; i < length; i++ {
		hash ^= uint32(s[i])
		hash *= 16777619
	}
	return hash
}

func NewFunction() ObjFunction {
	fn := new(ObjFunction)
	fn.Arity = 0 // actually not necessary in glox
	fn.Name = "" // actually not necessary in glox
	chunk.InitChunk(&fn.Chun)
	return *fn
}

func NewNative(function NativeFn) ObjNative {
	native := ObjNative{Type_: OBJ_NATIVE, Function: function}
	return native
}
