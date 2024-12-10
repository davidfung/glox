package main

func writeChunk(chunk *Chunk, code uint8) {
	chunk.code = append(chunk.code, code)
}
