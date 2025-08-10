package ast

import (
	"fmt"
)

type BoolExpr struct {
	Value *Value
}

func NewBoolExpr(value bool) *BoolExpr {
	return &BoolExpr{
		Value: NewValue(value),
	}
}

func (n *BoolExpr) Eval() *Value {
	return n.Value
}

func (n *BoolExpr) Print() string {
	return fmt.Sprintf("%t", n.Value.Bool())
}
