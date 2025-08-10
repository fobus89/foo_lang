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
	output := formatValue(val.Any())

	if !n.isPrint {
		fmt.Println(output)
	} else {
		fmt.Print(output)
	}

	return nil
}

func formatValue(v any) string {
	switch arr := v.(type) {
	case []any:
		if len(arr) == 0 {
			return "[]"
		}
		result := "["
		for i, item := range arr {
			if i > 0 {
				result += ", "
			}
			result += fmt.Sprintf("%v", item)
		}
		result += "]"
		return result
	default:
		return fmt.Sprintf("%v", v)
	}
}
