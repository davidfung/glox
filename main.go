package main

import (
	"fmt"
	"os"
)

const DEBUG_TRACE_EXECUTION = true

func repl() {
}

func runFile(path string) {
}

func main() {
	initVM()

	if len(os.Args) == 1 {
		repl()
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		fmt.Fprintln(os.Stderr, "Usage: glox [path]")
	}

	freeVM()
}
