package ast

// MultiReturnExpr represents multiple return values: return a, b, c
type MultiReturnExpr struct {
	Values []Expr // Multiple expressions to return
}

func NewMultiReturnExpr(values []Expr) *MultiReturnExpr {
	return &MultiReturnExpr{Values: values}
}

func (m *MultiReturnExpr) Eval() *Value {
	// Evaluate all return values
	var results []*Value
	for _, expr := range m.Values {
		results = append(results, expr.Eval())
	}
	
	// Create a special return value that contains multiple values
	multiValue := NewValue(results)
	multiValue.SetReturn(true)
	
	return multiValue
}