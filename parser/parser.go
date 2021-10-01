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

	p.readToken()
	p.readToken()

	return p
}

func (p *Parser) readToken() {
	p.curToken = p.nextToken
	p.nextToken = p.l.Next()
}

func (p *Parser) Parse() *ast.Program {
	return nil
}
