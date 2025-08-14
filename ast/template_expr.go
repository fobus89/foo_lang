package ast

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"regexp"
	"strings"
)

// TemplateExpr представляет шаблон с интерполяцией ${expr}
type TemplateExpr struct {
	Template string // Шаблон с ${} интерполяциями
}

func NewTemplateExpr(template string) *TemplateExpr {
	return &TemplateExpr{Template: template}
}

func (t *TemplateExpr) Eval() *value.Value {
	// Обрабатываем шаблон и заменяем ${...} на вычисленные значения
	result := t.processTemplate()
	
	// Возвращаем обработанную строку
	return value.NewValue(result)
}

// processTemplate обрабатывает шаблон и заменяет ${...} интерполяции
func (t *TemplateExpr) processTemplate() string {
	template := t.Template
	
	// Регулярное выражение для поиска ${...}
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	
	// Заменяем все ${...} на вычисленные значения
	result := re.ReplaceAllStringFunc(template, func(match string) string {
		// Убираем ${ и }
		expr := strings.TrimSpace(match[2 : len(match)-1])
		
		// Обрабатываем специальные конструкции
		if strings.HasPrefix(expr, "for ") {
			return t.processForLoop(expr)
		} else if strings.HasPrefix(expr, "if ") {
			return t.processIfStatement(expr)
		} else {
			// Обычная интерполяция переменной или выражения
			return t.evaluateExpression(expr)
		}
	})
	
	return result
}

// evaluateExpression вычисляет выражение и возвращает его строковое представление
func (t *TemplateExpr) evaluateExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	
	// Пытаемся получить значение из scope
	if val, found := scope.GlobalScope.Get(expr); found {
		if val != nil {
			return t.valueToString(val)
		}
	}
	
	// Обрабатываем доступ к свойствам (например, structParam.Name)
	if strings.Contains(expr, ".") {
		return t.evaluatePropertyAccess(expr)
	}
	
	// Обрабатываем вызовы методов (например, structParam.getName())
	if strings.Contains(expr, "(") && strings.Contains(expr, ")") {
		return t.evaluateMethodCall(expr)
	}
	
	// Если не можем вычислить, возвращаем как есть
	return expr
}

// evaluatePropertyAccess обрабатывает доступ к свойствам объектов
func (t *TemplateExpr) evaluatePropertyAccess(expr string) string {
	parts := strings.Split(expr, ".")
	if len(parts) < 2 {
		return expr
	}
	
	// Получаем базовый объект
	objName := parts[0]
	property := parts[1]
	
	if obj, found := scope.GlobalScope.Get(objName); found && obj != nil {
		// Если это TypeInfo, получаем его свойство
		if typeInfo, ok := obj.Any().(*TypeInfo); ok {
			switch property {
			case "Name":
				return typeInfo.Name
			case "Kind":
				return typeInfo.Kind
			case "Fields":
				// Возвращаем список полей для циклов
				fields := typeInfo.GetFields()
				if len(fields) > 0 {
					var fieldNames []string
					for _, field := range fields {
						fieldNames = append(fieldNames, field.Name)
					}
					return strings.Join(fieldNames, ", ")
				}
				return ""
			case "Methods":
				// TypeInfo не имеет поля Methods, возвращаем пустую строку
				return ""
			}
		}
		
		// Если это обычный объект/map
		if objValue, ok := obj.Any().(map[string]*value.Value); ok {
			if propValue, exists := objValue[property]; exists {
				return t.valueToString(propValue)
			}
		}
	}
	
	return expr
}

// evaluateMethodCall обрабатывает вызовы методов
func (t *TemplateExpr) evaluateMethodCall(expr string) string {
	// Упрощенная обработка вызовов методов
	// В будущем можно добавить полный парсер выражений
	
	// Пока возвращаем как есть
	return expr
}

// valueToString конвертирует value.Value в строку
func (t *TemplateExpr) valueToString(val *value.Value) string {
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
	case *TypeInfo:
		return v.Name
	default:
		return fmt.Sprintf("%v", v)
	}
}

// processForLoop обрабатывает цикл for внутри шаблона
// Формат: for field in structParam.Fields { template }
func (t *TemplateExpr) processForLoop(expr string) string {
	// Парсим for цикл
	// Упрощенная версия: for varName in collection { body }
	
	// Извлекаем части цикла
	forParts := strings.Fields(expr)
	if len(forParts) < 4 || forParts[0] != "for" || forParts[2] != "in" {
		return fmt.Sprintf("/* invalid for loop syntax: %s */", expr)
	}
	
	varName := forParts[1]
	collectionExpr := forParts[3]
	
	// Ищем тело цикла в фигурных скобках
	bodyStart := strings.Index(expr, "{")
	bodyEnd := strings.LastIndex(expr, "}")
	
	if bodyStart == -1 || bodyEnd == -1 || bodyEnd <= bodyStart {
		return fmt.Sprintf("/* for loop body not found: %s */", expr)
	}
	
	body := strings.TrimSpace(expr[bodyStart+1 : bodyEnd])
	
	// Получаем коллекцию для итерации
	collection := t.getCollection(collectionExpr)
	if collection == nil {
		return fmt.Sprintf("/* collection not found: %s */", collectionExpr)
	}
	
	// Сохраняем текущее состояние scope
	scope.GlobalScope.Push()
	defer scope.GlobalScope.Pop()
	
	var results []string
	
	// Итерируемся по коллекции
	for _, item := range collection {
		// Устанавливаем переменную цикла
		scope.GlobalScope.Set(varName, item)
		
		// Обрабатываем тело цикла
		processedBody := t.processTemplateString(body)
		results = append(results, processedBody)
	}
	
	return strings.Join(results, "\n")
}

// getCollection получает коллекцию для итерации из выражения
func (t *TemplateExpr) getCollection(expr string) []*value.Value {
	// Обрабатываем доступ к свойствам (например, structParam.Fields)
	if strings.Contains(expr, ".") {
		parts := strings.Split(expr, ".")
		if len(parts) == 2 {
			objName := parts[0]
			property := parts[1]
			
			if obj, found := scope.GlobalScope.Get(objName); found && obj != nil {
				if typeInfo, ok := obj.Any().(*TypeInfo); ok {
					switch property {
					case "Fields":
						fields := typeInfo.GetFields()
						if len(fields) > 0 {
							var collection []*value.Value
							for _, field := range fields {
								collection = append(collection, value.NewValue(field))
							}
							return collection
						}
					case "Methods":
						// TypeInfo не имеет поля Methods, возвращаем nil
						return nil
					}
				}
			}
		}
	}
	
	// Пытаемся получить массив из scope
	if val, found := scope.GlobalScope.Get(expr); found && val != nil {
		if arr, ok := val.Any().([]*value.Value); ok {
			return arr
		}
	}
	
	return nil
}

// processTemplateString обрабатывает строку шаблона (используется в циклах)
func (t *TemplateExpr) processTemplateString(template string) string {
	// Создаем новый TemplateExpr для обработки вложенного шаблона
	nestedTemplate := NewTemplateExpr(template)
	result := nestedTemplate.Eval()
	
	if result != nil && result.Any() != nil {
		return fmt.Sprintf("%v", result.Any())
	}
	
	return template
}

// processIfStatement обрабатывает условные операторы в шаблоне
// Формат: if condition { trueTemplate } else { falseTemplate }
func (t *TemplateExpr) processIfStatement(expr string) string {
	// Упрощенная обработка if операторов
	// В будущем можно добавить полный парсер условий
	
	// Ищем части if выражения
	ifStart := strings.Index(expr, "if ")
	if ifStart == -1 {
		return fmt.Sprintf("/* invalid if syntax: %s */", expr)
	}
	
	// Ищем тело условия
	bodyStart := strings.Index(expr, "{")
	bodyEnd := strings.Index(expr, "}")
	
	if bodyStart == -1 || bodyEnd == -1 || bodyEnd <= bodyStart {
		return fmt.Sprintf("/* if body not found: %s */", expr)
	}
	
	// Извлекаем условие
	condition := strings.TrimSpace(expr[3:bodyStart]) // убираем "if "
	
	// Извлекаем тело
	body := strings.TrimSpace(expr[bodyStart+1 : bodyEnd])
	
	// Вычисляем условие
	if t.evaluateCondition(condition) {
		return t.processTemplateString(body)
	}
	
	// Ищем else часть
	elseStart := strings.Index(expr[bodyEnd:], "else")
	if elseStart != -1 {
		elseBodyStart := strings.Index(expr[bodyEnd+elseStart:], "{")
		elseBodyEnd := strings.LastIndex(expr, "}")
		
		if elseBodyStart != -1 && elseBodyEnd != -1 {
			elseBody := strings.TrimSpace(expr[bodyEnd+elseStart+elseBodyStart+1 : elseBodyEnd])
			return t.processTemplateString(elseBody)
		}
	}
	
	return ""
}

// evaluateCondition вычисляет логическое условие
func (t *TemplateExpr) evaluateCondition(condition string) bool {
	condition = strings.TrimSpace(condition)
	
	// Простые логические значения
	if condition == "true" {
		return true
	}
	if condition == "false" {
		return false
	}
	
	// Проверка переменных из scope
	if val, found := scope.GlobalScope.Get(condition); found && val != nil {
		if b, ok := val.Any().(bool); ok {
			return b
		}
		// Любое непустое значение считается true
		return val.Any() != nil
	}
	
	// Простые сравнения (можно расширить)
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			return t.evaluateExpression(left) == t.evaluateExpression(right)
		}
	}
	
	// По умолчанию false
	return false
}

// processCodeBlock обрабатывает блок кода - генерирует код, не выполняет его
func (t *TemplateExpr) processCodeBlock(codeBlock string) string {
	// Обрабатываем интерполяции внутри блока кода
	processedCode := t.processTemplateString(codeBlock)
	
	// Возвращаем обработанный код как часть сгенерированного кода
	return processedCode
}