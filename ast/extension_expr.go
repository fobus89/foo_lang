package ast

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
)

// ExtensionMethodInfo представляет информацию о методе расширения
type ExtensionMethodInfo struct {
	Name         string
	Params       []string
	Defaults     []Expr
	Body         Expr
	GenericParams []string
	ParamTypes   []string
	ReturnType   string
}

// ExtensionExpr представляет расширение типа новыми методами
// extension string {
//     fn isPalindrome() -> bool { ... }
// }
type ExtensionExpr struct {
	TypeName string                 // Имя типа для расширения (string, int, etc)
	Methods  []*ExtensionMethodInfo // Методы расширения
}

func (e *ExtensionExpr) Eval() *value.Value {
	// Регистрируем методы расширения для типа
	for _, method := range e.Methods {
		// Создаем функцию-обертку, которая будет принимать this как первый параметр
		extensionMethod := &ExtensionMethodWrapper{
			Method: method,
			TypeName: e.TypeName,
		}
		
		// Регистрируем метод в системе типов
		value.RegisterExtensionMethod(e.TypeName, method.Name, extensionMethod)
	}
	
	// Extension не возвращает значение, просто регистрирует методы
	return value.NewValue(nil)
}

func (e *ExtensionExpr) String() string {
	return fmt.Sprintf("extension %s { %d methods }", e.TypeName, len(e.Methods))
}

// ExtensionMethodWrapper оборачивает метод расширения для вызова с правильным контекстом
type ExtensionMethodWrapper struct {
	Method   *ExtensionMethodInfo
	TypeName string
}

func (w *ExtensionMethodWrapper) Call(receiver *value.Value, args []*value.Value) *value.Value {
	// Сохраняем текущую область видимости и создаем новую
	scope.GlobalScope.Push()
	defer scope.GlobalScope.Pop()
	
	// Добавляем this в область видимости
	scope.GlobalScope.Set("this", receiver)
	
	// Добавляем параметры метода
	for i, param := range w.Method.Params {
		if i < len(args) {
			scope.GlobalScope.Set(param, args[i])
		} else if w.Method.Defaults != nil && i < len(w.Method.Defaults) && w.Method.Defaults[i] != nil {
			// Используем значение по умолчанию
			defaultValue := w.Method.Defaults[i].Eval()
			scope.GlobalScope.Set(param, defaultValue)
		} else {
			scope.GlobalScope.Set(param, value.NewValue(nil))
		}
	}
	
	// Выполняем тело метода
	result := w.Method.Body.Eval()
	
	// Убираем флаг возврата, если он установлен
	if result != nil && result.IsReturn() {
		result.SetReturn(false)
	}
	
	return result
}