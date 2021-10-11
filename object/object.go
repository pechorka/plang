package object

import (
	"strconv"
)

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
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
