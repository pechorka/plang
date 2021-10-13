package evaluator

import "github.com/pechorka/plang/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(`wrong number of arguments. got=%d, want=%d`,
					len(args), 1)
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		}},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(`wrong number of arguments. got=%d, want=%d`,
					len(args), 1)
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			if len(arr.Elements) == 0 {
				return NULL
			}
			return arr.Elements[0]
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(`wrong number of arguments. got=%d, want=%d`,
					len(args), 1)
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
			if len(arr.Elements) == 0 {
				return NULL
			}
			return arr.Elements[len(arr.Elements)-1]
		},
	},
}
