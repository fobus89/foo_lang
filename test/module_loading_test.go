package test

import (
	"os"
	"path/filepath"
	"foo_lang/parser"
	"foo_lang/scope"
	"foo_lang/modules"
	"foo_lang/value"
	"foo_lang/ast"
	"testing"
)

// Helper function to create parse function for tests
func createParseFunc() modules.ParseFunc {
	return func(code string) []modules.Expr {
		exprs := parser.NewParser(code).ParseWithoutScopeInit()
		result := make([]modules.Expr, len(exprs))
		for i, expr := range exprs {
			result[i] = expr
		}
		return result
	}
}

func TestModuleLoading(t *testing.T) {
	InitTestEnvironment()
	
	// Set up global parse function
	parseFunc := createParseFunc()
	ast.SetGlobalParseFunc(parseFunc)

	// Create a temporary test module file
	tempDir, err := os.MkdirTemp("", "foo_test_modules")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create math_utils.foo module
	mathModulePath := filepath.Join(tempDir, "math_utils.foo")
	mathModuleContent := `
export fn add(a, b) {
    return a + b
}

export fn multiply(x, y) {
    return x * y
}

export let PI = 3.14159
`
	err = os.WriteFile(mathModulePath, []byte(mathModuleContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test module: %v", err)
	}

	// Test full import: import "./math_utils.foo"  
	err = modules.ImportModule(mathModulePath, "./", []string{}, "", parseFunc)
	if err != nil {
		t.Fatalf("Failed to import module: %v", err)
	}

	// Test that imported functions work
	testCode := `
let sum = add(5, 3)
let product = multiply(4, 6)
let piValue = PI
`
	
	exprs := parser.NewParser(testCode).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Verify imported items
	sumVal, ok := scope.GlobalScope.Get("sum")
	if !ok {
		t.Errorf("sum not found after import")
		return
	}
	if sumVal.Int64() != 8 {
		t.Errorf("expected sum to be 8, got %d", sumVal.Int64())
	}

	productVal, ok := scope.GlobalScope.Get("product")
	if !ok {
		t.Errorf("product not found after import")
		return
	}
	if productVal.Int64() != 24 {
		t.Errorf("expected product to be 24, got %d", productVal.Int64())
	}

	piVal, ok := scope.GlobalScope.Get("PI")
	if !ok {
		t.Errorf("PI not found after import")
		return
	}
	if piVal.Float64() != 3.14159 {
		t.Errorf("expected PI to be 3.14159, got %f", piVal.Float64())
	}
}

func TestSelectiveImport(t *testing.T) {
	InitTestEnvironment()
	
	// Set up global parse function
	parseFunc := createParseFunc()
	ast.SetGlobalParseFunc(parseFunc)

	// Create a temporary test module file
	tempDir, err := os.MkdirTemp("", "foo_test_modules")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create utils.foo module
	utilsModulePath := filepath.Join(tempDir, "utils.foo")
	utilsModuleContent := `
export fn greet(name) {
    return "Hello " + name
}

export fn double(x) {
    return x * 2
}

export let VERSION = "1.0.0"
export let DEBUG = true
`
	err = os.WriteFile(utilsModulePath, []byte(utilsModuleContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test module: %v", err)
	}

	// Test selective import: import { greet, VERSION } from "./utils.foo" 
	err = modules.ImportModule(utilsModulePath, "./", []string{"greet", "VERSION"}, "", parseFunc)
	if err != nil {
		t.Fatalf("Failed to import module: %v", err)
	}

	// Test that only selected items are imported
	testCode := `
let greeting = greet("World")
let version = VERSION
`
	
	exprs := parser.NewParser(testCode).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Verify selective import worked
	greetingVal, ok := scope.GlobalScope.Get("greeting")
	if !ok {
		t.Errorf("greeting not found after selective import")
		return
	}
	if greetingVal.String() != "Hello World" {
		t.Errorf("expected 'Hello World', got %s", greetingVal.String())
	}

	versionVal, ok := scope.GlobalScope.Get("version")
	if !ok {
		t.Errorf("version not found after selective import")
		return
	}
	if versionVal.String() != "1.0.0" {
		t.Errorf("expected '1.0.0', got %s", versionVal.String())
	}

	// Verify that non-imported items are not available
	_, doubleExists := scope.GlobalScope.Get("double")
	if doubleExists {
		t.Errorf("double should not be imported with selective import")
	}

	_, debugExists := scope.GlobalScope.Get("DEBUG")
	if debugExists {
		t.Errorf("DEBUG should not be imported with selective import")
	}
}

func TestAliasImport(t *testing.T) {
	InitTestEnvironment()
	
	// Set up global parse function
	parseFunc := createParseFunc()
	ast.SetGlobalParseFunc(parseFunc)

	// Create a temporary test module file
	tempDir, err := os.MkdirTemp("", "foo_test_modules")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create calculator.foo module
	calcModulePath := filepath.Join(tempDir, "calculator.foo")
	calcModuleContent := `
export fn add(a, b) {
    return a + b
}

export fn subtract(a, b) {
    return a - b
}

export let NAME = "Calculator"
`
	err = os.WriteFile(calcModulePath, []byte(calcModuleContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test module: %v", err)
	}

	// Test alias import: import * as Calc from "./calculator.foo" 
	err = modules.ImportModule(calcModulePath, "./", []string{}, "Calc", parseFunc)
	if err != nil {
		t.Fatalf("Failed to import module: %v", err)
	}

	// Verify that module is available under alias
	calcVal, ok := scope.GlobalScope.Get("Calc")
	if !ok {
		t.Errorf("Calc alias not found after import")
		return
	}

	// Verify it's an object containing the exports
	calcObj, ok := calcVal.Any().(map[string]*value.Value)
	if !ok {
		t.Errorf("Calc should be a module object, got %T", calcVal.Any())
		return
	}

	// Check that expected exports are present
	if _, exists := calcObj["add"]; !exists {
		t.Errorf("add function not found in Calc module")
	}
	if _, exists := calcObj["subtract"]; !exists {
		t.Errorf("subtract function not found in Calc module")
	}
	if _, exists := calcObj["NAME"]; !exists {
		t.Errorf("NAME variable not found in Calc module")
	}
}

func TestModuleCaching(t *testing.T) {
	InitTestEnvironment()
	
	// Set up global parse function
	parseFunc := createParseFunc()
	ast.SetGlobalParseFunc(parseFunc)

	// Create a temporary test module file
	tempDir, err := os.MkdirTemp("", "foo_test_modules")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create simple.foo module without state
	simpleModulePath := filepath.Join(tempDir, "simple.foo")
	simpleModuleContent := `
export fn double(x) {
    return x * 2
}

export let MESSAGE = "Hello from module"
`
	err = os.WriteFile(simpleModulePath, []byte(simpleModuleContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test module: %v", err)
	}

	// Import module first time
	err = modules.ImportModule(simpleModulePath, "./", []string{}, "", parseFunc)
	if err != nil {
		t.Fatalf("Failed to import module: %v", err)
	}
	
	// Check initial import worked
	doubleVal, ok := scope.GlobalScope.Get("double")
	if !ok {
		t.Errorf("double function not found after first import")
		return
	}
	
	// Import module second time - should use cached version
	err = modules.ImportModule(simpleModulePath, "./", []string{}, "", parseFunc)
	if err != nil {
		t.Fatalf("Failed to import module second time: %v", err)
	}
	
	// Verify module is still accessible
	doubleVal2, ok := scope.GlobalScope.Get("double")
	if !ok {
		t.Errorf("double function not found after second import")
		return
	}
	
	// Both should be the same function reference (from cache)
	if doubleVal != doubleVal2 {
		t.Errorf("Module caching failed - different function instances")
	}
}