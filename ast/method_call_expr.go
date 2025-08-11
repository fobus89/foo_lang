package ast

import (
	"fmt"
	"strconv"
	"strings"
	"foo_lang/value"
	"foo_lang/scope"
)

// StructObject представляет экземпляр структуры
type StructObject struct {
	TypeInfo *TypeInfo
	Fields   map[string]*Value
}

// NewStructObject создает новый экземпляр структуры
func NewStructObject(typeInfo *TypeInfo, fields map[string]*Value) *StructObject {
	return &StructObject{
		TypeInfo: typeInfo,
		Fields:   fields,
	}
}

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
	
	// Методы для экземпляров структур
	if structObj, ok := obj.Any().(*StructObject); ok {
		// Проверяем все зарегистрированные интерфейсы для этого типа
		for _, interfaceName := range getImplementedInterfaces(structObj.TypeInfo.Name) {
			if impl := GetImplementation(structObj.TypeInfo.Name, interfaceName); impl != nil {
				// Ищем метод в реализации
				for _, method := range impl.Methods {
					if method.FuncName == m.MethodName {
						// Вычисляем аргументы
						args := make([]*Value, len(m.Args))
						for i, arg := range m.Args {
							args[i] = arg.Eval()
						}
						
						// Вызываем метод с 'this' контекстом
						return callMethodWithContext(method, obj, args)
					}
				}
			}
		}
		
		// Если метод не найден в интерфейсах, возможно это доступ к полю
		panic("method '" + m.MethodName + "' not found for struct " + structObj.TypeInfo.Name)
	}
	
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
		case "pop":
			if len(m.Args) != 0 {
				panic("pop() expects no arguments")
			}
			if len(arr) == 0 {
				panic("pop() called on empty array")
			}
			return NewValue(arr[len(arr)-1])
		case "slice":
			if len(m.Args) != 2 {
				panic("slice() expects exactly 2 arguments")
			}
			start := int(m.Args[0].Eval().Int64())
			end := int(m.Args[1].Eval().Int64())
			if start < 0 || end > len(arr) || start > end {
				panic("slice() index out of bounds")
			}
			return NewValue(arr[start:end])
		case "map":
			if len(m.Args) != 1 {
				panic("map() expects exactly 1 argument (function)")
			}
			
			// Получаем функцию из аргумента
			fnArg := m.Args[0].Eval()
			
			// Проверяем что это вызываемый объект
			if callable, ok := fnArg.Any().(Callable); ok {
				result := make([]any, len(arr))
				
				// Применяем функцию к каждому элементу
				for i, item := range arr {
					// Вызываем функцию с текущим элементом
					argValue := NewValue(item)
					mappedValue := callable.Call([]*Value{argValue})
					result[i] = mappedValue.Any()
				}
				
				return NewValue(result)
			}
			
			panic("map() argument must be a function")
		case "filter":
			if len(m.Args) != 1 {
				panic("filter() expects exactly 1 argument (predicate function)")
			}
			
			// Получаем функцию из аргумента
			fnArg := m.Args[0].Eval()
			
			// Проверяем что это вызываемый объект
			if callable, ok := fnArg.Any().(Callable); ok {
				var result []any
				
				// Фильтруем элементы через переданную функцию
				for _, item := range arr {
					// Вызываем функцию с текущим элементом
					argValue := NewValue(item)
					shouldInclude := callable.Call([]*Value{argValue})
					
					// Если функция вернула true, добавляем элемент
					if shouldInclude.Bool() {
						result = append(result, item)
					}
				}
				
				return NewValue(result)
			}
			
			panic("filter() argument must be a function")
		case "reduce":
			if len(m.Args) != 2 {
				panic("reduce() expects 2 arguments (initial value, reducer function)")
			}
			
			initialValue := m.Args[0].Eval()
			fnArg := m.Args[1].Eval()
			
			// Проверяем что второй аргумент - функция
			if callable, ok := fnArg.Any().(Callable); ok {
				accumulator := initialValue
				
				// Применяем reducer функцию к каждому элементу
				for _, item := range arr {
					// Вызываем функцию с аккумулятором и текущим элементом
					accValue := accumulator
					itemValue := NewValue(item)
					accumulator = callable.Call([]*Value{accValue, itemValue})
				}
				
				return accumulator
			}
			
			panic("reduce() second argument must be a function")
		}
	}
	
	// Методы для строк
	if str, ok := obj.Any().(string); ok {
		switch m.MethodName {
		case "length":
			if len(m.Args) != 0 {
				panic("string.length() expects no arguments")
			}
			return NewValue(int64(len(str)))
		case "charAt":
			if len(m.Args) != 1 {
				panic("string.charAt() expects exactly 1 argument")
			}
			index := int(m.Args[0].Eval().Int64())
			if index < 0 || index >= len(str) {
				panic("string.charAt() index out of bounds")
			}
			return NewValue(string(str[index]))
		case "substring":
			if len(m.Args) != 2 {
				panic("string.substring() expects exactly 2 arguments")
			}
			start := int(m.Args[0].Eval().Int64())
			end := int(m.Args[1].Eval().Int64())
			if start < 0 {
				start = 0
			}
			if end > len(str) {
				end = len(str)
			}
			if start > end {
				start = end
			}
			return NewValue(str[start:end])
		case "toUpper":
			if len(m.Args) != 0 {
				panic("string.toUpper() expects no arguments")
			}
			return NewValue(strings.ToUpper(str))
		case "toLower":
			if len(m.Args) != 0 {
				panic("string.toLower() expects no arguments")
			}
			return NewValue(strings.ToLower(str))
		}
	}
	
	// Методы для int64
	if num, ok := obj.Any().(int64); ok {
		switch m.MethodName {
		case "toString":
			if len(m.Args) != 0 {
				panic("int.toString() expects no arguments")
			}
			return NewValue(strconv.FormatInt(num, 10))
		case "abs":
			if len(m.Args) != 0 {
				panic("int.abs() expects no arguments")
			}
			if num < 0 {
				return NewValue(-num)
			}
			return NewValue(num)
		case "toFloat":
			if len(m.Args) != 0 {
				panic("int.toFloat() expects no arguments")
			}
			return NewValue(float64(num))
		}
	}
	
	// Методы для int (32-bit)
	if num, ok := obj.Any().(int); ok {
		switch m.MethodName {
		case "toString":
			if len(m.Args) != 0 {
				panic("int.toString() expects no arguments")
			}
			return NewValue(strconv.Itoa(num))
		case "abs":
			if len(m.Args) != 0 {
				panic("int.abs() expects no arguments")
			}
			if num < 0 {
				return NewValue(-num)
			}
			return NewValue(num)
		case "toFloat":
			if len(m.Args) != 0 {
				panic("int.toFloat() expects no arguments")
			}
			return NewValue(float64(num))
		}
	}
	
	// Методы для float64 (все числа в foo_lang)
	if num, ok := obj.Any().(float64); ok {
		switch m.MethodName {
		case "toString":
			if len(m.Args) != 0 {
				panic("number.toString() expects no arguments")
			}
			return NewValue(strconv.FormatFloat(num, 'f', -1, 64))
		case "round":
			if len(m.Args) != 0 {
				panic("number.round() expects no arguments")
			}
			return NewValue(float64(int(num + 0.5)))
		case "floor":
			if len(m.Args) != 0 {
				panic("number.floor() expects no arguments")
			}
			return NewValue(float64(int(num)))
		case "ceil":
			if len(m.Args) != 0 {
				panic("number.ceil() expects no arguments")
			}
			if num == float64(int(num)) {
				return NewValue(num)
			}
			return NewValue(float64(int(num) + 1))
		case "toInt":
			if len(m.Args) != 0 {
				panic("number.toInt() expects no arguments")
			}
			return NewValue(int64(num))
		// Добавляем методы которые были для int
		case "abs":
			if len(m.Args) != 0 {
				panic("number.abs() expects no arguments")
			}
			if num < 0 {
				return NewValue(-num)
			}
			return NewValue(num)
		case "toFloat":
			if len(m.Args) != 0 {
				panic("number.toFloat() expects no arguments")
			}
			return NewValue(num)  // уже float64
		case "isInteger":
			if len(m.Args) != 0 {
				panic("number.isInteger() expects no arguments")
			}
			return NewValue(num == float64(int64(num)))
		}
	}
	
	// Методы для bool
	if b, ok := obj.Any().(bool); ok {
		switch m.MethodName {
		case "toString":
			if len(m.Args) != 0 {
				panic("bool.toString() expects no arguments")
			}
			return NewValue(strconv.FormatBool(b))
		case "not":
			if len(m.Args) != 0 {
				panic("bool.not() expects no arguments")
			}
			return NewValue(!b)
		}
	}
	
	// Проверяем extension методы
	typeName := value.GetValueTypeName(obj)
	if extensionMethod, ok := value.GetExtensionMethod(typeName, m.MethodName); ok {
		if wrapper, ok := extensionMethod.(*ExtensionMethodWrapper); ok {
			// Вычисляем аргументы
			args := make([]*Value, len(m.Args))
			for i, arg := range m.Args {
				args[i] = arg.Eval()
			}
			// Вызываем extension метод с receiver (this) как первым аргументом
			return wrapper.Call(obj, args)
		}
	}
	
	// Проверяем interface методы
	if typeInfo, ok := obj.Any().(*TypeInfo); ok {
		// Для объекта структуры, ищем методы интерфейса
		if impl := GetImplementation(typeInfo.Name, ""); impl != nil {
			// Ищем метод в реализации
			for _, method := range impl.Methods {
				if method.FuncName == m.MethodName {
					// Вычисляем аргументы
					args := make([]*Value, len(m.Args))
					for i, arg := range m.Args {
						args[i] = arg.Eval()
					}
					
					// Создаем временную область видимости для метода
					// TODO: Нужно установить 'this' как текущий объект
					
					// Пока что вызываем метод как обычную функцию  
					closure := NewTypedClosure(method.FuncName, method.Params, method.Body)
					return closure.Call(args)
				}
			}
		}
	}
	
	// Если объект - это экземпляр структуры (не TypeInfo), попробуем найти его тип
	if typeInfo := getObjectTypeInfo(obj); typeInfo != nil {
		// Проверяем все зарегистрированные интерфейсы для этого типа
		for _, interfaceName := range getImplementedInterfaces(typeInfo.Name) {
			if impl := GetImplementation(typeInfo.Name, interfaceName); impl != nil {
				// Ищем метод в реализации
				for _, method := range impl.Methods {
					if method.FuncName == m.MethodName {
						// Вычисляем аргументы
						args := make([]*Value, len(m.Args))
						for i, arg := range m.Args {
							args[i] = arg.Eval()
						}
						
						// Вызываем метод с 'this' контекстом
						return callMethodWithContext(method, obj, args)
					}
				}
			}
		}
	}
	
	// Отладка: покажем тип объекта
	objType := fmt.Sprintf("%T", obj.Any())
	panic("method '" + m.MethodName + "' not supported on type: " + objType)
}

// getObjectTypeInfo пытается определить TypeInfo для объекта
func getObjectTypeInfo(obj *Value) *TypeInfo {
	// Если объект - это StructObject, получаем его TypeInfo
	if structObj, ok := obj.Any().(*StructObject); ok {
		return structObj.TypeInfo
	}
	return nil
}

// getImplementedInterfaces возвращает список интерфейсов, реализованных типом
func getImplementedInterfaces(typeName string) []string {
	var interfaces []string
	for interfaceName := range interfaceImplementations[typeName] {
		interfaces = append(interfaces, interfaceName)
	}
	return interfaces
}

// callMethodWithContext вызывает метод интерфейса с установленным 'this' контекстом
func callMethodWithContext(method *TypedFuncStatement, thisObj *Value, args []*Value) *Value {
	// Создаем временную область видимости
	scope.GlobalScope.Push()
	defer scope.GlobalScope.Pop()
	
	// Устанавливаем 'this' в области видимости
	scope.GlobalScope.Set("this", thisObj)
	
	// Устанавливаем параметры метода
	for i, param := range method.Params {
		if i < len(args) {
			scope.GlobalScope.Set(param.Name, args[i])
		} else if param.Default != nil {
			scope.GlobalScope.Set(param.Name, param.Default.Eval())
		} else {
			panic(fmt.Sprintf("missing required argument: %s", param.Name))
		}
	}
	
	// Выполняем тело метода
	if bodyStm, ok := method.Body.(*BodyExpr); ok {
		for _, stmt := range bodyStm.Statments {
			result := stmt.Eval()
			if result != nil && result.IsReturn() {
				return result
			}
		}
		return NewValue(nil)
	} else {
		return method.Body.Eval()
	}
}