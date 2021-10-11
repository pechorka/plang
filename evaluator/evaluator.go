package evaluator

import (
	"fmt"

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
		switch res := result.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil {
			switch result.Type() {
			case object.RETURN_VALUE_OBJ, object.ERROR_OBJ:
				return result
			}
		}
	}
	return result
}

func evalPrefixExpression(prefix *ast.PrefixExpression) object.Object {
	right := Eval(prefix.Right)
	if isError(right) {
		return right
	}
	switch prefix.Operator {
	case "!":
		return evalBangExpression(right)
	case "-":
		return evalMinusExpression(right)
	default:
		return newError("unknown operator: %s%s", prefix.Operator, right.Type())
	}
}

func evalInfixExpression(infix *ast.InfixExpression) object.Object {
	left := Eval(infix.Left)
	if isError(left) {
		return left
	}
	right := Eval(infix.Right)
	if isError(right) {
		return right
	}
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(infix.Operator, left, right)
	case infix.Operator == "==":
		return boolToBooleanObject(left == right)
	case infix.Operator == "!=":
		return boolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), infix.Operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), infix.Operator, right.Type())
	}
}

func evalIfExpression(ifExpr *ast.IfExpression) object.Object {
	cond := Eval(ifExpr.Condition)
	if isError(cond) {
		return cond
	}
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
	if isError(val) {
		return val
	}
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
		return newError("unknown operator: -%s", right.Type())
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
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
