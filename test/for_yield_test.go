package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"testing"
)

// Helper function to compare numeric values
func compareNumbers(a, b any) bool {
	var aVal, bVal float64
	
	switch v := a.(type) {
	case float64:
		aVal = v
	case int64:
		aVal = float64(v)
	case int:
		aVal = float64(v)
	default:
		return false
	}
	
	switch v := b.(type) {
	case float64:
		bVal = v
	case int64:
		bVal = float64(v)
	case int:
		bVal = float64(v)
	default:
		return false
	}
	
	return aVal == bVal
}

func TestForYieldBasic(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `
		const arr = for let i = 0; i < 5; i++ {
			yield i
		}
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("arr")
	if !ok {
		t.Errorf("expected arr to be defined")
		return
	}

	result := val.Any().([]any)
	
	if len(result) != 5 {
		t.Errorf("expected 5 elements, got %d", len(result))
		return
	}

	for i := 0; i < 5; i++ {
		if !compareNumbers(result[i], i) {
			t.Errorf("expected element %d to be %d, got %v", i, i, result[i])
		}
	}
}

func TestForYieldWithCondition(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `
		const arr = for let i = 0; i < 10; i++ {
			if i > 2 && i < 5 {
				yield i
			}
		}
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("arr")
	if !ok {
		t.Errorf("expected arr to be defined")
		return
	}

	result := val.Any().([]any)
	
	if len(result) != 2 {
		t.Errorf("expected 2 elements, got %d", len(result))
		return
	}

	if !compareNumbers(result[0], 3) || !compareNumbers(result[1], 4) {
		t.Errorf("expected [3, 4], got %v", result)
	}
}

func TestForYieldEmpty(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `
		const empty = for let i = 0; i < 5; i++ {
			if i > 10 {
				yield i
			}
		}
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("empty")
	if !ok {
		t.Errorf("expected empty to be defined")
		return
	}

	result := val.Any().([]any)

	if len(result) != 0 {
		t.Errorf("expected empty array, got %v", result)
	}
}

func TestForYieldWithBreak(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `
		const arr = for let i = 0; i < 10; i++ {
			if i > 4 {
				break
			}
			yield i
		}
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("arr")
	if !ok {
		t.Errorf("expected arr to be defined")
		return
	}

	result := val.Any().([]any)
	
	if len(result) != 5 {
		t.Errorf("expected 5 elements, got %d", len(result))
		return
	}

	for i := 0; i < 5; i++ {
		if !compareNumbers(result[i], i) {
			t.Errorf("expected element %d to be %d, got %v", i, i, result[i])
		}
	}
}

func TestForYieldExpression(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `
		const squares = for let i = 1; i <= 4; i++ {
			yield i * i
		}
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("squares")
	if !ok {
		t.Errorf("expected squares to be defined")
		return
	}

	result := val.Any().([]any)
	
	if len(result) != 4 {
		t.Errorf("expected 4 elements, got %d", len(result))
		return
	}

	expectedSquares := []int{1, 4, 9, 16}
	for i, expected := range expectedSquares {
		if !compareNumbers(result[i], expected) {
			t.Errorf("expected element %d to be %d, got %v", i, expected, result[i])
		}
	}
}