package compiler

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/davidfung/glox/chunk"
	"github.com/davidfung/glox/scanner"
	"github.com/davidfung/glox/value"
)

type Parser struct {
	current   scanner.Token
	previous  scanner.Token
	hadError  bool
	panicMode bool
}

var parser Parser
var compilingChunk *chunk.Chunk

func currentChunk() *chunk.Chunk {
	return compilingChunk
}

func errorAt(token scanner.Token, message string) {
	if parser.panicMode {
		return
	}
	parser.panicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", token.Line)

	if token.Type == scanner.TOKEN_EOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.Type == scanner.TOKEN_ERROR {
		// Nothing.
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", (*token.Source)[token.Start:token.Length])
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
		parser.current = scanner.ScanToken()
		if parser.current.Type != scanner.TOKEN_ERROR {
			break
		}
		// TOFIX: not sure why passing the current current token tax to errorAtCurrent
		errorAtCurrent((*parser.current.Source)[parser.current.Start : parser.current.Start+parser.current.Length])
	}
}

func consume(typ scanner.TokenType, message string) {
	if parser.current.Type == typ {
		parser_advance()
		return
	}

	errorAtCurrent(message)
}

func emitByte(byte_ uint8) {
	chunk.WriteChunk(currentChunk(), byte_, parser.previous.Line)
}

func emitBytes(byte1 uint8, byte2 uint8) {
	emitByte(byte1)
	emitByte(byte2)
}

func emitReturn() {
	emitByte(chunk.OP_RETURN)
}

func makeConstant(value value.Value) uint8 {
	constant := chunk.AddConstant(currentChunk(), value)
	if constant > math.MaxUint8 {
		error("Too many constants in one chunk.")
		return 0
	}
	return uint8(constant)
}

func emitConstant(value value.Value) {
	emitBytes(chunk.OP_CONSTANT, makeConstant(value))
}

func endCompiler() {
	emitReturn()
}

func parser_number() {
	beg := parser.previous.Start
	end := parser.previous.Start + parser.previous.Length
	s := (*parser.previous.Source)[beg:end]
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		error(err.Error())
	}
	emitConstant(value.Value(val))
}

func expression() {
	// What goes here?
}

func Compile(source *string, chunk *chunk.Chunk) bool {
	scanner.InitScanner(source)
	compilingChunk = chunk

	parser.hadError = false
	parser.panicMode = false

	parser_advance()
	expression()
	consume(scanner.TOKEN_EOF, "Expect end of expression.")
	endCompiler()
	return !parser.hadError
}
