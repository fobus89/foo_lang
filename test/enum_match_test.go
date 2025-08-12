package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"testing"
)

func TestEnumDefinition(t *testing.T) {
	InitTestEnvironment()

	const code = `
		enum Color { RED, GREEN, BLUE }
		let red = Color.RED
		let green = Color.GREEN
		let blue = Color.BLUE
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test enum values
	redVal, ok := scope.GlobalScope.Get("red")
	if !ok {
		t.Errorf("red not found")
		return
	}

	if redVal.Int64() != 0 {
		t.Errorf("expected RED = 0, got %d", redVal.Int64())
	}

	greenVal, ok := scope.GlobalScope.Get("green")
	if !ok {
		t.Errorf("green not found")
		return
	}

	if greenVal.Int64() != 1 {
		t.Errorf("expected GREEN = 1, got %d", greenVal.Int64())
	}

	blueVal, ok := scope.GlobalScope.Get("blue")
	if !ok {
		t.Errorf("blue not found")
		return
	}

	if blueVal.Int64() != 2 {
		t.Errorf("expected BLUE = 2, got %d", blueVal.Int64())
	}
}

func TestEnumAccess(t *testing.T) {
	InitTestEnvironment()

	const code = `
		enum Status { PENDING, RUNNING, COMPLETED, FAILED }
		let currentStatus = Status.RUNNING
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	statusVal, ok := scope.GlobalScope.Get("currentStatus")
	if !ok {
		t.Errorf("currentStatus not found")
		return
	}

	if statusVal.Int64() != 1 {
		t.Errorf("expected RUNNING = 1, got %d", statusVal.Int64())
	}
}

func TestMatchExpression(t *testing.T) {
	InitTestEnvironment()

	const code = `
		let x = 2
		let result = match (x) {
			1 => "one",
			2 => "two", 
			3 => "three",
			_ => "other"
		}
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	resultVal, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	if resultVal.String() != "two" {
		t.Errorf("expected 'two', got %s", resultVal.String())
	}
}

func TestMatchWithDefault(t *testing.T) {
	InitTestEnvironment()

	const code = `
		let x = 99
		let result = match (x) {
			1 => "one",
			2 => "two",
			_ => "default case"
		}
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	resultVal, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	if resultVal.String() != "default case" {
		t.Errorf("expected 'default case', got %s", resultVal.String())
	}
}

func TestMatchWithBlocks(t *testing.T) {
	InitTestEnvironment()

	const code = `
		let score = 80
		let grade = match (score) {
			90 => {
				let temp = "excellent"
				temp
			},
			80 => {
				let temp = "good" 
				temp
			},
			_ => {
				let temp = "needs improvement"
				temp
			}
		}
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	gradeVal, ok := scope.GlobalScope.Get("grade")
	if !ok {
		t.Errorf("grade not found")
		return
	}

	if gradeVal.String() != "good" {
		t.Errorf("expected 'good', got %s", gradeVal.String())
	}
}

func TestMatchWithEnums(t *testing.T) {
	InitTestEnvironment()

	const code = `
		enum Direction { NORTH, SOUTH, EAST, WEST }
		let dir = Direction.NORTH
		let message = match (dir) {
			0 => "Going North",
			1 => "Going South", 
			2 => "Going East",
			3 => "Going West",
			_ => "Unknown direction"
		}
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	messageVal, ok := scope.GlobalScope.Get("message")
	if !ok {
		t.Errorf("message not found")
		return
	}

	if messageVal.String() != "Going North" {
		t.Errorf("expected 'Going North', got %s", messageVal.String())
	}
}

func TestConditionalExpression(t *testing.T) {
	InitTestEnvironment()

	const code = `
		let age = 20
		let status = age >= 18 ? "adult" : "minor"
		let number = 5
		let parity = number % 2 == 0 ? "even" : "odd"
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test age condition
	statusVal, ok := scope.GlobalScope.Get("status")
	if !ok {
		t.Errorf("status not found")
		return
	}

	if statusVal.String() != "adult" {
		t.Errorf("expected 'adult', got %s", statusVal.String())
	}

	// Test parity condition
	parityVal, ok := scope.GlobalScope.Get("parity")
	if !ok {
		t.Errorf("parity not found")
		return
	}

	if parityVal.String() != "odd" {
		t.Errorf("expected 'odd', got %s", parityVal.String())
	}
}

func TestChainedMethodCalls(t *testing.T) {
	InitTestEnvironment()

	const code = `
		let data = { 
			numbers: [1, 2, 3, 4, 5] 
		}
		let length = data.numbers.length()
		let newArray = data.numbers.push(6)
		let newLength = newArray.length()
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test property access with method call
	lengthVal, ok := scope.GlobalScope.Get("length")
	if !ok {
		t.Errorf("length not found")
		return
	}

	if lengthVal.Int64() != 5 {
		t.Errorf("expected 5, got %d", lengthVal.Int64())
	}

	// Test chained operations
	newLengthVal, ok := scope.GlobalScope.Get("newLength")
	if !ok {
		t.Errorf("newLength not found")
		return
	}

	if newLengthVal.Int64() != 6 {
		t.Errorf("expected 6, got %d", newLengthVal.Int64())
	}
}