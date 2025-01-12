package main

type Scanner struct {
	start   int
	current int
	line    int
}

var scanner Scanner

func initScanner(source string) {
	scanner.start = 0
	scanner.current = 0
	scanner.line = 1
}

func compile(source string) {
	initScanner(source)
}
