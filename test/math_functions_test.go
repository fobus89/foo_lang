package test

import (
	"math"
	"testing"
	"foo_lang/parser"
	"foo_lang/scope"
	"foo_lang/builtin"
	"foo_lang/ast"
	"foo_lang/modules"
)

func TestBasicMathFunctions(t *testing.T) {
	// Устанавливаем глобальную parseFunc для import
	parseFunc := func(code string) []modules.Expr {
		exprs := parser.NewParser(code).ParseWithoutScopeInit()
		result := make([]modules.Expr, len(exprs))
		for i, expr := range exprs {
			result[i] = expr
		}
		return result
	}
	ast.SetGlobalParseFunc(parseFunc)
	
	// Инициализируем тестовое окружение
	InitTestEnvironment(
		builtin.InitializeMathFunctions,
	)

	tests := []struct {
		name     string
		code     string
		expected float64
		delta    float64  // tolerance for floating point comparison
	}{
		{"sin", "let result = sin(1.5708)", 1.0, 0.0001},
		{"cos", "let result = cos(0)", 1.0, 0.0001},
		{"sqrt", "let result = sqrt(16)", 4.0, 0.0001},
		{"abs positive", "let result = abs(5.5)", 5.5, 0.0001},
		{"abs negative", "let result = abs(-3.2)", 3.2, 0.0001},
		{"pow", "let result = pow(2, 3)", 8.0, 0.0001},
		{"floor", "let result = floor(5.7)", 5.0, 0.0001},
		{"ceil", "let result = ceil(5.2)", 6.0, 0.0001},
		{"round", "let result = round(5.6)", 6.0, 0.0001},
		{"min", "let result = min(3, 7)", 3.0, 0.0001},
		{"max", "let result = max(3, 7)", 7.0, 0.0001},
		{"log", "let result = log(2.718281828)", 1.0, 0.0001},
		{"log10", "let result = log10(100)", 2.0, 0.0001},
		{"exp", "let result = exp(1)", math.E, 0.0001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сбрасываем окружение для каждого теста
			InitTestEnvironment(
				builtin.InitializeMathFunctions,
			)
			
			exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
			for _, expr := range exprs {
				expr.Eval()
			}

			result, ok := scope.GlobalScope.Get("result")
			if !ok {
				t.Errorf("result variable not found")
				return
			}

			actual := result.Float64()
			if math.Abs(actual-tt.expected) > tt.delta {
				t.Errorf("expected %.6f, got %.6f", tt.expected, actual)
			}
		})
	}
}

func TestMathFunctionErrors(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment(
		builtin.InitializeMathFunctions,
	)

	errorTests := []struct {
		name string
		code string
		expectedError string
	}{
		{"sqrt negative", "sqrt(-1)", "sqrt() argument must be non-negative"},
		{"log negative", "log(-1)", "log() argument must be positive"},
		{"log zero", "log(0)", "log() argument must be positive"},
		{"log10 negative", "log10(-1)", "log10() argument must be positive"},
		{"sin too many args", "sin(1, 2)", "sin() takes exactly 1 argument"},
		{"pow wrong args", "pow(2)", "pow() takes exactly 2 arguments"},
		{"min wrong args", "min(1)", "min() takes exactly 2 arguments"},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					errorStr := r.(string)
					if errorStr != tt.expectedError {
						t.Errorf("expected error '%s', got '%s'", tt.expectedError, errorStr)
					}
				} else {
					t.Errorf("expected panic with error '%s', but no panic occurred", tt.expectedError)
				}
			}()

			exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
			for _, expr := range exprs {
				expr.Eval()
			}
		})
	}
}

func TestMathFunctionIntegration(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment(
		builtin.InitializeMathFunctions,
	)

	// Test using math functions in user-defined functions
	const code = `
	fn distance(x1, y1, x2, y2) {
		let dx = x2 - x1
		let dy = y2 - y1
		return sqrt(pow(dx, 2) + pow(dy, 2))
	}
	
	let result = distance(0, 0, 3, 4)
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	result, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	expected := 5.0
	actual := result.Float64()
	if math.Abs(actual-expected) > 0.0001 {
		t.Errorf("expected %.6f, got %.6f", expected, actual)
	}
}

func TestTrigonometricFunctions(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment(
		builtin.InitializeMathFunctions,
	)

	// Test trigonometric identities
	const code = `
	let pi = 3.14159265359
	let sinPiOver2 = sin(pi / 2)
	let cosZero = cos(0)
	let tanPiOver4 = tan(pi / 4)
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// sin(π/2) should be close to 1
	sinVal, ok := scope.GlobalScope.Get("sinPiOver2")
	if !ok {
		t.Errorf("sinPiOver2 not found")
		return
	}
	if math.Abs(sinVal.Float64()-1.0) > 0.001 {
		t.Errorf("sin(π/2) expected ~1.0, got %.6f", sinVal.Float64())
	}

	// cos(0) should be 1
	cosVal, ok := scope.GlobalScope.Get("cosZero")
	if !ok {
		t.Errorf("cosZero not found")
		return
	}
	if math.Abs(cosVal.Float64()-1.0) > 0.001 {
		t.Errorf("cos(0) expected 1.0, got %.6f", cosVal.Float64())
	}

	// tan(π/4) should be close to 1
	tanVal, ok := scope.GlobalScope.Get("tanPiOver4")
	if !ok {
		t.Errorf("tanPiOver4 not found")
		return
	}
	if math.Abs(tanVal.Float64()-1.0) > 0.01 {
		t.Errorf("tan(π/4) expected ~1.0, got %.6f", tanVal.Float64())
	}
}