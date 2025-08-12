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
	// Сначала вычисляем аргументы
	evalArgs := make([]*Value, len(f.args))
	for i, arg := range f.args {
		evalArgs[i] = arg.Eval()
	}
	
	// Если функция зарегистрирована как перегружаемая, используем систему перегрузки
	if IsOverloadedMethod(f.funcName) {
		argTypes := GetArgTypesFromValues(evalArgs)
		
		// Пытаемся разрешить перегрузку
		overloadedFunc, err := ResolveMethodOverload(f.funcName, argTypes)
		if err != nil {
			panic(fmt.Sprintf("Overload resolution failed for '%s': %v", f.funcName, err))
		}
		
		// Вызываем найденную перегрузку
		return overloadedFunc.Call(evalArgs)
	}
	
	// Иначе проверяем обычную функцию (это поддерживает параметры по умолчанию)
	val, ok := scope.GlobalScope.Get(f.funcName)
	if ok {
		// Проверяем на TypedClosure (поддерживает параметры по умолчанию)
		if typedClosure, ok := val.Any().(*TypedClosure); ok {
			// Вызываем типизированное замыкание
			return typedClosure.Call(evalArgs)
		}
	}
	
	// Если перегрузок нет, используем обычную логику
	if !ok {
		panic("function '" + f.funcName + "' is not defined")
	}
	
	// Пробуем найти Callable объект (может быть FuncStatment или встроенная функция)
	if callable, ok := val.Any().(Callable); ok {
		// Вызываем функцию
		return callable.Call(evalArgs)
	}

	// Проверяем на Go-функцию (встроенные функции)
	if goFunc, ok := val.Any().(func([]*Value) *Value); ok {
		return goFunc(evalArgs)
	}

	// Старый код для совместимости с FuncStatment
	fnStatment, ok := val.Any().(*FuncStatment)
	if !ok {
		panic("'" + f.funcName + "' is not a function")
	}

	return fnStatment.Call(evalArgs)
}

func (f *FuncCallExpr) String() string {
	return fmt.Sprintf("%s()", f.funcName)
}
