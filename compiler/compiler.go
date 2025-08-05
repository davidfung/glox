package compiler

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/davidfung/glox/chunk"
	"github.com/davidfung/glox/debugger"
	"github.com/davidfung/glox/object"
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

type ParseFn func()

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

var parser Parser
var compilingChunk *chunk.Chunk
var rules []ParseRule

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

func emitByte[B chunk.Byte](byte_ B) {
	chunk.WriteChunk(currentChunk(), byte_, parser.previous.Line)
}

func emitBytes[B1 chunk.Byte, B2 chunk.Byte](byte1 B1, byte2 B2) {
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
	if debugger.DEBUG_PRINT_CODE {
		if !parser.hadError {
			debugger.DisassembleChunk(currentChunk(), "code")
		}
	}
}

func binary() {
	operatorType := parser.previous.Type
	rule := getRule(operatorType)
	parsePrecedence(rule.precedence + 1)

	switch operatorType {
	case scanner.TOKEN_BANG_EQUAL:
		emitBytes(chunk.OP_EQUAL, chunk.OP_NOT)
	case scanner.TOKEN_EQUAL_EQUAL:
		emitByte(chunk.OP_EQUAL)
	case scanner.TOKEN_GREATER:
		emitByte(chunk.OP_GREATER)
	case scanner.TOKEN_GREATER_EQUAL:
		emitBytes(chunk.OP_LESS, chunk.OP_NOT)
	case scanner.TOKEN_LESS:
		emitByte(chunk.OP_LESS)
	case scanner.TOKEN_LESS_EQUAL:
		emitBytes(chunk.OP_GREATER, chunk.OP_NOT)
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

func literal() {
	switch parser.previous.Type {
	case scanner.TOKEN_FALSE:
		emitByte(chunk.OP_FALSE)
	case scanner.TOKEN_NIL:
		emitByte(chunk.OP_NIL)
	case scanner.TOKEN_TRUE:
		emitByte(chunk.OP_TRUE)
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
	emitConstant(value.NUMBER_VAL(val))
}

func str() {
	// Create a string object, wrap it in a Value, and stuff
	// the value into the constant table.
	emitConstant(value.OBJ_VAL(object.CopyString(parser.previous.Source, parser.previous.Start+1, parser.previous.Length-2)))
}

func unary() {
	operatorType := parser.previous.Type

	// Compile the operand.
	parsePrecedence(PREC_UNARY)

	// Emit the operator instruction.
	switch operatorType {
	case scanner.TOKEN_BANG:
		emitByte(chunk.OP_NOT)
	case scanner.TOKEN_MINUS:
		emitByte(chunk.OP_NEGATE)
	default: // Unreachable
		return
	}
}

func parsePrecedence(precedence Precedence) {
	advance()
	prefixRule := getRule(parser.previous.Type).prefix
	if prefixRule == nil {
		error("Expect expression.")
		return
	}

	prefixRule()

	for precedence <= getRule(parser.current.Type).precedence {
		advance()
		infixRule := getRule(parser.previous.Type).infix
		infixRule()
	}
}

func getRule(tokenType scanner.TokenType) ParseRule {
	return rules[tokenType]
}

func initCompiler() {
	rules = []ParseRule{
		scanner.TOKEN_LEFT_PAREN:    {grouping, nil, PREC_NONE},
		scanner.TOKEN_RIGHT_PAREN:   {nil, nil, PREC_NONE},
		scanner.TOKEN_LEFT_BRACE:    {nil, nil, PREC_NONE},
		scanner.TOKEN_RIGHT_BRACE:   {nil, nil, PREC_NONE},
		scanner.TOKEN_COMMA:         {nil, nil, PREC_NONE},
		scanner.TOKEN_DOT:           {nil, nil, PREC_NONE},
		scanner.TOKEN_MINUS:         {unary, binary, PREC_TERM},
		scanner.TOKEN_PLUS:          {nil, binary, PREC_TERM},
		scanner.TOKEN_SEMICOLON:     {nil, nil, PREC_NONE},
		scanner.TOKEN_SLASH:         {nil, binary, PREC_FACTOR},
		scanner.TOKEN_STAR:          {nil, binary, PREC_FACTOR},
		scanner.TOKEN_BANG:          {unary, nil, PREC_NONE},
		scanner.TOKEN_BANG_EQUAL:    {nil, binary, PREC_EQUALITY},
		scanner.TOKEN_EQUAL:         {nil, nil, PREC_NONE},
		scanner.TOKEN_EQUAL_EQUAL:   {nil, binary, PREC_EQUALITY},
		scanner.TOKEN_GREATER:       {nil, binary, PREC_COMPARISON},
		scanner.TOKEN_GREATER_EQUAL: {nil, binary, PREC_COMPARISON},
		scanner.TOKEN_LESS:          {nil, binary, PREC_COMPARISON},
		scanner.TOKEN_LESS_EQUAL:    {nil, binary, PREC_COMPARISON},
		scanner.TOKEN_IDENTIFIER:    {nil, nil, PREC_NONE},
		scanner.TOKEN_STRING:        {str, nil, PREC_NONE},
		scanner.TOKEN_NUMBER:        {number, nil, PREC_NONE},
		scanner.TOKEN_AND:           {nil, nil, PREC_NONE},
		scanner.TOKEN_CLASS:         {nil, nil, PREC_NONE},
		scanner.TOKEN_ELSE:          {nil, nil, PREC_NONE},
		scanner.TOKEN_FALSE:         {literal, nil, PREC_NONE},
		scanner.TOKEN_FOR:           {nil, nil, PREC_NONE},
		scanner.TOKEN_FUN:           {nil, nil, PREC_NONE},
		scanner.TOKEN_IF:            {nil, nil, PREC_NONE},
		scanner.TOKEN_NIL:           {literal, nil, PREC_NONE},
		scanner.TOKEN_OR:            {nil, nil, PREC_NONE},
		scanner.TOKEN_PRINT:         {nil, nil, PREC_NONE},
		scanner.TOKEN_RETURN:        {nil, nil, PREC_NONE},
		scanner.TOKEN_SUPER:         {nil, nil, PREC_NONE},
		scanner.TOKEN_THIS:          {nil, nil, PREC_NONE},
		scanner.TOKEN_TRUE:          {literal, nil, PREC_NONE},
		scanner.TOKEN_VAR:           {nil, nil, PREC_NONE},
		scanner.TOKEN_WHILE:         {nil, nil, PREC_NONE},
		scanner.TOKEN_ERROR:         {nil, nil, PREC_NONE},
		scanner.TOKEN_EOF:           {nil, nil, PREC_NONE},
	}
}

func Compile(source *string, chunk *chunk.Chunk) bool {
	scanner.InitScanner(source)
	compilingChunk = chunk

	parser.hadError = false
	parser.panicMode = false

	initCompiler()
	advance()
	expression()
	consume(scanner.TOKEN_EOF, "Expect end of expression.")
	endCompiler()
	return !parser.hadError
}
