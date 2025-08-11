package ast

import (
	"foo_lang/scope"
	"foo_lang/value"
)

// TypeNameExpr представляет прямое обращение к типу по имени (User, int, string, etc)
// Используется в вызовах макросов: @macro(User) вместо @macro(type(User))
type TypeNameExpr struct {
	TypeName string
}

func NewTypeNameExpr(typeName string) *TypeNameExpr {
	return &TypeNameExpr{TypeName: typeName}
}

func (t *TypeNameExpr) Eval() *value.Value {
	// ТОЧНО КОПИРУЕМ логику из TypeExpr.Eval() для совместимости
	
	// Сначала проверяем специальный ключ для TypeInfo у enum
	if typeInfoValue, found := scope.GlobalScope.Get(t.TypeName + "__TypeInfo"); found {
		return typeInfoValue
	}
	
	// Получаем информацию о типе из scope
	typeValue, found := scope.GlobalScope.Get(t.TypeName)
	if !found {
		// Если не найден пользовательский тип, проверяем примитивные типы
		switch t.TypeName {
		case "int", "integer":
			return value.NewValue(NewPrimitiveTypeInfo("int"))
		case "float", "double":
			return value.NewValue(NewPrimitiveTypeInfo("float"))
		case "string":
			return value.NewValue(NewPrimitiveTypeInfo("string"))
		case "bool", "boolean":
			return value.NewValue(NewPrimitiveTypeInfo("bool"))
		case "array":
			return value.NewValue(NewPrimitiveTypeInfo("array"))
		default:
			panic("type '" + t.TypeName + "' not found")
		}
	}
	
	// Если это уже TypeInfo, возвращаем как есть
	if _, ok := typeValue.Any().(*TypeInfo); ok {
		return typeValue
	}
	
	// ТОЧНО КОПИРУЕМ поведение TypeExpr: возвращаем значение как есть
	return typeValue
}