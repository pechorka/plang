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
		return evalProgram(n)
	case *ast.BlockStatement:
		return evalBlockStatement(n)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.PrefixExpression:
		return evalPrefixExpression(n)
	case *ast.InfixExpression:
		return evalInfixExpression(n)
	case *ast.IfExpression:
		return evalIfExpression(n)
	case *ast.ReturnStatement:
		return evalReturn(n)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return boolToBooleanObject(n.Value)
	}

	return NULL
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
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

func evalInfixExpression(infix *ast.InfixExpression) object.Object {
	left := Eval(infix.Left)
	right := Eval(infix.Right)
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(infix.Operator, left, right)
	case infix.Operator == "==":
		return boolToBooleanObject(left == right)
	case infix.Operator == "!=":
		return boolToBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalIfExpression(ifExpr *ast.IfExpression) object.Object {
	cond := Eval(ifExpr.Condition)
	if isTruthy(cond) {
		return Eval(ifExpr.Then)
	}
	if ifExpr.Else != nil {
		return Eval(ifExpr.Else)
	}

	return NULL
}

func evalReturn(retExpr *ast.ReturnStatement) object.Object {
	val := Eval(retExpr.Value)
	return &object.ReturnValue{Value: val}
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

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	// comparison
	case "<":
		return boolToBooleanObject(leftValue < rightValue)
	case ">":
		return boolToBooleanObject(leftValue > rightValue)
	case "==":
		return boolToBooleanObject(leftValue == rightValue)
	case "!=":
		return boolToBooleanObject(leftValue != rightValue)
	default:
		return NULL
	}
}

func boolToBooleanObject(b bool) object.Object {
	if b {
		return TRUE
	}

	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
