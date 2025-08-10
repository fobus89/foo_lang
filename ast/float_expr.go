package ast

type FloatExpr struct {
	Value *Value
}

func NewFloat64Expr(value float64) *FloatExpr {
	return &FloatExpr{Value: NewValue(value)}
}

func (n *FloatExpr) Eval() *Value {
	return n.Value
}
