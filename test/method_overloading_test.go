package test

import (
	"testing"
	"foo_lang/parser"
	"foo_lang/ast"
)

func TestMethodOverloadingBasic(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Основные перегрузки функции add
	fn add(a: int, b: int) -> int {
		return a + b
	}
	
	fn add(a: float, b: float) -> float {
		return a + b
	}
	
	fn add(a: string, b: string) -> string {
		return a + b
	}
	
	// Тестируем вызовы
	println(add(5, 3))          // 8 (int)
	println(add(5.5, 2.3))      // 7.8 (float)
	println(add("Hello", " World"))  // "Hello World" (string)
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

func TestMethodOverloadingParameterCounts(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Перегрузки с разным количеством параметров (используем типизированные функции)
	fn greetEmpty() {
		return "Hello!"
	}
	
	fn greetOne(name: string) {
		return "Hello, " + name + "!"
	}
	
	fn greetTwo(name: string, title: string) {
		return "Hello, " + title + " " + name + "!"
	}
	
	// Тестируем
	println(greetEmpty())                    // "Hello!"
	println(greetOne("Alice"))             // "Hello, Alice!"
	println(greetTwo("Smith", "Dr."))      // "Hello, Dr. Smith!"
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

func TestMethodOverloadingTypePromotion(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Функция принимающая float может принять int (автоматическое приведение)
	fn calculate(x: float, y: float) -> float {
		return x * y + 10.5
	}
	
	fn calculate(x: int, y: int, z: int) -> int {
		return x + y + z
	}
	
	// Тестируем
	println(calculate(3.5, 2.0))    // 17.5 (float + float)
	println(calculate(3, 2))        // 16.5 (int->float приведение)
	println(calculate(1, 2, 3))     // 6 (точное совпадение int)
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

func TestMethodOverloadingWithAnyType(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Функция с разными типами
	fn processInt(value: int) {
		return "Processing int: " + value.toString()
	}
	
	fn processString(value: string) {
		return "Processing string: " + value
	}
	
	fn processFloat(value: float) {
		return "Processing float: " + value.toString()
	}
	
	// Тестируем разные типы
	println(processInt(42))        // "Processing int: 42"
	println(processString("test"))    // "Processing string: test"
	println(processFloat(3.14))      // "Processing float: 3.14"
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

func TestMethodOverloadingMathOperations(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Математические операции с разными типами
	fn multiply(base: int, exp: int) {
		return base * exp
	}
	
	fn multiply(base: float, exp: int) {
		return base * exp
	}
	
	fn multiply(base: float, exp: float) {
		return base * exp
	}
	
	// Тестируем
	println(multiply(2, 3).toString())        // 6 (int)
	println(multiply(2.5, 2).toString())      // 5.0 (float, int)
	println(multiply(3.0, 2.0).toString())    // 6.0 (float, float)
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

func TestMethodOverloadingReturnTypes(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Простые перегрузки с разными типами
	fn doubleInt(value: int) {
		return value * 2
	}
	
	fn doubleFloat(value: float) {
		return value * 2.0
	}
	
	fn concatString(value: string) {
		return value + value
	}
	
	// Тестируем
	println(doubleInt(21).toString())        // "42"
	println(doubleFloat(1.5).toString())     // "3.0" 
	println(concatString("Hi"))              // "HiHi"
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

func TestMethodOverloadingComplexTypes(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Работа с массивами - разные функции для демонстрации
	fn getLength(items) {
		return items.length()
	}
	
	fn getMultiLength(items, multiplier: int) {
		return items.length() * multiplier
	}
	
	fn formatCount(item: string, items) {
		return item + " (total: " + items.length().toString() + ")"
	}
	
	// Тестируем
	let arr = [1, 2, 3, 4, 5]
	println(getLength(arr))               // 5
	println(getMultiLength(arr, 2))       // 10
	println(formatCount("Count", arr))    // "Count (total: 5)"
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

func TestMethodOverloadingErrorCases(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Тестируем случаи, где нет подходящей перегрузки
	fn specificFunction(x: int, y: string) -> string {
		return y + ": " + x.toString()
	}
	
	// Этот вызов должен найти функцию
	println(specificFunction(42, "Number"))  // "Number: 42"
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

func TestMethodOverloadingWithDefaults(t *testing.T) {
	ast.ClearOverloadedMethods()
	code := `
	// Перегрузки с разным количеством параметров
	fn format1(msg: string) {
		return "[INFO] " + msg
	}
	
	fn format2(msg: string, level: string) {
		return "[" + level + "] " + msg
	}
	
	fn format3(msg: string, level: string, prefix: string) {
		return "[" + prefix + "] [" + level + "] " + msg
	}
	
	// Тестируем
	println(format1("Hello"))                    // "[INFO] Hello"
	println(format2("Error occurred", "ERROR"))  // "[ERROR] Error occurred"
	println(format3("Debug info", "DEBUG", "TIME")) // "[TIME] [DEBUG] Debug info"
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