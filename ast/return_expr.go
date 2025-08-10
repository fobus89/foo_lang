package ast

type ReturnExpr struct {
	Expr Expr
}

func NewReturnExpr(expr Expr) *ReturnExpr {
	return &ReturnExpr{Expr: expr}
}

func (r *ReturnExpr) Eval() *Value {
	if r.Expr == nil {
		result := NewValue(nil)
		result.SetReturn(true)
		return result
	}
	result := r.Expr.Eval()
	if result == nil {
		result = NewValue(nil)
	}
	result.SetReturn(true)
	return result
}
