package ast

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"strings"
	"regexp"
)

// GenerateExpr представляет блок generate {} для шаблонной генерации кода в макросах
type GenerateExpr struct {
	Template string // Шаблон с ${} интерполяциями
}

func NewGenerateExpr(template string) *GenerateExpr {
	return &GenerateExpr{Template: template}
}

func (g *GenerateExpr) Eval() *value.Value {
	// Обрабатываем шаблон и заменяем ${...} на реальные значения
	processedCode := g.processTemplate()
	
	// TODO: Нужен способ динамически парсить и выполнять сгенерированный код
	// Пока просто выводим сгенерированный код для отладки
	fmt.Println("=== Generated code ===")
	fmt.Println(processedCode)
	fmt.Println("======================")
	
	// Возвращаем nil пока не реализован полный механизм
	return value.NewValue(nil)
}

// processTemplate обрабатывает шаблон и заменяет ${...} интерполяции
func (g *GenerateExpr) processTemplate() string {
	template := g.Template
	
	// Регулярное выражение для поиска ${...}
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	
	// Заменяем все ${...} на вычисленные значения
	result := re.ReplaceAllStringFunc(template, func(match string) string {
		// Убираем ${ и }
		expr := match[2:len(match)-1]
		
		// Обрабатываем специальные конструкции
		if strings.HasPrefix(expr, "for ") {
			return g.processForLoop(expr)
		} else if strings.HasPrefix(expr, "if ") {
			return g.processIfStatement(expr)
		} else {
			// Обычная интерполяция переменной или выражения
			return g.evaluateExpression(expr)
		}
	})
	
	return result
}

// evaluateExpression вычисляет выражение и возвращает его строковое представление
func (g *GenerateExpr) evaluateExpression(expr string) string {
	// Пытаемся получить значение из scope
	if val, found := scope.GlobalScope.Get(expr); found {
		if val != nil {
			return fmt.Sprintf("%v", val.Any())
		}
	}
	
	// Обрабатываем доступ к свойствам (например, structParam.Name)
	if strings.Contains(expr, ".") {
		parts := strings.Split(expr, ".")
		if len(parts) == 2 {
			// Получаем объект из scope
			if obj, found := scope.GlobalScope.Get(parts[0]); found && obj != nil {
				// Если это TypeInfo, получаем его свойство
				if typeInfo, ok := obj.Any().(*TypeInfo); ok {
					if parts[1] == "Name" {
						return typeInfo.Name
					} else if parts[1] == "Kind" {
						return typeInfo.Kind
					}
				}
			}
		}
	}
	
	// Если не можем вычислить, возвращаем как есть в ${}
	return "${" + expr + "}"
}

// processForLoop обрабатывает цикл for внутри шаблона
func (g *GenerateExpr) processForLoop(expr string) string {
	// Упрощенная обработка for циклов
	// Формат: for field in structParam.Fields { ... }
	
	// Извлекаем переменную, коллекцию и тело
	// Это упрощенная версия, в реальной реализации нужен полный парсер
	
	// Пока возвращаем заглушку
	return "/* for loop processing not yet implemented */"
}

// processIfStatement обрабатывает условные операторы в шаблоне
func (g *GenerateExpr) processIfStatement(expr string) string {
	// Упрощенная обработка if операторов
	// Пока возвращаем заглушку
	return "/* if statement processing not yet implemented */"
}