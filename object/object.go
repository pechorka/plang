package object

import (
	"strconv"
)

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	BOOLEAN_OBJ ObjectType = "BOOLEAN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string {
	return strconv.FormatBool(i.Value)
}

func (i *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
