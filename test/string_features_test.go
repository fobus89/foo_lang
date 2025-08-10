package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"testing"
)

func TestStringInterpolation(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		let name = "World"
		let age = 25
		let greeting = "Hello ${name}!"
		let info = "Name: ${name}, Age: ${age}"
	`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test simple interpolation
	greetingVal, ok := scope.GlobalScope.Get("greeting")
	if !ok {
		t.Errorf("greeting not found")
		return
	}

	if greetingVal.String() != "Hello World!" {
		t.Errorf("expected 'Hello World!', got %s", greetingVal.String())
	}

	// Test multiple interpolation
	infoVal, ok := scope.GlobalScope.Get("info")
	if !ok {
		t.Errorf("info not found")
		return
	}

	if infoVal.String() != "Name: World, Age: 25" {
		t.Errorf("expected 'Name: World, Age: 25', got %s", infoVal.String())
	}
}

func TestStringInterpolationWithExpressions(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		let x = 10
		let y = 5
		let result = "${x} + ${y} = ${x + y}"
		let complex = "Result: ${x * 2 + y}"
	`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test expression interpolation
	resultVal, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	if resultVal.String() != "10 + 5 = 15" {
		t.Errorf("expected '10 + 5 = 15', got %s", resultVal.String())
	}

	// Test complex expression
	complexVal, ok := scope.GlobalScope.Get("complex")
	if !ok {
		t.Errorf("complex not found")
		return
	}

	if complexVal.String() != "Result: 25" {
		t.Errorf("expected 'Result: 25', got %s", complexVal.String())
	}
}

func TestStringInterpolationWithArraysAndObjects(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		let arr = [1, 2, 3]
		let obj = { name: "test", value: 42 }
		let arrayStr = "Array: ${arr}"
		let objStr = "Object name: ${obj.name}"
		let methodStr = "Length: ${arr.length()}"
	`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test array interpolation
	arrayStrVal, ok := scope.GlobalScope.Get("arrayStr")
	if !ok {
		t.Errorf("arrayStr not found")
		return
	}

	if arrayStrVal.String() != "Array: [1, 2, 3]" {
		t.Errorf("expected 'Array: [1, 2, 3]', got %s", arrayStrVal.String())
	}

	// Test object property interpolation
	objStrVal, ok := scope.GlobalScope.Get("objStr")
	if !ok {
		t.Errorf("objStr not found")
		return
	}

	if objStrVal.String() != "Object name: test" {
		t.Errorf("expected 'Object name: test', got %s", objStrVal.String())
	}

	// Test method call interpolation
	methodStrVal, ok := scope.GlobalScope.Get("methodStr")
	if !ok {
		t.Errorf("methodStr not found")
		return
	}

	if methodStrVal.String() != "Length: 3" {
		t.Errorf("expected 'Length: 3', got %s", methodStrVal.String())
	}
}

func TestComments(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		// This is a single line comment
		let x = 42 // Comment at end of line
		
		/*
		This is a 
		multi-line comment
		*/
		let y = /* inline comment */ 100
		
		/* Another multi-line
		   comment with different formatting */
		let z = x + y
	`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Comments should not affect execution
	xVal, ok := scope.GlobalScope.Get("x")
	if !ok {
		t.Errorf("x not found")
		return
	}

	if xVal.Int64() != 42 {
		t.Errorf("expected 42, got %d", xVal.Int64())
	}

	yVal, ok := scope.GlobalScope.Get("y")
	if !ok {
		t.Errorf("y not found")
		return
	}

	if yVal.Int64() != 100 {
		t.Errorf("expected 100, got %d", yVal.Int64())
	}

	zVal, ok := scope.GlobalScope.Get("z")
	if !ok {
		t.Errorf("z not found")
		return
	}

	if zVal.Int64() != 142 {
		t.Errorf("expected 142, got %d", zVal.Int64())
	}
}

func TestStringConcatenation(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		let str1 = "Hello"
		let str2 = "World"
		let result = str1 + " " + str2
		
		let arr = [1, 2, 3]
		let arrStr = "Array: " + arr
		
		let num = 42
		let numStr = "Number: " + num
	`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test string concatenation
	resultVal, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	if resultVal.String() != "Hello World" {
		t.Errorf("expected 'Hello World', got %s", resultVal.String())
	}

	// Test string + array concatenation
	arrStrVal, ok := scope.GlobalScope.Get("arrStr")
	if !ok {
		t.Errorf("arrStr not found")
		return
	}

	if arrStrVal.String() != "Array: [1, 2, 3]" {
		t.Errorf("expected 'Array: [1, 2, 3]', got %s", arrStrVal.String())
	}

	// Test string + number concatenation
	numStrVal, ok := scope.GlobalScope.Get("numStr")
	if !ok {
		t.Errorf("numStr not found")
		return
	}

	if numStrVal.String() != "Number: 42" {
		t.Errorf("expected 'Number: 42', got %s", numStrVal.String())
	}
}