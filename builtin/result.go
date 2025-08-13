package builtin

import (
	"foo_lang/ast"
	"foo_lang/scope"
	"foo_lang/value"
)

// InitializeResultFunctions инициализирует Result функции Ok и Err
func InitializeResultFunctions(globalScope *scope.ScopeStack) {
	// Ok функция - создает Ok(value)
	okFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewValue("Error: Ok() expects exactly 1 argument")
		}
		
		// Конвертируем value.Value в ast.Value 
		astValue := ast.NewValue(args[0].Any())
		result := ast.NewResultOk(astValue)
		
		// Конвертируем обратно в value.Value
		return value.NewValue(result)
	}
	globalScope.Set("Ok", value.NewValue(okFunc))
	
	// Err функция - создает Err(error)
	errFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewValue("Error: Err() expects exactly 1 argument")
		}
		
		// Конвертируем value.Value в ast.Value
		astValue := ast.NewValue(args[0].Any())
		result := ast.NewResultErr(astValue)
		
		// Конвертируем обратно в value.Value
		return value.NewValue(result)
	}
	globalScope.Set("Err", value.NewValue(errFunc))
}

// InitializeResultMethods инициализирует методы для Result типов
func InitializeResultMethods() {
	// Эти методы уже реализованы в method_call_expr.go
	// для типа *ast.ResultValue:
	// - isOk()
	// - isErr() 
	// - unwrap()
	// - unwrapOr(default)
}