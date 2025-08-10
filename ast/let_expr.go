package ast

import "foo_lang/scope"

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
	if scope.GlobalScope.Has(n.name) {
		panic("variable " + n.name + " is already defined")
	}

	val := n.expr.Eval()
	scope.GlobalScope.Set(n.name, val)

	return nil
}
