package ast

import "foo_lang/value"

// IndexExpr представляет индексацию массива или объекта (arr[index])
type IndexExpr struct {
	Object Expr // Объект или массив
	Index  Expr // Индекс
}

func NewIndexExpr(object, index Expr) *IndexExpr {
	return &IndexExpr{
		Object: object,
		Index:  index,
	}
}

func (i *IndexExpr) Eval() *Value {
	obj := i.Object.Eval()
	idx := i.Index.Eval()
	
	// Для массивов
	if arr, ok := obj.Any().([]any); ok {
		// Индекс должен быть числом
		if !idx.IsNumber() {
			panic("array index must be an integer")
		}
		
		index := idx.Int()
		
		// Проверка границ
		if index < 0 || index >= len(arr) {
			panic("array index out of bounds")
		}
		
		return value.NewValue(arr[index])
	}
	
	// Для объектов (строковый индекс)
	if objMap, ok := obj.Any().(map[string]*Value); ok {
		if !idx.IsString() {
			panic("object property name must be a string")
		}
		
		key := idx.String()
		if val, exists := objMap[key]; exists {
			return val
		} else {
			return value.NewValue(nil) // undefined
		}
	}
	
	panic("cannot index non-array, non-object value")
}