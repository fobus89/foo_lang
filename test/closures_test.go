package test

import (
	"testing"
	"foo_lang/parser"
	"foo_lang/scope"
	"foo_lang/builtin"
	"foo_lang/ast"
	"foo_lang/modules"
)

func TestBasicClosure(t *testing.T) {
	// Устанавливаем global parse function
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

	// Тест базового замыкания - функция захватывает внешнюю переменную
	const code = `
	let x = 10
	
	fn inner() {
		return x + 5  // x захвачена из внешней области
	}
	
	let result = inner()
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

	if result.Int64() != 15 {
		t.Errorf("expected 15, got %d", result.Int64())
	}
}

func TestClosureModification(t *testing.T) {
	// Устанавливаем global parse function
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

	// Тест замыкания с изменением захваченной переменной
	const code = `
	let counter = 0
	
	fn increment() {
		counter = counter + 1
		return counter
	}
	
	let first = increment()
	let second = increment()
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	first, ok := scope.GlobalScope.Get("first")
	if !ok {
		t.Errorf("first not found")
		return
	}

	second, ok := scope.GlobalScope.Get("second")
	if !ok {
		t.Errorf("second not found")
		return
	}

	if first.Int64() != 1 {
		t.Errorf("expected first = 1, got %d", first.Int64())
	}

	if second.Int64() != 2 {
		t.Errorf("expected second = 2, got %d", second.Int64())
	}
}

func TestNestedClosures(t *testing.T) {
	// Устанавливаем global parse function
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

	// Тест вложенных замыканий
	const code = `
	let x = 10
	
	fn outer() {
		let y = 20
		
		fn inner() {
			return x + y  // захватывает и x, и y
		}
		
		return inner()
	}
	
	let result = outer()
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

	if result.Int64() != 30 {
		t.Errorf("expected 30, got %d", result.Int64())
	}
}

func TestClosureWithParameters(t *testing.T) {
	// Устанавливаем global parse function
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

	// Тест замыкания с параметрами и захваченными переменными
	const code = `
	let multiplier = 10
	
	fn createMultiplier() {
		fn multiply(x) {
			return x * multiplier  // захватывает multiplier
		}
		return multiply
	}
	
	let mult = createMultiplier()
	let result = mult(5)
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

	if result.Int64() != 50 {
		t.Errorf("expected 50, got %d", result.Int64())
	}
}

func TestClosureWithMathFunctions(t *testing.T) {
	// Устанавливаем global parse function
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

	// Тест замыканий с математическими функциями
	const code = `
	let radius = 5
	
	fn calculateArea() {
		let pi = 3.14159
		return pi * pow(radius, 2)  // использует math функцию и захваченный radius
	}
	
	let area = calculateArea()
	`

	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	area, ok := scope.GlobalScope.Get("area")
	if !ok {
		t.Errorf("area not found")
		return
	}

	expected := 3.14159 * 25 // π * r²
	actual := area.Float64()
	
	if actual < expected-0.001 || actual > expected+0.001 {
		t.Errorf("expected area ≈ %.5f, got %.5f", expected, actual)
	}
}