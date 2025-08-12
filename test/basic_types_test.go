package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"testing"
)

func TestBasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{"Integer", "let x = 42", int64(42)},
		{"Float", "let x = 3.14", 3.14},
		{"String", `let x = "hello"`, "hello"},
		{"Boolean True", "let x = true", true},
		{"Boolean False", "let x = false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitTestEnvironment()
			
			exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
			for _, expr := range exprs {
				expr.Eval()
			}

			val, ok := scope.GlobalScope.Get("x")
			if !ok {
				t.Errorf("variable x not found")
				return
			}

			result := val.Any()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestArithmeticOperators(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int64
	}{
		{"Addition", "let x = 5 + 3", 8},
		{"Subtraction", "let x = 10 - 4", 6},
		{"Multiplication", "let x = 6 * 7", 42},
		{"Division", "let x = 15 / 3", 5},
		{"Modulo", "let x = 17 % 5", 2},
		{"Complex", "let x = 2 + 3 * 4", 14},
		{"Parentheses", "let x = (2 + 3) * 4", 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitTestEnvironment()
			
			exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
			for _, expr := range exprs {
				expr.Eval()
			}

			val, ok := scope.GlobalScope.Get("x")
			if !ok {
				t.Errorf("variable x not found")
				return
			}

			if val.Int64() != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, val.Int64())
			}
		})
	}
}

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"AND true", "let x = true && true", true},
		{"AND false", "let x = true && false", false},
		{"OR true", "let x = true || false", true},
		{"OR false", "let x = false || false", false},
		{"NOT true", "let x = !false", true},
		{"NOT false", "let x = !true", false},
		{"Complex", "let x = (true || false) && !false", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitTestEnvironment()
			
			exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
			for _, expr := range exprs {
				expr.Eval()
			}

			val, ok := scope.GlobalScope.Get("x")
			if !ok {
				t.Errorf("variable x not found")
				return
			}

			if val.Bool() != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val.Bool())
			}
		})
	}
}

func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"Equal", "let x = 5 == 5", true},
		{"Not Equal", "let x = 5 != 3", true},
		{"Greater", "let x = 7 > 5", true},
		{"Less", "let x = 3 < 8", true},
		{"Greater Equal", "let x = 5 >= 5", true},
		{"Less Equal", "let x = 4 <= 6", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitTestEnvironment()
			
			exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
			for _, expr := range exprs {
				expr.Eval()
			}

			val, ok := scope.GlobalScope.Get("x")
			if !ok {
				t.Errorf("variable x not found")
				return
			}

			if val.Bool() != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val.Bool())
			}
		})
	}
}