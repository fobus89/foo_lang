package ast

type LiteralString struct {
	Value string
}

func NewLiteralString(value string) *LiteralString {
	return &LiteralString{Value: value}
}

func (l *LiteralString) Eval() *Value {
	return NewValue(l.Value)
}
