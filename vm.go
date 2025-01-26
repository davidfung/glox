package main

import "fmt"

const STACK_MAX = 256

type VM struct {
	chunk    *Chunk
	ip       int
	stack    [STACK_MAX]Value
	stackTop int
}

const (
	INTERPRET_OK = iota
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

func initVM() {
	resetStack()
}

func freeVM() {
}

func push(value Value) {
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func pop() Value {
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

func interpret(source *string) int {
	fmt.Println("Interpreting...")
	fmt.Println(*source)
	compile(source)
	return INTERPRET_OK
}

func run() int {
	readByte := func() uint8 {
		instruction := vm.chunk.code[vm.ip]
		vm.ip++
		return instruction
	}

	readConstant := func() Value {
		return vm.chunk.constants.values[readByte()]
	}

	for {
		if DEBUG_TRACE_EXECUTION {
			fmt.Printf("         ")
			for i := 0; i < vm.stackTop; i++ {
				fmt.Printf("[")
				printValue(vm.stack[i])
				fmt.Printf("]")
			}
			fmt.Printf("\n")
			disassembleInstruction(vm.chunk, vm.ip)
		}

		instruction := readByte()
		switch instruction {
		case OP_CONSTANT:
			constant := readConstant()
			push(constant)
		case OP_ADD:
			binary_op(BINARY_OP_ADD)
		case OP_SUBTRACT:
			binary_op(BINARY_OP_SUBTRACT)
		case OP_MULTIPLY:
			binary_op(BINARY_OP_MULTIPLY)
		case OP_DIVIDE:
			binary_op(BINARY_OP_DIVIDE)
		case OP_NEGATE:
			push(-pop())
		case OP_RETURN:
			printValue(pop())
			fmt.Printf("\n")
			return INTERPRET_OK
		}
	}
}
