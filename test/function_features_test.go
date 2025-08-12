package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"testing"
)

func TestFunctionDefaultParameters(t *testing.T) {
	InitTestEnvironment()

	const code = `
		fn greet(name, greeting = "Hello") {
			return greeting + " " + name
		}
		
		fn calculate(x, y = 10, z = 5) {
			return x + y + z
		}
		
		let result1 = greet("John")
		let result2 = greet("Jane", "Hi")
		let result3 = calculate(5)
		let result4 = calculate(5, 15)
		let result5 = calculate(1, 2, 3)
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test default parameter usage
	result1Val, ok := scope.GlobalScope.Get("result1")
	if !ok {
		t.Errorf("result1 not found")
		return
	}
	if result1Val.String() != "Hello John" {
		t.Errorf("expected 'Hello John', got %s", result1Val.String())
	}

	// Test explicit parameter
	result2Val, ok := scope.GlobalScope.Get("result2")
	if !ok {
		t.Errorf("result2 not found")
		return
	}
	if result2Val.String() != "Hi Jane" {
		t.Errorf("expected 'Hi Jane', got %s", result2Val.String())
	}

	// Test multiple default parameters
	result3Val, ok := scope.GlobalScope.Get("result3")
	if !ok {
		t.Errorf("result3 not found")
		return
	}
	if result3Val.Int64() != 20 { // 5 + 10 + 5
		t.Errorf("expected 20, got %d", result3Val.Int64())
	}

	// Test partial default parameters
	result4Val, ok := scope.GlobalScope.Get("result4")
	if !ok {
		t.Errorf("result4 not found")
		return
	}
	if result4Val.Int64() != 25 { // 5 + 15 + 5
		t.Errorf("expected 25, got %d", result4Val.Int64())
	}

	// Test all explicit parameters
	result5Val, ok := scope.GlobalScope.Get("result5")
	if !ok {
		t.Errorf("result5 not found")
		return
	}
	if result5Val.Int64() != 6 { // 1 + 2 + 3
		t.Errorf("expected 6, got %d", result5Val.Int64())
	}
}

func TestFunctionComplexDefaultParameters(t *testing.T) {
	InitTestEnvironment()

	const code = `
		let defaultValue = 42
		
		fn test(a, b = defaultValue * 2) {
			return a + b
		}
		
		let result = test(8)
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test expression as default parameter
	resultVal, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}
	if resultVal.Int64() != 92 { // 8 + (42 * 2) = 8 + 84
		t.Errorf("expected 92, got %d", resultVal.Int64())
	}
}

func TestFunctionMissingRequiredParameter(t *testing.T) {
	InitTestEnvironment()

	const code = `
		fn requiresTwo(a, b) {
			return a + b
		}
		
		requiresTwo(5) // Should panic
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()

	// This should panic with missing required argument
	defer func() {
		if r := recover(); r != nil {
			errorStr := r.(string)
			if errorStr != "missing required argument: b" {
				t.Errorf("expected 'missing required argument: b', got: %s", errorStr)
			}
		} else {
			t.Errorf("expected panic due to missing required argument")
		}
	}()

	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestMultipleReturnValues(t *testing.T) {
	InitTestEnvironment()

	const code = `
		fn divmod(a, b) {
			return a / b, a % b
		}
		
		fn getNameAge() {
			return "John", 25
		}
		
		let quotient, remainder = divmod(17, 5)
		let name, age = getNameAge()
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test divmod function
	quotientVal, ok := scope.GlobalScope.Get("quotient")
	if !ok {
		t.Errorf("quotient not found")
		return
	}
	if quotientVal.Int64() != 3 { // 17 / 5 = 3
		t.Errorf("expected 3, got %d", quotientVal.Int64())
	}

	remainderVal, ok := scope.GlobalScope.Get("remainder")
	if !ok {
		t.Errorf("remainder not found")
		return
	}
	if remainderVal.Int64() != 2 { // 17 % 5 = 2
		t.Errorf("expected 2, got %d", remainderVal.Int64())
	}

	// Test getNameAge function
	nameVal, ok := scope.GlobalScope.Get("name")
	if !ok {
		t.Errorf("name not found")
		return
	}
	if nameVal.String() != "John" {
		t.Errorf("expected 'John', got %s", nameVal.String())
	}

	ageVal, ok := scope.GlobalScope.Get("age")
	if !ok {
		t.Errorf("age not found")
		return
	}
	if ageVal.Int64() != 25 {
		t.Errorf("expected 25, got %d", ageVal.Int64())
	}
}

func TestMultipleReturnWithSingleValue(t *testing.T) {
	InitTestEnvironment()

	const code = `
		fn singleReturn() {
			return 42
		}
		
		// Assign single return to multiple variables
		let a, b = singleReturn()
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// First variable should get the value
	aVal, ok := scope.GlobalScope.Get("a")
	if !ok {
		t.Errorf("a not found")
		return
	}
	if aVal.Int64() != 42 {
		t.Errorf("expected 42, got %d", aVal.Int64())
	}

	// Second variable should get nil
	bVal, ok := scope.GlobalScope.Get("b")
	if !ok {
		t.Errorf("b not found")
		return
	}
	if bVal.Any() != nil {
		t.Errorf("expected nil, got %v", bVal.Any())
	}
}