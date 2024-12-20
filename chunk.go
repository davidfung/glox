package main

const (
	OP_RETURN = iota
)

type Chunk struct {
	code      []uint8
	constants ValueArray
}

func initChunk(chunk *Chunk) {
	chunk.code = nil
	initValueArray(&chunk.constants)
}

func writeChunk(chunk *Chunk, code uint8) {
	chunk.code = append(chunk.code, code)
}

func addConstant(chunk *Chunk, value Value) int {
	writeValueArray(&chunk.constants, value)
	return len(chunk.constants.values) - 1
}

func freeChunk(chunk *Chunk) {
	freeValueArrary(&chunk.constants)
	initChunk(chunk)
}
