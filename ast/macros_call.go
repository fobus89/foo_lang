package ast

import (
	"fmt"
)

type MacrosCallExpr struct {
	funcName string
	args     []Expr
}

func NewMacrosCallExpr(funcName string, args []Expr) *MacrosCallExpr {

	f := &MacrosCallExpr{
		funcName: funcName,
		args:     args,
	}

	return f
}

func (f *MacrosCallExpr) Eval() *Value {

	val := Container[f.funcName]
	{
		if val == nil {
			return nil
		}
	}

	fnStatment, ok := val.Any().(*FuncStatment)
	{
		if !ok {
			panic("not a function")
		}
	}

	bodyStm := fnStatment.body.(*BodyExpr)

	expected := len(fnStatment.args)
	passed := len(f.args)

	if passed > expected {
		panic(fmt.Sprintf("too many arguments: expected %d, got %d", expected, passed))
	}

	for i, arg := range fnStatment.args {
		for name, expr := range arg {
			if i < len(f.args) {
				value := f.args[i].Eval()
				Container[name] = value
			} else if expr != nil {
				defaultValue := expr.Eval()
				Container[name] = defaultValue
			} else {
				panic(fmt.Sprintf("missing required argument: %s", name))
			}
		}
	}

	for _, stm := range bodyStm.Statments {
		if stm == nil {
			continue
		}

		if _, ok := stm.(*ReturnExpr); ok {
			return stm.Eval()
		}

	}

	return nil
}
