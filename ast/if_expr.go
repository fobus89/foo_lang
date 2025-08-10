package ast


type IfExpr struct {
	Condition []Expr
	Then      []Expr
	Else      Expr
}

func NewIfExpr(conditions []Expr, then []Expr, elseExpr Expr) *IfExpr {
	return &IfExpr{Condition: conditions, Then: then, Else: elseExpr}
}

func (i *IfExpr) Eval() *Value {

	for index, cond := range i.Condition {
		if cond.Eval().Bool() {
			body := i.Then[index].(*BodyExpr)
			{
				var result *Value
				for _, statment := range body.Statments {
					switch stm := statment.(type) {
					case *ReturnExpr:
						val := stm.Eval()
						val.SetReturn(true)
						return val
					case *BreakExpr:
						return stm.Eval()
					default:
						result = stm.Eval()
						if result != nil && (result.IsReturn() || result.IsYield() || result.IsBreak()) {
							return result
						}
					}
				}
				return result
			}
		}
	}

	if i.Else == nil {
		return nil
	}

	return i.Else.Eval()
}
