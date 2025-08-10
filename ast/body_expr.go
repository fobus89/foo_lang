package ast

type BodyExpr struct {
	Statments []Expr
}

func NewBodyStatment(stratments []Expr) *BodyExpr {
	return &BodyExpr{Statments: stratments}
}

func (b *BodyExpr) Eval() *Value {
	var result *Value

	for _, stm := range b.Statments {
		result = stm.Eval()
	}

	return result
}
