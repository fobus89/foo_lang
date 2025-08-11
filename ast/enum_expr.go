package ast

import (
	"foo_lang/scope"
	"foo_lang/value"
)

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
	// Создаём объект с enum значениями для обратной совместимости
	enumObj := make(map[string]*Value)
	
	for i, value := range e.Values {
		enumObj[value] = NewValue(int64(i))
	}
	
	// Создаем TypeInfo для enum
	enumTypeInfo := NewEnumTypeInfo(e.Name, e.Values)
	
	// Сохраняем enum как объект для обычного использования
	scope.GlobalScope.Set(e.Name, NewValue(enumObj))
	
	// ДОПОЛНИТЕЛЬНО: Сохраняем TypeInfo для использования в type() и макросах
	scope.GlobalScope.Set(e.Name+"__TypeInfo", value.NewValue(enumTypeInfo))
	
	return value.NewValue(enumTypeInfo)
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