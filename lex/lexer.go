package lex

import "go/ast"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	s            ast.Stmt
}

type Lexeme struct {
	Type    Token
	Literal string
	Name    string
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Lexeme {
	var tok Lexeme

	l.skipWhitespace()

	switch l.ch {
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(ASSIGN, ch, l.ch)
		} else {
			tok = newToken(COLON, l.ch)
		}
	case '=':
		tok = newToken(EQL, l.ch)
	case '+':
		tok = newToken(ADD, l.ch)
	case '-':
		tok = newToken(SUBT, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(NEQ, ch, l.ch)
		} else {
			tok = newToken(NOT, l.ch)
		}
	case '/':
		tok = newToken(DIV, l.ch)
	case '*':
		tok = newToken(MULT, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(LEQ, ch, l.ch)
		} else {
			tok = newToken(LSS, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(GEQ, ch, l.ch)
		} else {
			tok = newToken(GTR, l.ch)
		}
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case '[':
		tok = newToken(LBRACK, l.ch)
	case ']':
		tok = newToken(RBRACK, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = Lookup(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
	}

	return l.input[position:l.position]
}

func newToken(token Token, chars ...byte) Lexeme {
	literal := ""
	for _, ch := range chars {
		literal += string(ch)
	}
	return Lexeme{
		Type:    token,
		Literal: literal,
		Name:    token.String(),
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
