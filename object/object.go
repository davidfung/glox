package object

type ObjType int

const (
	OBJ_STRING ObjType = iota
)

type Obj struct {
	Type_ ObjType
	val   any
}

func CopyString(s *string, start int, len int) Obj {
	d := (*s)[start : start+len]
	o := Obj{Type_: OBJ_STRING, val: d}
	return o
}
