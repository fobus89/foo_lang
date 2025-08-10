package ast

type VarExpr struct {
	name string
	expr Expr
}

func NewVarExpr(name string, expr Expr) *VarExpr {
	return &VarExpr{
		name: name,
		expr: expr,
	}
}

func (n *VarExpr) Eval() *Value {

	val, ok := Container[n.name]
	{
		if !ok {
			panic("variable " + n.name + " is not defined")
		}
		if n.expr != nil && val.IsConst() {
			panic("constant " + n.name + " is already defined")
		}
	}

	if n.expr != nil {
		tmp := n.expr.Eval()
		Container[n.name] = tmp
		return tmp
	}

	return val
}
