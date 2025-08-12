package ast

// NullExpr представляет null литерал
type NullExpr struct {}

func NewNullExpr() *NullExpr {
	return &NullExpr{}
}

func (n *NullExpr) Eval() *Value {
	// Null возвращает специальное значение с nil
	return NewValue(nil)
}