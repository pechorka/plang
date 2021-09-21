package lexer

import (
	"bufio"
	"io"

	"github.com/pechorka/plang/token"
)

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
	c, _, err := l.r.ReadRune()
	if err != nil {
		l.err = err
		return false
	}

	var tt token.Type
	switch c {
	case '=':
		tt = token.ASSIGN
	case '+':
		tt = token.PLUS
	case ',':
		tt = token.COMMA
	case ';':
		tt = token.SEMICOLON
	case '(':
		tt = token.LPAREN
	case ')':
		tt = token.RPAREN
	case '{':
		tt = token.LBRACE
	case '}':
		tt = token.RBRACE
	}

	l.t = token.Token{
		Type:    tt,
		Literal: string(c),
	}

	return true
}

func (l *Lexer) Token() token.Token {
	return l.t
}

func (l *Lexer) Err() error {
	return l.err
}
