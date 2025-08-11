package ast

import (
	"fmt"
	"foo_lang/scope"
)

type FuncCallExpr struct {
	funcName string
	args     []Expr
}

func NewFuncCallExpr(funcName string, args []Expr) *FuncCallExpr {

	f := &FuncCallExpr{
		funcName: funcName,
		args:     args,
	}

	return f
}

func (f *FuncCallExpr) Eval() *Value {
	val, ok := scope.GlobalScope.Get(f.funcName)
	if !ok {
		panic("function '" + f.funcName + "' is not defined")
	}

	// Проверяем на TypedClosure
	if typedClosure, ok := val.Any().(*TypedClosure); ok {
		// Вычисляем аргументы
		evalArgs := make([]*Value, len(f.args))
		for i, arg := range f.args {
			evalArgs[i] = arg.Eval()
		}
		
		// Вызываем типизированное замыкание
		return typedClosure.Call(evalArgs)
	}
	
	// Пробуем найти Callable объект (может быть FuncStatment или встроенная функция)
	if callable, ok := val.Any().(Callable); ok {
		// Вычисляем аргументы
		evalArgs := make([]*Value, len(f.args))
		for i, arg := range f.args {
			evalArgs[i] = arg.Eval()
		}
		
		// Вызываем функцию
		return callable.Call(evalArgs)
	}

	// Старый код для совместимости с FuncStatment
	fnStatment, ok := val.Any().(*FuncStatment)
	if !ok {
		panic("'" + f.funcName + "' is not a function")
	}

	// Вычисляем аргументы
	evalArgs := make([]*Value, len(f.args))
	for i, arg := range f.args {
		evalArgs[i] = arg.Eval()
	}

	return fnStatment.Call(evalArgs)
}

func (f *FuncCallExpr) String() string {
	return fmt.Sprintf("%s()", f.funcName)
}
