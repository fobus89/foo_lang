package ast

type BreakExpr struct {
	Expr Expr
}

func NewBreakExpr(expr Expr) *BreakExpr {
	return &BreakExpr{Expr: expr}
}

func (r *BreakExpr) Eval() *Value {
	// Break always returns nil but marks the value as a break
	result := NewValue(nil)
	result.SetBreak(true)
	return result
}
