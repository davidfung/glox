package main

import "fmt"

func disassembleChunk(chunk *Chunk, name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < len(chunk.code); {
		offset = disassembleInstruction(chunk, offset)
	}
}

func constantInstruction(name string, chunk *Chunk, offset int) int {
	constant := chunk.code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	printValue(chunk.constants.values[constant])
	fmt.Println()
	return offset + 2
}

func disassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && chunk.lines[offset] == chunk.lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", chunk.lines[offset])
	}

	instruction := chunk.code[offset]
	switch instruction {
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", chunk, offset)
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
