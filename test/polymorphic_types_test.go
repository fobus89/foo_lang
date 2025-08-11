package test

import (
	"bytes"
	"foo_lang/parser"
	"foo_lang/scope"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPolymorphicTypeMethods(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "isStruct method",
			code: `
				struct Person {
					name: string,
					age: int
				}
				
				let personType = type(Person)
				println(personType.isStruct())
				println(personType.isEnum())
				println(personType.isPrimitive())
			`,
			output: "true\nfalse\nfalse",
		},
		{
			name: "isPrimitive method",
			code: `
				let intType = type(int)
				let stringType = type(string)
				
				println(intType.isPrimitive())
				println(intType.isStruct())
				println(stringType.isPrimitive())
			`,
			output: "true\nfalse\ntrue",
		},
		{
			name: "typeof with polymorphic methods",
			code: `
				let x = 42
				let str = "hello"
				let obj = {name: "test", value: 123}
				
				let xType = typeof(x)
				let strType = typeof(str)
				let objType = typeof(obj)
				
				println(xType.isPrimitive())
				println(strType.isPrimitive())
				println(objType.isStruct())
			`,
			output: "true\ntrue\ntrue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput3(func() {
				scope.GlobalScope = scope.NewScopeStack()
				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})
			expected := strings.TrimSpace(tt.output)
			actual := strings.TrimSpace(result)
			
			if actual != expected {
				t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, actual)
			}
		})
	}
}

func TestTypeConversions(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "struct to struct conversion",
			code: `
				struct User {
					name: string,
					id: int
				}
				
				let userType = type(User)
				let structType = userType.toStruct()
				println(structType.Name)
				println(structType.Kind)
			`,
			output: "User\nstruct",
		},
		{
			name: "primitive conversions",
			code: `
				let intType = type(int)
				// Примитивные типы не могут быть конвертированы в struct
				println(intType.isPrimitive())
			`,
			output: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput3(func() {
				scope.GlobalScope = scope.NewScopeStack()
				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})
			expected := strings.TrimSpace(tt.output)
			actual := strings.TrimSpace(result)
			
			if actual != expected {
				t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, actual)
			}
		})
	}
}

func TestPolymorphicMacros(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "polymorphic type analysis macro",
			code: `
				struct Product {
					name: string,
					price: float
				}
				
				macro analyzeAnyType(someType) {
					if someType.isStruct() {
						println("Found struct: " + someType.Name)
					} else if someType.isPrimitive() {
						println("Found primitive: " + someType.Name)
					} else {
						println("Unknown type: " + someType.Kind)
					}
				}
				
				let productType = type(Product)
				let intType = type(int)
				
				@analyzeAnyType(productType)
				@analyzeAnyType(intType)
			`,
			output: "Found struct: Product\nFound primitive: int",
		},
		{
			name: "conditional code generation",
			code: `
				struct Config {
					host: string,
					port: int
				}
				
				macro generateForType(typeInfo) {
					if typeInfo.isStruct() {
						println("fn new" + typeInfo.Name + "() { return {} }")
					} else if typeInfo.isPrimitive() {
						println("fn default" + typeInfo.Name + "() { return nil }")
					}
				}
				
				let configType = type(Config)
				let stringType = type(string)
				
				@generateForType(configType)
				@generateForType(stringType)
			`,
			output: "fn newConfig() { return {} }\nfn defaultstring() { return nil }",
		},
		{
			name: "typeof in macro",
			code: `
				macro processValue(value) {
					let valueType = typeof(value)
					if valueType.isPrimitive() {
						println("Processing primitive: " + value)
					} else if valueType.isStruct() {
						println("Processing object")
					}
				}
				
				@processValue(42)
				@processValue("hello")
			`,
			output: "Processing primitive: 42\nProcessing primitive: hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput3(func() {
				scope.GlobalScope = scope.NewScopeStack()
				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})
			expected := strings.TrimSpace(tt.output)
			actual := strings.TrimSpace(result)
			
			if actual != expected {
				t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, actual)
			}
		})
	}
}

func TestTypeConversionErrors(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		shouldPanic bool
		errorMsg    string
	}{
		{
			name: "primitive to struct conversion error",
			code: `
				let intType = type(int)
				let structType = intType.toStruct()
			`,
			shouldPanic: true,
			errorMsg:    "cannot convert primitive to struct",
		},
		{
			name: "struct to enum conversion error",
			code: `
				struct User { name: string }
				let userType = type(User)
				let enumType = userType.toEnum()
			`,
			shouldPanic: true,
			errorMsg:    "cannot convert struct to enum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldPanic {
						t.Errorf("Code panicked unexpectedly: %v", r)
					} else if !strings.Contains(r.(string), tt.errorMsg) {
						t.Errorf("Expected error containing '%s', got '%v'", tt.errorMsg, r)
					}
				} else if tt.shouldPanic {
					t.Errorf("Expected panic but code executed successfully")
				}
			}()
			
			scope.GlobalScope = scope.NewScopeStack()
			exprs := parser.NewParser(tt.code).Parse()
			for _, expr := range exprs {
				expr.Eval()
			}
		})
	}
}

func TestUniversalTypeAnalysis(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "universal type analyzer macro",
			code: `
				struct Item { id: int }
				
				macro universalAnalyzer(anyType) {
					println("=== Type Analysis ===")
					println("Name: " + anyType.Name)
					println("Kind: " + anyType.Kind)
					println("isStruct: " + anyType.isStruct())
					println("isPrimitive: " + anyType.isPrimitive())
				}
				
				let itemType = type(Item)
				@universalAnalyzer(itemType)
			`,
			output: "=== Type Analysis ===\nName: Item\nKind: struct\nisStruct: true\nisPrimitive: false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput3(func() {
				scope.GlobalScope = scope.NewScopeStack()
				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})
			expected := strings.TrimSpace(tt.output)
			actual := strings.TrimSpace(result)
			
			if actual != expected {
				t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, actual)
			}
		})
	}
}

// captureOutput3 captures stdout output from a function (renamed to avoid conflict)
func captureOutput3(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}