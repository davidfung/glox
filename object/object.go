package object

import (
	"fmt"

	"github.com/davidfung/glox/chunk"
)

type ObjType int

const (
	_ ObjType = iota
	OBJ_FUNCTION
	OBJ_STRING
)

type Obj struct {
	Type_ ObjType
	Val   any
}

type ObjFunction struct {
	Type_ ObjType
	arity int
	Chun  chunk.Chunk
	Name  ObjString
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
	fmt.Printf("<fn %s>\n", function.Name)
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

func NewFunction() *ObjFunction {
	fn := new(ObjFunction)
	fn.arity = 0 // actually not necessary in glox
	fn.Name = "" // actually not necessary in glox
	chunk.InitChunk(&fn.Chun)
	return fn
}
