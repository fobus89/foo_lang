package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"foo_lang/ast"
	"testing"
)

func TestBasicImportParsing(t *testing.T) {
	InitTestEnvironment()

	const code = `import "./math.foo"`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The import should parse successfully
	expr := exprs[0]
	
	// Just check that we parsed an ImportExpr (don't actually execute it)
	importExpr, ok := expr.(*ast.ImportExpr)
	if !ok {
		t.Errorf("expected ImportExpr, got %T", expr)
		return
	}
	
	// Verify the import path was parsed correctly
	if importExpr.Path != "./math.foo" {
		t.Errorf("expected path './math.foo', got '%s'", importExpr.Path)
	}

	// Import parsing test complete - actual loading is tested in module_loading_test.go
}

func TestSelectiveImportParsing(t *testing.T) {
	InitTestEnvironment()

	const code = `import { add, subtract } from "./math.foo"`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The selective import should parse successfully
	expr := exprs[0]
	
	// Just check that we parsed an ImportExpr (don't actually execute it)
	importExpr, ok := expr.(*ast.ImportExpr)
	if !ok {
		t.Errorf("expected ImportExpr, got %T", expr)
		return
	}
	
	// Verify the imported items were parsed correctly
	if len(importExpr.ImportedItems) != 2 {
		t.Errorf("expected 2 imported items, got %d", len(importExpr.ImportedItems))
	}
	if importExpr.ImportedItems[0] != "add" || importExpr.ImportedItems[1] != "subtract" {
		t.Errorf("expected ['add', 'subtract'], got %v", importExpr.ImportedItems)
	}
}

func TestAliasImportParsing(t *testing.T) {
	InitTestEnvironment()

	const code = `import * as Math from "./math.foo"`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The alias import should parse successfully
	expr := exprs[0]
	
	// Just check that we parsed an ImportExpr (don't actually execute it)
	importExpr, ok := expr.(*ast.ImportExpr)
	if !ok {
		t.Errorf("expected ImportExpr, got %T", expr)
		return
	}
	
	// Verify the alias was parsed correctly
	if importExpr.AliasName != "Math" {
		t.Errorf("expected alias 'Math', got '%s'", importExpr.AliasName)
	}
	if !importExpr.ImportedAll {
		t.Errorf("expected ImportedAll to be true for alias import")
	}
}

func TestBasicExportParsing(t *testing.T) {
	InitTestEnvironment()

	const code = `
		export fn add(a, b) {
			return a + b
		}
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The export should parse and execute successfully
	expr := exprs[0]
	result := expr.Eval()
	
	// Export statements return the result of their declaration
	if result != nil {
		t.Errorf("expected export to return nil, got %v", result)
	}

	// Check if function was defined
	if _, ok := scope.GlobalScope.Get("add"); !ok {
		t.Errorf("exported function 'add' not found in scope")
	}

	// Check if export was registered
	if _, ok := scope.GlobalScope.Get("__export_add"); !ok {
		t.Errorf("export 'add' not registered")
	}
}

func TestExportVariable(t *testing.T) {
	InitTestEnvironment()

	const code = `export let PI = 3.14159`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The export should parse and execute successfully
	expr := exprs[0]
	result := expr.Eval()
	
	// Export statements return nil
	if result != nil {
		t.Errorf("expected export to return nil, got %v", result)
	}

	// Check if variable was defined
	if val, ok := scope.GlobalScope.Get("PI"); !ok {
		t.Errorf("exported variable 'PI' not found in scope")
	} else if val.Float64() != 3.14159 {
		t.Errorf("expected PI = 3.14159, got %f", val.Float64())
	}

	// Check if export was registered
	if _, ok := scope.GlobalScope.Get("__export_PI"); !ok {
		t.Errorf("export 'PI' not registered")
	}
}

func TestExportConstant(t *testing.T) {
	InitTestEnvironment()

	const code = `export const MAX_SIZE = 100`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The export should parse and execute successfully
	expr := exprs[0]
	result := expr.Eval()
	
	// Export statements return nil
	if result != nil {
		t.Errorf("expected export to return nil, got %v", result)
	}

	// Check if constant was defined
	if val, ok := scope.GlobalScope.Get("MAX_SIZE"); !ok {
		t.Errorf("exported constant 'MAX_SIZE' not found in scope")
	} else if val.Int64() != 100 {
		t.Errorf("expected MAX_SIZE = 100, got %d", val.Int64())
	}

	// Check if export was registered
	if _, ok := scope.GlobalScope.Get("__export_MAX_SIZE"); !ok {
		t.Errorf("export 'MAX_SIZE' not registered")
	}
}

func TestExportEnum(t *testing.T) {
	InitTestEnvironment()

	const code = `export enum Status { PENDING, RUNNING, COMPLETED }`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The export should parse and execute successfully
	expr := exprs[0]
	result := expr.Eval()
	
	// Export statements return nil
	if result != nil {
		t.Errorf("expected export to return nil, got %v", result)
	}

	// Check if enum was defined (enum creates namespace)
	if _, ok := scope.GlobalScope.Get("Status"); !ok {
		t.Errorf("exported enum 'Status' not found in scope")
	}

	// Check if export was registered
	if _, ok := scope.GlobalScope.Get("__export_Status"); !ok {
		t.Errorf("export 'Status' not registered")
	}
}

// TODO: Tests for actual module loading and importing will be added 
// when the full module system is implemented