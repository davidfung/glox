package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/davidfung/glox/vm"
)

const versionMajor = 19
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
		source := input.Text()
		vm.Interpret(&source)
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
	result := vm.Interpret(&source)
	if result == vm.INTERPRET_COMPILE_ERROR {
		os.Exit(65)
	}
	if result == vm.INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
	}

}

func printVersion() {
	fmt.Printf("glox version %d.%d.%d\n", versionMajor, versionMinor, versionPatch)
}

func main() {
	printVersion()

	vm.InitVM()

	if len(os.Args) == 1 {
		repl()
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		fmt.Fprintln(os.Stderr, "Usage: glox [path]")
	}

	vm.FreeVM()
}
