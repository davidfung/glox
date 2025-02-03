package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const DEBUG_TRACE_EXECUTION = true

const versionMajor = 0
const versionMinor = 3
const versionPatch = 0

func repl() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Println()
	fmt.Println("Type ctrl-d to exit.")
	for {
		fmt.Printf("> ")
		if !input.Scan() {
			break
		}
		source := input.Text() //TODO why &(input.Text()) does not work?
		interpret(&source)
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
	result := interpret(&source)
	if result == INTERPRET_COMPILE_ERROR {
		os.Exit(65)
	}
	if result == INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
	}

}

func printVersion() {
	fmt.Printf("glox version %d.%d.%d\n", versionMajor, versionMinor, versionPatch)
}

func main() {
	printVersion()

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
