package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"foo_lang/scope"
	"foo_lang/value"
)

// Module represents a loaded module
type Module struct {
	Path     string                     // File path of the module
	Exports  map[string]*value.Value    // Exported variables/functions
	Loaded   bool                       // Whether the module has been loaded
	Scope    *scope.ScopeStack         // Module's own scope
}

// ModuleCache stores loaded modules to prevent re-loading
var ModuleCache = make(map[string]*Module)

// ParseFunc represents a function that can parse code and return AST expressions
type ParseFunc func(string) []Expr

// Expr interface for AST expressions (to avoid circular import)
type Expr interface {
	Eval() *value.Value
}

// LoadModule loads a module from the given path
func LoadModule(modulePath string, parseFunc ParseFunc) (*Module, error) {
	// Normalize the path
	absPath, err := filepath.Abs(modulePath)
	if err != nil {
		return nil, fmt.Errorf("invalid module path: %s", modulePath)
	}
	
	// Check if already loaded
	if module, exists := ModuleCache[absPath]; exists {
		return module, nil
	}
	
	// Read the module file
	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read module file: %s", err)
	}
	
	// Create new module
	module := &Module{
		Path:    absPath,
		Exports: make(map[string]*value.Value),
		Loaded:  false,
		Scope:   scope.NewScopeStack(), // Each module has its own scope
	}
	
	// Save current global scope
	originalScope := scope.GlobalScope
	
	// Set module scope as global temporarily
	scope.GlobalScope = module.Scope
	
	// Parse and execute the module
	exprs := parseFunc(string(content))
	for _, expr := range exprs {
		expr.Eval()
	}
	
	// Collect exported items (items with "__export_" prefix)
	for name, val := range module.Scope.GetAll() {
		if strings.HasPrefix(name, "__export_") {
			exportName := strings.TrimPrefix(name, "__export_")
			module.Exports[exportName] = val
		}
	}
	
	// Restore original scope
	scope.GlobalScope = originalScope
	
	module.Loaded = true
	ModuleCache[absPath] = module
	
	return module, nil
}

// ResolveModulePath resolves relative module paths
func ResolveModulePath(currentFile, importPath string) string {
	if filepath.IsAbs(importPath) {
		return importPath
	}
	
	// Handle relative paths
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		dir := filepath.Dir(currentFile)
		return filepath.Join(dir, importPath)
	}
	
	// For now, treat all other paths as relative to current directory
	return importPath
}

// ImportModule imports items from a module into current scope
func ImportModule(modulePath string, currentFile string, importedItems []string, alias string, parseFunc ParseFunc) error {
	// Resolve the path
	resolvedPath := ResolveModulePath(currentFile, modulePath)
	
	// Load the module
	module, err := LoadModule(resolvedPath, parseFunc)
	if err != nil {
		return err
	}
	
	// Import items into current scope
	if alias != "" {
		// Import as alias: import * as ModuleName from "./module"
		moduleObj := make(map[string]*value.Value)
		for name, val := range module.Exports {
			moduleObj[name] = val
		}
		scope.GlobalScope.Set(alias, value.NewValue(moduleObj))
	} else if len(importedItems) == 0 {
		// Import all: import "./module"
		for name, val := range module.Exports {
			scope.GlobalScope.Set(name, val)
		}
	} else {
		// Selective import: import { item1, item2 } from "./module"
		for _, itemName := range importedItems {
			if val, exists := module.Exports[itemName]; exists {
				scope.GlobalScope.Set(itemName, val)
			} else {
				return fmt.Errorf("module %s does not export '%s'", modulePath, itemName)
			}
		}
	}
	
	return nil
}