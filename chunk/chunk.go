package chunk

import "github.com/davidfung/glox/value"

type OpCode uint8

type Byte interface {
	uint8 | OpCode
}

const (
	OP_CONSTANT OpCode = iota
	OP_NIL
	OP_TRUE
	OP_FALSE
	OP_POP
	OP_GET_LOCAL
	OP_SET_LOCAL
	OP_GET_GLOBAL
	OP_DEFINE_GLOBAL
	OP_SET_GLOBAL
	OP_EQUAL
	OP_GREATER
	OP_LESS
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NOT
	OP_NEGATE
	OP_PRINT
	OP_JUMP
	OP_JUMP_IF_FALSE
	OP_RETURN
)

type Chunk struct {
	Code      []uint8
	Lines     []int
	Constants value.ValueArray
}

func InitChunk(chun *Chunk) {
	chun.Code = nil
	chun.Lines = nil
	value.InitValueArray(&chun.Constants)
}

func WriteChunk[B Byte](chun *Chunk, code B, line int) {
	chun.Code = append(chun.Code, uint8(code))
	chun.Lines = append(chun.Lines, line)
}

func AddConstant(chun *Chunk, val value.Value) int {
	value.WriteValueArray(&chun.Constants, val)
	return len(chun.Constants.Values) - 1
}

func FreeChunk(chun *Chunk) {
	value.FreeValueArrary(&chun.Constants)
	InitChunk(chun)
}
