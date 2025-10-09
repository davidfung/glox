package object

type ObjType int

const (
	OBJ_STRING ObjType = iota
)

type Obj struct {
	Type_ ObjType
	Val   any
}

type ObjString string

// Given a segment of a string, return a object whose value is a string.
// It can be used to convert a scanner token into a string object.
func CopyString(s *string, start int, length int) Obj {
	d := ObjString((*s)[start : start+length])
	o := Obj{Type_: OBJ_STRING, Val: d}
	return o
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
