package vm

import (
	"fmt"

	"github.com/davidfung/glox/chunk"
	"github.com/davidfung/glox/compiler"
	"github.com/davidfung/glox/debugger"
	"github.com/davidfung/glox/value"
)

const STACK_MAX = 256

type VM struct {
	chunk    *chunk.Chunk
	ip       int
	stack    [STACK_MAX]value.Value
	stackTop int
}

type InterpretResult int

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

const (
	BINARY_OP_ADD = iota
	BINARY_OP_SUBTRACT
	BINARY_OP_MULTIPLY
	BINARY_OP_DIVIDE
)

var vm VM

func resetStack() {
	vm.stackTop = 0
}

func InitVM() {
	resetStack()
}

func FreeVM() {
}

func push(value value.Value) {
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func pop() value.Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func binary_op(op int) {
	b := pop()
	a := pop()
	switch op {
	case BINARY_OP_ADD:
		push(a + b)
	case BINARY_OP_SUBTRACT:
		push(a - b)
	case BINARY_OP_MULTIPLY:
		push(a * b)
	case BINARY_OP_DIVIDE:
		push(a / b)
	}
}

func Interpret(source *string) InterpretResult {
	var chun chunk.Chunk
	chunk.InitChunk(&chun)

	if !compiler.Compile(source, &chun) {
		chunk.FreeChunk(&chun)
		return INTERPRET_COMPILE_ERROR
	}

	vm.chunk = &chun
	vm.ip = 0

	result := run()

	chunk.FreeChunk(&chun)
	return result
}

func run() InterpretResult {
	readByte := func() uint8 {
		instruction := vm.chunk.Code[vm.ip]
		vm.ip++
		return instruction
	}

	readConstant := func() value.Value {
		return vm.chunk.Constants.Values[readByte()]
	}

	for {
		if debugger.DEBUG_TRACE_EXECUTION {
			fmt.Printf("         ")
			for i := 0; i < vm.stackTop; i++ {
				fmt.Printf("[")
				value.PrintValue(vm.stack[i])
				fmt.Printf("]")
			}
			fmt.Printf("\n")
			debugger.DisassembleInstruction(vm.chunk, vm.ip)
		}

		instruction := readByte()
		switch instruction {
		case chunk.OP_CONSTANT:
			constant := readConstant()
			push(constant)
		case chunk.OP_ADD:
			binary_op(BINARY_OP_ADD)
		case chunk.OP_SUBTRACT:
			binary_op(BINARY_OP_SUBTRACT)
		case chunk.OP_MULTIPLY:
			binary_op(BINARY_OP_MULTIPLY)
		case chunk.OP_DIVIDE:
			binary_op(BINARY_OP_DIVIDE)
		case chunk.OP_NEGATE:
			push(-pop())
		case chunk.OP_RETURN:
			value.PrintValue(pop())
			fmt.Printf("\n")
			return INTERPRET_OK
		}
	}
}
