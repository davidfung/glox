package value

type ValueType int

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
	VAL_OBJ
)

type Value struct {
	Type_ ValueType // promote to export field?
	Val   any       // promote to export field?
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
