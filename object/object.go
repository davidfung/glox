package object

type ObjType int

const (
	OBJ_STRING ObjType = iota
)

type Obj struct {
	Type_ ObjType
	val   any
}
