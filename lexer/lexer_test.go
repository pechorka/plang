package lexer

import (
	"io"
	"strings"
	"testing"

	"github.com/pechorka/plang/token"
)

type lexerResult struct {
	expectedType    token.Type
	expectedLiteral string
}

func TestNext_basic(t *testing.T) {
	input := `=+(){},;`
	tests := []lexerResult{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	testLexer(t, input, tests)
}

func TestNext_full(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	   let add = fn(x, y) {
		 x + y;
	};
	   let result = add(five, ten);
	   !-/*5;
	   5 < 10 > 5;
`
	tests := []lexerResult{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
	}

	testLexer(t, input, tests)
}

func TestNext_invalid(t *testing.T) {
	input := `@let`
	l := New(strings.NewReader(input))

	isNext := l.Next()
	if isNext {
		t.Fatalf("next should be false")
	}

	err := l.Err()

	if err != ErrInvalidToken {
		t.Fatalf("error should %v, not the %v", ErrInvalidToken, err)
	}
}

func testLexer(t *testing.T, input string, tests []lexerResult) {
	t.Helper()
	l := New(strings.NewReader(input))
	for i, res := range tests {
		isNext := l.Next()
		if !isNext {
			t.Fatalf("tests[%d] - next token should exist, but got err %v", i, l.Err())
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
