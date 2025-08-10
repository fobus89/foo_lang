package ast

type YieldExpr struct {
	Expr Expr
}

func NewYieldExpr(expr Expr) *YieldExpr {
	return &YieldExpr{Expr: expr}
}

func (r *YieldExpr) Eval() *Value {
	if r.Expr == nil {
		return nil
	}
	val := r.Expr.Eval()
	// Mark the value as a yield value
	result := NewValue(val.Any())
	result.SetYield(true)
	return result
}
