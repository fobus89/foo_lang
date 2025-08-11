package ast

import "foo_lang/value"

type Value = value.Value

var NewValue = value.NewValue

type Expr interface {
	Eval() *Value
}
