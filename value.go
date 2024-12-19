package main

type Value float64

type ValueArray struct {
	values []Value
}

func initValueArray(array ValueArray) {
	array.values = nil
}

func writeValueArray(array ValueArray, value Value) {
	array.values = append(array.values, value)
}

func freeValueArrary(array ValueArray) {
	initValueArray(array)
}
