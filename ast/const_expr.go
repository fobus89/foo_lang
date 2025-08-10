package ast

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

	if _, ok := Container[n.name]; ok {
		panic("constant " + n.name + " is already defined")
	}

	if n.expr == nil {
		return
	}

	val := n.expr.Eval()

	Container[n.name] = val
}
