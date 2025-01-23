package main

const (
	TOKEN_LEFT_PAREN = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR

	// One or two character tokens
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL

	// Literals

	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FOR
	TOKEN_FUN
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_ERROR
	TOKEN_EOF
)

type Token struct {
	type_  int
	start  int
	length int
	line   int
}

type Scanner struct {
	source  string
	start   int
	current int
	end     int
	line    int
}

var scanner Scanner

func initScanner(source string) {
	scanner.source = source
	scanner.start = 0
	scanner.current = 0
	scanner.end = len(source) - 1
	scanner.line = 1
}

func scanToken() {
	scanner.start = scanner.current
	if isAtEnd() {
		return makeToken(TOKEN_EOF)
	}
	return errorToken("Unexpected character.")
}

func isAtEnd() bool {
	return scanner.start >= scanner.end
}

func makeToken(typ int) Token {
	var token Token
	token.type_ = typ
	token.start = scanner.start
	token.length = scanner.current - scanner.start
	token.line = scanner.line
	return token
}

func errorToken(message string) Token {
	var token Token
	token.type_ = TOKEN_ERROR
	token.start = scanner.start
	token.length = scanner.current - scanner.start
	token.line = scanner.line
	return token
}
