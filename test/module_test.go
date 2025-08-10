package test

import (
	"foo_lang/parser"
	"foo_lang/scope"
	"testing"
)

func TestBasicImportParsing(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `import "./math.foo"`
	
	exprs := parser.NewParser(code).Parse()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The import should parse successfully
	expr := exprs[0]
	result := expr.Eval()
	
	// Import statements return nil
	if result != nil {
		t.Errorf("expected import to return nil, got %v", result)
	}

	// Check if import was registered (placeholder test)
	if val, ok := scope.GlobalScope.Get("__import_./math.foo"); !ok {
		t.Errorf("import not registered in scope")
	} else if val.String() != "loaded" {
		t.Errorf("expected 'loaded', got %s", val.String())
	}
}

func TestSelectiveImportParsing(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `import { add, subtract } from "./math.foo"`
	
	exprs := parser.NewParser(code).Parse()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The selective import should parse successfully
	expr := exprs[0]
	result := expr.Eval()
	
	// Import statements return nil
	if result != nil {
		t.Errorf("expected import to return nil, got %v", result)
	}
}

func TestAliasImportParsing(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `import * as Math from "./math.foo"`
	
	exprs := parser.NewParser(code).Parse()
	if len(exprs) != 1 {
		t.Errorf("expected 1 expression, got %d", len(exprs))
		return
	}

	// The alias import should parse successfully
	expr := exprs[0]
	result := expr.Eval()
	
	// Import statements return nil
	if result != nil {
		t.Errorf("expected import to return nil, got %v", result)
	}
}

func TestBasicExportParsing(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()

	const code = `
		export fn add(a, b) {
			return a + b
		}
	`
	
	exprs := parser.NewParser(code).Parse()
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
	scope.GlobalScope = scope.NewScopeStack()

	const code = `export let PI = 3.14159`
	
	exprs := parser.NewParser(code).Parse()
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
	scope.GlobalScope = scope.NewScopeStack()

	const code = `export const MAX_SIZE = 100`
	
	exprs := parser.NewParser(code).Parse()
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
	scope.GlobalScope = scope.NewScopeStack()

	const code = `export enum Status { PENDING, RUNNING, COMPLETED }`
	
	exprs := parser.NewParser(code).Parse()
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