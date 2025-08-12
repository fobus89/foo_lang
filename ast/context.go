package ast

import "foo_lang/modules"

// Global parse function registry to avoid circular imports
var GlobalParseFunc modules.ParseFunc

// Global current file context for imports
var CurrentFileContext string

// SetGlobalParseFunc sets the parse function for import statements
func SetGlobalParseFunc(parseFunc modules.ParseFunc) {
	GlobalParseFunc = parseFunc
}

// SetCurrentFileContext устанавливает контекст текущего файла для импортов
func SetCurrentFileContext(filePath string) {
	CurrentFileContext = filePath
}

// GetCurrentFileContext возвращает контекст текущего файла
func GetCurrentFileContext() string {
	return CurrentFileContext
}