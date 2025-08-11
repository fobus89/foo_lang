package ast

import "foo_lang/modules"

// Global parse function registry to avoid circular imports
var GlobalParseFunc modules.ParseFunc

// SetGlobalParseFunc sets the parse function for import statements
func SetGlobalParseFunc(parseFunc modules.ParseFunc) {
	GlobalParseFunc = parseFunc
}