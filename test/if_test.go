package test

import (
	"foo_lang/parser"
	"testing"
)

func TestIf(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `if true {1} else {2}`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()

	for _, expr := range exprs {
		value := expr.Eval()

		if value.String() != "1" {
			t.Errorf("expected 1, got %s", value.String())
		}
	}

}

func TestIfElse(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `if false {1} else {2}`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()

	for _, expr := range exprs {
		value := expr.Eval()

		if value.String() != "2" {
			t.Errorf("expected 2, got %s", value.String())
		}
	}

}

func TestIfExpression(t *testing.T) {
	// Clear scope and create new one
	InitTestEnvironment()

	const code = `if 1+2/2 {1} else {2}`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()

	for _, expr := range exprs {
		value := expr.Eval()

		if value.String() != "1" {
			t.Errorf("expected 2, got %s", value.String())
		}
	}

}
