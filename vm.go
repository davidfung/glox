package main

type VM struct {
	chunk *Chunk
	ip    *uint8
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
	vm.ip = &vm.chunk.code[0]
	return INTERPRET_OK //TOFIX run()
}
