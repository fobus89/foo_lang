package ast

import "foo_lang/scope"

// ImportExpr represents different types of import statements:
// import "./module.foo"
// import { func1, var1 } from "./module.foo"
// import * as ModuleName from "./module.foo"
type ImportExpr struct {
	Path         string   // Path to the module file
	ImportedAll  bool     // true for import * as Name
	AliasName    string   // Name for import * as Name, or empty
	ImportedItems []string // List of specific items to import, empty for import all
}

func NewImportExpr(path string) *ImportExpr {
	return &ImportExpr{
		Path:        path,
		ImportedAll: true, // Default: import everything
	}
}

func NewSelectiveImportExpr(path string, items []string) *ImportExpr {
	return &ImportExpr{
		Path:          path,
		ImportedAll:   false,
		ImportedItems: items,
	}
}

func NewAliasImportExpr(path string, alias string) *ImportExpr {
	return &ImportExpr{
		Path:        path,
		ImportedAll: true,
		AliasName:   alias,
	}
}

func (i *ImportExpr) Eval() *Value {
	// Import statements don't return values, they modify the current scope
	// We need to use the modules system to load and import
	
	// For now, create a simple placeholder until we implement the full module system
	// This will be replaced with actual module loading logic
	
	// Mark as imported in global scope (temporary solution)
	scope.GlobalScope.Set("__import_"+i.Path, NewValue("loaded"))
	
	// TODO: Implement actual module loading using modules.ImportModule()
	// err := modules.ImportModule(i.Path, "current_file", i.ImportedItems, i.AliasName)
	// if err != nil {
	//     panic(err.Error())
	// }
	
	return nil
}