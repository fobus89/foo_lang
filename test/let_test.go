package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"testing"
)

func TestLet(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `
		let x = 1
		let y = 2
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()

	for _, expr := range exprs {
		expr.Eval()
	}

	valX, okX := scope.GlobalScope.Get("x")
	valY, okY := scope.GlobalScope.Get("y")

	if !okX || !okY {
		t.Errorf("expected x and y to be defined")
		return
	}

	if valX.Int64() != 1 {
		t.Errorf("expected 1, got %d", valX.Int64())
	}

	if valY.Int64() != 2 {
		t.Errorf("expected 2, got %d", valY.Int64())
	}
}

func TestLetExpression(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `
		let x = 1
		let y = 2
		let z = x + y
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()

	for _, expr := range exprs {
		expr.Eval()
	}

	valZ, okZ := scope.GlobalScope.Get("z")

	if !okZ {
		t.Errorf("expected z to be defined")
		return
	}

	if valZ.Int64() != 3 {
		t.Errorf("expected 3, got %d", valZ.Int64())
	}
}
