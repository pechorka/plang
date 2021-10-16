package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch n := node.(type) {
	case *Program:
		for i, statement := range n.Statements {
			n.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		n.Expression, _ = Modify(n.Expression, modifier).(Expression)
	}
	return modifier(node)
}
