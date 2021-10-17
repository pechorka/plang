package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch n := node.(type) {
	case *Program:
		for i, statement := range n.Statements {
			n.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *BlockStatement:
		for i := range n.Statements {
			n.Statements[i], _ = Modify(n.Statements[i], modifier).(Statement)
		}
	case *ExpressionStatement:
		n.Expression, _ = Modify(n.Expression, modifier).(Expression)
	case *InfixExpression:
		n.Left, _ = Modify(n.Left, modifier).(Expression)
		n.Right, _ = Modify(n.Right, modifier).(Expression)
	case *PrefixExpression:
		n.Right, _ = Modify(n.Right, modifier).(Expression)
	case *IndexExpression:
		n.Left, _ = Modify(n.Left, modifier).(Expression)
		n.Index, _ = Modify(n.Index, modifier).(Expression)
	case *IfExpression:
		n.Condition, _ = Modify(n.Condition, modifier).(Expression)
		n.Then, _ = Modify(n.Then, modifier).(*BlockStatement)
		if n.Else != nil {
			n.Else, _ = Modify(n.Else, modifier).(*BlockStatement)
		}
	case *ReturnStatement:
		n.Value, _ = Modify(n.Value, modifier).(Expression)
	case *LetStatement:
		n.Value, _ = Modify(n.Value, modifier).(Expression)
	case *FnExpression:
		for i := range n.Params {
			n.Params[i], _ = Modify(n.Params[i], modifier).(*Identifier)
		}
		n.Body, _ = Modify(n.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i := range n.Elements {
			n.Elements[i], _ = Modify(n.Elements[i], modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for k, v := range n.Pairs {
			key, _ := Modify(k, modifier).(Expression)
			val, _ := Modify(v, modifier).(Expression)
			newPairs[key] = val
		}
		n.Pairs = newPairs
	}
	return modifier(node)
}
