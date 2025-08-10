package ast

import "foo_lang/scope"

type VarExpr struct {
	Name string
	expr Expr
}

func NewVarExpr(name string, expr Expr) *VarExpr {
	return &VarExpr{
		Name: name,
		expr: expr,
	}
}

func (n *VarExpr) Eval() *Value {
	val, ok := scope.GlobalScope.Get(n.Name)
	if !ok {
		panic("variable " + n.Name + " is not defined")
	}
	
	if n.expr != nil && val.IsConst() {
		panic("cannot assign to constant " + n.Name)
	}

	if n.expr != nil {
		// Обновление существующей переменной
		tmp := n.expr.Eval()
		scope.GlobalScope.Update(n.Name, tmp)
		return tmp
	}

	return val
}
