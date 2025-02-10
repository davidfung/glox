package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type Parser struct {
	current   Token
	previous  Token
	hadError  bool
	panicMode bool
}

var parser Parser
var compilingChunk *Chunk

func currentChunk() *Chunk {
	return compilingChunk
}

func errorAt(token Token, message string) {
	if parser.panicMode {
		return
	}
	parser.panicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", token.line)

	if token.type_ == TOKEN_EOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.type_ == TOKEN_ERROR {
		// Nothing.
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", (*token.source)[token.start:token.length])
	}

	fmt.Fprintf(os.Stderr, ": %s\n", message)
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
		if parser.current.type_ != TOKEN_ERROR {
			break
		}
		// TOFIX: not sure why passing the current current token tax to errorAtCurrent
		errorAtCurrent((*parser.current.source)[parser.current.start : parser.current.start+parser.current.length])
	}
}

func consume(typ TokenType, message string) {
	if parser.current.type_ == typ {
		parser_advance()
		return
	}

	errorAtCurrent(message)
}

func emitByte(byte_ uint8) {
	writeChunk(currentChunk(), byte_, parser.previous.line)
}

func emitBytes(byte1 uint8, byte2 uint8) {
	emitByte(byte1)
	emitByte(byte2)
}

func emitReturn() {
	emitByte(OP_RETURN)
}

func makeConstant(value Value) uint8 {
	constant := addConstant(currentChunk(), value)
	if constant > math.MaxUint8 {
		error("Too many constants in one chunk.")
		return 0
	}
	return uint8(constant)
}

func emitConstant(value Value) {
	emitBytes(OP_CONSTANT, makeConstant(value))
}

func endCompiler() {
	emitReturn()
}

func parser_number() {
	beg := parser.previous.start
	end := parser.previous.start + parser.previous.length
	s := (*parser.previous.source)[beg:end]
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		error(err.Error())
	}
	emitConstant(Value(value))
}

func expression() {
	// What goes here?
}

func compile(source *string, chunk *Chunk) bool {
	initScanner(source)
	compilingChunk = chunk

	parser.hadError = false
	parser.panicMode = false

	parser_advance()
	expression()
	consume(TOKEN_EOF, "Expect end of expression.")
	endCompiler()
	return !parser.hadError
}
