package ast

type BreakExpr struct {
	Expr Expr
}

func NewBreakExpr(expr Expr) *BreakExpr {
	return &BreakExpr{Expr: expr}
}

func (r *BreakExpr) Eval() *Value {
	if r.Expr == nil {
		return nil
	}
	return r.Expr.Eval()
}
