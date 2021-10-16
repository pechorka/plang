package ast

import (
	"reflect"
	"testing"

	"github.com/pechorka/plang/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestModify(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }
	turnOneIntoTwo := func(node Node) Node {
		i, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}
		if i.Value != 1 {
			return node
		}
		i.Value = 2
		return i
	}

	tests := []struct {
		input    Node
		expected Node
	}{
		{
			one(),
			two(),
		},

		{
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: one()},
				},
			},
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: two()},
				},
			},
		},
	}

	for _, tt := range tests {
		modified := Modify(tt.input, turnOneIntoTwo)
		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. got=%#v, want=%#v",
				modified, tt.expected)
		}
	}
}
