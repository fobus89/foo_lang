package ast

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
)

// Macro представляет определение макроса
type Macro struct {
	Name   string
	Params []string
	Body   Expr // AST тело макроса
}

// MacroDefExpr представляет определение макроса
type MacroDefExpr struct {
	Name   string
	Params []string
	Body   Expr
}

func NewMacroDefExpr(name string, params []string, body Expr) *MacroDefExpr {
	return &MacroDefExpr{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (m *MacroDefExpr) Eval() *value.Value {
	// Сохраняем макрос в текущей области видимости
	macro := &Macro{
		Name:   m.Name,
		Params: m.Params,
		Body:   m.Body,
	}
	
	// Сохраняем макрос как специальное значение в scope
	scope.GlobalScope.Set(m.Name, value.NewValue(macro))
	
	// Возвращаем nil, так как определение макроса не производит значения
	return value.NewValue(nil)
}

// MacroCallExpr представляет вызов макроса
type MacroCallExpr struct {
	Name string
	Args []Expr
}

func NewMacroCallExpr(name string, args []Expr) *MacroCallExpr {
	return &MacroCallExpr{
		Name: name,
		Args: args,
	}
}

func (m *MacroCallExpr) Eval() *value.Value {
	// Получаем макрос из scope
	macroValue, found := scope.GlobalScope.Get(m.Name)
	if !found || macroValue == nil {
		panic(fmt.Sprintf("macro '%s' not found", m.Name))
	}
	
	macro, ok := macroValue.Any().(*Macro)
	if !ok {
		panic(fmt.Sprintf("'%s' is not a macro", m.Name))
	}
	
	// Проверяем количество аргументов
	if len(m.Args) != len(macro.Params) {
		panic(fmt.Sprintf("macro '%s' expects %d arguments, got %d", 
			m.Name, len(macro.Params), len(m.Args)))
	}
	
	// Создаем новую область видимости для макроса
	scope.GlobalScope.Push()
	defer scope.GlobalScope.Pop()
	
	// Связываем параметры с аргументами
	for i, param := range macro.Params {
		// Вычисляем аргументы перед передачей в макрос
		// В будущем можно добавить возможность передавать AST для метапрограммирования
		argValue := m.Args[i].Eval()
		scope.GlobalScope.Set(param, argValue)
	}
	
	// Выполняем тело макроса
	// В будущем здесь должна быть более сложная логика для манипуляции AST
	result := macro.Body.Eval()
	
	return result
}

// QuoteExpr представляет оператор quote для создания AST
type QuoteExpr struct {
	Expr Expr
}

func NewQuoteExpr(expr Expr) *QuoteExpr {
	return &QuoteExpr{Expr: expr}
}

func (q *QuoteExpr) Eval() *value.Value {
	// Quote возвращает AST как значение, не выполняя его
	return value.NewValue(q.Expr)
}

// UnquoteExpr представляет оператор unquote для вставки значений в quoted код
type UnquoteExpr struct {
	Expr Expr
}

func NewUnquoteExpr(expr Expr) *UnquoteExpr {
	return &UnquoteExpr{Expr: expr}
}

func (u *UnquoteExpr) Eval() *value.Value {
	// Unquote выполняет выражение и возвращает результат
	return u.Expr.Eval()
}

// ExpandExpr представляет оператор для раскрытия макроса в AST
type ExpandExpr struct {
	Expr Expr
}

func NewExpandExpr(expr Expr) *ExpandExpr {
	return &ExpandExpr{Expr: expr}
}

func (e *ExpandExpr) Eval() *value.Value {
	// Раскрываем макрос и возвращаем результирующий AST
	result := e.Expr.Eval()
	
	// Если результат - это AST, выполняем его
	if expr, ok := result.Any().(Expr); ok {
		return expr.Eval()
	}
	
	return result
}