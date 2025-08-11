package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"testing"
)

func TestAnonymousFunctionArrow(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()
	
	input := `
	let add = fn(x, y) => x + y
	let result = add(5, 3)
	println(result)
	`

	program := parser.NewParser(input).Parse()

	if len(program) != 3 {
		t.Errorf("expected 3 statements, got %d", len(program))
		return
	}

	// Выполняем программу
	for _, stmt := range program {
		if stmt != nil {
			stmt.Eval()
		}
	}
}

func TestAnonymousFunctionBlock(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()
	
	input := `
	let multiply = fn(x, y) {
		return x * y
	}
	let result = multiply(4, 6)
	println(result)
	`

	program := parser.NewParser(input).Parse()

	if len(program) != 3 {
		t.Errorf("expected 3 statements, got %d", len(program))
		return
	}

	// Выполняем программу
	for _, stmt := range program {
		if stmt != nil {
			stmt.Eval()
		}
	}
}

func TestAnonymousFunctionWithDefaults(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()
	
	input := `
	let greet = fn(name, prefix = "Hello") => prefix + ", " + name + "!"
	let result1 = greet("World")
	let result2 = greet("Alice", "Hi")
	println(result1)
	println(result2)
	`

	program := parser.NewParser(input).Parse()

	if len(program) != 5 {
		t.Errorf("expected 5 statements, got %d", len(program))
		return
	}

	// Выполняем программу
	for _, stmt := range program {
		if stmt != nil {
			stmt.Eval()
		}
	}
}

func TestAnonymousFunctionWithClosures(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()
	
	input := `
	let counter = 0
	let increment = fn() {
		counter = counter + 1
		return counter
	}
	
	let val1 = increment()
	let val2 = increment()
	println(val1)
	println(val2)
	`

	program := parser.NewParser(input).Parse()

	if len(program) != 6 {
		t.Errorf("expected 6 statements, got %d", len(program))
		return
	}

	// Выполняем программу
	for _, stmt := range program {
		if stmt != nil {
			stmt.Eval()
		}
	}
}

func TestAnonymousFunctionAssignedToVariable(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()
	
	input := `
	let add = fn(a, b) => a + b
	let multiply = fn(a, b) => a * b
	let subtract = fn(a, b) => a - b
	
	let sum = add(10, 5)
	let product = multiply(3, 4)
	let difference = subtract(8, 3)
	
	println(sum)
	println(product)
	println(difference)
	`

	program := parser.NewParser(input).Parse()

	if len(program) != 9 {
		t.Errorf("expected 9 statements, got %d", len(program))
		return
	}

	// Выполняем программу
	for _, stmt := range program {
		if stmt != nil {
			stmt.Eval()
		}
	}
}