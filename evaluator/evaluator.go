package evaluator

import (
	"github.com/pechorka/plang/ast"
	"github.com/pechorka/plang/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalStatements(n.Statements)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: n.Value,
		}
	case *ast.Boolean:
		return boolToBooleanObject(n.Value)
	}

	return NULL
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}

func boolToBooleanObject(b bool) object.Object {
	if b {
		return TRUE
	}

	return FALSE
}
