package value

import "fmt"

type Value float64

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
