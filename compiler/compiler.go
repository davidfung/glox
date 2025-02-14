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

type Precedence int

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

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

func advance() {
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
		advance()
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

func binary() {
	operatorType := parser.previous.Type
	rule := getRule(operatorType)
	parsePrecedence(rule.precedence + 1)

	switch operatorType {
	case scanner.TOKEN_PLUS:
		emitByte(chunk.OP_ADD)
	case scanner.TOKEN_MINUS:
		emitByte(chunk.OP_SUBTRACT)
	case scanner.TOKEN_STAR:
		emitByte(chunk.OP_MULTIPLY)
	case scanner.TOKEN_SLASH:
		emitByte(chunk.OP_DIVIDE)
	default:
		return // Unreachable.
	}
}

func grouping() {
	expression()
	consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func expression() {
	parsePrecedence(PREC_ASSIGNMENT)
}

func number() {
	beg := parser.previous.Start
	end := parser.previous.Start + parser.previous.Length
	s := (*parser.previous.Source)[beg:end]
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		error(err.Error())
	}
	emitConstant(value.Value(val))
}

func unary() {
	operatorType := parser.previous.Type

	// Compile the operand.
	parsePrecedence(PREC_UNARY)

	// Emit the operator instruction.
	switch operatorType {
	case scanner.TOKEN_MINUS:
		emitByte(chunk.OP_NEGATE)
	default: // Unreachable
		return
	}
}

func parsePrecedence(precedence Precedence) {
	// What goes here?
}

func Compile(source *string, chunk *chunk.Chunk) bool {
	scanner.InitScanner(source)
	compilingChunk = chunk

	parser.hadError = false
	parser.panicMode = false

	advance()
	expression()
	consume(scanner.TOKEN_EOF, "Expect end of expression.")
	endCompiler()
	return !parser.hadError
}
