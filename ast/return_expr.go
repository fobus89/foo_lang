package ast

type ReturnExpr struct {
	Expr Expr
}

func NewReturnExpr(expr Expr) *ReturnExpr {
	return &ReturnExpr{Expr: expr}
}

func (r *ReturnExpr) Eval() *Value {
	if r.Expr == nil {
		return nil
	}
	return r.Expr.Eval()
}
