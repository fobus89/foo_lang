package ast

import (
	"fmt"
	"foo_lang/scope"
)

// SafeVarAccess безопасно получает переменную из scope
func SafeVarAccess(name string) *Value {
	if value, exists := scope.GlobalScope.Get(name); exists {
		return value
	}
	
	return NewReferenceError(name, "variable")
}

// SafeFunctionCall безопасно вызывает функцию
func SafeFunctionCall(funcName string, args []*Value) *Value {
	// Проверяем существование функции
	funcValue, exists := scope.GlobalScope.Get(funcName)
	if !exists {
		return NewReferenceError(funcName, "function")
	}
	
	// Проверяем что это функция
	switch fn := funcValue.Any().(type) {
	case func([]*Value) *Value:
		// Встроенная функция - вызываем напрямую
		return fn(args)
		
	case *Closure:
		// Замыкание - проверяем аргументы
		expected := len(fn.args)
		got := len(args)
		
		if got != expected {
			return NewArgumentError(expected, got, funcName)
		}
		
		return fn.Call(args)
		
	default:
		err := NewErrorInfo(TypeError,
			fmt.Sprintf("'%s' is not a function", funcName),
			E102_UNDEFINED_FUNCTION).
			WithContext(funcName).
			WithSuggestion("Check that the name refers to a function")
		
		result := NewResultErr(NewValue(err))
		return NewValue(result)
	}
}

// SafeArrayAccess безопасно получает элемент массива
func SafeArrayAccess(array *Value, index *Value) *Value {
	// Проверяем что это массив
	arr, ok := array.Any().([]interface{})
	if !ok {
		return NewTypeError("array", fmt.Sprintf("%T", array.Any()), "array access")
	}
	
	// Проверяем что индекс - число
	idx, ok := index.Any().(int64)
	if !ok {
		return NewTypeError("integer", fmt.Sprintf("%T", index.Any()), "array index")
	}
	
	// Проверяем границы
	if idx < 0 || int(idx) >= len(arr) {
		return NewIndexError(int(idx), len(arr))
	}
	
	// Возвращаем элемент
	element := NewValue(arr[idx])
	result := NewResultOk(element)
	return NewValue(result)
}

// SafeObjectAccess безопасно получает свойство объекта
func SafeObjectAccess(object *Value, property *Value) *Value {
	// Проверяем что это объект
	obj, ok := object.Any().(map[string]interface{})
	if !ok {
		return NewTypeError("object", fmt.Sprintf("%T", object.Any()), "property access")
	}
	
	// Проверяем что свойство - строка
	prop, ok := property.Any().(string)
	if !ok {
		return NewTypeError("string", fmt.Sprintf("%T", property.Any()), "property name")
	}
	
	// Проверяем существование свойства
	if value, exists := obj[prop]; exists {
		element := NewValue(value)
		result := NewResultOk(element)
		return NewValue(result)
	}
	
	// Свойство не найдено
	err := NewErrorInfo(AttributeError,
		fmt.Sprintf("Property '%s' does not exist", prop),
		"E404").
		WithSuggestion(fmt.Sprintf("Check the spelling of '%s'", prop)).
		WithSuggestion("Use Object.keys() to see available properties")
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

// SafeArithmetic безопасно выполняет арифметические операции
func SafeArithmetic(left, right *Value, operation string) *Value {
	// Проверяем что оба операнда - числа
	leftNum, leftOk := getNumericValue(left)
	rightNum, rightOk := getNumericValue(right)
	
	if !leftOk {
		return NewTypeError("number", fmt.Sprintf("%T", left.Any()), operation)
	}
	
	if !rightOk {
		return NewTypeError("number", fmt.Sprintf("%T", right.Any()), operation)
	}
	
	var result float64
	switch operation {
	case "+":
		result = leftNum + rightNum
	case "-":
		result = leftNum - rightNum
	case "*":
		result = leftNum * rightNum
	case "/":
		if rightNum == 0 {
			return NewValueError("Division by zero", "arithmetic operation")
		}
		result = leftNum / rightNum
	case "%":
		if rightNum == 0 {
			return NewValueError("Modulo by zero", "arithmetic operation")
		}
		result = float64(int64(leftNum) % int64(rightNum))
	default:
		return NewValueError(fmt.Sprintf("Unknown operation: %s", operation), "arithmetic")
	}
	
	// Возвращаем результат
	value := NewValue(result)
	resultOk := NewResultOk(value)
	return NewValue(resultOk)
}

// getNumericValue извлекает числовое значение из Value
func getNumericValue(v *Value) (float64, bool) {
	switch val := v.Any().(type) {
	case int64:
		return float64(val), true
	case float64:
		return val, true
	case int:
		return float64(val), true
	case float32:
		return float64(val), true
	default:
		return 0, false
	}
}

// SafeMethodCall безопасно вызывает метод объекта
func SafeMethodCall(receiver *Value, methodName string, args []*Value) *Value {
	// Проверяем что у объекта есть данный метод
	switch receiverType := receiver.Any().(type) {
	case []interface{}:
		// Методы массива
		return safeArrayMethod(receiver, methodName, args)
		
	case map[string]interface{}:
		// Методы объекта
		return safeObjectMethod(receiver, methodName, args)
		
	case string:
		// Методы строки
		return safeStringMethod(receiver, methodName, args)
		
	case int64, float64:
		// Методы числа
		return safeNumberMethod(receiver, methodName, args)
		
	case bool:
		// Методы булева
		return safeBoolMethod(receiver, methodName, args)
		
	default:
		err := NewErrorInfo(AttributeError,
			fmt.Sprintf("Method '%s' not supported on type: %T", methodName, receiverType),
			"E405").
			WithContext(methodName).
			WithSuggestion("Check the available methods for this type")
		
		result := NewResultErr(NewValue(err))
		return NewValue(result)
	}
}

// safeArrayMethod безопасно вызывает методы массива
func safeArrayMethod(array *Value, methodName string, args []*Value) *Value {
	arr := array.Any().([]interface{})
	
	switch methodName {
	case "length":
		if len(args) != 0 {
			return NewArgumentError(0, len(args), "array.length")
		}
		value := NewValue(int64(len(arr)))
		result := NewResultOk(value)
		return NewValue(result)
		
	case "push":
		if len(args) != 1 {
			return NewArgumentError(1, len(args), "array.push")
		}
		// Здесь должна быть реализация push
		arr = append(arr, args[0].Any())
		value := NewValue(arr)
		result := NewResultOk(value)
		return NewValue(result)
		
	case "pop":
		if len(args) != 0 {
			return NewArgumentError(0, len(args), "array.pop")
		}
		if len(arr) == 0 {
			return NewEmptyContainerError("pop", "array")
		}
		// Возвращаем последний элемент
		lastElement := NewValue(arr[len(arr)-1])
		result := NewResultOk(lastElement)
		return NewValue(result)
		
	default:
		err := NewErrorInfo(AttributeError,
			fmt.Sprintf("Array has no method '%s'", methodName),
			"E405").
			WithSuggestion("Available methods: length, push, pop, slice, map, filter, reduce")
		
		result := NewResultErr(NewValue(err))
		return NewValue(result)
	}
}

// safeStringMethod безопасно вызывает методы строки
func safeStringMethod(str *Value, methodName string, args []*Value) *Value {
	s := str.Any().(string)
	
	switch methodName {
	case "length":
		if len(args) != 0 {
			return NewArgumentError(0, len(args), "string.length")
		}
		value := NewValue(int64(len(s)))
		result := NewResultOk(value)
		return NewValue(result)
		
	case "charAt":
		if len(args) != 1 {
			return NewArgumentError(1, len(args), "string.charAt")
		}
		
		idx, ok := args[0].Any().(int64)
		if !ok {
			return NewTypeError("integer", fmt.Sprintf("%T", args[0].Any()), "string.charAt index")
		}
		
		if idx < 0 || int(idx) >= len(s) {
			return NewIndexError(int(idx), len(s))
		}
		
		char := string(s[idx])
		value := NewValue(char)
		result := NewResultOk(value)
		return NewValue(result)
		
	default:
		err := NewErrorInfo(AttributeError,
			fmt.Sprintf("String has no method '%s'", methodName),
			"E405").
			WithSuggestion("Available methods: length, charAt, substring, toUpper, toLower")
		
		result := NewResultErr(NewValue(err))
		return NewValue(result)
	}
}

// safeNumberMethod и safeBoolMethod для полноты
func safeNumberMethod(num *Value, methodName string, args []*Value) *Value {
	err := NewErrorInfo(AttributeError,
		fmt.Sprintf("Number has no method '%s'", methodName),
		"E405").
		WithSuggestion("Available methods: toString, abs, round, floor, ceil")
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

func safeBoolMethod(b *Value, methodName string, args []*Value) *Value {
	err := NewErrorInfo(AttributeError,
		fmt.Sprintf("Boolean has no method '%s'", methodName),
		"E405").
		WithSuggestion("Available methods: toString, not")
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

func safeObjectMethod(obj *Value, methodName string, args []*Value) *Value {
	err := NewErrorInfo(AttributeError,
		fmt.Sprintf("Object has no method '%s'", methodName),
		"E405").
		WithSuggestion("Available methods depend on object type")
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}