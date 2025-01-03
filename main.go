package main

const DEBUG_TRACE_EXECUTION = true

func main() {
	initVM()

	var chunk Chunk
	initChunk(&chunk)

	var constant int

	constant = addConstant(&chunk, 1.2)
	writeChunk(&chunk, OP_CONSTANT, 123)
	writeChunk(&chunk, uint8(constant), 123)

	constant = addConstant(&chunk, 3.4)
	writeChunk(&chunk, OP_CONSTANT, 123)
	writeChunk(&chunk, uint8(constant), 123)

	writeChunk(&chunk, OP_ADD, 123)

	constant = addConstant(&chunk, 5.6)
	writeChunk(&chunk, OP_CONSTANT, 123)
	writeChunk(&chunk, uint8(constant), 123)

	writeChunk(&chunk, OP_DIVIDE, 123)
	writeChunk(&chunk, OP_NEGATE, 123)

	writeChunk(&chunk, OP_RETURN, 123) //TODO can we just pass chunk instead of &chunk?

	disassembleChunk(&chunk, "test chunk")
	interpret(&chunk)
	freeVM()
	freeChunk(&chunk)
}
