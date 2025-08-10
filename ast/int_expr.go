package ast

type IntExpr struct {
	Value *Value
}

func NewInt64Expr(value int64) *IntExpr {
	return &IntExpr{Value: NewValue(value)}
}

func (n *IntExpr) Eval() *Value {
	return n.Value
}
