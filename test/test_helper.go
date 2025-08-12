package test

import (
	"foo_lang/builtin"
	"foo_lang/parser"
	"foo_lang/scope"
	"foo_lang/value"
)

// InitTestEnvironment инициализирует тестовое окружение через парсер
// Поддерживает старые сигнатуры для совместимости
func InitTestEnvironment(initFuncs ...interface{}) {
	// Создаем парсер с пустым кодом для инициализации scope
	p := parser.NewParser("")
	
	// Парсим пустой код для установки GlobalScope
	p.Parse()
	
	// Инициализируем функции, если они переданы
	for _, fn := range initFuncs {
		switch f := fn.(type) {
		case func(*scope.ScopeStack):
			f(scope.GlobalScope)
		case func(builtin.ScopeStack):
			f(scope.GlobalScope)
		case func(interface{ Set(name string, val *value.Value) }):
			f(scope.GlobalScope)
		default:
			// Игнорируем неизвестные типы для совместимости
		}
	}
}

// InitWithChannels инициализирует окружение с функциями каналов
func InitWithChannels() {
	InitTestEnvironment()
	builtin.InitializeMathFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	builtin.InitializeChannelFunctions(scope.GlobalScope)
}

// InitWithMath инициализирует окружение с математическими функциями
func InitWithMath() {
	InitTestEnvironment()
	builtin.InitializeMathFunctions(scope.GlobalScope)
}

// InitWithAll инициализирует окружение со всеми функциями
func InitWithAll() {
	InitTestEnvironment()
	builtin.InitializeMathFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	builtin.InitializeFilesystemFunctions(scope.GlobalScope)
	builtin.InitializeHttpFunctions(scope.GlobalScope)
	builtin.InitializeChannelFunctions(scope.GlobalScope)
	builtin.InitializeTimeFunctions(scope.GlobalScope)
	builtin.InitializeCryptoFunctions(scope.GlobalScope)
	builtin.InitializeRegexFunctions(scope.GlobalScope)
	builtin.InitializeSyncFunctions(scope.GlobalScope)
}