package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pechorka/plang/token"
)

type Lexer struct {
	r           bufio.Reader
	currentRune rune
	nextRune    rune
}

func New(r io.Reader) *Lexer {
	l := &Lexer{
		r: *bufio.NewReader(r),
	}
	l.readRune()
	l.readRune()
	return l
}

func NewFromString(in string) *Lexer {
	return New(strings.NewReader(in))
}

func (l *Lexer) Next() (tok token.Token) {
	l.skipWhitespace()

	switch l.currentRune {
	case '=':
		switch l.nextRune {
		case '=':
			tok.Type = token.EQ
			tok.Literal = "=="
			l.readRune() // skip picked rune
		default:
			tok = l.newToken(token.ASSIGN)
		}
	case '+':
		tok = l.newToken(token.PLUS)
	case '-':
		tok = l.newToken(token.MINUS)
	case ',':
		tok = l.newToken(token.COMMA)
	case ';':
		tok = l.newToken(token.SEMICOLON)
	case '(':
		tok = l.newToken(token.LPAREN)
	case ')':
		tok = l.newToken(token.RPAREN)
	case '{':
		tok = l.newToken(token.LBRACE)
	case '}':
		tok = l.newToken(token.RBRACE)
	case '[':
		tok = l.newToken(token.LBRACKET)
	case ']':
		tok = l.newToken(token.RBRACKET)
	case '!':
		switch l.nextRune {
		case '=':
			tok.Type = token.NOT_EQ
			tok.Literal = "!="
			l.readRune()
		default:
			tok = l.newToken(token.BANG)
		}
	case '*':
		tok = l.newToken(token.ASTERISK)
	case '/':
		tok = l.newToken(token.SLASH)
	case '<':
		tok = l.newToken(token.LT)
	case '>':
		tok = l.newToken(token.GT)
	case '"':
		l.readRune() // skip quote
		tok = l.readString()
	case 0:
		tok.Type = token.EOF
	case utf8.RuneError:
		tok = l.newToken(token.INVALID)
	default:
		return l.multiRuneToken() // return early to avoid l.readRune()
	}

	l.readRune()

	return tok
}

func (l *Lexer) readRune() {
	l.currentRune = l.nextRune
	var err error
	l.nextRune, _, err = l.r.ReadRune() // TODO handle error
	if err == io.EOF {
		l.nextRune = 0
	}
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.currentRune) {
		l.readRune()
	}
}

func (l *Lexer) newToken(tt token.Type) token.Token {
	return token.Token{
		Type:    tt,
		Literal: string(l.currentRune),
	}
}

func (l *Lexer) multiRuneToken() token.Token {
	switch {
	case isLetter(l.currentRune):
		return l.readIdent()
	case unicode.IsDigit(l.currentRune):
		return l.readNumber()
	default:
		return l.newToken(token.INVALID)
	}
}

func (l *Lexer) readIdent() (tok token.Token) {
	var buf strings.Builder
	for isLetter(l.currentRune) {
		buf.WriteRune(l.currentRune)
		l.readRune()
	}
	tok.Literal = buf.String()
	tok.Type = lookupIdentType(tok.Literal)
	return tok
}

func (l *Lexer) readNumber() (tok token.Token) {
	var buf strings.Builder
	for unicode.IsDigit(l.currentRune) {
		buf.WriteRune(l.currentRune)
		l.readRune()
	}
	tok.Literal = buf.String()
	tok.Type = token.INT
	return tok
}

func (l *Lexer) readString() (tok token.Token) {
	var buf strings.Builder
	for l.currentRune != '"' && l.currentRune != 0 {
		buf.WriteRune(l.currentRune)
		l.readRune()
	}
	tok.Literal = buf.String()
	tok.Type = token.STRING
	return tok
}

func lookupIdentType(ident string) token.Type {
	switch ident {
	case "fn":
		return token.FUNCTION
	case "let":
		return token.LET
	case "if":
		return token.IF
	case "else":
		return token.ELSE
	case "true":
		return token.TRUE
	case "false":
		return token.FALSE
	case "return":
		return token.RETURN
	default:
		return token.IDENT
	}
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}
