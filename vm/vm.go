package vm

import (
	"fmt"
	"os"

	"github.com/davidfung/glox/chunk"
	"github.com/davidfung/glox/compiler"
	"github.com/davidfung/glox/debugger"
	"github.com/davidfung/glox/object"
	"github.com/davidfung/glox/objval"
	"github.com/davidfung/glox/table"
	"github.com/davidfung/glox/value"
)

const STACK_MAX = 256

type VM struct {
	chunk    *chunk.Chunk
	ip       int
	stack    [STACK_MAX]value.Value
	stackTop int
	globals  table.Table
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
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)

	instruction := vm.ip - 1
	line := vm.chunk.Lines[instruction]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	resetStack()
}

func InitVM() {
	resetStack()
	table.InitTable(&vm.globals)
}

func FreeVM() {
	table.FreeTable(&vm.globals)
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
	return objval.IS_NIL(val) || objval.IS_BOOL(val) && !objval.AS_BOOL(val)
}

func concatenate() InterpretResult {
	b := objval.AS_STRING(pop())
	a := objval.AS_STRING(pop())
	c := a + b
	o := object.Obj{Type_: object.OBJ_STRING, Val: c}
	v := objval.OBJ_VAL(o)
	push(v)
	return INTERPRET_OK
}

func binary_op(op BinaryOp) InterpretResult {
	if !objval.IS_NUMBER(peek(0)) || !objval.IS_NUMBER(peek(1)) {
		runtimeError("Operands must be numbers.")
		return INTERPRET_RUNTIME_ERROR
	}
	b := objval.AS_NUMBER(pop())
	a := objval.AS_NUMBER(pop())
	switch op {
	case BINARY_OP_ADD:
		push(objval.NUMBER_VAL(a + b))
	case BINARY_OP_SUBTRACT:
		push(objval.NUMBER_VAL(a - b))
	case BINARY_OP_MULTIPLY:
		push(objval.NUMBER_VAL(a * b))
	case BINARY_OP_DIVIDE:
		push(objval.NUMBER_VAL(a / b))
	case BINARY_OP_GREATER:
		push(objval.BOOL_VAL(a > b))
	case BINARY_OP_LESS:
		push(objval.BOOL_VAL(a < b))
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

	readShort := func() uint16 {
		vm.ip += 2
		var x uint16 = (uint16(vm.chunk.Code[vm.ip-2]) << 8) | uint16((vm.chunk.Code[vm.ip-1]))
		return x
	}

	readString := func() object.ObjString {
		return objval.AS_STRING(readConstant())
	}

	for {
		if debugger.DEBUG_TRACE_EXECUTION {
			fmt.Printf("         ")
			for i := 0; i < vm.stackTop; i++ {
				fmt.Printf("[")
				objval.PrintValue(vm.stack[i])
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
			push(objval.NIL_VAL())
		case chunk.OP_TRUE:
			push(objval.BOOL_VAL(true))
		case chunk.OP_FALSE:
			push(objval.BOOL_VAL(false))
		case chunk.OP_POP:
			pop()
		case chunk.OP_GET_LOCAL:
			// Load the value from the local index and then
			// push it on top of the stack where later
			// instructions can find it.
			slot := readByte()
			push(vm.stack[slot])
		case chunk.OP_SET_LOCAL:
			// Take the assigned value from the top of the
			// stack and stores it in the stack slot corresponding
			// to the local variable.  Note that it does not
			// pop the value from the stack because assignment is
			// an expression, and every expression produces a value.
			// The value of an assignment expression is the assigned
			// value itself, so the VM just leaves the value on the
			// stack.
			slot := readByte()
			vm.stack[slot] = peek(0)
		case chunk.OP_GET_GLOBAL:
			name := readString()
			val, ok := table.TableGet(&vm.globals, name)
			if !ok {
				runtimeError("Undefined variable '%s'.", name)
				return INTERPRET_RUNTIME_ERROR
			}
			push(val)
		case chunk.OP_DEFINE_GLOBAL:
			name := readString()
			table.TableSet(&vm.globals, name, peek(0))
			pop()
		case chunk.OP_SET_GLOBAL:
			name := readString()
			if table.TableSet(&vm.globals, name, peek(0)) {
				// Lox doesn't support implicit variable declaration
				table.TableDelete(&vm.globals, name)
				runtimeError("Undefined variable '%s'.", name)
				return INTERPRET_RUNTIME_ERROR
			}
		case chunk.OP_EQUAL:
			a := pop()
			b := pop()
			push(objval.BOOL_VAL(objval.ValuesEqual(a, b)))
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
			if objval.IS_STRING(peek(0)) && objval.IS_STRING(peek(1)) {
				result = concatenate()
			} else if objval.IS_NUMBER(peek(0)) && objval.IS_NUMBER(peek(1)) {
				result = binary_op(BINARY_OP_ADD)
			} else {
				runtimeError("Operands must be two numbers or two strings.")
				return INTERPRET_RUNTIME_ERROR
			}
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
			push(objval.BOOL_VAL(isFalsey(pop())))
		case chunk.OP_NEGATE:
			if !objval.IS_NUMBER(peek(0)) {
				runtimeError("Operand must be a number")
			}
			push(objval.NUMBER_VAL(-objval.AS_NUMBER(pop())))
		case chunk.OP_PRINT:
			objval.PrintValue(pop())
			fmt.Println()
		case chunk.OP_JUMP_IF_FALSE:
			offset := readShort()
			if isFalsey(peek(0)) {
				vm.ip += int(offset)
			}
		case chunk.OP_RETURN:
			// Exit interpreter.
			return INTERPRET_OK
		}
	}
}
