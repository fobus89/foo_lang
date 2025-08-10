package ast

// MemberExpr представляет доступ к полю объекта (object.property)
type MemberExpr struct {
	Object   Expr
	Property string
}

func NewMemberExpr(object Expr, property string) *MemberExpr {
	return &MemberExpr{
		Object:   object,
		Property: property,
	}
}

func (m *MemberExpr) Eval() *Value {
	obj := m.Object.Eval()
	
	// Проверяем, что объект - это словарь
	if objMap, ok := obj.Any().(map[string]*Value); ok {
		if value, exists := objMap[m.Property]; exists {
			return value
		}
		panic("property '" + m.Property + "' does not exist")
	}
	
	panic("cannot access property of non-object")
}