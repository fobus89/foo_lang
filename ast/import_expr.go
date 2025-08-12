package ast

import (
	"foo_lang/modules"
)

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
	// Use the modules system to load and import
	
	// Используем контекст текущего файла, установленный парсером
	currentFile := GetCurrentFileContext()
	if currentFile == "" {
		currentFile = "./" // Fallback для обратной совместимости
	}
	
	if GlobalParseFunc == nil {
		panic("GlobalParseFunc not set - cannot import modules")
	}
	
	err := modules.ImportModule(i.Path, currentFile, i.ImportedItems, i.AliasName, GlobalParseFunc)
	if err != nil {
		panic(err.Error())
	}
	
	return nil
}