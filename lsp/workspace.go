package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"foo_lang/ast"
	"foo_lang/parser"
	"foo_lang/scope"
)

// Workspace управляет документами и анализом кода
type Workspace struct {
	documents map[string]*Document
	symbols   map[string][]*Symbol
	parsers   map[string]*parser.Parser     // AST парсеры для каждого документа
	scopes    map[string]*scope.ScopeStack // Области видимости
}

// Document представляет открытый документ
type Document struct {
	URI     string
	Content string
	Version int
	Lines   []string
}

// Symbol представляет символ в коде (функция, переменная, etc.)
type Symbol struct {
	Name          string
	Kind          SymbolKind
	Range         Range
	Detail        string
	Type          string
	Signature     string
	Documentation string
}

// SymbolKind виды символов
type SymbolKind int

const (
	SymbolKindVariable SymbolKind = iota
	SymbolKindFunction
	SymbolKindStruct
	SymbolKindInterface
	SymbolKindEnum
	SymbolKindMacro
	SymbolKindKeyword
	SymbolKindBuiltin
)

// MethodInfo информация о методах
type MethodInfo struct {
	Name          string
	Signature     string
	Documentation string
}

// NewWorkspace создает новый workspace
func NewWorkspace() *Workspace {
	return &Workspace{
		documents: make(map[string]*Document),
		symbols:   make(map[string][]*Symbol),
		parsers:   make(map[string]*parser.Parser),
		scopes:    make(map[string]*scope.ScopeStack),
	}
}

// AddDocument добавляет документ в workspace
func (w *Workspace) AddDocument(uri, content string) {
	doc := &Document{
		URI:     uri,
		Content: content,
		Version: 0,
		Lines:   strings.Split(content, "\n"),
	}
	
	w.documents[uri] = doc
	
	// Создаем AST парсер и scope для документа
	p := parser.NewParser(content)
	w.parsers[uri] = p
	w.scopes[uri] = scope.NewScopeStack()
	
	w.analyzeDocumentWithAST(doc)
}

// UpdateDocument обновляет содержимое документа
func (w *Workspace) UpdateDocument(uri, content string) {
	if doc, exists := w.documents[uri]; exists {
		doc.Content = content
		doc.Version++
		doc.Lines = strings.Split(content, "\n")
		
		// Пересоздаем AST парсер
		p := parser.NewParser(content)
		w.parsers[uri] = p
		
		w.analyzeDocumentWithAST(doc)
	}
}

// analyzeDocumentWithAST анализирует документ используя настоящий AST парсер
func (w *Workspace) analyzeDocumentWithAST(doc *Document) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error analyzing document %s with AST: %v", doc.URI, r)
			// Fallback to regex analysis
			w.analyzeDocument(doc)
		}
	}()
	
	log.Printf("🔍 Analyzing document with AST: %s", doc.URI)
	
	parser := w.parsers[doc.URI]
	if parser == nil {
		log.Printf("No parser found for document: %s", doc.URI)
		w.analyzeDocument(doc) // Fallback
		return
	}
	
	// Парсим весь документ в AST
	expressions := parser.Parse()
	if expressions == nil {
		log.Printf("Failed to parse program for document: %s", doc.URI)
		w.analyzeDocument(doc) // Fallback
		return
	}
	
	// Извлекаем символы из AST
	symbols := w.extractSymbolsFromAST(expressions, doc.URI)
	w.symbols[doc.URI] = symbols
	
	log.Printf("✅ Found %d symbols with AST parser: %s", len(symbols), doc.URI)
}

// extractSymbolsFromAST извлекает символы из AST узлов
func (w *Workspace) extractSymbolsFromAST(expressions []ast.Expr, uri string) []*Symbol {
	symbols := []*Symbol{}
	
	// Проходим по всем выражениям
	for _, expr := range expressions {
		// Простое извлечение символов - можно расширить позже
		if expr != nil {
			// Пока что добавляем примитивный анализ
			symbol := &Symbol{
				Name: "parsed_symbol",
				Kind: SymbolKindVariable,
				Range: Range{
					Start: Position{Line: 0, Character: 0},
					End:   Position{Line: 0, Character: 10},
				},
				Detail:        "ast_parsed",
				Type:          "unknown",
				Documentation: "Symbol extracted from AST",
			}
			symbols = append(symbols, symbol)
		}
	}
	
	return symbols
}

// inferTypeFromAST выводит тип из AST узла
func (w *Workspace) inferTypeFromAST(expr ast.Expr) string {
	if expr == nil {
		return "unknown"
	}
	
	// Простая реализация - можно расширить позже
	return "ast_type"
}

// analyzeDocument анализирует документ и извлекает символы (REGEX анализ с extension methods)
func (w *Workspace) analyzeDocument(doc *Document) {
	log.Printf("Analyzing document: %s", doc.URI)
	
	symbols := []*Symbol{}
	lines := doc.Lines
	
	// Отслеживаем extension блоки
	currentExtension := ""
	
	for lineNum, line := range lines {
		// Анализируем extension блоки: extension TypeName {
		if matches := regexp.MustCompile(`extension\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\{`).FindStringSubmatch(line); matches != nil {
			currentExtension = matches[1]
			log.Printf("Found extension for type: %s", currentExtension)
			continue
		}
		
		// Если мы внутри extension блока, анализируем методы
		if currentExtension != "" {
			// Конец extension блока
			if strings.Contains(line, "}") && !strings.Contains(line, "fn") {
				log.Printf("End of extension for type: %s", currentExtension)
				currentExtension = ""
				continue
			}
			
			// Анализируем методы внутри extension: fn methodName() -> returnType
			if matches := regexp.MustCompile(`fn\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\)(?:\s*->\s*([^{]+))?`).FindStringSubmatch(line); matches != nil {
				methodName := matches[1]
				params := matches[2]
				returnType := "void"
				if matches[3] != "" {
					returnType = strings.TrimSpace(matches[3])
				}
				
				symbol := &Symbol{
					Name: fmt.Sprintf("%s.%s", currentExtension, methodName),
					Kind: SymbolKindFunction,
					Range: Range{
						Start: Position{Line: lineNum, Character: strings.Index(line, methodName)},
						End:   Position{Line: lineNum, Character: strings.Index(line, methodName) + len(methodName)},
					},
					Detail:        fmt.Sprintf("extension method for %s", currentExtension),
					Type:          returnType,
					Signature:     fmt.Sprintf("fn %s(%s) -> %s", methodName, params, returnType),
					Documentation: fmt.Sprintf("Extension method %s for type %s", methodName, currentExtension),
				}
				symbols = append(symbols, symbol)
				log.Printf("Found extension method: %s.%s", currentExtension, methodName)
				continue
			}
		}
		// Анализируем переменные: let name = value
		if matches := regexp.MustCompile(`let\s+([a-zA-Z_][a-zA-Z0-9_]*)`).FindStringSubmatch(line); matches != nil {
			symbol := &Symbol{
				Name: matches[1],
				Kind: SymbolKindVariable,
				Range: Range{
					Start: Position{Line: lineNum, Character: strings.Index(line, matches[1])},
					End:   Position{Line: lineNum, Character: strings.Index(line, matches[1]) + len(matches[1])},
				},
				Detail:        "variable",
				Type:          w.inferType(line),
				Documentation: fmt.Sprintf("Variable %s", matches[1]),
			}
			symbols = append(symbols, symbol)
		}
		
		// Анализируем константы: const name = value
		if matches := regexp.MustCompile(`const\s+([a-zA-Z_][a-zA-Z0-9_]*)`).FindStringSubmatch(line); matches != nil {
			symbol := &Symbol{
				Name: matches[1],
				Kind: SymbolKindVariable,
				Range: Range{
					Start: Position{Line: lineNum, Character: strings.Index(line, matches[1])},
					End:   Position{Line: lineNum, Character: strings.Index(line, matches[1]) + len(matches[1])},
				},
				Detail:        "constant",
				Type:          w.inferType(line),
				Documentation: fmt.Sprintf("Constant %s", matches[1]),
			}
			symbols = append(symbols, symbol)
		}
		
		// Анализируем функции: fn name(params) -> returnType
		if matches := regexp.MustCompile(`fn\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\)(?:\s*->\s*([^{]+))?`).FindStringSubmatch(line); matches != nil {
			params := matches[2]
			returnType := "void"
			if matches[3] != "" {
				returnType = strings.TrimSpace(matches[3])
			}
			
			symbol := &Symbol{
				Name: matches[1],
				Kind: SymbolKindFunction,
				Range: Range{
					Start: Position{Line: lineNum, Character: strings.Index(line, matches[1])},
					End:   Position{Line: lineNum, Character: strings.Index(line, matches[1]) + len(matches[1])},
				},
				Detail:        "function",
				Type:          returnType,
				Signature:     fmt.Sprintf("fn %s(%s) -> %s", matches[1], params, returnType),
				Documentation: fmt.Sprintf("Function %s with parameters (%s) returning %s", matches[1], params, returnType),
			}
			symbols = append(symbols, symbol)
		}
		
		// Анализируем структуры: struct name
		if matches := regexp.MustCompile(`struct\s+([a-zA-Z_][a-zA-Z0-9_]*)`).FindStringSubmatch(line); matches != nil {
			symbol := &Symbol{
				Name: matches[1],
				Kind: SymbolKindStruct,
				Range: Range{
					Start: Position{Line: lineNum, Character: strings.Index(line, matches[1])},
					End:   Position{Line: lineNum, Character: strings.Index(line, matches[1]) + len(matches[1])},
				},
				Detail:        "struct",
				Type:          "struct",
				Documentation: fmt.Sprintf("Struct %s", matches[1]),
			}
			symbols = append(symbols, symbol)
		}
		
		// Анализируем интерфейсы: interface name
		if matches := regexp.MustCompile(`interface\s+([a-zA-Z_][a-zA-Z0-9_]*)`).FindStringSubmatch(line); matches != nil {
			symbol := &Symbol{
				Name: matches[1],
				Kind: SymbolKindInterface,
				Range: Range{
					Start: Position{Line: lineNum, Character: strings.Index(line, matches[1])},
					End:   Position{Line: lineNum, Character: strings.Index(line, matches[1]) + len(matches[1])},
				},
				Detail:        "interface",
				Type:          "interface",
				Documentation: fmt.Sprintf("Interface %s", matches[1]),
			}
			symbols = append(symbols, symbol)
		}
	}
	
	w.symbols[doc.URI] = symbols
	log.Printf("Found %d symbols in document: %s", len(symbols), doc.URI)
}

// inferType выводит тип из строки кода
func (w *Workspace) inferType(line string) string {
	// Явная типизация: let name: type = value
	if matches := regexp.MustCompile(`:\s*([a-zA-Z_][a-zA-Z0-9_|?]*)`).FindStringSubmatch(line); matches != nil {
		return matches[1]
	}
	
	// Вывод из значения
	if strings.Contains(line, "\"") {
		return "string"
	}
	if regexp.MustCompile(`=\s*\d+\.\d+`).MatchString(line) {
		return "float"
	}
	if regexp.MustCompile(`=\s*\d+`).MatchString(line) {
		return "int"
	}
	if regexp.MustCompile(`=\s*(true|false)`).MatchString(line) {
		return "bool"
	}
	if strings.Contains(line, "[") && strings.Contains(line, "]") {
		return "array"
	}
	if strings.Contains(line, "{") && strings.Contains(line, "}") {
		return "object"
	}
	
	return "unknown"
}

// GetCompletions возвращает автодополнения для позиции
func (w *Workspace) GetCompletions(uri string, pos Position) []map[string]interface{} {
	log.Printf("🔍 Getting completions for position %d:%d in %s", pos.Line, pos.Character, uri)
	
	// Получаем текущую строку для анализа контекста
	currentLine := w.getCurrentLine(uri, pos.Line)
	log.Printf("Current line: '%s'", currentLine)
	
	// Проверяем, набирает ли пользователь после точки (extension methods)
	if pos.Character > 0 {
		// Получаем текст до курсора
		textBeforeCursor := ""
		if pos.Character <= len(currentLine) {
			textBeforeCursor = currentLine[:pos.Character]
		}
		
		// Ищем паттерн "variableName." 
		if matches := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\.$`).FindStringSubmatch(textBeforeCursor); matches != nil {
			variableName := matches[1]
			log.Printf("🎯 Found dot completion for variable: %s", variableName)
			
			// Определяем тип переменной
			variableType := w.getVariableType(uri, variableName)
			log.Printf("Variable type: %s", variableType)
			
			// Возвращаем только extension методы для этого типа
			return w.getExtensionMethodsForType(uri, variableType)
		}
	}
	
	completions := []map[string]interface{}{}
	
	// Ключевые слова языка
	keywords := []string{
		"let", "const", "fn", "struct", "enum", "interface", "impl", "extension",
		"if", "else", "for", "match", "return", "yield", "break", "continue",
		"async", "await", "sleep", "Promise",
		"import", "export", "from", "as",
		"macro", "quote", "unquote", "typeof", "type",
		"true", "false", "null",
		"int", "float", "string", "bool",
	}
	
	for _, keyword := range keywords {
		completions = append(completions, map[string]interface{}{
			"label": keyword,
			"kind":  CompletionItemKindKeyword,
		})
	}
	
	// Встроенные функции
	builtins := map[string]string{
		"print":      "fn print(value) -> void",
		"println":    "fn println(value) -> void",
		"sin":        "fn sin(x: float) -> float",
		"cos":        "fn cos(x: float) -> float",
		"sqrt":       "fn sqrt(x: float) -> float",
		"abs":        "fn abs(x: number) -> number",
		"readFile":   "fn readFile(path: string) -> string",
		"writeFile":  "fn writeFile(path: string, content: string) -> string",
		"httpGet":    "fn httpGet(url: string) -> string",
		"httpPost":   "fn httpPost(url: string, data: string) -> string",
		"newChannel": "fn newChannel(size: int, name: string) -> Channel",
		"send":       "fn send(ch: Channel, value: any) -> void",
		"receive":    "fn receive(ch: Channel) -> any",
		"sha256Hash": "fn sha256Hash(data: string) -> string",
		"base64Encode": "fn base64Encode(data: string) -> string",
		"regexMatch": "fn regexMatch(pattern: string, text: string) -> bool",
		"regexReplace": "fn regexReplace(pattern: string, text: string, replacement: string) -> string",
		"now":        "fn now() -> Time",
		"timeFormat": "fn timeFormat(time: Time, format: string) -> string",
		"sleep":      "fn sleep(ms: int) -> void",
		"typeof":     "fn typeof(value: any) -> string",
		"jsonParse":  "fn jsonParse(json: string) -> any",
		"jsonStringify": "fn jsonStringify(value: any) -> string",
		// Встроенные методы
		"toString":   "fn toString() -> string",
		"length":     "fn length() -> int",
		"push":       "fn push(value: any) -> void",
		"pop":        "fn pop() -> any",
	}
	
	for name, signature := range builtins {
		completions = append(completions, map[string]interface{}{
			"label":  name,
			"kind":   CompletionItemKindFunction,
			"detail": signature,
		})
	}
	
	// Пользовательские символы из текущего файла
	if symbols, exists := w.symbols[uri]; exists {
		for _, symbol := range symbols {
			kind := CompletionItemKindVariable
			if symbol.Kind == SymbolKindFunction {
				kind = CompletionItemKindFunction
			} else if symbol.Kind == SymbolKindStruct {
				kind = CompletionItemKindStruct
			} else if symbol.Kind == SymbolKindInterface {
				kind = CompletionItemKindInterface
			}
			
			completions = append(completions, map[string]interface{}{
				"label":  symbol.Name,
				"kind":   kind,
				"detail": symbol.Detail,
			})
		}
	}
	
	log.Printf("✅ Returning %d completions", len(completions))
	return completions
}

// GetHover возвращает информацию при наведении
func (w *Workspace) GetHover(uri string, pos Position) *Hover {
	if symbols, exists := w.symbols[uri]; exists {
		for _, symbol := range symbols {
			if w.isPositionInRange(pos, symbol.Range) {
				content := MarkupContent{
					Kind:  "markdown",
					Value: fmt.Sprintf("**%s** `%s`\n\n%s", symbol.Name, symbol.Type, symbol.Documentation),
				}
				
				if symbol.Signature != "" {
					content.Value = fmt.Sprintf("```foo\n%s\n```\n\n%s", symbol.Signature, symbol.Documentation)
				}
				
				return &Hover{
					Contents: content,
					Range:    &symbol.Range,
				}
			}
		}
	}
	
	return nil
}

// GetDefinitions возвращает определения для символа
func (w *Workspace) GetDefinitions(uri string, pos Position) []Location {
	// Базовая реализация - можно расширить
	return []Location{}
}

// GetDiagnostics возвращает диагностику для документа
func (w *Workspace) GetDiagnostics(uri string) []Diagnostic {
	diagnostics := []Diagnostic{}
	
	if doc, exists := w.documents[uri]; exists {
		for lineNum, line := range doc.Lines {
			// Проверка незакрытых скобок
			openBraces := strings.Count(line, "{")
			closeBraces := strings.Count(line, "}")
			if openBraces != closeBraces && strings.TrimSpace(line) != "" {
				diagnostics = append(diagnostics, Diagnostic{
					Range: Range{
						Start: Position{Line: lineNum, Character: 0},
						End:   Position{Line: lineNum, Character: len(line)},
					},
					Severity: DiagnosticSeverityWarning,
					Message:  "Mismatched braces in line",
				})
			}
			
			// Проверка undefined функций
			if matches := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`).FindAllStringSubmatch(line, -1); matches != nil {
				for _, match := range matches {
					funcName := match[1]
					if !w.isSymbolDefined(uri, funcName) && !w.isBuiltinFunction(funcName) {
						diagnostics = append(diagnostics, Diagnostic{
							Range: Range{
								Start: Position{Line: lineNum, Character: strings.Index(line, funcName)},
								End:   Position{Line: lineNum, Character: strings.Index(line, funcName) + len(funcName)},
							},
							Severity: DiagnosticSeverityError,
							Message:  fmt.Sprintf("Undefined function '%s'", funcName),
						})
					}
				}
			}
		}
	}
	
	return diagnostics
}

// isPositionInRange проверяет находится ли позиция в диапазоне
func (w *Workspace) isPositionInRange(pos Position, r Range) bool {
	if pos.Line < r.Start.Line || pos.Line > r.End.Line {
		return false
	}
	if pos.Line == r.Start.Line && pos.Character < r.Start.Character {
		return false
	}
	if pos.Line == r.End.Line && pos.Character > r.End.Character {
		return false
	}
	return true
}

// isSymbolDefined проверяет определен ли символ в документе
func (w *Workspace) isSymbolDefined(uri, name string) bool {
	if symbols, exists := w.symbols[uri]; exists {
		for _, symbol := range symbols {
			if symbol.Name == name {
				return true
			}
		}
	}
	return false
}

// isBuiltinFunction проверяет является ли функция встроенной
func (w *Workspace) isBuiltinFunction(name string) bool {
	builtins := []string{
		"print", "println", "sin", "cos", "sqrt", "abs", "min", "max",
		"strlen", "charAt", "substring", "readFile", "writeFile",
		"httpGet", "httpPost", "newChannel", "send", "receive",
		"sha256Hash", "base64Encode", "regexMatch", "regexReplace",
		"now", "timeFormat", "sleep", "typeof", "jsonParse", "jsonStringify",
		// Также добавляем встроенные методы
		"toString", "length", "push", "pop", "slice", "charAt", "toUpper", "toLower",
	}
	
	for _, builtin := range builtins {
		if name == builtin {
			return true
		}
	}
	
	return false
}

// getCurrentLine возвращает текущую строку из документа
func (w *Workspace) getCurrentLine(uri string, lineNum int) string {
	if doc, exists := w.documents[uri]; exists {
		if lineNum >= 0 && lineNum < len(doc.Lines) {
			return doc.Lines[lineNum]
		}
	}
	return ""
}

// getVariableType определяет тип переменной по её имени
func (w *Workspace) getVariableType(uri, variableName string) string {
	log.Printf("🔍 Looking for variable type: %s", variableName)
	
	// Ищем переменную среди символов
	if symbols, exists := w.symbols[uri]; exists {
		for _, symbol := range symbols {
			if symbol.Name == variableName {
				log.Printf("Found variable %s with type: %s", variableName, symbol.Type)
				return symbol.Type
			}
		}
	}
	
	// Ищем в коде документа (fallback)
	if doc, exists := w.documents[uri]; exists {
		for _, line := range doc.Lines {
			// Ищем объявление переменной: let variableName: type = ...
			pattern := fmt.Sprintf(`let\s+%s\s*:\s*([a-zA-Z_][a-zA-Z0-9_]*)`, regexp.QuoteMeta(variableName))
			if matches := regexp.MustCompile(pattern).FindStringSubmatch(line); matches != nil {
				log.Printf("Found typed variable declaration: %s -> %s", variableName, matches[1])
				return matches[1]
			}
			
			// Ищем создание структуры: let variableName = StructName{...}
			pattern = fmt.Sprintf(`let\s+%s\s*=\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*\{`, regexp.QuoteMeta(variableName))
			if matches := regexp.MustCompile(pattern).FindStringSubmatch(line); matches != nil {
				log.Printf("Found struct instantiation: %s -> %s", variableName, matches[1])
				return matches[1]
			}
		}
	}
	
	log.Printf("Could not determine type for variable: %s", variableName)
	return "unknown"
}

// getExtensionMethodsForType возвращает extension методы для конкретного типа
func (w *Workspace) getExtensionMethodsForType(uri, typeName string) []map[string]interface{} {
	log.Printf("🎯 Getting extension methods for type: %s", typeName)
	
	completions := []map[string]interface{}{}
	
	// Встроенные методы для разных типов
	builtinMethods := w.getBuiltinMethodsForType(typeName)
	completions = append(completions, builtinMethods...)
	
	// Extension методы из символов
	if symbols, exists := w.symbols[uri]; exists {
		for _, symbol := range symbols {
			// Ищем extension методы в формате "TypeName.methodName"
			if strings.HasPrefix(symbol.Name, typeName+".") {
				methodName := strings.TrimPrefix(symbol.Name, typeName+".")
				log.Printf("Found extension method: %s for type %s", methodName, typeName)
				
				completions = append(completions, map[string]interface{}{
					"label":           methodName,
					"kind":            CompletionItemKindMethod,
					"detail":          symbol.Detail,
					"documentation":   symbol.Documentation,
					"insertText":      methodName,
					"sortText":        "0" + methodName, // Приоритет extension методам
				})
			}
		}
	}
	
	log.Printf("✅ Found %d methods for type %s", len(completions), typeName)
	return completions
}

// getBuiltinMethodsForType возвращает встроенные методы для типа
func (w *Workspace) getBuiltinMethodsForType(typeName string) []map[string]interface{} {
	completions := []map[string]interface{}{}
	
	switch typeName {
	case "string":
		methods := map[string]string{
			"length":    "fn length() -> int",
			"charAt":    "fn charAt(index: int) -> string",
			"substring": "fn substring(start: int, end: int) -> string",
			"toUpper":   "fn toUpper() -> string",
			"toLower":   "fn toLower() -> string",
			"toString":  "fn toString() -> string",
		}
		for name, signature := range methods {
			completions = append(completions, map[string]interface{}{
				"label":  name,
				"kind":   CompletionItemKindMethod,
				"detail": signature,
			})
		}
		
	case "int", "float":
		methods := map[string]string{
			"toString": "fn toString() -> string",
			"abs":      "fn abs() -> " + typeName,
		}
		if typeName == "float" {
			methods["round"] = "fn round() -> int"
			methods["floor"] = "fn floor() -> int"
			methods["ceil"] = "fn ceil() -> int"
		}
		for name, signature := range methods {
			completions = append(completions, map[string]interface{}{
				"label":  name,
				"kind":   CompletionItemKindMethod,
				"detail": signature,
			})
		}
		
	case "array":
		methods := map[string]string{
			"length": "fn length() -> int",
			"push":   "fn push(value: any) -> void",
			"pop":    "fn pop() -> any",
			"slice":  "fn slice(start: int, end: int) -> array",
		}
		for name, signature := range methods {
			completions = append(completions, map[string]interface{}{
				"label":  name,
				"kind":   CompletionItemKindMethod,
				"detail": signature,
			})
		}
		
	default:
		// Для всех типов доступен toString
		completions = append(completions, map[string]interface{}{
			"label":  "toString",
			"kind":   CompletionItemKindMethod,
			"detail": "fn toString() -> string",
		})
	}
	
	return completions
}