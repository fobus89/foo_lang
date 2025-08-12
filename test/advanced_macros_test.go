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

func TestStructDefinition(t *testing.T) {
	tests := []struct {
		name string
		code string

		checks func(*testing.T)
	}{
		{
			name: "simple struct definition",
			code: `
				struct Person {
					name: string,
					age: int
				}
			`,
			checks: func(t *testing.T) {
				// Проверяем, что тип Person создан в scope
				personType, found := scope.GlobalScope.Get("Person")
				if !found {
					t.Error("Person type not found in scope")
					return
				}

				// Дополнительные проверки можно добавить здесь
				if personType == nil {
					t.Error("Person type is nil")
				}
			},
		},
		{
			name: "struct with multiple field types",
			code: `
				struct User {
					id: int,
					name: string,
					email: string,
					active: bool
				}
			`,
			checks: func(t *testing.T) {
				userType, found := scope.GlobalScope.Get("User")
				if !found {
					t.Error("User type not found in scope")
				}
				if userType == nil {
					t.Error("User type is nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope.GlobalScope = scope.NewScopeStack()

			exprs := parser.NewParser(tt.code).Parse()
			for _, expr := range exprs {
				expr.Eval()
			}

			if tt.checks != nil {
				tt.checks(t)
			}
		})
	}
}

func TestTypeofExpression(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "typeof integer",
			code: `
				let x = 42
				let t = typeof(x)
				println(t.String())
			`,
			output: "int",
		},
		{
			name: "typeof string",
			code: `
				let s = "hello"
				let t = typeof(s)
				println(t.String())
			`,
			output: "string",
		},
		{
			name: "typeof boolean",
			code: `
				let b = true
				let t = typeof(b)
				println(t.String())
			`,
			output: "bool",
		},
		{
			name: "typeof object",
			code: `
				let obj = {name: "test", value: 123}
				let t = typeof(obj)
				println(t.String())
			`,
			output: "struct object { value: int, name: string }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput2(func() {
				scope.GlobalScope = scope.NewScopeStack()
				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})
			expected := strings.TrimSpace(tt.output)
			actual := strings.TrimSpace(result)

			// Специальная обработка для typeof_object - порядок полей может меняться
			if tt.name == "typeof object" {
				// Проверяем, что содержит нужные части
				if !strings.Contains(actual, "struct object") ||
					!strings.Contains(actual, "name: string") ||
					!strings.Contains(actual, "value: int") {
					t.Errorf("Expected struct with name: string and value: int, got: %s", actual)
				}
			} else if actual != expected {
				t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, actual)
			}
		})
	}
}

func TestTypeExpression(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "type expression with struct",
			code: `
				struct Person {
					name: string,
					age: int
				}
				
				let personType = type(Person)
				println(personType.Name)
				println(personType.Kind)
			`,
			output: "Person\nstruct",
		},
		{
			name: "primitive types",
			code: `
				let intType = type(int)
				let stringType = type(string)
				println(intType.String())
				println(stringType.String())
			`,
			output: "int\nstring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput2(func() {
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

func TestAdvancedMacrosWithTypes(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "macro with type parameter",
			code: `
				struct User {
					name: string,
					id: int
				}
				
				macro showTypeName(structType) {
					println("Type: " + structType.Name)
				}
				
				let userType = type(User)
				@showTypeName(userType)
			`,
			output: "Type: User",
		},
		{
			name: "macro generating code",
			code: `
				struct Person {
					name: string,
					age: int
				}
				
				macro generateGetter(structType, fieldName) {
					println("fn get" + fieldName + "(obj) {")
					println("    return obj." + fieldName)
					println("}")
				}
				
				let personType = type(Person)
				@generateGetter(personType, "name")
			`,
			output: "fn getname(obj) {\n    return obj.name\n}",
		},
		{
			name: "type introspection",
			code: `
				struct Product {
					name: string,
					price: float,
					inStock: bool
				}
				
				macro analyzeType(structType) {
					println("Analyzing type: " + structType.Name)
					println("Kind: " + structType.Kind)
					if structType.Kind == "struct" {
						println("This is a struct type")
					}
				}
				
				let productType = type(Product)
				@analyzeType(productType)
			`,
			output: "Analyzing type: Product\nKind: struct\nThis is a struct type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput2(func() {
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

func TestMacroCodeGeneration(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "generate setter macro",
			code: `
				struct Config {
					host: string,
					port: int
				}
				
				macro generateSetter(structType, fieldName) {
					println("fn set" + fieldName + "(obj, value) {")
					println("    obj." + fieldName + " = value")
					println("}")
				}
				
				let configType = type(Config)
				@generateSetter(configType, "port")
			`,
			output: "fn setport(obj, value) {\n    obj.port = value\n}",
		},
		{
			name: "multiple code generation",
			code: `
				struct Item {
					name: string
				}
				
				macro generateAccessors(structType, fieldName) {
					println("// Getter")
					println("fn get" + fieldName + "() { return this." + fieldName + " }")
					println("// Setter") 
					println("fn set" + fieldName + "(val) { this." + fieldName + " = val }")
				}
				
				let itemType = type(Item)
				@generateAccessors(itemType, "name")
			`,
			output: "// Getter\nfn getname() { return this.name }\n// Setter\nfn setname(val) { this.name = val }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput2(func() {
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

// captureOutput2 captures stdout output from a function (renamed to avoid conflict)
func captureOutput2(f func()) string {
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
