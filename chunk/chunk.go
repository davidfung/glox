package chunk

import "github.com/davidfung/glox/value"

const (
	OP_CONSTANT = iota
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NEGATE
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

func WriteChunk(chun *Chunk, code uint8, line int) {
	chun.Code = append(chun.Code, code)
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
