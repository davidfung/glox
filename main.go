package main

import "fmt"

func main() {
	fmt.Println("the beginning of glox starts here...")
	var chunk Chunk
	writeChunk(&chunk, OP_RETURN) //TODO can we just pass chunk instead of &chunk?
	//TODO write disassembleChunk(&chunk, "test chunk")
}
