package ast

import "foo_lang/scope"


type ForExpr struct {
	InitExpr      Expr
	ConditionExpr Expr
	StepExpr      Expr
	BodyExpr      Expr
}

func NewForExpr(init, condition, step, body Expr) *ForExpr {
	return &ForExpr{
		InitExpr:      init,
		ConditionExpr: condition,
		StepExpr:      step,
		BodyExpr:      body,
	}
}

func (f *ForExpr) Eval() *Value {
	// Создаём локальную область видимости для цикла
	scope.GlobalScope.Push()
	defer scope.GlobalScope.Pop()

	statments := f.BodyExpr.(*BodyExpr).Statments
	var yield []any

	f.InitExpr.Eval()
	for {
		if !f.ConditionExpr.Eval().Bool() {
			break
		}
		for _, statment := range statments {
			switch stm := statment.(type) {
			case *ReturnExpr:
				val := stm.Eval()
				val.SetReturn(true)
				return val
			case *BreakExpr:
				// Break out of the loop
				goto endLoop
			case *YieldExpr:
				yield = append(yield, stm.Eval().Any())
			default:
				val := stm.Eval()
				if val != nil {
					if val.IsReturn() {
						return val
					}
					if val.IsBreak() {
						goto endLoop
					}
					if val.IsYield() {
						yield = append(yield, val.Any())
					}
				}
			}
		}
		f.StepExpr.Eval()
	}

endLoop:
	return NewValue(yield)
}
