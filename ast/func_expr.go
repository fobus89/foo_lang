package ast

import (
	"fmt"
	"foo_lang/scope"
)

type FuncStatment struct {
	funcName string
	args     []map[string]Expr
	body     Expr
	isMacro  bool
}

func NewFuncStatment(funcName string, args []map[string]Expr, body Expr, isMacro bool) *FuncStatment {
	f := &FuncStatment{
		funcName: funcName,
		args:     args,
		body:     body,
		isMacro:  isMacro,
	}

	// Создаем замыкание для захвата переменных из текущей области видимости
	closure := NewClosure(funcName, args, body, isMacro)
	
	// Регистрируем замыкание в области видимости
	scope.GlobalScope.Set(funcName, NewValue(closure))

	return f
}

func (f *FuncStatment) Name() string {
	return f.funcName
}

func (f *FuncStatment) Params() []string {
	var args []string
	{
		for _, v := range f.args {
			for k := range v {
				args = append(args, k)
			}
		}
	}
	return args
}

func (f *FuncStatment) Eval() *Value {
	// Function definitions don't return values, they register the function in scope
	// The function is already registered in NewFuncStatment constructor
	return nil
}

func (f *FuncStatment) IsMacro() bool {
	return f.isMacro
}

func (f *FuncStatment) String() string {
	return fmt.Sprintf("func %s(%s) { %s }", f.funcName, f.args, f.body)
}

// Реализация интерфейса Callable для FuncStatment
func (f *FuncStatment) Call(args []*Value) *Value {
	bodyStm := f.body.(*BodyExpr)

	expected := len(f.args)
	passed := len(args)

	if passed > expected {
		panic(fmt.Sprintf("too many arguments: expected %d, got %d", expected, passed))
	}

	// Создаем новую область видимости для функции с проверкой рекурсии
	err := scope.GlobalScope.PushFunction()
	if err != nil {
		panic(err.Error())
	}
	defer scope.GlobalScope.PopFunction()

	// Устанавливаем параметры функции в локальной области
	for i, arg := range f.args {
		for name, expr := range arg {
			if i < len(args) {
				scope.GlobalScope.Set(name, args[i])
			} else if expr != nil {
				defaultValue := expr.Eval()
				scope.GlobalScope.Set(name, defaultValue)
			} else {
				panic(fmt.Sprintf("missing required argument: %s", name))
			}
		}
	}

	// Выполняем тело функции
	for _, stm := range bodyStm.Statments {
		if stm == nil {
			continue
		}

		result := stm.Eval()
		
		// Проверяем на return
		if result != nil && result.IsReturn() {
			return result
		}
	}

	return nil
}
