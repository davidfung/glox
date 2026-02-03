package debugger

import (
	"fmt"

	"github.com/davidfung/glox/chunk"
	"github.com/davidfung/glox/objval"
)

const DEBUG_PRINT_CODE = true
const DEBUG_TRACE_EXECUTION = (0 == 1)

func DisassembleChunk(chun *chunk.Chunk, name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < len(chun.Code); {
		offset = DisassembleInstruction(chun, offset)
	}
}

func constantInstruction(name string, chun *chunk.Chunk, offset int) int {
	constant := chun.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	objval.PrintValue(chun.Constants.Values[constant])
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

	instruction := chunk.OpCode(chun.Code[offset])
	switch instruction {
	case chunk.OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", chun, offset)
	case chunk.OP_NIL:
		return simpleInstruction("OP_NIL", offset)
	case chunk.OP_TRUE:
		return simpleInstruction("OP_TRUE", offset)
	case chunk.OP_FALSE:
		return simpleInstruction("OP_FALSE", offset)
	case chunk.OP_POP:
		return simpleInstruction("OP_POP", offset)
	case chunk.OP_GET_LOCAL:
		return byteInstruction("OP_GET_LOCAL", chun, offset)
	case chunk.OP_SET_LOCAL:
		return byteInstruction("OP_SET_LOCAL", chun, offset)
	case chunk.OP_GET_GLOBAL:
		return constantInstruction("OP_GET_GLOBAL", chun, offset)
	case chunk.OP_DEFINE_GLOBAL:
		return constantInstruction("OP_DEFINE_GLOBAL", chun, offset)
	case chunk.OP_SET_GLOBAL:
		return constantInstruction("OP_SET_GLOBAL", chun, offset)
	case chunk.OP_EQUAL:
		return simpleInstruction("OP_EQUAL", offset)
	case chunk.OP_GREATER:
		return simpleInstruction("OP_GREATER", offset)
	case chunk.OP_LESS:
		return simpleInstruction("OP_LESS", offset)
	case chunk.OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case chunk.OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset)
	case chunk.OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case chunk.OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	case chunk.OP_NOT:
		return simpleInstruction("OP_NOT", offset)
	case chunk.OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
	case chunk.OP_PRINT:
		return simpleInstruction("OP_PRINT", offset)
	case chunk.OP_JUMP:
		return jumpInstruction("OP_JUMP", 1, chun, offset)
	case chunk.OP_JUMP_IF_FALSE:
		return jumpInstruction("OP_JUMP_IF_FALSE", 1, chun, offset)
	case chunk.OP_LOOP:
		return jumpInstruction("OP_LOOP", -1, chun, offset)
	case chunk.OP_CALL:
		return byteInstruction("OP_CALL", chun, offset)
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

func byteInstruction(name string, chun *chunk.Chunk, offset int) int {
	slot := chun.Code[offset+1]
	fmt.Printf("%-16s %4d\n", name, slot)
	return offset + 2
}

func jumpInstruction(name string, sign int, chun *chunk.Chunk, offset int) int {
	var jump uint16 = uint16(chun.Code[offset+1]) << 8
	jump |= uint16(chun.Code[offset+2])
	fmt.Printf("%-16s %4d -> %d\n", name, offset, offset+3+sign*int(jump))
	return offset + 3
}
