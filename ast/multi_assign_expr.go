package ast

import "foo_lang/scope"

// MultiAssignExpr represents multiple assignment: let a, b = func()
type MultiAssignExpr struct {
	Names []string // Variable names to assign to
	Expr  Expr     // Expression that returns multiple values
}

func NewMultiAssignExpr(names []string, expr Expr) *MultiAssignExpr {
	return &MultiAssignExpr{
		Names: names,
		Expr:  expr,
	}
}

func (m *MultiAssignExpr) Eval() *Value {
	// Evaluate the expression that should return multiple values
	result := m.Expr.Eval()
	
	if result == nil {
		panic("expression returned nil, expected multiple values")
	}
	
	// Check if result contains multiple values (slice of *Value)
	if values, ok := result.Any().([]*Value); ok {
		// Assign each value to corresponding variable
		for i, name := range m.Names {
			if i < len(values) {
				scope.GlobalScope.Set(name, values[i])
			} else {
				// If fewer values than names, assign nil
				scope.GlobalScope.Set(name, NewValue(nil))
			}
		}
	} else {
		// Single value case - assign to first variable, rest get nil
		scope.GlobalScope.Set(m.Names[0], result)
		for i := 1; i < len(m.Names); i++ {
			scope.GlobalScope.Set(m.Names[i], NewValue(nil))
		}
	}
	
	return nil
}