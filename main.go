package main

func main() {
	initVM()

	var chunk Chunk
	initChunk(&chunk)

	constant := addConstant(&chunk, 1.2)
	writeChunk(&chunk, OP_CONSTANT, 123)
	writeChunk(&chunk, uint8(constant), 123)

	writeChunk(&chunk, OP_RETURN, 123) //TODO can we just pass chunk instead of &chunk?

	disassembleChunk(&chunk, "test chunk")
	interpret(&chunk)
	freeVM()
	freeChunk(&chunk)
}
