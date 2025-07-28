package vm

import (
	"fmt"
	"os"

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

type BinaryOp int

const (
	BINARY_OP_ADD BinaryOp = iota
	BINARY_OP_SUBTRACT
	BINARY_OP_MULTIPLY
	BINARY_OP_DIVIDE
	BINARY_OP_GREATER
	BINARY_OP_LESS
)

var vm VM

func resetStack() {
	vm.stackTop = 0
}

func runtimeError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args)
	fmt.Fprintln(os.Stderr)

	instruction := vm.ip - 1
	line := vm.chunk.Lines[instruction]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	resetStack()
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

func peek(distance int) value.Value {
	return vm.stack[vm.stackTop-1-distance]
}

func isFalsey(val value.Value) bool {
	return value.IS_NIL(val) || value.IS_BOOL(val) && !value.AS_BOOL(val)
}

func binary_op(op BinaryOp) InterpretResult {
	if !value.IS_NUMBER(peek(0)) || !value.IS_NUMBER(peek(1)) {
		runtimeError("Operands must be numbers.")
		return INTERPRET_RUNTIME_ERROR
	}
	b := value.AS_NUMBER(pop())
	a := value.AS_NUMBER(pop())
	switch op {
	case BINARY_OP_ADD:
		push(value.NUMBER_VAL(a + b))
	case BINARY_OP_SUBTRACT:
		push(value.NUMBER_VAL(a - b))
	case BINARY_OP_MULTIPLY:
		push(value.NUMBER_VAL(a * b))
	case BINARY_OP_DIVIDE:
		push(value.NUMBER_VAL(a / b))
	case BINARY_OP_GREATER:
		push(value.BOOL_VAL(a > b))
	case BINARY_OP_LESS:
		push(value.BOOL_VAL(a < b))
	}
	return INTERPRET_OK
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
	var result InterpretResult

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

		instruction := chunk.OpCode(readByte())
		switch instruction {
		case chunk.OP_CONSTANT:
			constant := readConstant()
			push(constant)
		case chunk.OP_NIL:
			push(value.NIL_VAL())
		case chunk.OP_TRUE:
			push(value.BOOL_VAL(true))
		case chunk.OP_FALSE:
			push(value.BOOL_VAL(false))
		case chunk.OP_EQUAL:
			a := pop()
			b := pop()
			push(value.BOOL_VAL(value.ValuesEqual(a, b)))
		case chunk.OP_GREATER:
			result = binary_op(BINARY_OP_GREATER)
			if result != INTERPRET_OK {
				return result
			}
		case chunk.OP_LESS:
			result = binary_op(BINARY_OP_LESS)
			if result != INTERPRET_OK {
				return result
			}
		case chunk.OP_ADD:
			result = binary_op(BINARY_OP_ADD)
			if result != INTERPRET_OK {
				return result
			}
		case chunk.OP_SUBTRACT:
			result = binary_op(BINARY_OP_SUBTRACT)
			if result != INTERPRET_OK {
				return result
			}
		case chunk.OP_MULTIPLY:
			result = binary_op(BINARY_OP_MULTIPLY)
			if result != INTERPRET_OK {
				return result
			}
		case chunk.OP_DIVIDE:
			result = binary_op(BINARY_OP_DIVIDE)
			if result != INTERPRET_OK {
				return result
			}
		case chunk.OP_NOT:
			push(value.BOOL_VAL(isFalsey(pop())))
		case chunk.OP_NEGATE:
			if !value.IS_NUMBER(peek(0)) {
				runtimeError("Operand must be a number")
			}
			push(value.NUMBER_VAL(-value.AS_NUMBER(pop())))
		case chunk.OP_RETURN:
			value.PrintValue(pop())
			fmt.Printf("\n")
			return INTERPRET_OK
		}
	}
}
