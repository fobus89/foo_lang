package ast

import "strings"

type LiteralString struct {
	Value string
}

func NewLiteralString(value string) *LiteralString {
	// Process escape sequences
	value = strings.ReplaceAll(value, `\n`, "\n")
	value = strings.ReplaceAll(value, `\t`, "\t")
	value = strings.ReplaceAll(value, `\r`, "\r")
	value = strings.ReplaceAll(value, `\\`, "\\")
	value = strings.ReplaceAll(value, `\"`, "\"")
	return &LiteralString{Value: value}
}

func (l *LiteralString) Eval() *Value {
	return NewValue(l.Value)
}
