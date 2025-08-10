package ast

import "foo_lang/scope"

// EnumExpr представляет определение enum
type EnumExpr struct {
	Name   string
	Values []string
}

func NewEnumExpr(name string, values []string) *EnumExpr {
	return &EnumExpr{
		Name:   name,
		Values: values,
	}
}

func (e *EnumExpr) Eval() *Value {
	// Создаём объект с enum значениями
	enumObj := make(map[string]*Value)
	
	for i, value := range e.Values {
		enumObj[value] = NewValue(int64(i))
	}
	
	// Сохраняем enum в scope
	scope.GlobalScope.Set(e.Name, NewValue(enumObj))
	
	return nil
}

// EnumValueExpr представляет доступ к значению enum (Color.RED)
type EnumValueExpr struct {
	EnumName string
	Value    string
}

func NewEnumValueExpr(enumName, value string) *EnumValueExpr {
	return &EnumValueExpr{
		EnumName: enumName,
		Value:    value,
	}
}

func (e *EnumValueExpr) Eval() *Value {
	// Получаем enum из scope
	enumVal, ok := scope.GlobalScope.Get(e.EnumName)
	if !ok {
		panic("enum '" + e.EnumName + "' is not defined")
	}
	
	// Проверяем, что это объект
	if enumObj, ok := enumVal.Any().(map[string]*Value); ok {
		if value, exists := enumObj[e.Value]; exists {
			return value
		}
		panic("enum value '" + e.Value + "' does not exist in " + e.EnumName)
	}
	
	panic("'" + e.EnumName + "' is not an enum")
}