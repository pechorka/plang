package parser

import (
	"github.com/pechorka/plang/ast"
	"github.com/pechorka/plang/lexer"
	"github.com/pechorka/plang/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	nextToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// fill cur and next token
	p.readToken()
	p.readToken()

	return p
}

func (p *Parser) readToken() {
	p.curToken = p.nextToken
	p.nextToken = p.l.Next()
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
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
	for p.curToken.Type != token.SEMICOLON {
		p.readToken()
	}

	return &stmt
}

func (p *Parser) isNextToken(tt token.Type) bool {
	if p.nextToken.Type == tt {
		p.readToken()
		return true
	}

	return false
}
