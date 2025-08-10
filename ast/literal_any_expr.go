package ast

type LiteralAny struct {
	Value any
}

func NewLiteralAny(value any) *LiteralAny {
	return &LiteralAny{Value: value}
}

func (l *LiteralAny) Eval() *Value {
	return NewValue(l.Value)
}
