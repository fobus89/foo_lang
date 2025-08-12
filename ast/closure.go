package ast

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
)

// Closure представляет замыкание - функцию с захваченными переменными
type Closure struct {
	funcName     string
	args         []map[string]Expr
	body         Expr
	capturedVars map[string]*value.Value  // Захваченные переменные
	isMacro      bool
}

// TypedClosure представляет типизированное замыкание
type TypedClosure struct {
	funcName     string
	params       []FuncParam
	body         Expr
	returnType   string                   // Ожидаемый тип возвращаемого значения
	capturedVars map[string]*value.Value  // Захваченные переменные
}

// NewClosure создает новое замыкание с захватом переменных из текущей области
func NewClosure(funcName string, args []map[string]Expr, body Expr, isMacro bool) *Closure {
	// Захватываем все переменные из текущей области видимости
	capturedVars := make(map[string]*value.Value)
	
	// Анализируем тело функции для поиска свободных переменных
	freeVars := findFreeVariables(body, args)
	
	// Захватываем значения свободных переменных
	for varName := range freeVars {
		if val, exists := scope.GlobalScope.Get(varName); exists {
			// Создаем копию значения для захвата
			capturedVars[varName] = value.NewValue(val.Any())
		}
	}
	
	return &Closure{
		funcName:     funcName,
		args:         args,
		body:         body,
		capturedVars: capturedVars,
		isMacro:      isMacro,
	}
}

// findFreeVariables анализирует AST и находит свободные переменные
func findFreeVariables(expr Expr, funcArgs []map[string]Expr) map[string]bool {
	freeVars := make(map[string]bool)
	
	// Пока что упростим: захватываем все переменные из текущей области видимости
	// TODO: Улучшить анализ AST для более точного определения свободных переменных
	allVars := scope.GlobalScope.GetAll()
	
	// Исключаем параметры функции из захваченных переменных
	localVars := make(map[string]bool)
	for _, arg := range funcArgs {
		for name := range arg {
			localVars[name] = true
		}
	}
	
	// Захватываем все остальные переменные
	for name := range allVars {
		if !localVars[name] {
			freeVars[name] = true
		}
	}
	
	return freeVars
}

// Name возвращает имя функции
func (c *Closure) Name() string {
	return c.funcName
}

// Call выполняет замыкание с захваченными переменными
func (c *Closure) Call(args []*Value) *Value {

	expected := len(c.args)
	passed := len(args)

	if passed > expected {
		panic(fmt.Sprintf("too many arguments: expected %d, got %d", expected, passed))
	}

	// Создаем новую область видимости для функции
	err := scope.GlobalScope.PushFunction()
	if err != nil {
		panic(err.Error())
	}
	defer scope.GlobalScope.PopFunction()

	// Восстанавливаем захваченные переменные в новой области
	for name, val := range c.capturedVars {
		scope.GlobalScope.Set(name, val)
	}

	// Устанавливаем параметры функции
	for i, arg := range c.args {
		for name, expr := range arg {
			if i < len(args) {
				scope.GlobalScope.Set(name, args[i])
			} else if expr != nil {
				defaultValue := expr.Eval()
				scope.GlobalScope.Set(name, defaultValue)
			} else {
				panic(fmt.Sprintf("missing required argument: %s", name))
			}
		}
	}

	// Выполняем тело функции
	if bodyStm, ok := c.body.(*BodyExpr); ok {
		// Если тело - блок выражений
		for _, stm := range bodyStm.Statments {
			if stm == nil {
				continue
			}

			result := stm.Eval()
			
			// Проверяем на return
			if result != nil && result.IsReturn() {
				return result
			}
		}
	} else {
		// Если тело - одно выражение (для стрелочных функций)
		result := c.body.Eval()
		if result != nil {
			// Для одиночных выражений автоматически возвращаем результат
			retResult := NewValue(result.Any())
			retResult.SetReturn(true)
			return retResult
		}
	}

	return nil
}

// String возвращает строковое представление замыкания
func (c *Closure) String() string {
	return fmt.Sprintf("closure %s(%v) with %d captured vars", c.funcName, c.args, len(c.capturedVars))
}

// IsMacro проверяет, является ли замыкание макросом
func (c *Closure) IsMacro() bool {
	return c.isMacro
}

// GetCapturedVars возвращает захваченные переменные (для отладки)
func (c *Closure) GetCapturedVars() map[string]*value.Value {
	return c.capturedVars
}

// NewTypedClosure создает новое типизированное замыкание
func NewTypedClosure(funcName string, params []FuncParam, body Expr, returnType string) *TypedClosure {
	// Захватываем все переменные из текущей области видимости
	capturedVars := make(map[string]*value.Value)
	
	// Пока используем простое решение - не захватываем переменные для типизированных функций
	// В будущем можно добавить анализ свободных переменных
	
	return &TypedClosure{
		funcName:     funcName,
		params:       params,
		body:         body,
		returnType:   returnType,
		capturedVars: capturedVars,
	}
}

// Name возвращает имя типизированной функции
func (tc *TypedClosure) Name() string {
	return tc.funcName
}

// Call вызывает типизированное замыкание  
func (tc *TypedClosure) Call(args []*Value) *Value {
	// Проверяем количество аргументов
	requiredArgs := 0
	for _, param := range tc.params {
		if param.Default == nil {
			requiredArgs++
		}
	}
	
	if len(args) < requiredArgs {
		panic(fmt.Sprintf("function '%s' requires at least %d arguments, got %d", tc.funcName, requiredArgs, len(args)))
	}
	
	if len(args) > len(tc.params) {
		panic(fmt.Sprintf("function '%s' accepts at most %d arguments, got %d", tc.funcName, len(tc.params), len(args)))
	}
	
	// Создаем новую область видимости
	scope.GlobalScope.Push()
	defer scope.GlobalScope.Pop()
	
	// Восстанавливаем захваченные переменные
	for name, val := range tc.capturedVars {
		scope.GlobalScope.Set(name, val)
	}
	
	// Устанавливаем параметры функции с проверкой типов
	for i, param := range tc.params {
		var argValue *Value
		
		if i < len(args) {
			argValue = args[i]
			
			// Проверяем тип параметра, если указан
			if param.TypeName != "" {
				if err := validateFunctionParameterType(argValue, param.TypeName); err != nil {
					panic(fmt.Sprintf("function '%s' parameter '%s': %s", tc.funcName, param.Name, err.Error()))
				}
			}
		} else if param.Default != nil {
			// Используем значение по умолчанию
			argValue = param.Default.Eval()
		} else {
			panic(fmt.Sprintf("missing required argument: %s", param.Name))
		}
		
		scope.GlobalScope.Set(param.Name, argValue)
	}
	
	// Выполняем тело функции
	var result *Value
	if bodyStm, ok := tc.body.(*BodyExpr); ok {
		for _, stmt := range bodyStm.Statments {
			result = stmt.Eval()
			if result != nil && result.IsReturn() {
				break
			}
		}
		if result == nil || (result != nil && !result.IsReturn()) {
			result = NewValue(nil)
		}
	} else {
		result = tc.body.Eval()
	}
	
	// Проверяем тип возвращаемого значения, если он указан
	if tc.returnType != "" && result != nil {
		if err := tc.validateReturnType(result); err != nil {
			panic(fmt.Sprintf("function '%s' return type error: %s", tc.funcName, err.Error()))
		}
	}
	
	return result
}

// validateReturnType проверяет соответствие типа возвращаемого значения ожидаемому
func (tc *TypedClosure) validateReturnType(returnValue *Value) error {
	// Извлекаем значение из return-обертки, если оно есть
	actualValue := returnValue
	if returnValue.IsReturn() {
		actualValue = NewValue(returnValue.Any())
	}
	
	return validateFunctionParameterType(actualValue, tc.returnType)
}

// validateFunctionParameterType проверяет соответствие типа аргумента ожидаемому типу (включая Union типы)
func validateFunctionParameterType(argValue *Value, expectedTypeName string) error {
	// Используем универсальную функцию валидации, которая поддерживает Union типы
	return validateVariableType(argValue, expectedTypeName)
}