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