package main

import "fmt"

type VM struct {
	chunk *Chunk
	ip    int
}

const (
	INTERPRET_OK = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

var vm VM

func initVM() {
}

func freeVM() {
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
			disassembleInstruction(vm.chunk, vm.ip)
		}

		instruction := readByte()
		switch instruction {
		case OP_CONSTANT:
			constant := readConstant()
			printValue(constant)
			fmt.Printf("\n")
		case OP_RETURN:
			return INTERPRET_OK
		}
	}
}
