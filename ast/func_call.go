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

	fnStatment, ok := val.Any().(*FuncStatment)
	if !ok {
		panic("'" + f.funcName + "' is not a function")
	}

	bodyStm := fnStatment.body.(*BodyExpr)

	expected := len(fnStatment.args)
	passed := len(f.args)

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
	for i, arg := range fnStatment.args {
		for name, expr := range arg {
			if i < len(f.args) {
				value := f.args[i].Eval()
				scope.GlobalScope.Set(name, value)
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

func (f *FuncCallExpr) String() string {
	return fmt.Sprintf("%s()", f.funcName)
}
