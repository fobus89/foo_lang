package ast

type UnaryOpExpr struct {
	Op    rune
	Count int
	Expr  Expr
}

func NewUnaryOpExpr(op rune, expr Expr, count int) *UnaryOpExpr {
	return &UnaryOpExpr{Op: op, Expr: expr, Count: count}
}

func (u *UnaryOpExpr) Eval() *Value {
	switch u.Op {
	case '-':
		return NewValue(-u.Expr.Eval().Float64())
	case '!':
		if u.Count%2 == 1 {
			return NewValue(!u.Expr.Eval().Bool())
		}
		return NewValue(u.Expr.Eval().Bool())
	default:
		return u.Expr.Eval()
	}
}
