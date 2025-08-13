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
		
		// 🔥 ИСПРАВЛЕНИЕ: Создаем новую область видимости для каждой итерации
		// Это изолирует переменные, объявленные с let внутри цикла
		scope.GlobalScope.Push()
		
		iterationBreak := false
		for _, statment := range statments {
			switch stm := statment.(type) {
			case *ReturnExpr:
				val := stm.Eval()
				val.SetReturn(true)
				scope.GlobalScope.Pop() // Очищаем область итерации перед return
				return val
			case *BreakExpr:
				// Break out of the loop immediately
				iterationBreak = true
				break
			case *YieldExpr:
				// Проверяем, если уже был break, не выполняем yield
				if !iterationBreak {
					yield = append(yield, stm.Eval().Any())
				}
			default:
				val := stm.Eval()
				if val != nil {
					if val.IsReturn() {
						scope.GlobalScope.Pop() // Очищаем область итерации перед return
						return val
					}
					if val.IsBreak() {
						iterationBreak = true
						break
					}
					if val.IsYield() {
						// Проверяем, если уже был break, не выполняем yield
						if !iterationBreak {
							yield = append(yield, val.Any())
						}
					}
				}
			}
			
			// Если был break, прекращаем обработку остальных statement в этой итерации
			if iterationBreak {
				break
			}
		}
		
		// Очищаем область видимости итерации
		scope.GlobalScope.Pop()
		
		// Проверяем, нужно ли выйти из цикла
		if iterationBreak {
			break
		}
		
		f.StepExpr.Eval()
	}

	return NewValue(yield)
}
