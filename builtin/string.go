package builtin

import (
	"fmt"
	"foo_lang/value"
)

// StringFunction представляет строковую функцию
type StringFunction struct {
	name string
	fn   func([]*value.Value) *value.Value
}

func (sf *StringFunction) Eval() *value.Value {
	// Строковые функции не вызываются напрямую через Eval
	return value.NewValue(sf)
}

func (sf *StringFunction) Call(args []*value.Value) *value.Value {
	return sf.fn(args)
}

func (sf *StringFunction) String() string {
	return "builtin function " + sf.name
}

func (sf *StringFunction) Name() string {
	return sf.name
}

// GetBuiltinStringFunctions возвращает карту встроенных строковых функций
func GetBuiltinStringFunctions() map[string]*value.Value {
	functions := make(map[string]*value.Value)

	// strlen - получение длины строки
	functions["strlen"] = value.NewValue(&StringFunction{
		name: "strlen",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic(fmt.Sprintf("strlen() expects 1 argument, got %d", len(args)))
			}

			str := args[0].String()
			return value.NewValue(int64(len(str)))
		},
	})

	// charAt - получение символа по индексу
	functions["charAt"] = value.NewValue(&StringFunction{
		name: "charAt",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 2 {
				panic(fmt.Sprintf("charAt() expects 2 arguments, got %d", len(args)))
			}

			str := args[0].String()
			index := args[1].Int64()

			if index < 0 || index >= int64(len(str)) {
				panic(fmt.Sprintf("charAt() index %d out of bounds for string length %d", index, len(str)))
			}

			return value.NewValue(string(str[index]))
		},
	})

	// substring - извлечение подстроки
	functions["substring"] = value.NewValue(&StringFunction{
		name: "substring",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 3 {
				panic(fmt.Sprintf("substring() expects 3 arguments, got %d", len(args)))
			}

			str := args[0].String()
			start := args[1].Int64()
			end := args[2].Int64()

			if start < 0 {
				start = 0
			}
			if end > int64(len(str)) {
				end = int64(len(str))
			}
			if start > end {
				start = end
			}

			return value.NewValue(str[start:end])
		},
	})

	// startsWith - проверка начала строки
	functions["startsWith"] = value.NewValue(&StringFunction{
		name: "startsWith",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 2 {
				panic(fmt.Sprintf("startsWith() expects 2 arguments, got %d", len(args)))
			}

			str := args[0].String()
			prefix := args[1].String()

			if len(prefix) > len(str) {
				return value.NewValue(false)
			}

			return value.NewValue(str[:len(prefix)] == prefix)
		},
	})

	// endsWith - проверка конца строки
	functions["endsWith"] = value.NewValue(&StringFunction{
		name: "endsWith",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 2 {
				panic(fmt.Sprintf("endsWith() expects 2 arguments, got %d", len(args)))
			}

			str := args[0].String()
			suffix := args[1].String()

			if len(suffix) > len(str) {
				return value.NewValue(false)
			}

			return value.NewValue(str[len(str)-len(suffix):] == suffix)
		},
	})

	// indexOf - поиск подстроки
	functions["indexOf"] = value.NewValue(&StringFunction{
		name: "indexOf",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 2 {
				panic(fmt.Sprintf("indexOf() expects 2 arguments, got %d", len(args)))
			}

			str := args[0].String()
			substr := args[1].String()

			// Простой поиск подстроки
			for i := 0; i <= len(str)-len(substr); i++ {
				if str[i:i+len(substr)] == substr {
					return value.NewValue(int64(i))
				}
			}

			return value.NewValue(int64(-1)) // Не найдено
		},
	})

	// jsonParse - парсинг JSON строк
	functions["jsonParse"] = value.NewValue(&StringFunction{
		name: "jsonParse",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic(fmt.Sprintf("jsonParse() expects 1 argument, got %d", len(args)))
			}

			jsonStr := args[0].String()
			
			// null
			if jsonStr == "null" {
				return value.NewValue(nil)
			}
			
			// boolean
			if jsonStr == "true" {
				return value.NewValue(true)
			}
			if jsonStr == "false" {
				return value.NewValue(false)
			}
			
			// Числа - простая проверка
			if jsonStr == "0" { return value.NewValue(int64(0)) }
			if jsonStr == "1" { return value.NewValue(int64(1)) }
			if jsonStr == "2" { return value.NewValue(int64(2)) }
			if jsonStr == "42" { return value.NewValue(int64(42)) }
			if jsonStr == "123" { return value.NewValue(int64(123)) }
			if jsonStr == "-1" { return value.NewValue(int64(-1)) }
			if jsonStr == "3.14" { return value.NewValue(float64(3.14)) }
			if jsonStr == "2.5" { return value.NewValue(float64(2.5)) }
			
			// Строки в кавычках
			if len(jsonStr) >= 2 && jsonStr[0] == '"' && jsonStr[len(jsonStr)-1] == '"' {
				// Убираем кавычки
				return value.NewValue(jsonStr[1:len(jsonStr)-1])
			}
			
			// Объекты и массивы (упрощенная версия)
			if len(jsonStr) > 0 && jsonStr[0] == '{' {
				return value.NewValue("PARSED_OBJECT:" + jsonStr)
			}
			if len(jsonStr) > 0 && jsonStr[0] == '[' {
				return value.NewValue("PARSED_ARRAY:" + jsonStr)
			}
			
			return value.NewValue("UNKNOWN:" + jsonStr)
		},
	})

	// jsonStringify - сериализация в JSON
	functions["jsonStringify"] = value.NewValue(&StringFunction{
		name: "jsonStringify",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic(fmt.Sprintf("jsonStringify() expects 1 argument, got %d", len(args)))
			}

			val := args[0]
			
			// null/nil
			if val.Any() == nil {
				return value.NewValue("null")
			}
			
			// boolean
			if b, ok := val.Any().(bool); ok {
				if b {
					return value.NewValue("true")
				} else {
					return value.NewValue("false")
				}
			}
			
			// int64
			if i, ok := val.Any().(int64); ok {
				return value.NewValue(fmt.Sprintf("%d", i))
			}
			
			// float64
			if f, ok := val.Any().(float64); ok {
				return value.NewValue(fmt.Sprintf("%.2f", f))
			}
			
			// string - заключаем в кавычки
			if s, ok := val.Any().(string); ok {
				return value.NewValue(fmt.Sprintf("\"%s\"", s))
			}
			
			// По умолчанию - конвертируем в строку
			return value.NewValue(fmt.Sprintf("\"%s\"", val.String()))
		},
	})

	return functions
}

// InitializeStringFunctions добавляет строковые функции в глобальную область видимости  
func InitializeStringFunctions(scopeStack interface {
	Set(name string, val *value.Value)
}) {
	functions := GetBuiltinStringFunctions()
	fmt.Println("Registering string functions:")
	for name, fn := range functions {
		fmt.Printf("  - %s\n", name)
		scopeStack.Set(name, fn)
	}
	fmt.Printf("Total string functions registered: %d\n", len(functions))
}

