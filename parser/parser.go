package parser

import (
	"fmt"

	"github.com/pechorka/plang/ast"
	"github.com/pechorka/plang/lexer"
	"github.com/pechorka/plang/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	nextToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// fill cur and next token
	p.readToken()
	p.readToken()

	return p
}

func (p *Parser) Parse() *ast.Program {
	var prog ast.Program
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.readToken()
	}
	return &prog
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) readToken() {
	p.curToken = p.nextToken
	p.nextToken = p.l.Next()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	}

	return nil
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := ast.LetStatement{
		Token: p.curToken,
	}

	if !p.isNextToken(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.isNextToken(token.ASSIGN) {
		return nil
	}

	// TODO: add value parsing
	if !p.skipUntilSemicolon() {
		return nil
	}

	return &stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := ast.ReturnStatement{
		Token: p.curToken,
	}

	if !p.skipUntilSemicolon() {
		return nil
	}

	return &stmt
}

func (p *Parser) isNextToken(tt token.Type) bool {
	if p.nextToken.Type == tt {
		p.readToken()
		return true
	}
	p.appendErrorf("expect next token to be %s, got %s instead", tt, p.nextToken.Type)
	return false
}

func (p *Parser) skipUntilSemicolon() bool {
	for p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.EOF {
		p.readToken()
	}
	if p.curToken.Type == token.EOF {
		p.appendErrorf("no semicolon after statement")
		return false
	}

	return true
}

func (p *Parser) appendErrorf(text string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(text, args...))
}
