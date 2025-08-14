package ast

import (
	"fmt"
	"regexp"
	"strings"
	"foo_lang/scope"
	"foo_lang/value"
)

// CompileTimeIfExpr представляет compile-time условие $if
// Выполняется во время макроса, а не в runtime
type CompileTimeIfExpr struct {
	Condition Expr // Условие, проверяемое в compile-time
	ThenBody  Expr // Код, генерируемый если условие истинно
	ElseBody  Expr // Код, генерируемый если условие ложно (опционально)
}

func NewCompileTimeIfExpr(condition, thenBody, elseBody Expr) *CompileTimeIfExpr {
	return &CompileTimeIfExpr{
		Condition: condition,
		ThenBody:  thenBody,
		ElseBody:  elseBody,
	}
}

func (c *CompileTimeIfExpr) Eval() *value.Value {
	// Вычисляем условие в compile-time
	condition := c.Condition.Eval()
	
	// Проверяем результат условия
	if condition.IsTruthy() {
		// Если условие истинно, возвращаем тело then
		return c.ThenBody.Eval()
	} else if c.ElseBody != nil {
		// Если условие ложно и есть else, возвращаем тело else
		return c.ElseBody.Eval()
	}
	
	// Если условие ложно и нет else, возвращаем пустую строку
	return value.NewValue("")
}

// CompileTimeForExpr представляет compile-time цикл $for
// Выполняется во время макроса для генерации повторяющегося кода
type CompileTimeForExpr struct {
	Iterator   string // Имя переменной итератора
	Collection Expr   // Коллекция для итерации
	Body       Expr   // Тело цикла, генерируемое для каждого элемента
}

func NewCompileTimeForExpr(iterator string, collection, body Expr) *CompileTimeForExpr {
	return &CompileTimeForExpr{
		Iterator:   iterator,
		Collection: collection,
		Body:       body,
	}
}

func (c *CompileTimeForExpr) Eval() *value.Value {
	// Вычисляем коллекцию
	collection := c.Collection.Eval()
	
	var generatedCode []string
	
	// Создаем новую область видимости для цикла
	scope.GlobalScope.Push()
	defer scope.GlobalScope.Pop()
	
	// Итерируемся по коллекции
	// Пока упростим - будем работать только со строками и числами
	if collection.IsString() {
		// Итерируемся по символам строки
		str := collection.String()
		for i, char := range str {
			// Устанавливаем переменную итератора
			scope.GlobalScope.Set(c.Iterator, value.NewValue(string(char)))
			scope.GlobalScope.Set("index", value.NewValue(int64(i)))
			
			// Генерируем код для текущего элемента
			result := c.Body.Eval()
			if result.IsString() {
				// Обрабатываем интерполяцию ${переменная_итератора} вручную
				processedString := c.processStringInterpolation(result.String(), c.Iterator, string(char))
				generatedCode = append(generatedCode, processedString)
			}
		}
	} else if collection.IsInt64() {
		// Итерируемся от 0 до значения числа
		count := collection.Int64()
		for i := int64(0); i < count; i++ {
			// Устанавливаем переменную итератора ПЕРЕД выполнением Body
			scope.GlobalScope.Set(c.Iterator, value.NewValue(i))
			
			
			// Специальная обработка для compile-time интерполяций
			processedCode := c.processBodyWithInterpolation()
			if processedCode != "" {
				generatedCode = append(generatedCode, processedCode)
			}
		}
	}
	
	// Возвращаем объединенный сгенерированный код
	return value.NewValue(strings.Join(generatedCode, " "))
}

// processInterpolation обрабатывает ${переменная} интерполяции в строке
func (c *CompileTimeForExpr) processInterpolation(text string, varName string, varValue int64) string {
	// Регулярное выражение для поиска ${varName}
	pattern := fmt.Sprintf(`\$\{%s\}`, regexp.QuoteMeta(varName))
	re := regexp.MustCompile(pattern)
	
	// Заменяем ${varName} на значение
	result := re.ReplaceAllString(text, fmt.Sprintf("%d", varValue))
	
	return result
}

// processStringInterpolation обрабатывает ${переменная} интерполяции в строке для строковых значений
func (c *CompileTimeForExpr) processStringInterpolation(text string, varName string, varValue string) string {
	// Регулярное выражение для поиска ${varName}
	pattern := fmt.Sprintf(`\$\{%s\}`, regexp.QuoteMeta(varName))
	re := regexp.MustCompile(pattern)
	
	// Заменяем ${varName} на значение
	result := re.ReplaceAllString(text, varValue)
	
	return result
}

// getBodyCode возвращает исходный код тела цикла как строку без выполнения
// Пока используем упрощенный подход - получаем результат как строку
func (c *CompileTimeForExpr) getBodyCode() string {
	// Упрощенное решение: выполняем Body и получаем результат
	// В будущем можно улучшить для получения исходного кода без выполнения
	result := c.Body.Eval()
	if result.IsString() {
		return result.String()
	}
	return fmt.Sprintf("%v", result.Any())
}

// processBodyWithInterpolation обрабатывает Body с учетом compile-time интерполяций
func (c *CompileTimeForExpr) processBodyWithInterpolation() string {
	// Проверяем, является ли Body оберткой BodyExpr
	if bodyExpr, ok := c.Body.(*BodyExpr); ok {
		// Обрабатываем все statements в BodyExpr
		var results []string
		for _, stmt := range bodyExpr.Statments {
			// Обрабатываем PrintExpr (println)
			if printExpr, isPrintExpr := stmt.(*PrintExpr); isPrintExpr {
				if printExpr.Expr != nil {
					// Обрабатываем выражение специально для compile-time
					processedArg := c.processArgumentWithInterpolation(printExpr.Expr)
					
					// Добавляем сгенерированный код
					results = append(results, fmt.Sprintf("println(\"%s\")", processedArg))
				}
			}
			
			// Обрабатываем FuncCallExpr (обычные вызовы функций)
			if funcCall, isFuncCall := stmt.(*FuncCallExpr); isFuncCall {
				if funcCall.funcName == "println" && len(funcCall.args) == 1 {
					// Получаем аргумент println
					arg := funcCall.args[0]
					
					// Обрабатываем аргумент специально для compile-time
					processedArg := c.processArgumentWithInterpolation(arg)
					
					// Добавляем сгенерированный код
					results = append(results, fmt.Sprintf("println(\"%s\")", processedArg))
				}
			}
		}
		if len(results) > 0 {
			return strings.Join(results, "; ")
		}
	}
	
	// Проверяем, является ли Body вызовом функции (например, println)
	if funcCall, ok := c.Body.(*FuncCallExpr); ok {
		if funcCall.funcName == "println" && len(funcCall.args) == 1 {
			// Получаем аргумент println
			arg := funcCall.args[0]
			
			// Обрабатываем аргумент специально для compile-time
			processedArg := c.processArgumentWithInterpolation(arg)
			
			// Возвращаем сгенерированный код
			return fmt.Sprintf("println(\"%s\")", processedArg)
		}
	}
	
	// Для других случаев используем стандартную обработку
	result := c.Body.Eval()
	if result.IsString() {
		return result.String()
	}
	return fmt.Sprintf("%v", result.Any())
}

// processArgumentWithInterpolation обрабатывает аргументы функций с интерполяциями
func (c *CompileTimeForExpr) processArgumentWithInterpolation(arg Expr) string {
	// Если это StringFormatExpr, обрабатываем его части
	if stringFormat, ok := arg.(*StringFormatExpr); ok {
		return c.processStringFormatWithInterpolation(stringFormat)
	}
	
	// Если это простая строка, проверяем на возможные необработанные интерполяции
	if literalString, ok := arg.(*LiteralString); ok {
		// Проверяем, содержит ли строка имя переменной итератора
		text := literalString.Value
		if strings.Contains(text, c.Iterator) {
			// Пытаемся восстановить интерполяцию и обработать её
			restoredText := c.restoreInterpolation(text)
			return c.processManualInterpolation(restoredText)
		}
		return text
	}
	
	// Для других типов выполняем стандартную обработку
	result := arg.Eval()
	if result.IsString() {
		return result.String()
	}
	return fmt.Sprintf("%v", result.Any())
}

// restoreInterpolation пытается восстановить ${...} интерполяцию в строке
func (c *CompileTimeForExpr) restoreInterpolation(text string) string {
	// Простой подход: заменяем все вхождения имени переменной на ${имя_переменной}
	// Это работает, потому что мы знаем, что в compile-time контексте 
	// ${counter} превращается в просто "counter"
	restored := strings.ReplaceAll(text, c.Iterator, "${"+c.Iterator+"}")
	
	return restored
}

// processManualInterpolation обрабатывает восстановленную интерполяцию
func (c *CompileTimeForExpr) processManualInterpolation(text string) string {
	// Регулярное выражение для поиска ${...}
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	
	// Заменяем все ${...} на вычисленные значения
	result := re.ReplaceAllStringFunc(text, func(match string) string {
		// Убираем ${ и }
		varName := strings.TrimSpace(match[2 : len(match)-1])
		
		// Пытаемся получить значение из scope
		if val, found := scope.GlobalScope.Get(varName); found {
			if val != nil {
				return c.valueToString(val)
			}
		}
		
		// Если переменная не найдена, возвращаем имя переменной
		return varName
	})
	
	return result
}

// processStringFormatWithInterpolation обрабатывает StringFormatExpr с compile-time интерполяциями
func (c *CompileTimeForExpr) processStringFormatWithInterpolation(format *StringFormatExpr) string {
	var result strings.Builder
	
	// Обрабатываем каждую часть StringFormatExpr
	for _, part := range format.exprs {
		switch p := part.(type) {
		case *LiteralString:
			result.WriteString(p.Value)
		case *VarExpr:
			// Для переменных ищем значение в текущем scope
			if val, found := scope.GlobalScope.Get(p.Name); found {
				result.WriteString(c.valueToString(val))
			} else {
				// Если переменная не найдена, возвращаем её имя
				result.WriteString(p.Name)
			}
		default:
			// Для других выражений пытаемся выполнить их
			res := part.Eval()
			if res != nil {
				result.WriteString(c.valueToString(res))
			}
		}
	}
	
	return result.String()
}

// valueToString конвертирует value.Value в строку (дублирует метод из RawStringExpr)
func (c *CompileTimeForExpr) valueToString(val *value.Value) string {
	if val == nil || val.Any() == nil {
		return ""
	}
	
	switch v := val.Any().(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// CompileTimeLetExpr представляет compile-time переменную $let
// Переменная существует только во время выполнения макроса
type CompileTimeLetExpr struct {
	Name string // Имя переменной
	Expr Expr   // Выражение для вычисления значения
}

func NewCompileTimeLetExpr(name string, expr Expr) *CompileTimeLetExpr {
	return &CompileTimeLetExpr{
		Name: name,
		Expr: expr,
	}
}

func (c *CompileTimeLetExpr) Eval() *value.Value {
	// Вычисляем значение выражения
	val := c.Expr.Eval()
	
	// Сохраняем переменную в текущей области видимости
	scope.GlobalScope.Set(c.Name, val)
	
	// Compile-time переменные не генерируют код, возвращаем пустую строку
	return value.NewValue("")
}

// CompileTimeWhileExpr представляет compile-time цикл $while
// Выполняется во время макроса, пока условие истинно
type CompileTimeWhileExpr struct {
	Condition Expr // Условие цикла
	Body      Expr // Тело цикла
}

func NewCompileTimeWhileExpr(condition, body Expr) *CompileTimeWhileExpr {
	return &CompileTimeWhileExpr{
		Condition: condition,
		Body:      body,
	}
}

func (c *CompileTimeWhileExpr) Eval() *value.Value {
	var generatedCode []string
	
	// Создаем новую область видимости для цикла
	scope.GlobalScope.Push()
	defer scope.GlobalScope.Pop()
	
	// Выполняем цикл пока условие истинно
	for {
		// Проверяем условие
		condition := c.Condition.Eval()
		if !condition.IsTruthy() {
			break
		}
		
		// Выполняем тело цикла
		result := c.Body.Eval()
		if result.IsString() {
			generatedCode = append(generatedCode, result.String())
		}
		
		// Защита от бесконечного цикла
		if len(generatedCode) > 1000 {
			break
		}
	}
	
	// Возвращаем объединенный сгенерированный код
	return value.NewValue(strings.Join(generatedCode, " "))
}