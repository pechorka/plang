package evaluator

import "github.com/pechorka/plang/object"

var builtins = map[string]*object.Builtin{"len": {
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError(`wrong number of arguments. got=%d, want=%d`,
				len(args), 1)
		}
		str, ok := args[0].(*object.String)
		if !ok {
			return newError("argument to `len` not supported, got %s", args[0].Type())
		}
		return &object.Integer{Value: int64(len(str.Value))}
	}},
}
