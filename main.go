package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const DEBUG_TRACE_EXECUTION = true

func repl() {
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("> ")
		if !input.Scan() {
			break
		}
		interpret(input.Text())
	}
	fmt.Println("terminating...")
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func runFile(path string) {
	source := readFile(path)
	result := interpret(source)
	if result == INTERPRET_COMPILE_ERROR {
		os.Exit(65)
	}
	if result == INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
	}

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
