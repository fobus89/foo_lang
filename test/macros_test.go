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

func TestMacros(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "simple macro definition and call",
			code: `
				macro greet(name) {
					println("Hello, " + name + "!")
				}
				
				@greet("World")
			`,
			output: "Hello, World!",
		},
		{
			name: "macro with multiple arguments",
			code: `
				macro add(a, b) {
					println(a + b)
				}
				
				@add(5, 3)
			`,
			output: "8",
		},
		{
			name: "debug macro",
			code: `
				macro debug(expr) {
					println("DEBUG: " + expr)
				}
				
				let x = 42
				@debug(x)
			`,
			output: "DEBUG: 42",
		},
		{
			name: "assert macro",
			code: `
				macro assert(condition, message) {
					if !condition {
						println("ASSERTION FAILED: " + message)
					}
				}
				
				@assert(5 > 3, "5 should be greater than 3")
				@assert(2 > 5, "2 should be greater than 5")
			`,
			output: "ASSERTION FAILED: 2 should be greater than 5",
		},
		{
			name: "macro with complex body",
			code: `
				macro repeat(n, text) {
					for let i = 0; i < n; i++ {
						println(text + " " + i)
					}
				}
				
				@repeat(3, "Item")
			`,
			output: "Item 0\nItem 1\nItem 2",
		},
		{
			name: "nested macro calls",
			code: `
				macro double(x) {
					println(x * 2)
				}
				
				macro quadruple(x) {
					@double(x)
					@double(x)
				}
				
				@quadruple(5)
			`,
			output: "10\n10",
		},
		{
			name: "macro with expression evaluation",
			code: `
				macro eval(expr) {
					println("Result: " + expr)
				}
				
				let a = 10
				let b = 5
				@eval(a + b * 2)
			`,
			output: "Result: 20",
		},
		{
			name: "macro scope isolation",
			code: `
				let x = 100
				
				macro changeX() {
					let x = 50
					println("Inside macro: " + x)
				}
				
				@changeX()
				println("Outside macro: " + x)
			`,
			output: "Inside macro: 50\nOutside macro: 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput(func() {
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

func TestMacroErrors(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		shouldPanic bool
		errorMsg    string
	}{
		{
			name: "undefined macro",
			code: `
				@undefinedMacro()
			`,
			shouldPanic: true,
			errorMsg:    "macro 'undefinedMacro' not found",
		},
		{
			name: "wrong number of arguments",
			code: `
				macro test(a, b) {
					println(a + b)
				}
				
				@test(1)
			`,
			shouldPanic: true,
			errorMsg:    "macro 'test' expects 2 arguments, got 1",
		},
		{
			name: "too many arguments",
			code: `
				macro test(a) {
					println(a)
				}
				
				@test(1, 2, 3)
			`,
			shouldPanic: true,
			errorMsg:    "macro 'test' expects 1 arguments, got 3",
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

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		output string
	}{
		{
			name: "quote expression",
			code: `
				let expr = quote(5 + 3)
				println("Quoted expression stored")
			`,
			output: "Quoted expression stored",
		},
		{
			name: "unquote expression",
			code: `
				let x = 10
				let result = unquote(x * 2)
				println(result)
			`,
			output: "20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput(func() {
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

// captureOutput captures stdout output from a function
func captureOutput(f func()) string {
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