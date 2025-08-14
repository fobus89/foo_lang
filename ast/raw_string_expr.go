package ast

import (
	"fmt"
	"regexp"
	"strings"
	"foo_lang/scope"
	"foo_lang/value"
)

// RawStringExpr представляет строку с необработанными интерполяциями ${...}
// Используется в compile-time контексте, где интерполяции должны обрабатываться позже
type RawStringExpr struct {
	RawText string // Исходный текст с ${...} интерполяциями
}

func NewRawStringExpr(rawText string) *RawStringExpr {
	return &RawStringExpr{RawText: rawText}
}

func (r *RawStringExpr) Eval() *value.Value {
	// Обрабатываем интерполяции на этапе выполнения
	result := r.processInterpolations()
	return value.NewValue(result)
}

// processInterpolations обрабатывает ${...} интерполяции с текущим scope
func (r *RawStringExpr) processInterpolations() string {
	text := r.RawText
	
	// Регулярное выражение для поиска ${...}
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	
	// Заменяем все ${...} на вычисленные значения
	result := re.ReplaceAllStringFunc(text, func(match string) string {
		// Убираем ${ и }
		varName := strings.TrimSpace(match[2 : len(match)-1])
		
		// Пытаемся получить значение из scope
		if val, found := scope.GlobalScope.Get(varName); found {
			if val != nil {
				return r.valueToString(val)
			}
		}
		
		// Если переменная не найдена, возвращаем имя переменной
		return varName
	})
	
	return result
}

// valueToString конвертирует value.Value в строку
func (r *RawStringExpr) valueToString(val *value.Value) string {
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