package object

import (
	"strconv"
)

type Type string

const (
	INTEGER_OBJ Type = "INTEGER"
	BOOLEAN_OBJ Type = "BOOLEAN"
	NULL_OBJ    Type = "NULL"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i *Integer) Type() Type {
	return INTEGER_OBJ
}

type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string {
	return strconv.FormatBool(i.Value)
}

func (i *Boolean) Type() Type {
	return BOOLEAN_OBJ
}

type Null struct {
}

func (i *Null) Inspect() string {
	return "null"
}

func (i *Null) Type() Type {
	return NULL_OBJ
}
