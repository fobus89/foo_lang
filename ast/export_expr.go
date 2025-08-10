package ast

import "foo_lang/scope"

// ExportExpr represents export statements:
// export fn add(a, b) { return a + b }
// export let PI = 3.14159
// export enum Color { RED, GREEN, BLUE }
type ExportExpr struct {
	Declaration Expr // The declaration being exported (FuncExpr, LetExpr, EnumExpr, etc.)
	Name        string // Name of the exported item
}

func NewExportExpr(declaration Expr, name string) *ExportExpr {
	return &ExportExpr{
		Declaration: declaration,
		Name:        name,
	}
}

func (e *ExportExpr) Eval() *Value {
	// Execute the declaration first
	result := e.Declaration.Eval()
	
	// Mark this item as exported in the module's export table
	// For now, store in a special scope with "__export_" prefix
	if e.Name != "" {
		// Get the actual value from global scope and mark it as exported
		if val, exists := scope.GlobalScope.Get(e.Name); exists {
			scope.GlobalScope.Set("__export_"+e.Name, val)
		}
	}
	
	return result
}