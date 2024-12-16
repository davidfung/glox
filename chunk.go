package main

import "fmt"

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

func disassembleChunk(chunk *Chunk, name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < len(chunk.code); {
		offset = disassembleInstruction(chunk, offset)
	}
}

func disassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	instruction := chunk.code[offset]
	switch instruction {
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}
