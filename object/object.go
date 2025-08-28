package object

type ObjType int

const (
	OBJ_STRING ObjType = iota
)

type Obj struct {
	Type_ ObjType
	Val   any
}

type ObjString struct {
	Chars string
	hash  uint32
}

func CopyString(s *string, start int, length int) Obj {
	d := (*s)[start : start+length]
	h := hashString(d, len(d))
	o := Obj{Type_: OBJ_STRING, Val: ObjString{Chars: d, hash: h}}
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
