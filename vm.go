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

func interpret(chunk *Chunk) int {
	vm.chunk = chunk
	vm.ip = 0
	return run()
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
		case OP_RETURN:
			printValue(pop())
			fmt.Printf("\n")
			return INTERPRET_OK
		}
	}
}
