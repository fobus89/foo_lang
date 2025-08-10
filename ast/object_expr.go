package ast

// ObjectExpr представляет объект как словарь ключ-значение
type ObjectExpr struct {
	Fields map[string]Expr
}

func NewObjectExpr(fields map[string]Expr) *ObjectExpr {
	return &ObjectExpr{Fields: fields}
}

func (o *ObjectExpr) Eval() *Value {
	result := make(map[string]*Value)
	for key, expr := range o.Fields {
		result[key] = expr.Eval()
	}
	return NewValue(result)
}