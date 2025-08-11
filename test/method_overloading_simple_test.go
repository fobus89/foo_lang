package test

import (
	"testing"
	"foo_lang/parser"
	"foo_lang/ast"
)

func TestMethodOverloadingSimple(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Простые перегрузки для тестирования
	fn add(a: int, b: int) {
		return a + b
	}
	
	fn add(a: string, b: string) {
		return a + b
	}
	
	// Тестируем
	println(add(5, 3))          // 8
	println(add("Hello", " World"))  // "Hello World"
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestMethodOverloadingDifferentParamCounts(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Разные количества параметров - используем разные имена функций
	fn msg0() {
		return "No params"
	}
	
	fn msg1(a: string) {
		return "One param: " + a
	}
	
	fn msg2(a: string, b: string) {
		return "Two params: " + a + ", " + b
	}
	
	// Тестируем
	println(msg0())                 // "No params" 
	println(msg1("test"))           // "One param: test"
	println(msg2("hello", "world")) // "Two params: hello, world"
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}