package main

const (
	OP_RETURN = iota
)

type Chunk struct {
	code []uint8
}

func initChunk(chunk *Chunk) {
	chunk.code = nil
}

func writeChunk(chunk *Chunk, code uint8) {
	chunk.code = append(chunk.code, code)
}

func freeChunk(chunk *Chunk) {
	initChunk(chunk)
}
