package main

const (
	OP_CONSTANT = iota
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NEGATE
	OP_RETURN
)

type Chunk struct {
	code      []uint8
	lines     []int
	constants ValueArray
}

func initChunk(chunk *Chunk) {
	chunk.code = nil
	chunk.lines = nil
	initValueArray(&chunk.constants)
}

func writeChunk(chunk *Chunk, code uint8, line int) {
	chunk.code = append(chunk.code, code)
	chunk.lines = append(chunk.lines, line)
}

func addConstant(chunk *Chunk, value Value) int {
	writeValueArray(&chunk.constants, value)
	return len(chunk.constants.values) - 1
}

func freeChunk(chunk *Chunk) {
	freeValueArrary(&chunk.constants)
	initChunk(chunk)
}
