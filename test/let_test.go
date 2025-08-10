package test

import (
	"foo_lang/parser"
	"foo_lang/value"
	"testing"
)

var Container = value.Container

func TestLet(t *testing.T) {
	const code = `
		let x = 1
		let y = 2
	`

	exprs := parser.NewParser(code).Parse()

	for _, expr := range exprs {
		expr.Eval()
	}

	valX := Container["x"]
	valY := Container["y"]

	if valX == nil || valY == nil {
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
	const code = `
		let x = 1
		let y = 2
		let z = x + y
	`

	exprs := parser.NewParser(code).Parse()

	for _, expr := range exprs {
		expr.Eval()
	}

	valZ := Container["z"]

	if valZ == nil {
		t.Errorf("expected z to be defined")
		return
	}

	if valZ.Int64() != 3 {
		t.Errorf("expected 3, got %d", valZ.Int64())
	}
}
