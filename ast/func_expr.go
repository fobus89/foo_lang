package ast

import (
	"fmt"
	"foo_lang/scope"
)

type FuncStatment struct {
	funcName string
	args     []map[string]Expr
	body     Expr
	isMacro  bool
}

func NewFuncStatment(funcName string, args []map[string]Expr, body Expr, isMacro bool) *FuncStatment {
	f := &FuncStatment{
		funcName: funcName,
		args:     args,
		body:     body,
		isMacro:  isMacro,
	}

	scope.GlobalScope.Set(funcName, NewValue(f))

	return f
}

func (f *FuncStatment) Name() string {
	return f.funcName
}

func (f *FuncStatment) Params() []string {
	var args []string
	{
		for _, v := range f.args {
			for k := range v {
				args = append(args, k)
			}
		}
	}
	return args
}

func (f *FuncStatment) Eval() *Value {
	// Function definitions don't return values, they register the function in scope
	// The function is already registered in NewFuncStatment constructor
	return nil
}

func (f *FuncStatment) IsMacro() bool {
	return f.isMacro
}

func (f *FuncStatment) String() string {
	return fmt.Sprintf("func %s(%s) { %s }", f.funcName, f.args, f.body)
}
