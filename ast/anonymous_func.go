package ast

import (
	"fmt"
	"foo_lang/scope"
)

// AnonymousFunc представляет анонимную функцию (лямбду)
type AnonymousFunc struct {
	args []map[string]Expr  // параметры функции
	body Expr                // тело функции
}

// NewAnonymousFunc создает новую анонимную функцию
func NewAnonymousFunc(args []map[string]Expr, body Expr) *AnonymousFunc {
	return &AnonymousFunc{
		args: args,
		body: body,
	}
}

// Eval возвращает замыкание для анонимной функции
func (af *AnonymousFunc) Eval() *Value {
	// Создаем замыкание для анонимной функции
	closure := NewClosure("(anonymous)", af.args, af.body, false)
	return NewValue(closure)
}

// Name возвращает имя функции (для интерфейса Callable)
func (af *AnonymousFunc) Name() string {
	return "(anonymous)"
}

// Call выполняет анонимную функцию
func (af *AnonymousFunc) Call(args []*Value) *Value {
	expected := len(af.args)
	passed := len(args)

	if passed > expected {
		panic(fmt.Sprintf("too many arguments: expected %d, got %d", expected, passed))
	}

	// Создаем новую область видимости для функции
	err := scope.GlobalScope.PushFunction()
	if err != nil {
		panic(err.Error())
	}
	defer scope.GlobalScope.PopFunction()

	// Устанавливаем параметры функции
	for i, arg := range af.args {
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
	if bodyStm, ok := af.body.(*BodyExpr); ok {
		// Если тело - блок выражений
		for _, stm := range bodyStm.Statments {
			if stm == nil {
				continue
			}
			result := stm.Eval()
			if result != nil && result.IsReturn() {
				return result
			}
		}
	} else {
		// Если тело - одно выражение (для стрелочных функций)
		result := af.body.Eval()
		if result != nil {
			// Для одиночных выражений автоматически возвращаем результат
			retResult := NewValue(result.Any())
			retResult.SetReturn(true)
			return retResult
		}
	}

	return nil
}

// String возвращает строковое представление
func (af *AnonymousFunc) String() string {
	return fmt.Sprintf("(anonymous function with %d args)", len(af.args))
}