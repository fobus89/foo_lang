package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"strings"
	"testing"
)

func TestFunctionDefinition(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	const code = `
		fn add(a, b) {
			return a + b
		}
		let result = add(5, 3)
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	if val.Int64() != 8 {
		t.Errorf("expected 8, got %d", val.Int64())
	}
}

func TestFunctionWithoutReturn(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	const code = `
		fn setGlobal(val) {
			let x = val * 2
		}
		setGlobal(5)
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Function should execute without error even without return
	// This test just ensures no panic occurs
}

func TestRecursiveFunction(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	const code = `
		fn factorial(n) {
			if n <= 1 {
				return 1
			}
			return n * factorial(n - 1)
		}
		let result = factorial(5)
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	if val.Int64() != 120 {
		t.Errorf("expected 120, got %d", val.Int64())
	}
}

func TestMutualRecursion(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	const code = `
		fn isEven(n) {
			if n == 0 {
				return true
			}
			return isOdd(n - 1)
		}
		
		fn isOdd(n) {
			if n == 0 {
				return false
			}
			return isEven(n - 1)
		}
		
		let even4 = isEven(4)
		let odd5 = isOdd(5)
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	even4Val, ok := scope.GlobalScope.Get("even4")
	if !ok {
		t.Errorf("even4 not found")
		return
	}

	if !even4Val.Bool() {
		t.Errorf("expected isEven(4) = true, got %t", even4Val.Bool())
	}

	odd5Val, ok := scope.GlobalScope.Get("odd5")
	if !ok {
		t.Errorf("odd5 not found")
		return
	}

	if !odd5Val.Bool() {
		t.Errorf("expected isOdd(5) = true, got %t", odd5Val.Bool())
	}
}

func TestRecursionDepthLimit(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	const code = `
		fn infiniteRecursion(n) {
			return infiniteRecursion(n + 1)
		}
		infiniteRecursion(1)
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	
	// This should panic with recursion limit exceeded
	defer func() {
		if r := recover(); r != nil {
			errorStr := r.(string)
			if !strings.Contains(errorStr, "maximum recursion depth exceeded") {
				t.Errorf("expected recursion depth error, got: %s", errorStr)
			}
		} else {
			t.Errorf("expected panic due to recursion limit")
		}
	}()

	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestFunctionScope(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	const code = `
		let globalVar = "global"
		
		fn testScope(param) {
			let localVar = "local"
			return param + localVar
		}
		
		let result = testScope("param")
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Global variable should still exist
	globalVal, ok := scope.GlobalScope.Get("globalVar")
	if !ok {
		t.Errorf("globalVar not found")
		return
	}

	if globalVal.String() != "global" {
		t.Errorf("expected 'global', got %s", globalVal.String())
	}

	// Local variable should not exist in global scope
	if _, ok := scope.GlobalScope.Get("localVar"); ok {
		t.Errorf("localVar should not be in global scope")
	}

	// Function result should be correct
	resultVal, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	if resultVal.String() != "paramlocal" {
		t.Errorf("expected 'paramlocal', got %s", resultVal.String())
	}
}

// TODO: Default parameters are not yet implemented
// func TestFunctionWithDefaultParameters(t *testing.T) {
//	scope.GlobalScope = scope.NewScopeStack()
//
//	const code = `
//		fn greet(name, greeting = "Hello") {
//			return greeting + " " + name
//		}
//		
//		let result1 = greet("John")
//		let result2 = greet("Jane", "Hi")
//	`
//	
//	exprs := parser.NewParser(code).ParseWithoutScopeInit()
//	for _, expr := range exprs {
//		expr.Eval()
//	}
//
//	// Test default parameter usage
//	result1Val, ok := scope.GlobalScope.Get("result1")
//	if !ok {
//		t.Errorf("result1 not found")
//		return
//	}
//
//	if result1Val.String() != "Hello John" {
//		t.Errorf("expected 'Hello John', got %s", result1Val.String())
//	}
//
//	// Test explicit parameter
//	result2Val, ok := scope.GlobalScope.Get("result2")
//	if !ok {
//		t.Errorf("result2 not found")
//		return
//	}
//
//	if result2Val.String() != "Hi Jane" {
//		t.Errorf("expected 'Hi Jane', got %s", result2Val.String())
//	}
// }