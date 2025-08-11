package ast

// MethodCallExpr представляет вызов метода объекта (object.method(args))
type MethodCallExpr struct {
	Object     Expr
	MethodName string
	Args       []Expr
}

func NewMethodCallExpr(object Expr, methodName string, args []Expr) *MethodCallExpr {
	return &MethodCallExpr{
		Object:     object,
		MethodName: methodName,
		Args:       args,
	}
}

func (m *MethodCallExpr) Eval() *Value {
	obj := m.Object.Eval()
	
	// Методы для TypeInfo
	if typeInfo, ok := obj.Any().(*TypeInfo); ok {
		switch m.MethodName {
		case "String":
			if len(m.Args) != 0 {
				panic("String() expects no arguments")
			}
			return NewValue(typeInfo.String())
		case "GetFieldNames":
			if len(m.Args) != 0 {
				panic("GetFieldNames() expects no arguments")
			}
			names := typeInfo.GetFieldNames()
			values := make([]*Value, len(names))
			for i, name := range names {
				values[i] = NewValue(name)
			}
			return NewValue(values)
		case "GetFieldType":
			if len(m.Args) != 1 {
				panic("GetFieldType() expects exactly 1 argument")
			}
			fieldName := m.Args[0].Eval().String()
			fieldType := typeInfo.GetFieldType(fieldName)
			if fieldType == nil {
				return NewValue(nil)
			}
			return NewValue(fieldType)
		case "HasField":
			if len(m.Args) != 1 {
				panic("HasField() expects exactly 1 argument")
			}
			fieldName := m.Args[0].Eval().String()
			return NewValue(typeInfo.HasField(fieldName))
		// Полиморфные методы проверки типов
		case "isStruct":
			if len(m.Args) != 0 {
				panic("isStruct() expects no arguments")
			}
			return NewValue(typeInfo.Kind == "struct")
		case "isFunction":
			if len(m.Args) != 0 {
				panic("isFunction() expects no arguments")
			}
			return NewValue(typeInfo.Kind == "function")
		case "isEnum":
			if len(m.Args) != 0 {
				panic("isEnum() expects no arguments")
			}
			return NewValue(typeInfo.Kind == "enum")
		case "isPrimitive":
			if len(m.Args) != 0 {
				panic("isPrimitive() expects no arguments")
			}
			return NewValue(typeInfo.Kind == "primitive")
		// Методы преобразования
		case "toStruct":
			if len(m.Args) != 0 {
				panic("toStruct() expects no arguments")
			}
			if typeInfo.Kind != "struct" {
				panic("cannot convert " + typeInfo.Kind + " to struct")
			}
			return NewValue(typeInfo)
		case "toFunction":
			if len(m.Args) != 0 {
				panic("toFunction() expects no arguments")
			}
			if typeInfo.Kind != "function" {
				panic("cannot convert " + typeInfo.Kind + " to function")
			}
			return NewValue(typeInfo)
		case "toEnum":
			if len(m.Args) != 0 {
				panic("toEnum() expects no arguments")
			}
			if typeInfo.Kind != "enum" {
				panic("cannot convert " + typeInfo.Kind + " to enum")
			}
			return NewValue(typeInfo)
		}
	}
	
	// Методы для Result типа
	if result, ok := obj.Any().(*ResultValue); ok {
		switch m.MethodName {
		case "isOk":
			if len(m.Args) != 0 {
				panic("isOk() expects no arguments")
			}
			return NewValue(result.IsOk())
		case "isErr":
			if len(m.Args) != 0 {
				panic("isErr() expects no arguments")
			}
			return NewValue(result.IsErr())
		case "unwrap":
			if len(m.Args) != 0 {
				panic("unwrap() expects no arguments")
			}
			return result.Unwrap()
		case "unwrapOr":
			if len(m.Args) != 1 {
				panic("unwrapOr() expects exactly 1 argument")
			}
			defaultValue := m.Args[0].Eval()
			return result.UnwrapOr(defaultValue)
		}
	}
	
	// Методы для массивов
	if arr, ok := obj.Any().([]any); ok {
		switch m.MethodName {
		case "push":
			if len(m.Args) != 1 {
				panic("push() expects exactly 1 argument")
			}
			newItem := m.Args[0].Eval().Any()
			newArr := append(arr, newItem)
			return NewValue(newArr)
		case "length":
			if len(m.Args) != 0 {
				panic("length() expects no arguments")
			}
			return NewValue(int64(len(arr)))
		}
	}
	
	// Для других объектов - пока просто ошибка
	panic("method '" + m.MethodName + "' not supported on this type")
}