package ast

import (
	"foo_lang/scope"
)

// StructInstanceExpr представляет создание экземпляра структуры: TypeName{field: value, ...}
type StructInstanceExpr struct {
	TypeName string
	Fields   map[string]Expr
}

func NewStructInstanceExpr(typeName string, fields map[string]Expr) *StructInstanceExpr {
	return &StructInstanceExpr{
		TypeName: typeName,
		Fields:   fields,
	}
}

func (s *StructInstanceExpr) Eval() *Value {
	// Получаем информацию о типе структуры из глобального scope
	typeValue, found := scope.GlobalScope.Get(s.TypeName)
	if !found {
		panic("unknown struct type: " + s.TypeName)
	}
	
	// Проверяем, что это действительно TypeInfo структуры
	typeInfo, ok := typeValue.Any().(*TypeInfo)
	if !ok || typeInfo.Kind != "struct" {
		panic(s.TypeName + " is not a struct type")
	}
	
	// Создаем экземпляр структуры
	fields := make(map[string]*Value)
	
	// Вычисляем значения полей
	for fieldName, fieldExpr := range s.Fields {
		// Проверяем, что поле существует в определении структуры
		if _, exists := typeInfo.Fields[fieldName]; !exists {
			panic("field '" + fieldName + "' does not exist in struct " + s.TypeName)
		}
		
		// Вычисляем значение поля
		fields[fieldName] = fieldExpr.Eval()
	}
	
	// Устанавливаем значения по умолчанию для незаданных полей
	for fieldName := range typeInfo.Fields {
		if _, exists := fields[fieldName]; !exists {
			// Пока что используем nil для незаданных полей
			fields[fieldName] = NewValue(nil)
		}
	}
	
	// Создаем объект структуры
	structObj := NewStructObject(typeInfo, fields)
	
	return NewValue(structObj)
}