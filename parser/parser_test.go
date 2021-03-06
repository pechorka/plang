package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/pechorka/plang/ast"
	"github.com/pechorka/plang/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for i, tt := range tests {
		l := lexer.NewFromString(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrorsForTest(t, p, i)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
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

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for i, tt := range tests {
		l := lexer.NewFromString(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrorsForTest(t, p, i)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}

		if !testLiteralExpression(t, returnStmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	stmt := getExpressionStmt(t, input)

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	stmt := getExpressionStmt(t, input)

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"foobar";`

	stmt := getExpressionStmt(t, input)

	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "foobar" {
		t.Errorf("literal.Value not %s. got=%s", "foobar", literal.Value)
	}
	if literal.TokenLiteral() != "foobar" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "foobar",
			literal.TokenLiteral())
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input            string
		expectedValue    bool
		expectedValueStr string
	}{
		{"true;", true, "true"},
		{"false;", false, "false"},
	}

	for _, tt := range tests {
		stmt := getExpressionStmt(t, tt.input)

		b, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if b.Value != tt.expectedValue {
			t.Errorf("literal.Value not %v. got=%v", tt.expectedValue, b.Value)
		}
		if b.TokenLiteral() != tt.expectedValueStr {
			t.Errorf("literal.TokenLiteral not %s. got=%s", tt.expectedValueStr,
				b.TokenLiteral())
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		stmt := getExpressionStmt(t, tt.input)

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		stmt := getExpressionStmt(t, tt.input)

		if !testInfixExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for i, tt := range tests {
		l := lexer.NewFromString(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrorsForTest(t, p, i)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x > y) { x }"

	stmt := getExpressionStmt(t, input)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if exp.TokenLiteral() != "if" {
		t.Errorf("exp.TokenLiteral not %s. got=%s", "if",
			exp.TokenLiteral())
	}

	if !testInfixExpression(t, exp.Condition, "x", ">", "y") {
		return
	}

	if len(exp.Then.Statements) != 1 {
		t.Fatalf("exp.Then has not enough statements. got=%d",
			len(exp.Then.Statements))
	}

	thenExprStmt, ok := exp.Then.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Then.Statements[0] not *ast.ExpressionStatement. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, thenExprStmt.Expression, "x") {
		return
	}

	if exp.Else != nil {
		t.Errorf("exp.Else was not nil, got=%T", exp.Else)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x > y) { x } else { y }"

	stmt := getExpressionStmt(t, input)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if exp.TokenLiteral() != "if" {
		t.Errorf("exp.TokenLiteral not %s. got=%s", "if",
			exp.TokenLiteral())
	}

	if !testInfixExpression(t, exp.Condition, "x", ">", "y") {
		return
	}

	if len(exp.Then.Statements) != 1 {
		t.Fatalf("exp.Then has not enough statements. got=%d",
			len(exp.Then.Statements))
	}

	thenExprStmt, ok := exp.Then.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Then.Statements[0] not *ast.ExpressionStatement. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, thenExprStmt.Expression, "x") {
		return
	}

	if exp.Else == nil {
		t.Errorf("exp.Else was nil")
	}

	if len(exp.Else.Statements) != 1 {
		t.Fatalf("exp.Else has not enough statements. got=%d",
			len(exp.Else.Statements))
	}

	elseExprStmt, ok := exp.Else.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Else.Statements[0] not *ast.ExpressionStatement. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, elseExprStmt.Expression, "y") {
		return
	}
}

func TestFnExpression(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	stmt := getExpressionStmt(t, input)

	exp, ok := stmt.Expression.(*ast.FnExpression)
	if !ok {
		t.Fatalf("exp not *ast.FnExpression. got=%T", stmt.Expression)
	}

	if exp.TokenLiteral() != "fn" {
		t.Errorf("exp.TokenLiteral not %s. got=%s", "fn",
			exp.TokenLiteral())
	}

	if len(exp.Params) != 2 {
		t.Fatalf("exp.Params has not %d identifiers. got=%d", 2, len(exp.Params))
	}

	if !testIdentifier(t, exp.Params[0], "x") {
		return
	}
	if !testIdentifier(t, exp.Params[1], "y") {
		return
	}

	if len(exp.Body.Statements) != 1 {
		t.Fatalf("exp.Then has not enough statements. got=%d",
			len(exp.Body.Statements))
	}

	sumExprStmt, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Body.Statements[0] not *ast.ExpressionStatement. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, sumExprStmt.Expression, "x", "+", "y") {
		return
	}
}

func TestFnParamsParsing(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedParams []string
	}{
		{name: "empty params", input: "fn() {};", expectedParams: []string{}},
		{name: "single param", input: "fn(x) {};", expectedParams: []string{"x"}},
		{name: "multiple params", input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := getExpressionStmt(t, tt.input)

			exp, ok := stmt.Expression.(*ast.FnExpression)
			if !ok {
				t.Fatalf("exp not *ast.FnExpression. got=%T", stmt.Expression)
			}

			if len(tt.expectedParams) != len(exp.Params) {
				t.Fatalf("wrong param list len: expected %d, got %d",
					len(tt.expectedParams), len(exp.Params))
			}

			for i, param := range exp.Params {
				if !testIdentifier(t, param, tt.expectedParams[i]) {
					return
				}
			}
		})
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	stmt := getExpressionStmt(t, input)

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp not *ast.CallExpression. got=%T", stmt.Expression)
	}

	if exp.TokenLiteral() != "(" {
		t.Errorf("exp.TokenLiteral not %q. got=%q", "(",
			exp.TokenLiteral())
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("exp.Arguments has not %d identifiers. got=%d", 3, len(exp.Arguments))
	}

	if !testLiteralExpression(t, exp.Arguments[0], 1) {
		return
	}

	if !testInfixExpression(t, exp.Arguments[1], 2, "*", 3) {
		return
	}

	if !testInfixExpression(t, exp.Arguments[2], 4, "+", 5) {
		return
	}
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		stmt := getExpressionStmt(t, tt.input)

		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	stmt := getExpressionStmt(t, input)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	stmt := getExpressionStmt(t, input)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	stmt := getExpressionStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`
	stmt := getExpressionStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not ast.IntegerLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`
	stmt := getExpressionStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 2 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not ast.Boolean. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"
	stmt := getExpressionStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	stmt := getExpressionStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
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

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, il ast.Expression, value bool) bool {
	boolean, ok := il.(*ast.Boolean)
	if !ok {
		t.Errorf("il not *ast.Boolean. got=%T", il)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %v. got=%v", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != strconv.FormatBool(value) {
		t.Errorf("integ.TokenLiteral not %v. got=%s", value,
			boolean.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func getExpressionStmt(t *testing.T, input string) *ast.ExpressionStatement {
	l := lexer.NewFromString(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	return stmt
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
}

func checkParserErrorsForTest(t *testing.T, p *Parser, i int) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors in test %d", len(errors), i)
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
