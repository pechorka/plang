package lexer

import (
	"io"
	"strings"
	"testing"

	"github.com/pechorka/plang/token"
)

func TestNextToken(t *testing.T) {

	input := `=+(){},;`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	l := New(strings.NewReader(input))
	for i, res := range tests {
		isNext := l.Next()
		if !isNext {
			t.Fatalf("next token should exist")
		}
		tok := l.Token()
		if tok.Type != res.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, res.expectedType, tok.Type)
		}
		if tok.Literal != res.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, res.expectedLiteral, tok.Literal)
		}
	}

	isNext := l.Next()
	if isNext {
		t.Fatalf("should be no next token")
	}

	err := l.Err()
	switch err {
	case nil:
		t.Fatalf("should be io.EOF")
	case io.EOF:
	default:
		t.Fatalf("enexpected error %v", err)
	}
}
