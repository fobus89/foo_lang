package ast

// Callable представляет объект, который может быть вызван как функция
type Callable interface {
	// Call выполняет вызов с переданными аргументами
	Call(args []*Value) *Value
	
	// Name возвращает имя функции (для отладки)
	Name() string
}