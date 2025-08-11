package ast

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
)

// MacroParam представляет типизированный параметр макроса
type MacroParam struct {
	Name     string
	TypeName string // "Type", "FnType", "StructType", "EnumType", или пустая строка для нетипизированного
}

// Macro представляет определение макроса с поддержкой macro-time execution
type Macro struct {
	Name         string
	Params       []MacroParam // Типизированные параметры
	MacroTime    []Expr       // Код, выполняющийся во время макро-времени
	CodeGenBody  Expr         // Expr блок для генерации кода
}

// MacroDefExpr представляет определение макроса с поддержкой macro-time и code-gen
type MacroDefExpr struct {
	Name        string
	Params      []MacroParam // Типизированные параметры
	MacroTime   []Expr       // Выполняется во время компиляции макроса
	CodeGenBody Expr         // Expr {} блок для генерации кода
}

func NewMacroDefExpr(name string, params []MacroParam, macroTime []Expr, codeGen Expr) *MacroDefExpr {
	return &MacroDefExpr{
		Name:        name,
		Params:      params,
		MacroTime:   macroTime,
		CodeGenBody: codeGen,
	}
}

// NewSimpleMacroDefExpr создает макрос с единым телом (для обратной совместимости)
func NewSimpleMacroDefExpr(name string, params []MacroParam, body Expr) *MacroDefExpr {
	return &MacroDefExpr{
		Name:        name,
		Params:      params,
		MacroTime:   nil,       // Нет macro-time кода
		CodeGenBody: body,      // Все тело как генерация кода
	}
}

// NewLegacyMacroDefExpr создает макрос из строковых параметров (для обратной совместимости) 
func NewLegacyMacroDefExpr(name string, paramNames []string, body Expr) *MacroDefExpr {
	params := make([]MacroParam, len(paramNames))
	for i, paramName := range paramNames {
		params[i] = MacroParam{Name: paramName, TypeName: ""} // Нетипизированный параметр
	}
	return NewSimpleMacroDefExpr(name, params, body)
}

func (m *MacroDefExpr) Eval() *value.Value {
	// Сохраняем макрос в текущей области видимости
	macro := &Macro{
		Name:        m.Name,
		Params:      m.Params,
		MacroTime:   m.MacroTime,
		CodeGenBody: m.CodeGenBody,
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
	
	// Связываем параметры с аргументами с проверкой типов
	for i, param := range macro.Params {
		// Вычисляем аргументы перед передачей в макрос
		argValue := m.Args[i].Eval()
		
		// Проверяем тип параметра, если указан
		if param.TypeName != "" {
			if err := validateMacroParameterType(argValue, param.TypeName); err != nil {
				panic(fmt.Sprintf("macro '%s' parameter '%s': %s", m.Name, param.Name, err.Error()))
			}
		}
		
		scope.GlobalScope.Set(param.Name, argValue)
	}
	
	// ФАЗА 1: Выполняем macro-time код
	if macro.MacroTime != nil {
		for _, stmt := range macro.MacroTime {
			stmt.Eval()
		}
	}
	
	// ФАЗА 2: Генерируем и выполняем код из Expr блока
	if macro.CodeGenBody != nil {
		result := macro.CodeGenBody.Eval()
		return result
	}
	
	// Если нет кода для генерации, возвращаем nil
	return value.NewValue(nil)
}

// validateMacroParameterType проверяет соответствие типа аргумента ожидаемому типу параметра макроса
func validateMacroParameterType(argValue *value.Value, expectedTypeName string) error {
	// Получаем TypeInfo из аргумента
	typeInfo, ok := argValue.Any().(*TypeInfo)
	if !ok {
		return fmt.Errorf("expected TypeInfo, got %T", argValue.Any())
	}
	
	switch expectedTypeName {
	case "Type":
		// Type принимает любой тип
		return nil
	case "FnType":
		if typeInfo.Kind != "function" {
			return fmt.Errorf("expected function type, got %s", typeInfo.Kind)
		}
		return nil
	case "StructType":
		if typeInfo.Kind != "struct" {
			return fmt.Errorf("expected struct type, got %s", typeInfo.Kind)
		}
		return nil
	case "EnumType":
		if typeInfo.Kind != "enum" {
			return fmt.Errorf("expected enum type, got %s", typeInfo.Kind)
		}
		return nil
	default:
		return fmt.Errorf("unknown type constraint: %s", expectedTypeName)
	}
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

// ExprBlockExpr представляет блок Expr {} для генерации кода в макросах
type ExprBlockExpr struct {
	Statements []Expr
}

func NewExprBlockExpr(statements []Expr) *ExprBlockExpr {
	return &ExprBlockExpr{Statements: statements}
}

func (e *ExprBlockExpr) Eval() *value.Value {
	// Выполняем все выражения в блоке для генерации кода
	var result *value.Value = value.NewValue(nil) // Инициализируем значением по умолчанию
	
	for _, stmt := range e.Statements {
		result = stmt.Eval()
		
		// Проверяем специальные флаги (если result не nil)
		if result != nil && (result.IsReturn() || result.IsBreak()) {
			break
		}
	}
	
	// Всегда возвращаем непустое значение
	if result == nil {
		return value.NewValue(nil)
	}
	
	return result
}