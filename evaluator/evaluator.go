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
	case *ast.PrefixExpression:
		return evalPrefixExpression(n)
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

func evalPrefixExpression(prefix *ast.PrefixExpression) object.Object {
	right := Eval(prefix.Right)
	switch prefix.Operator {
	case "!":
		return evalBangExpression(right)
	case "-":
		return evalMinusExpression(right)
	default:
		return NULL
	}
}

func evalBangExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return NULL
	default:
		return FALSE
	}
}

func evalMinusExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func boolToBooleanObject(b bool) object.Object {
	if b {
		return TRUE
	}

	return FALSE
}
