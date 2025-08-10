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
	
	// Пока что простая реализация - для массивов добавим встроенные методы
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