package parser

import (
	"testing"

	"github.com/pechorka/plang/ast"
	"github.com/pechorka/plang/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
   let x = 5;
   let y = 10;
   let foobar = 838383;
   `
	l := lexer.NewFromString(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

	}

}

func TestLetStatementsError(t *testing.T) {
	tests := []struct {
		inputWithError string
	}{
		{"let"},
		{"let x "},
		{"let x = "},
	}
	for i, tt := range tests {
		l := lexer.NewFromString(tt.inputWithError)
		p := New(l)
		p.Parse()
		if len(p.Errors()) == 0 {
			t.Fatalf("test[%d]: expected error", i)
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	t.Helper()
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
