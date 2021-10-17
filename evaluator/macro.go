package evaluator

import (
	"github.com/pechorka/plang/ast"
	"github.com/pechorka/plang/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}
	for i, statement := range program.Statements {
		if isMacroDefinition(statement) {
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		definitionIndex := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	letStatement, _ := stmt.(*ast.LetStatement)
	macroExpr, _ := letStatement.Value.(*ast.MacroExpression)

	obj := object.Macro{
		Parameters: macroExpr.Params,
		Body:       macroExpr.Body,
		Env:        env,
	}
	env.Set(letStatement.Name.Value, &obj)
}

func isMacroDefinition(stmt ast.Statement) bool {
	letStatement, ok := stmt.(*ast.LetStatement)
	if !ok {
		return false
	}
	_, ok = letStatement.Value.(*ast.MacroExpression)
	return ok
}
