package ast

import "foo_lang/scope"

type ConstExpr struct {
	name string
	expr Expr
}

func NewConstExpr(name string, expr Expr) *ConstExpr {
	c := &ConstExpr{
		name: name,
		expr: expr,
	}

	c.define()
	return c
}

func (n *ConstExpr) Eval() *Value {
	return nil
}

func (n *ConstExpr) define() {
	if scope.GlobalScope.Has(n.name) {
		panic("constant " + n.name + " is already defined")
	}

	if n.expr == nil {
		return
	}

	val := n.expr.Eval()
	val.SetConst(true)
	scope.GlobalScope.Set(n.name, val)
}
