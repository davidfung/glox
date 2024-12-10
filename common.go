package main

const (
	OP_RETURN = iota
)

type Chunk struct {
	code []uint8
}
