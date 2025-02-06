package main

import (
	"fmt"
	"os"
)

type Parser struct {
	current   Token
	previous  Token
	hadError  bool
	panicMode bool
}

var parser Parser

func errorAt(token Token, message string) {
	if parser.panicMode {
		return
	}
	parser.panicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", token.line)

	if token.typ == TOKEN_EOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.typ == TOKEN_ERROR {
		// Nothing.
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", (*token.source)[token.start:token.length])
	}

	fmt.Fprintf(os.Stderr, ": %s\n, message")
	parser.hadError = true
}

func error(message string) {
	errorAt(parser.previous, message)
}

func errorAtCurrent(message string) {
	errorAt(parser.current, message)
}

func parser_advance() {
	parser.previous = parser.current

	for {
		parser.current = scanToken()
		if parser.current.typ != TOKEN_ERROR {
			break
		}
		// TOFIX: not sure why passing the current current token tax to errorAtCurrent
		errorAtCurrent((*parser.current.source)[parser.current.start : parser.current.start+parser.current.length])
	}
}

func consume(typ TokenType, message string) {
	if parser.current.typ == typ {
		parser_advance()
		return
	}

	errorAtCurrent(message)
}

func compile(source *string, chunk *Chunk) bool {
	initScanner(source)

	parser.hadError = false
	parser.panicMode = false

	parser_advance()
	expression()
	consume(TOKEN_EOF, "Expect end of expression.")
	return !parser.hadError
}
