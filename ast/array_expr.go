package ast

type ArrayExpr struct {
	Elements []Expr
}

func NewArrayExpr(elements []Expr) *ArrayExpr {
	return &ArrayExpr{
		Elements: elements,
	}
}

func (a *ArrayExpr) Eval() *Value {
	var values []any
	
	for _, element := range a.Elements {
		values = append(values, element.Eval().Any())
	}
	
	return NewValue(values)
}