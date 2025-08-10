package ast

import (
	"fmt"
)

type PrintExpr struct {
	Expr    Expr
	isPrint bool
}

func NewPrintExpr(expr Expr, isPrint bool) *PrintExpr {
	return &PrintExpr{Expr: expr, isPrint: isPrint}
}

func (n *PrintExpr) Eval() *Value {

	if n.Expr == nil {
		return nil
	}

	val := n.Expr.Eval()
	output := FormatValue(val.Any())

	if !n.isPrint {
		fmt.Println(output)
	} else {
		fmt.Print(output)
	}

	return nil
}

func FormatValue(v any) string {
	switch val := v.(type) {
	case []any:
		if len(val) == 0 {
			return "[]"
		}
		result := "["
		for i, item := range val {
			if i > 0 {
				result += ", "
			}
			result += fmt.Sprintf("%v", item)
		}
		result += "]"
		return result
	case *ResultValue:
		if val.IsOk() {
			return fmt.Sprintf("Ok(%s)", FormatValue(val.GetValue().Any()))
		} else {
			return fmt.Sprintf("Err(%s)", FormatValue(val.GetValue().Any()))
		}
	default:
		return fmt.Sprintf("%v", v)
	}
}
