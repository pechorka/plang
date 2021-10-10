package parser

import (
	"fmt"
	"strconv"

	"github.com/pechorka/plang/ast"
	"github.com/pechorka/plang/lexer"
	"github.com/pechorka/plang/token"
)

// precedences of operators;
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)
type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	nextToken      token.Token
	errors         []string
	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		prefixParseFns: make(map[token.Type]prefixParseFn),
		infixParseFns:  make(map[token.Type]infixParseFn),
	}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFnExpression)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

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

	return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() ast.Statement {
	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}

	stmt := ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: exp,
	}

	if p.nextToken.Type == token.SEMICOLON { // semicolon is optional
		p.readToken()
	}

	return &stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.appendErrorf("no prefix func for %q token type", p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for p.nextToken.Type != token.SEMICOLON && precedence < p.nextTokenPrecedence() {
		infix := p.infixParseFns[p.nextToken.Type]
		if infix == nil {
			return leftExp
		}
		p.readToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.appendErrorf("cant parse %q as 64-bit integer", p.curToken.Literal)
		return nil
	}
	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: val,
	}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	var val bool
	switch p.curToken.Literal {
	case "true":
		val = true
	case "false":
		val = false
	default:
		p.appendErrorf("cant parse %q as boolean", p.curToken.Literal)
		return nil
	}
	return &ast.Boolean{
		Token: p.curToken,
		Value: val,
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.readToken() // consume token, so parse expression will analyze next token

	right := p.parseExpression(PREFIX)
	if right == nil {
		return nil
	}
	expr.Right = right
	return &expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curTokenPrecedence()
	p.readToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.readToken()
	exp := p.parseExpression(LOWEST)
	if !p.isNextToken(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExp := ast.IfExpression{
		Token: p.curToken,
	}

	if !p.isNextToken(token.LPAREN) {
		p.appendErrorf("invalid if expression: no ( after if")
		return nil
	}

	ifExp.Condition = p.parseExpression(LOWEST)

	if !p.isNextToken(token.LBRACE) {
		p.appendErrorf("invalid if expression: no { after condition")
		return nil
	}

	ifExp.Then = p.parseBlockStatement()

	if p.curToken.Type == token.ELSE {
		if !p.isNextToken(token.LBRACE) {
			p.appendErrorf("invalid if expression: no { after else")
			return nil
		}

		ifExp.Else = p.parseBlockStatement()
	}

	return &ifExp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStmt := ast.BlockStatement{
		Token: p.curToken,
	}

	p.readToken() // consume token.LBRACE

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			blockStmt.Statements = append(blockStmt.Statements, stmt)
		}
		p.readToken()
	}

	p.readToken() // consume token.RBRACE

	return &blockStmt
}

func (p *Parser) parseFnExpression() ast.Expression {
	fnExpr := ast.FnExpression{
		Token: p.curToken,
	}

	if !p.isNextToken(token.LPAREN) {
		p.appendErrorf("invalid fn expression: no ( after fn")
		return nil
	}

	fnExpr.Params = p.parseFnParams()

	if !p.isNextToken(token.LBRACE) {
		p.appendErrorf("invalid fn expression: no { after param list")
		return nil
	}

	fnExpr.Body = p.parseBlockStatement()

	return &fnExpr
}

func (p *Parser) parseFnParams() []*ast.Identifier {
	if p.nextToken.Type == token.RPAREN { // empty param list
		p.readToken()
		return nil
	}

	p.readToken() // consume left parenthesis

	var params []*ast.Identifier
	ident := p.parseIdentifier().(*ast.Identifier)
	params = append(params, ident)
	for p.nextToken.Type == token.COMMA {
		p.readToken()
		p.readToken()
		ident := p.parseIdentifier().(*ast.Identifier)
		params = append(params, ident)
	}

	if !p.isNextToken(token.RPAREN) {
		p.appendErrorf("expected right parenthesis after fn params")
		return nil
	}

	return params
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	callExp := ast.CallExpression{
		Token:    p.curToken,
		Function: left,
	}
	callExp.Arguments = p.parseCallArguments()
	return &callExp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	if p.nextToken.Type == token.RPAREN { // empty arg list
		p.readToken()
		return nil
	}
	p.readToken() // consume left parenthesis

	var args []ast.Expression
	args = append(args, p.parseExpression(LOWEST))
	for p.nextToken.Type == token.COMMA {
		p.readToken()
		p.readToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.isNextToken(token.RPAREN) {
		p.appendErrorf("expected right parenthesis after call arguments")
		return nil
	}

	return args
}

func (p *Parser) isNextToken(tt token.Type) bool {
	if p.nextToken.Type == tt {
		p.readToken()
		return true
	}
	p.appendErrorf("expect next token to be %q, got %q instead", tt, p.nextToken.Type)
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

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) curTokenPrecedence() int {
	return p.tokenPrecedence(p.curToken)
}

func (p *Parser) nextTokenPrecedence() int {
	return p.tokenPrecedence(p.nextToken)
}

func (p *Parser) tokenPrecedence(t token.Token) int {
	if p, ok := precedences[t.Type]; ok {
		return p
	}
	return LOWEST
}
