package ast

type MatchExpr struct {
	Value Expr
	Arms  []MatchArm
}

type MatchArm struct {
	Pattern Expr
	Body    Expr
}

//	let q =	match 1 > 2 && 1 || false {
//		1 => 1+22,
//		2 => ("b"),
//		true => true + 2,
//		false => false + 1,
//		_ => ("a"),
//	}
func NewMatchArm(pattern Expr, body Expr) MatchArm {
	return MatchArm{Pattern: pattern, Body: body}
}

func NewMatchExpr(value Expr, arms []MatchArm) *MatchExpr {
	return &MatchExpr{Value: value, Arms: arms}
}

func (m *MatchExpr) Eval() *Value {

	for _, arm := range m.Arms {
		value := m.Value.Eval()
		pattern := arm.Pattern.Eval()

		if pattern.IsString() && pattern.String() == "_" {
			return arm.Body.Eval()
		}

		if value.IsString() && pattern.IsString() && value.String() == pattern.String() {
			return arm.Body.Eval()
		} else if value.IsInt64() && pattern.IsInt64() && value.Int64() == pattern.Int64() {
			return arm.Body.Eval()
		} else if value.IsFloat64() && pattern.IsFloat64() && value.Float64() == pattern.Float64() {
			return arm.Body.Eval()
		} else if value.IsBool() && pattern.IsBool() && value.Bool() == pattern.Bool() {
			return arm.Body.Eval()
		}
	}

	return nil
}
