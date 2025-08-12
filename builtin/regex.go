package builtin

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"regexp"
	"strings"
)

// InitializeRegexFunctions инициализирует встроенные функции регулярных выражений
func InitializeRegexFunctions(globalScope *scope.ScopeStack) {

	// ============ ПОИСК И СОПОСТАВЛЕНИЕ ============

	// regexMatch - проверяет соответствие строки регулярному выражению
	regexMatchFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: regexMatch() requires 2 arguments (pattern, string)")
		}

		pattern, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexMatch() first argument must be string (pattern)")
		}

		text, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: regexMatch() second argument must be string")
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error compiling regex pattern: %v", err))
		}

		matched := regex.MatchString(text)
		return value.NewBool(matched)
	}
	globalScope.Set("regexMatch", value.NewValue(regexMatchFunc))

	// regexFind - находит первое совпадение регулярного выражения
	regexFindFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: regexFind() requires 2 arguments (pattern, string)")
		}

		pattern, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexFind() first argument must be string (pattern)")
		}

		text, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: regexFind() second argument must be string")
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error compiling regex pattern: %v", err))
		}

		result := regex.FindString(text)
		return value.NewString(result)
	}
	globalScope.Set("regexFind", value.NewValue(regexFindFunc))

	// regexFindAll - находит все совпадения регулярного выражения
	regexFindAllFunc := func(args []*value.Value) *value.Value {
		if len(args) < 2 || len(args) > 3 {
			return value.NewString("Error: regexFindAll() requires 2-3 arguments (pattern, string, [limit])")
		}

		pattern, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexFindAll() first argument must be string (pattern)")
		}

		text, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: regexFindAll() second argument must be string")
		}

		// Лимит совпадений (по умолчанию -1 = все)
		limit := -1
		if len(args) == 3 {
			if limitVal, ok := args[2].Any().(int64); ok {
				limit = int(limitVal)
			} else if limitVal, ok := args[2].Any().(float64); ok {
				limit = int(limitVal)
			} else {
				return value.NewString("Error: regexFindAll() third argument must be numeric (limit)")
			}
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error compiling regex pattern: %v", err))
		}

		matches := regex.FindAllString(text, limit)
		
		// Преобразуем в массив Value
		arrayValues := make([]*value.Value, len(matches))
		for i, match := range matches {
			arrayValues[i] = value.NewString(match)
		}
		
		return value.NewArray(arrayValues)
	}
	globalScope.Set("regexFindAll", value.NewValue(regexFindAllFunc))

	// ============ ЗАМЕНА ============

	// regexReplace - заменяет первое совпадение регулярного выражения
	regexReplaceFunc := func(args []*value.Value) *value.Value {
		if len(args) != 3 {
			return value.NewString("Error: regexReplace() requires 3 arguments (pattern, string, replacement)")
		}

		pattern, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexReplace() first argument must be string (pattern)")
		}

		text, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: regexReplace() second argument must be string")
		}

		replacement, ok := args[2].Any().(string)
		if !ok {
			return value.NewString("Error: regexReplace() third argument must be string (replacement)")
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error compiling regex pattern: %v", err))
		}

		// Заменяем только первое совпадение
		result := regex.ReplaceAllStringFunc(text, func(match string) string {
			return replacement
		})

		// Но на самом деле нужно заменить только первое
		firstMatch := regex.FindStringIndex(text)
		if firstMatch == nil {
			return value.NewString(text) // Нет совпадений
		}

		result = text[:firstMatch[0]] + replacement + text[firstMatch[1]:]
		return value.NewString(result)
	}
	globalScope.Set("regexReplace", value.NewValue(regexReplaceFunc))

	// regexReplaceAll - заменяет все совпадения регулярного выражения
	regexReplaceAllFunc := func(args []*value.Value) *value.Value {
		if len(args) != 3 {
			return value.NewString("Error: regexReplaceAll() requires 3 arguments (pattern, string, replacement)")
		}

		pattern, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexReplaceAll() first argument must be string (pattern)")
		}

		text, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: regexReplaceAll() second argument must be string")
		}

		replacement, ok := args[2].Any().(string)
		if !ok {
			return value.NewString("Error: regexReplaceAll() third argument must be string (replacement)")
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error compiling regex pattern: %v", err))
		}

		result := regex.ReplaceAllString(text, replacement)
		return value.NewString(result)
	}
	globalScope.Set("regexReplaceAll", value.NewValue(regexReplaceAllFunc))

	// ============ РАЗДЕЛЕНИЕ ============

	// regexSplit - разделяет строку по регулярному выражению
	regexSplitFunc := func(args []*value.Value) *value.Value {
		if len(args) < 2 || len(args) > 3 {
			return value.NewString("Error: regexSplit() requires 2-3 arguments (pattern, string, [limit])")
		}

		pattern, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexSplit() first argument must be string (pattern)")
		}

		text, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: regexSplit() second argument must be string")
		}

		// Лимит частей (по умолчанию -1 = все)
		limit := -1
		if len(args) == 3 {
			if limitVal, ok := args[2].Any().(int64); ok {
				limit = int(limitVal)
			} else if limitVal, ok := args[2].Any().(float64); ok {
				limit = int(limitVal)
			} else {
				return value.NewString("Error: regexSplit() third argument must be numeric (limit)")
			}
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error compiling regex pattern: %v", err))
		}

		parts := regex.Split(text, limit)
		
		// Преобразуем в массив Value
		arrayValues := make([]*value.Value, len(parts))
		for i, part := range parts {
			arrayValues[i] = value.NewString(part)
		}
		
		return value.NewArray(arrayValues)
	}
	globalScope.Set("regexSplit", value.NewValue(regexSplitFunc))

	// ============ ГРУППЫ ЗАХВАТА ============

	// regexGroups - извлекает группы захвата из первого совпадения
	regexGroupsFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: regexGroups() requires 2 arguments (pattern, string)")
		}

		pattern, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexGroups() first argument must be string (pattern)")
		}

		text, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: regexGroups() second argument must be string")
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error compiling regex pattern: %v", err))
		}

		matches := regex.FindStringSubmatch(text)
		if matches == nil {
			// Нет совпадений - возвращаем пустой массив
			return value.NewArray([]*value.Value{})
		}

		// Преобразуем группы в массив Value
		arrayValues := make([]*value.Value, len(matches))
		for i, match := range matches {
			arrayValues[i] = value.NewString(match)
		}
		
		return value.NewArray(arrayValues)
	}
	globalScope.Set("regexGroups", value.NewValue(regexGroupsFunc))

	// ============ ВАЛИДАЦИЯ И УТИЛИТЫ ============

	// regexValid - проверяет валидность регулярного выражения
	regexValidFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: regexValid() requires 1 argument (pattern)")
		}

		pattern, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexValid() argument must be string (pattern)")
		}

		_, err := regexp.Compile(pattern)
		return value.NewBool(err == nil)
	}
	globalScope.Set("regexValid", value.NewValue(regexValidFunc))

	// regexEscape - экранирует специальные символы в строке
	regexEscapeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: regexEscape() requires 1 argument (string)")
		}

		text, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexEscape() argument must be string")
		}

		escaped := regexp.QuoteMeta(text)
		return value.NewString(escaped)
	}
	globalScope.Set("regexEscape", value.NewValue(regexEscapeFunc))

	// ============ ДОПОЛНИТЕЛЬНЫЕ УТИЛИТЫ ============

	// regexCount - подсчитывает количество совпадений
	regexCountFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: regexCount() requires 2 arguments (pattern, string)")
		}

		pattern, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: regexCount() first argument must be string (pattern)")
		}

		text, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: regexCount() second argument must be string")
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error compiling regex pattern: %v", err))
		}

		matches := regex.FindAllString(text, -1)
		return value.NewInt64(int64(len(matches)))
	}
	globalScope.Set("regexCount", value.NewValue(regexCountFunc))

	// ============ ПРОСТЫЕ СТРОКОВЫЕ ФУНКЦИИ ============

	// Для совместимости добавим также простые строковые функции без регулярных выражений
	
	// stringSplit - простое разделение строки по подстроке
	stringSplitFunc := func(args []*value.Value) *value.Value {
		if len(args) < 2 || len(args) > 3 {
			return value.NewString("Error: stringSplit() requires 2-3 arguments (string, separator, [limit])")
		}

		text, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: stringSplit() first argument must be string")
		}

		separator, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: stringSplit() second argument must be string (separator)")
		}

		limit := -1
		if len(args) == 3 {
			if limitVal, ok := args[2].Any().(int64); ok {
				limit = int(limitVal)
			} else if limitVal, ok := args[2].Any().(float64); ok {
				limit = int(limitVal)
			} else {
				return value.NewString("Error: stringSplit() third argument must be numeric (limit)")
			}
		}

		var parts []string
		if limit == -1 {
			parts = strings.Split(text, separator)
		} else {
			parts = strings.SplitN(text, separator, limit)
		}
		
		// Преобразуем в массив Value
		arrayValues := make([]*value.Value, len(parts))
		for i, part := range parts {
			arrayValues[i] = value.NewString(part)
		}
		
		return value.NewArray(arrayValues)
	}
	globalScope.Set("stringSplit", value.NewValue(stringSplitFunc))
}