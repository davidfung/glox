package main

import "fmt"

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
		fmt.Printf("%2d '%s'\n", token.typ, string((*token.source)[token.start:token.start+token.length]))

		if token.typ == TOKEN_EOF {
			break
		}
	}
}
