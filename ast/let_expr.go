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

// GetName возвращает имя переменной
func (n *LetExpr) GetName() string {
	return n.name
}

// GetExpr возвращает выражение
func (n *LetExpr) GetExpr() Expr {
	return n.expr
}
