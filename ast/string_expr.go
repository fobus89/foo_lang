package ast

import "strings"

type StringFormatExpr struct {
	exprs []Expr
}

func NewStringFormatExpr(exprs []Expr) *StringFormatExpr {
	return &StringFormatExpr{
		exprs: exprs,
	}
}

func (s *StringFormatExpr) Eval() *Value {
	var result strings.Builder

	for _, part := range s.exprs {
		if part == nil {
			continue
		}
		res := part.Eval()
		if res == nil {
			continue
		}

		result.WriteString(part.Eval().String())
	}

	return NewValue(result.String())
}
