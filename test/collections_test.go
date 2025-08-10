package test

import (
	"foo_lang/ast"
	"foo_lang/parser"
	"foo_lang/scope"
	"testing"
)

func TestArrayLiterals(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `let arr = [1, 2, 3, 4, 5]`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("arr")
	if !ok {
		t.Errorf("array not found")
		return
	}

	arr, ok := val.Any().([]any)
	if !ok {
		t.Errorf("expected array, got %T", val.Any())
		return
	}

	if len(arr) != 5 {
		t.Errorf("expected length 5, got %d", len(arr))
	}

	for i, expected := range []int64{1, 2, 3, 4, 5} {
		if arr[i] != expected {
			t.Errorf("expected arr[%d] = %d, got %v", i, expected, arr[i])
		}
	}
}

func TestArrayIndexing(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		let arr = [10, 20, 30]
		let first = arr[0]
		let second = arr[1]
		let third = arr[2]
	`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	tests := []struct {
		name     string
		variable string
		expected int64
	}{
		{"First element", "first", 10},
		{"Second element", "second", 20},
		{"Third element", "third", 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := scope.GlobalScope.Get(tt.variable)
			if !ok {
				t.Errorf("variable %s not found", tt.variable)
				return
			}

			if val.Int64() != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, val.Int64())
			}
		})
	}
}

func TestArrayMethods(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		let arr = [1, 2, 3]
		let len = arr.length()
		let newArr = arr.push(4)
	`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test length
	lenVal, ok := scope.GlobalScope.Get("len")
	if !ok {
		t.Errorf("length not found")
		return
	}

	if lenVal.Int64() != 3 {
		t.Errorf("expected length 3, got %d", lenVal.Int64())
	}

	// Test push
	newArrVal, ok := scope.GlobalScope.Get("newArr")
	if !ok {
		t.Errorf("newArr not found")
		return
	}

	newArr, ok := newArrVal.Any().([]any)
	if !ok {
		t.Errorf("expected array, got %T", newArrVal.Any())
		return
	}

	if len(newArr) != 4 {
		t.Errorf("expected length 4, got %d", len(newArr))
	}

	if newArr[3] != int64(4) {
		t.Errorf("expected newArr[3] = 4, got %v", newArr[3])
	}
}

func TestObjectLiterals(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `let obj = { name: "John", age: 30, active: true }`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("obj")
	if !ok {
		t.Errorf("object not found")
		return
	}

	obj, ok := val.Any().(map[string]*ast.Value)
	if !ok {
		t.Errorf("expected object, got %T", val.Any())
		return
	}

	// Test string field
	if nameVal := obj["name"]; nameVal == nil || nameVal.String() != "John" {
		t.Errorf("expected obj.name = 'John', got %v", nameVal)
	}

	// Test integer field
	if ageVal := obj["age"]; ageVal == nil || ageVal.Int64() != 30 {
		t.Errorf("expected obj.age = 30, got %v", ageVal)
	}

	// Test boolean field
	if activeVal := obj["active"]; activeVal == nil || !activeVal.Bool() {
		t.Errorf("expected obj.active = true, got %v", activeVal)
	}
}

func TestObjectPropertyAccess(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		let person = { name: "Alice", age: 25 }
		let personName = person.name
		let personAge = person.age
	`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test name access
	nameVal, ok := scope.GlobalScope.Get("personName")
	if !ok {
		t.Errorf("personName not found")
		return
	}

	if nameVal.String() != "Alice" {
		t.Errorf("expected 'Alice', got %s", nameVal.String())
	}

	// Test age access
	ageVal, ok := scope.GlobalScope.Get("personAge")
	if !ok {
		t.Errorf("personAge not found")
		return
	}

	if ageVal.Int64() != 25 {
		t.Errorf("expected 25, got %d", ageVal.Int64())
	}
}

func TestObjectStringIndexing(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		let data = { key1: "value1", key2: 42 }
		let val1 = data["key1"]
		let val2 = data["key2"]
	`
	
	exprs := parser.NewParser(code).Parse()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test string value
	val1, ok := scope.GlobalScope.Get("val1")
	if !ok {
		t.Errorf("val1 not found")
		return
	}

	if val1.String() != "value1" {
		t.Errorf("expected 'value1', got %s", val1.String())
	}

	// Test integer value
	val2, ok := scope.GlobalScope.Get("val2")
	if !ok {
		t.Errorf("val2 not found")
		return
	}

	if val2.Int64() != 42 {
		t.Errorf("expected 42, got %d", val2.Int64())
	}
}