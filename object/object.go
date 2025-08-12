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

func CopyString(s *string, start int, len int) Obj {
	d := (*s)[start : start+len]
	o := Obj{Type_: OBJ_STRING, Val: ObjString{Chars: d}}
	return o
}
