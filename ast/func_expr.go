package ast

import (
	"fmt"
	"foo_lang/scope"
)

// FuncParam представляет типизированный параметр функции
type FuncParam struct {
	Name     string
	TypeName string // "int", "string", "float", "bool", или пустая строка для нетипизированного
	Default  Expr   // Значение по умолчанию (может быть nil)
}

type FuncStatment struct {
	funcName string
	args     []map[string]Expr
	body     Expr
	isMacro  bool
}

// TypedFuncStatement представляет функцию с типизированными параметрами
type TypedFuncStatement struct {
	FuncName string
	Params   []FuncParam
	Body     Expr
}

func NewTypedFuncStatement(funcName string, params []FuncParam, body Expr) *TypedFuncStatement {
	return &TypedFuncStatement{
		FuncName: funcName,
		Params:   params,
		Body:     body,
	}
}

func (f *TypedFuncStatement) Eval() *Value {
	// Преобразуем типизированные параметры в старый формат для совместимости с замыканиями
	args := make([]map[string]Expr, len(f.Params))
	for i, param := range f.Params {
		argMap := make(map[string]Expr)
		argMap[param.Name] = param.Default // Если Default nil, это будет nil
		args[i] = argMap
	}
	
	// Создаем замыкание с типизированными параметрами
	closure := NewTypedClosure(f.FuncName, f.Params, f.Body)
	
	// Регистрируем замыкание в области видимости
	scope.GlobalScope.Set(f.FuncName, NewValue(closure))
	
	return NewValue(nil)
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
