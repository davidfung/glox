package debugger

import (
	"fmt"

	"github.com/davidfung/glox/chunk"
	"github.com/davidfung/glox/value"
)

const DEBUG_TRACE_EXECUTION = true

func disassembleChunk(chun *chunk.Chunk, name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < len(chun.Code); {
		offset = DisassembleInstruction(chun, offset)
	}
}

func constantInstruction(name string, chun *chunk.Chunk, offset int) int {
	constant := chun.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	value.PrintValue(chun.Constants.Values[constant])
	fmt.Println()
	return offset + 2
}

func DisassembleInstruction(chun *chunk.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && chun.Lines[offset] == chun.Lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", chun.Lines[offset])
	}

	instruction := chun.Code[offset]
	switch instruction {
	case chunk.OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", chun, offset)
	case chunk.OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case chunk.OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset)
	case chunk.OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case chunk.OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	case chunk.OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
	case chunk.OP_RETURN:
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
