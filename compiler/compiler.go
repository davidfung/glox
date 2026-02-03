package compiler

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/davidfung/glox/chunk"
	"github.com/davidfung/glox/common"
	"github.com/davidfung/glox/debugger"
	"github.com/davidfung/glox/object"
	"github.com/davidfung/glox/objval"
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

type ParseFn func(canAssign bool)

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

type Local struct {
	name  scanner.Token
	depth int
}

type FunctionType int

const (
	_ FunctionType = iota
	TYPE_FUNCTION
	TYPE_SCRIPT
)

// We have a simple, flat array of all locals that are in scope
// during each point in the compilation process.  They are ordered
// in the array in the order that their declarations appear in the
// code. Since the instruction operand we’ll use to encode a
// local is a single byte, our VM has a hard limit on the number
// of locals that can be in scope at once. That means we can also
// give the locals array a fixed size.
type Compiler struct {
	enclosing *Compiler
	function  object.ObjFunction
	type_     FunctionType

	locals     [common.UINT8_COUNT]Local
	localCount int
	scopeDepth int
}

var parser Parser
var current *Compiler
var compilingChunk *chunk.Chunk
var rules []ParseRule

func currentChunk() *chunk.Chunk {
	return &current.function.Chun
}

// This function is a convenient function in glox, and is not in
// the original clox.
func currentIP() int {
	return len(currentChunk().Code)
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
		fmt.Fprintf(os.Stderr, " at '%s'", (*token.Source)[token.Start:token.Start+token.Length])
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
		// TOFIX: not sure why passing the current text to errorAtCurrent
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

// check next token
func check(typ scanner.TokenType) bool {
	return parser.current.Type == typ
}

// check next token, advance if match
func match(typ scanner.TokenType) bool {
	if !check(typ) {
		return false
	}
	advance()
	return true
}

func emitByte[B chunk.Byte](byte_ B) {
	chunk.WriteChunk(currentChunk(), byte_, parser.previous.Line)
}

func emitBytes[B1 chunk.Byte, B2 chunk.Byte](byte1 B1, byte2 B2) {
	emitByte(byte1)
	emitByte(byte2)
}

func emitLoop(loopStart int) {
	emitByte(chunk.OP_LOOP)

	offset := len(currentChunk().Code) - loopStart + 2
	if offset > common.UINT16_MAX {
		error("Loop body too large.")
	}

	emitByte(uint8((offset >> 8) & 0xff))
	emitByte(uint8(offset & 0xff))
}

func emitJump[B chunk.Byte](byte_ B) int {
	emitByte(byte_)
	emitByte(B(0xFF))
	emitByte(B(0xFF))
	return len(currentChunk().Code) - 2
}

func emitReturn() {
	emitByte(chunk.OP_NIL)
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

// This goes back into the bytecode and replaces the operand
// at the given location with the calculated jump offset.
// We call patchJump() right before we emit the next instruction
// that we want the jump to land on, so it uses the current
// bytecode count to determine how far to jump.
func patchJump(offset int) {
	// -2 to adjust for the bytecode for the jump offset itself.
	jump := len(currentChunk().Code) - offset - 2

	if jump > common.UINT16_MAX {
		error("Too much code to jump over.")
	}

	currentChunk().Code[offset] = uint8((jump >> 8) & 0xff)
	currentChunk().Code[offset+1] = uint8(jump & 0xff)
}

func initCompiler(compiler *Compiler, type_ FunctionType) {
	compiler.enclosing = current
	compiler.type_ = type_
	compiler.localCount = 0
	compiler.scopeDepth = 0
	compiler.function = object.NewFunction()
	current = compiler
	if type_ != TYPE_SCRIPT {
		source := parser.previous.Source
		start := parser.previous.Start
		length := parser.previous.Length
		name := (*source)[start : start+length]
		current.function.Name = object.ObjString(name)
	}

	// The compiler’s locals array keeps track of which stack
	// slots are associated with which local variables or temporaries.
	// The compiler implicitly claims stack slot zero for the
	// VM’s own internal use. We give it an empty name so that
	// the user can’t write an identifier that refers to it.
	var local Local
	local = current.locals[current.localCount]
	current.localCount++
	local.depth = 0
	local.name.Start = 0
	local.name.Length = 0
}

func endCompiler() object.ObjFunction {
	emitReturn()
	var function object.ObjFunction = current.function
	if debugger.DEBUG_PRINT_CODE {
		if !parser.hadError {
			name := function.Name
			if name == "" {
				name = "<script>"
			}
			debugger.DisassembleChunk(currentChunk(), string(name))
		}
	}
	current = current.enclosing
	return function
}

func beginScope() {
	current.scopeDepth++
}

func endScope() {
	current.scopeDepth--

	// When a block ends, we discard any variables declared
	// at the scope depth we just left by simply decrementing
	// the length of the array, and emit an OP_POP instruction
	// to pop them from the stack.
	for current.localCount > 0 &&
		current.locals[current.localCount-1].depth > current.scopeDepth {
		emitByte(chunk.OP_POP)
		current.localCount--
	}
}

func binary(canAssign bool) {
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

func call(canAssign bool) {
	argCount := argumentList()
	emitBytes(chunk.OP_CALL, argCount)
}

func literal(canAssign bool) {
	switch parser.previous.Type {
	case scanner.TOKEN_FALSE:
		emitByte(chunk.OP_FALSE)
	case scanner.TOKEN_NIL:
		emitByte(chunk.OP_NIL)
	case scanner.TOKEN_TRUE:
		emitByte(chunk.OP_TRUE)
	}
}

func grouping(canAssign bool) {
	expression()
	consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func expression() {
	parsePrecedence(PREC_ASSIGNMENT)
}

func ifStatement() {
	consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'if'.")
	expression()
	consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after condition.")

	thenJump := emitJump(chunk.OP_JUMP_IF_FALSE)
	emitByte(chunk.OP_POP)
	statement()

	elseJump := emitJump(chunk.OP_JUMP)

	patchJump(thenJump)
	emitByte(chunk.OP_POP)

	if match(scanner.TOKEN_ELSE) {
		statement()
	}
	patchJump(elseJump)
}

func block() {
	for !check(scanner.TOKEN_RIGHT_BRACE) && !check(scanner.TOKEN_EOF) {
		declaration()
	}
	consume(scanner.TOKEN_RIGHT_BRACE, "Expect '}' after block.")
}

func function(type_ FunctionType) {
	var compiler Compiler
	initCompiler(&compiler, type_)
	beginScope() // no need for a matching endScope()

	//fun() {}
	consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after function name.")
	if !check(scanner.TOKEN_RIGHT_PAREN) {
		for {
			current.function.Arity++
			if current.function.Arity > 255 {
				errorAtCurrent("Can't have more than 255 parameters.")
			}
			constant := parseVariable("Expect parameter name.")
			defineVariable(constant)
			if !match(scanner.TOKEN_COMMA) {
				break
			}
		}
	}
	consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after parameters.")
	consume(scanner.TOKEN_LEFT_BRACE, "Expect '{' after function body.")
	block()

	function := endCompiler()
	obj := object.Obj{Type_: object.OBJ_FUNCTION, Val: function}
	emitBytes(chunk.OP_CONSTANT, makeConstant(objval.OBJ_VAL(obj)))
}

func funDeclaration() {
	global := parseVariable("Expect function name.")
	markInitialized()
	function(TYPE_FUNCTION)
	defineVariable(global)
}

// The production of declaration grammar rule.
func varDeclaration() {
	global := parseVariable("Expect variable name.")

	if match(scanner.TOKEN_EQUAL) {
		expression()
	} else {
		emitByte(chunk.OP_NIL)
	}
	consume(scanner.TOKEN_SEMICOLON, "Expect ';' after variable declaration.")

	defineVariable(global)
}

func expressionStatement() {
	expression()
	consume(scanner.TOKEN_SEMICOLON, "Expect ';' after expression.")
	emitByte(chunk.OP_POP)
}

// As with implementing for loops in jlox/clox, we didn’t need to touch
// the runtime. It all gets compiled down to primitive control flow
// operations the VM already supports.
func forStatement() {
	beginScope()
	consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'for'.")

	// Initializer clause
	if match(scanner.TOKEN_SEMICOLON) {
		// No initializer
	} else if match(scanner.TOKEN_VAR) {
		varDeclaration()
	} else {
		// We call expressionStatement() instead of expression().
		// That looks for a semicolon, which we need here too, and
		// also emits an OP_POP instruction to discard the value.
		// We don’t want the initializer to leave anything on the
		// stack.
		expressionStatement()
	}

	loopStart := currentIP()

	// Condition clause
	exitJump := -1
	if !match(scanner.TOKEN_SEMICOLON) {
		expression()
		consume(scanner.TOKEN_SEMICOLON, "Expect ';' after loop condition.")

		// Jump out of the loop if the condition is false
		exitJump = emitJump(chunk.OP_JUMP_IF_FALSE)
		emitByte(chunk.OP_POP)
	}

	// Increment clause
	//
	// The last part is a little tricky. First, we emit a loop instruction.
	// This is the main loop that takes us back to the top of the for
	// loop—right before the condition expression if there is one. That
	// loop happens right after the increment, since the increment executes
	// at the end of each loop iteration.
	//
	// Then we change loopStart to point to the offset where the increment
	// expression begins. Later, when we emit the loop instruction after the
	// body statement, this will cause it to jump up to the increment expression
	// instead of the top of the loop like it does when there is no increment.
	// This is how we weave the increment in to run after the body.
	if !match(scanner.TOKEN_RIGHT_PAREN) {
		bodyJump := emitJump(chunk.OP_JUMP)
		incrementStart := currentIP()
		expression()
		emitByte(chunk.OP_POP)
		consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after for clauses.")

		emitLoop(loopStart)
		loopStart = incrementStart
		patchJump(bodyJump)
	}

	statement()
	emitLoop(loopStart)

	if exitJump != (-1) {
		patchJump(exitJump)
		emitByte(chunk.OP_POP) // pop the condition
	}

	endScope()
}

func printStatement() {
	expression()
	consume(scanner.TOKEN_SEMICOLON, "Expect ';' after value.")
	emitByte(chunk.OP_PRINT)
}

func returnStatement() {
	if current.type_ == TYPE_SCRIPT {
		error("Can't return from top-level code.")
	}

	if match(scanner.TOKEN_SEMICOLON) {
		emitReturn()
	} else {
		expression()
		consume(scanner.TOKEN_SEMICOLON, "Expect ';' after return value.")
		emitByte(chunk.OP_RETURN)
	}
}

func whileStatement() {
	loopStart := len(currentChunk().Code)
	consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'while'.")
	expression()
	consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after 'condition'.")

	exitJump := emitJump(chunk.OP_JUMP_IF_FALSE)
	emitByte(chunk.OP_POP)
	statement()
	emitLoop(loopStart)

	patchJump(exitJump)
	emitByte(chunk.OP_POP)
}

func synchronize() {
	parser.panicMode = false

	for parser.current.Type != scanner.TOKEN_EOF {
		if parser.previous.Type == scanner.TOKEN_SEMICOLON {
			return
		}
		switch parser.current.Type {
		case scanner.TOKEN_CLASS:
			return
		case scanner.TOKEN_FUN:
			return
		case scanner.TOKEN_VAR:
			return
		case scanner.TOKEN_FOR:
			return
		case scanner.TOKEN_IF:
			return
		case scanner.TOKEN_WHILE:
			return
		case scanner.TOKEN_PRINT:
			return
		case scanner.TOKEN_RETURN:
			return
		default:
		}
		advance()
	}
}

func declaration() {
	if match(scanner.TOKEN_FUN) {
		funDeclaration()
	} else if match(scanner.TOKEN_VAR) {
		varDeclaration()
	} else {
		statement()
	}
	if parser.panicMode {
		synchronize()
	}
}

func statement() {
	if match(scanner.TOKEN_PRINT) {
		printStatement()
	} else if match(scanner.TOKEN_FOR) {
		forStatement()
	} else if match(scanner.TOKEN_IF) {
		ifStatement()
	} else if match(scanner.TOKEN_RETURN) {
		returnStatement()
	} else if match(scanner.TOKEN_WHILE) {
		whileStatement()
	} else if match(scanner.TOKEN_LEFT_BRACE) {
		beginScope()
		block()
		endScope()
	} else {
		expressionStatement()
	}
}

func number(canAssign bool) {
	beg := parser.previous.Start
	end := parser.previous.Start + parser.previous.Length
	s := (*parser.previous.Source)[beg:end]
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		error(err.Error())
	}
	emitConstant(objval.NUMBER_VAL(val))
}

func or_(canAssign bool) {
	elseJump := emitJump(chunk.OP_JUMP_IF_FALSE)
	endJump := emitJump(chunk.OP_JUMP)

	patchJump(elseJump) // this jump is solely for the POP
	emitByte(chunk.OP_POP)

	parsePrecedence(PREC_OR)
	patchJump(endJump) // this jump is to go to the right operand expression
}

func str(canAssign bool) {
	// Create a string object, wrap it in a Value, and stuff
	// the value into the constant table.
	emitConstant(objval.OBJ_VAL(object.CopyString(parser.previous.Source, parser.previous.Start+1, parser.previous.Length-2)))
}

func namedVariable(token scanner.Token, canAssign bool) {
	// arg := identifierConstant(token)
	var getOp, setOp chunk.OpCode
	var arg int = resolveLocal(current, &token)
	if arg != (-1) {
		getOp = chunk.OP_GET_LOCAL
		setOp = chunk.OP_SET_LOCAL
	} else {
		arg = int(identifierConstant(token))
		getOp = chunk.OP_GET_GLOBAL
		setOp = chunk.OP_SET_GLOBAL
	}

	if canAssign && match(scanner.TOKEN_EQUAL) {
		expression()
		emitBytes(setOp, uint8(arg))
	} else {
		emitBytes(getOp, uint8(arg))
	}
}

func variable(canAssign bool) {
	namedVariable(parser.previous, canAssign)
}

func unary(canAssign bool) {
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

	canAssign := precedence <= PREC_ASSIGNMENT
	prefixRule(canAssign)

	for precedence <= getRule(parser.current.Type).precedence {
		advance()
		infixRule := getRule(parser.previous.Type).infix
		infixRule(canAssign)
	}

	// If assignment is allowed, and the equal sign still exists at this point,
	// it is an error because the equal sign should be already consumed.
	if canAssign && match(scanner.TOKEN_EQUAL) {
		error("Invalid assignment target.")
	}
}

// The token is the name of the identifier.
// Add a value in the constant table and return its index.
func identifierConstant(token scanner.Token) uint8 {
	strobj := object.CopyString(token.Source, token.Start, token.Length)
	return makeConstant(objval.OBJ_VAL(strobj))
}

func identifierEqual(a *scanner.Token, b *scanner.Token) bool {
	if a.Length != b.Length {
		return false
	}
	return (*a.Source)[a.Start:a.Start+a.Length] == (*b.Source)[b.Start:b.Start+b.Length]
}

// We walk the list of locals that are currently in scope. If one
// has the same name as the identifier token, the identifier must
// refer to that variable. We’ve found it! We walk the array backward
// so that we find the last declared variable with the identifier.
// That ensures that inner local variables correctly shadow locals
// with the same name in surrounding scopes.
//
// At runtime, we load and store locals using the stack slot index,
// so that’s what the compiler needs to calculate after it resolves
// the variable. Whenever a variable is declared, we append it to
// the locals array in Compiler. That means the first local variable
// is at index zero, the next one is at index one, and so on. In
// other words, the locals array in the compiler has the exact same
// layout as the VM’s stack will have at runtime. The variable’s index
// in the locals array is the same as its stack slot. How convenient!
//
// If we make it through the whole array without finding a variable
// with the given name, it must not be a local. In that case, we
// return -1 to signal that it wasn’t found and should be assumed to
// be a global variable instead.
func resolveLocal(compiler *Compiler, name *scanner.Token) int {
	for i := compiler.localCount - 1; i >= 0; i-- {
		local := &compiler.locals[i]
		if identifierEqual(name, &local.name) {
			if local.depth == (-1) {
				error("Can't read local variable in its own initializer.")
			}
			return i
		}
	}
	return -1
}

// Initializes the next available Local in the compiler's array
// of variables.  It stores the variable's name and the depth
// of the scope that owns the variable.
func addLocal(name scanner.Token) {
	if current.localCount == common.UINT8_COUNT {
		error("Too may local variables in function.")
		return
	}

	local := &current.locals[current.localCount]
	current.localCount++
	local.name = name
	local.depth = current.scopeDepth
}

// The function declareVariable() is where the compiler records
// the existence of the variable.
func declareVariable() {
	// We only do this for locals, so if we’re in the top-level global scope,
	// we just bail out. Because global variables are late bound, the compiler
	// doesn’t keep track of which declarations for them it has seen.
	if current.scopeDepth == 0 {
		return
	}

	// But for local variables, the compiler does need to remember that the
	// variable exists. That’s what declaring it does—it adds it to the compiler’s
	// list of variables in the current scope. We implement that using addLocal().

	// Local variables are appended to the array when they’re
	// declared, which means the current scope is always at
	// the end of the array. When we declare a new variable,
	// we start at the end and work backward, looking for an
	// existing variable with the same name. If we find one in
	// the current scope, we report the error. Otherwise, if
	// we reach the beginning of the array or a variable owned
	// by another scope, then we know we’ve checked all of the
	// existing variables in the scope.

	name := &parser.previous

	for i := current.localCount - 1; i >= 0; i-- {
		local := current.locals[i]
		if local.depth != -1 && local.depth < current.scopeDepth {
			break
		}
		if identifierEqual(name, &local.name) {
			error("Already a variable with this name in this scope.")
		}
	}

	addLocal(*name)
}

func parseVariable(errorMessage string) uint8 {
	consume(scanner.TOKEN_IDENTIFIER, errorMessage)

	declareVariable()
	// Exit function if we're in a local scope.  At runtime,
	// locals aren't looked up by name.  There's no need to
	// stuff the variable's name into the constant table, so
	// if the declaration is inside a local scope, we return
	// a dummy table index instead.
	if current.scopeDepth > 0 {
		return 0
	}

	return identifierConstant(parser.previous)
}

func markInitialized() {
	if current.scopeDepth == 0 {
		// When a top-level function declaration calls this
		// function, there is no local variable to mark initialized
		// because the function is bound to a global variable.
		return
	}
	current.locals[current.localCount-1].depth = current.scopeDepth
}

func defineVariable(global uint8) {
	// There is no code to create a local variable at runtime.
	// Think about what state the VM is in. It has already
	// executed the code for the variable’s initializer (or
	// the implicit nil if the user omitted an initializer),
	// and that value is sitting right on top of the stack as
	// the only remaining temporary. We also know that new
	// locals are allocated at the top of the stack ... right
	// where that value already is. Thus, there’s nothing to
	// do. The temporary simply becomes the local variable.
	// It doesn’t get much more efficient than that.
	if current.scopeDepth > 0 {
		markInitialized()
		return
	}
	emitBytes(chunk.OP_DEFINE_GLOBAL, global)
}

func argumentList() uint8 {
	var argCount uint8 = 0
	if !check(scanner.TOKEN_RIGHT_PAREN) {
		for {
			expression()
			if argCount == 255 {
				error("Can't have more than 255 arguments.")
			}
			argCount++
			if !match(scanner.TOKEN_COMMA) {
				break
			}
		}
	}
	consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after arguments")
	return argCount
}

// At the point this is called, the left-hand side expression
// has already been compiled. That means at runtime, its value
// will be on top of the stack. If that value is falsey, then
// we know the entire and must be false, so we skip the right
// operand and leave the left-hand side value as the result
// of the entire expression. Otherwise, we discard the
// left-hand value and evaluate the right operand which
// becomes the result of the whole and expression.
func and_(canAssign bool) {
	endJump := emitJump(chunk.OP_JUMP_IF_FALSE)

	emitByte(chunk.OP_POP)
	parsePrecedence(PREC_AND)

	patchJump(endJump)
}

func getRule(tokenType scanner.TokenType) ParseRule {
	return rules[tokenType]
}

func initParseRules() {
	rules = []ParseRule{
		scanner.TOKEN_LEFT_PAREN:    {grouping, call, PREC_CALL},
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
		scanner.TOKEN_IDENTIFIER:    {variable, nil, PREC_NONE},
		scanner.TOKEN_STRING:        {str, nil, PREC_NONE},
		scanner.TOKEN_NUMBER:        {number, nil, PREC_NONE},
		scanner.TOKEN_AND:           {nil, and_, PREC_AND},
		scanner.TOKEN_CLASS:         {nil, nil, PREC_NONE},
		scanner.TOKEN_ELSE:          {nil, nil, PREC_NONE},
		scanner.TOKEN_FALSE:         {literal, nil, PREC_NONE},
		scanner.TOKEN_FOR:           {nil, nil, PREC_NONE},
		scanner.TOKEN_FUN:           {nil, nil, PREC_NONE},
		scanner.TOKEN_IF:            {nil, nil, PREC_NONE},
		scanner.TOKEN_NIL:           {literal, nil, PREC_NONE},
		scanner.TOKEN_OR:            {nil, or_, PREC_OR},
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

func Compile(source *string) object.ObjFunction {
	scanner.InitScanner(source)
	var compiler Compiler
	initCompiler(&compiler, TYPE_SCRIPT)

	parser.hadError = false
	parser.panicMode = false

	initParseRules()
	advance()
	for !match(scanner.TOKEN_EOF) {
		declaration()
	}
	function := endCompiler()
	if parser.hadError {
		// Use Arity = -1 to denote a nil structure in Go.
		return object.ObjFunction{Arity: -1}
	} else {
		return function
	}
}
