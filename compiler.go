package main

import "fmt"

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
	line := -1
	for {
		token := scanToken()
		if token.line != line {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%2d '%.*s'\n", token.type, token.length, token.start) //TOFIX

		if token.type == TOKEN_EOF {
			break
		}
	}
}
