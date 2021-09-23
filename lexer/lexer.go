package lexer

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"

	"github.com/pechorka/plang/token"
)

var ErrInvalidToken = errors.New("invalid token")

type Lexer struct {
	r   bufio.Reader
	t   token.Token
	err error
}

func New(r io.Reader) *Lexer {
	return &Lexer{
		r: *bufio.NewReader(r),
	}
}

func (l *Lexer) Next() bool {
	var (
		r   rune = ' '
		err error
	)
	for unicode.IsSpace(r) { // to skip whitespace
		r, err = l.readRune()
		if err != nil {
			l.err = err
			return false
		}
	}

	switch r {
	case '=':
		l.t = l.newToken(token.ASSIGN, r)
	case '+':
		l.t = l.newToken(token.PLUS, r)
	case '-':
		l.t = l.newToken(token.MINUS, r)
	case ',':
		l.t = l.newToken(token.COMMA, r)
	case ';':
		l.t = l.newToken(token.SEMICOLON, r)
	case '(':
		l.t = l.newToken(token.LPAREN, r)
	case ')':
		l.t = l.newToken(token.RPAREN, r)
	case '{':
		l.t = l.newToken(token.LBRACE, r)
	case '}':
		l.t = l.newToken(token.RBRACE, r)
	case '!':
		l.t = l.newToken(token.BANG, r)
	case '*':
		l.t = l.newToken(token.ASTERISK, r)
	case '/':
		l.t = l.newToken(token.SLASH, r)
	case '<':
		l.t = l.newToken(token.LT, r)
	case '>':
		l.t = l.newToken(token.GT, r)
	default:
		l.t, err = l.multiRuneToken(r)
		if err != nil {
			l.err = err
			return false
		}
	}

	return true
}

func (l *Lexer) Token() token.Token {
	return l.t
}

func (l *Lexer) Err() error {
	return l.err
}

func (l *Lexer) readRune() (rune, error) {
	r, _, err := l.r.ReadRune()
	// fmt.Printf("rune %d, strrune(%s), error %v\n", r, string(r), err)
	return r, err
}

func (l *Lexer) unreadRune() {
	l.r.UnreadRune()
}

func (l *Lexer) newToken(tt token.Type, literal rune) token.Token {
	return token.Token{
		Type:    tt,
		Literal: string(literal),
	}
}

func (l *Lexer) multiRuneToken(r rune) (t token.Token, err error) {
	switch {
	case isLetter(r):
		return l.readIdent(r)
	case unicode.IsDigit(r):
		return l.readNumber(r)
	default:
		return t, ErrInvalidToken
	}
}

func (l *Lexer) readIdent(r rune) (t token.Token, err error) {
	var buf strings.Builder
	for isLetter(r) {
		buf.WriteRune(r)
		r, err = l.readRune()
		if err != nil {
			return t, err
		}
	}
	l.unreadRune()
	t.Literal = buf.String()
	t.Type = lookupIdentType(t.Literal)
	return t, nil
}

func (l *Lexer) readNumber(r rune) (t token.Token, err error) {
	var buf strings.Builder
	for unicode.IsDigit(r) {
		buf.WriteRune(r)
		r, err = l.readRune()
		if err != nil {
			return t, err
		}
	}
	l.unreadRune()
	t.Literal = buf.String()
	t.Type = token.INT
	return t, nil
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
