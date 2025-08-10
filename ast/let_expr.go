package ast

type LetExpr struct {
	name string
	expr Expr
}

func NewLetExpr(name string, expr Expr) *LetExpr {
	return &LetExpr{
		name: name,
		expr: expr,
	}
}

func (n *LetExpr) Eval() *Value {

	if _, ok := Container[n.name]; ok {
		panic("variable " + n.name + " is already defined")
	}

	Container[n.name] = n.expr.Eval()

	return nil
}
