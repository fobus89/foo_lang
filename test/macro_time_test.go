package test

import (
	"bytes"
	"foo_lang/parser"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMacroTimeExecution(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "simple macro-time execution",
			code: `
				macro testMacro() {
					println("Macro-time execution")
				}
				
				@testMacro()
			`,
			output: "Macro-time execution",
		},
		{
			name: "macro with Expr block",
			code: `
				macro testExprBlock() {
					println("Macro-time")
					
					Expr {
						println("Code generation")
					}
				}
				
				@testExprBlock()
			`,
			output: "Macro-time\nCode generation",
		},
		{
			name: "macro with type analysis",
			code: `
				struct User {
					name: string,
					age: int
				}
				
				macro analyzeType(typeInfo) {
					println("Analyzing: " + typeInfo.Name)
					println("Kind: " + typeInfo.Kind)
					
					Expr {
						if typeInfo.isStruct() {
							println("Generated for struct")
						}
					}
				}
				
				let userType = type(User)
				@analyzeType(userType)
			`,
			output: "Analyzing: User\nKind: struct\nGenerated for struct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput4(func() {
				InitTestEnvironment()
				exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
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

func TestMacroTimeWithPolymorphicTypes(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "polymorphic type checking in macros",
			code: `
				struct Product { name: string }
				
				macro smartGenerate(someType) {
					println("Type: " + someType.Name)
					
					if someType.isStruct() {
						println("Is struct: true")
					} else if someType.isPrimitive() {
						println("Is primitive: true")
					}
					
					Expr {
						println("Generated for: " + someType.Name)
					}
				}
				
				let productType = type(Product)
				let intType = type(int)
				
				@smartGenerate(productType)
				@smartGenerate(intType)
			`,
			output: "Type: Product\nIs struct: true\nGenerated for: Product\nType: int\nIs primitive: true\nGenerated for: int",
		},
		{
			name: "conditional code generation",
			code: `
				struct Config { host: string }
				
				macro conditionalGen(typeInfo) {
					println("Processing: " + typeInfo.Name)
					
					Expr {
						if typeInfo.isStruct() {
							println("fn create" + typeInfo.Name + "() { return {} }")
						} else {
							println("fn default" + typeInfo.Name + "() { return nil }")
						}
					}
				}
				
				let configType = type(Config)
				let stringType = type(string)
				
				@conditionalGen(configType)
				@conditionalGen(stringType)
			`,
			output: "Processing: Config\nfn createConfig() { return {} }\nProcessing: string\nfn defaultstring() { return nil }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput4(func() {
				InitTestEnvironment()
				exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
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

func TestMacroTimeAdvanced(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "complex macro-time metaprogramming",
			code: `
				struct Entity {
					id: int,
					name: string
				}
				
				macro generateCRUD(entityType) {
					println("=== MACRO-TIME ===")
					println("Entity: " + entityType.Name)
					
					if entityType.isStruct() {
						println("Generating CRUD operations")
					}
					
					Expr {
						println("=== GENERATED CODE ===")
						println("fn create" + entityType.Name + "() {}")
						println("fn update" + entityType.Name + "() {}")
						println("fn delete" + entityType.Name + "() {}")
					}
				}
				
				let entityType = type(Entity)
				@generateCRUD(entityType)
			`,
			output: "=== MACRO-TIME ===\nEntity: Entity\nGenerating CRUD operations\n=== GENERATED CODE ===\nfn createEntity() {}\nfn updateEntity() {}\nfn deleteEntity() {}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput4(func() {
				InitTestEnvironment()
				exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
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

// captureOutput4 captures stdout output from a function (renamed to avoid conflict)
func captureOutput4(f func()) string {
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