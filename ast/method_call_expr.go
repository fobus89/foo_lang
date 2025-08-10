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