package ast

type ConditionalExpr struct {
	Condition  Expr
	ThenBranch Expr
	ElseBranch Expr
}

func NewConditionalExpr(cond, thenExpr, elseExpr Expr) *ConditionalExpr {
	return &ConditionalExpr{
		Condition:  cond,
		ThenBranch: thenExpr,
		ElseBranch: elseExpr,
	}
}

func (c *ConditionalExpr) Eval() *Value {
	condVal := c.Condition.Eval()
	if condVal.Bool() {
		return c.ThenBranch.Eval()
	}
	return c.ElseBranch.Eval()
}
