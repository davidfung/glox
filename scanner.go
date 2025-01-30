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
	source *string
	typ    int
	start  int
	length int
	line   int
}

type Scanner struct {
	source  *string
	start   int
	current int
	line    int
}

var scanner Scanner

func initScanner(source *string) {
	scanner.source = source
	scanner.start = 0
	scanner.current = 0
	scanner.line = 1
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func scanToken() Token {
	skipWhitespace()
	scanner.start = scanner.current
	if isAtEnd() {
		return makeToken(TOKEN_EOF)
	}

	c := advance()
	if isDigit(c) {
		return number()
	}

	switch c {
	case '(':
		return makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return makeToken(TOKEN_SEMICOLON)
	case ',':
		return makeToken(TOKEN_COMMA)
	case '.':
		return makeToken(TOKEN_DOT)
	case '-':
		return makeToken(TOKEN_MINUS)
	case '+':
		return makeToken(TOKEN_PLUS)
	case '/':
		return makeToken(TOKEN_SLASH)
	case '*':
		return makeToken(TOKEN_STAR)
	case '!':
		if match('=') {
			return makeToken(TOKEN_BANG_EQUAL)
		} else {
			return makeToken(TOKEN_BANG)
		}
	case '=':
		if match('=') {
			return makeToken(TOKEN_EQUAL_EQUAL)
		} else {
			return makeToken(TOKEN_EQUAL)
		}
	case '<':
		if match('=') {
			return makeToken(TOKEN_LESS_EQUAL)
		} else {
			return makeToken(TOKEN_LESS)
		}
	case '>':
		if match('=') {
			return makeToken(TOKEN_GREATER_EQUAL)
		} else {
			return makeToken(TOKEN_GREATER)
		}
	case '"':
		return quotedString()
	}

	return errorToken("Unexpected character.")
}

func isAtEnd() bool {
	return scanner.current >= len(*scanner.source)
}

func advance() byte {
	scanner.current++
	return (*scanner.source)[scanner.current-1]
}

func match(expected byte) bool {
	if isAtEnd() {
		return false
	}
	if (*scanner.source)[scanner.current] != expected {
		return false
	}
	scanner.current++
	return true
}

func makeToken(typ int) Token {
	var token Token
	token.source = scanner.source
	token.typ = typ
	token.start = scanner.start
	token.length = scanner.current - scanner.start
	token.line = scanner.line

	if typ == TOKEN_EOF {
		s := ""
		token.source = &s
		token.start = 0
		token.length = len(s)
	}

	return token
}

func errorToken(msg string) Token {
	var token Token
	token.source = &msg
	token.typ = TOKEN_ERROR
	token.start = 0
	token.length = len(msg)
	token.line = scanner.line
	return token
}

func skipWhitespace() {
	for {
		if isAtEnd() {
			return
		}
		c := peek()
		switch c {
		case ' ':
			advance()
		case '\r':
			advance()
		case '\t':
			advance()
		case '\n':
			scanner.line++
			advance()
		case '/':
			if peekNext() == '/' {
				for peek() != '\n' && !isAtEnd() {
					advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func number() Token {
	for isDigit(peek()) {
		advance()
	}

	// Look for a fractional part.
	if peek() == '.' && isDigit(peekNext()) {
		// Consume the ".".
		advance()

		for isDigit(peek()) {
			advance()
		}
	}

	return makeToken(TOKEN_NUMBER)
}

func quotedString() Token {
	for !isAtEnd() && peek() != '"' {
		if peek() == '\n' {
			scanner.line++
		}
		advance()
	}

	if isAtEnd() {
		return errorToken("Unterminated string.")
	}

	advance() // the closing quote
	return makeToken(TOKEN_STRING)
}

func peek() byte {
	if isAtEnd() {
		return 0
	}
	return byte((*scanner.source)[scanner.current])
}

func peekNext() byte {
	if isAtEnd() {
		return 0
	}
	return byte((*scanner.source)[scanner.current+1])
}
