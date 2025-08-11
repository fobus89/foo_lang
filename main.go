package main

import (
	"fmt"
	"foo_lang/parser"
	"foo_lang/ast"
	"foo_lang/modules"
	"foo_lang/builtin"
	"foo_lang/scope"
	"os"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var mainFile []byte
	var err error
	
	if len(os.Args) > 1 {
		mainFile, err = os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
	} else {
		mainFile, _ = os.ReadFile("examples/main.foo")
	}

	// Set up global parse function for module imports
	parseFunc := func(code string) []modules.Expr {
		exprs := parser.NewParser(code).Parse()
		result := make([]modules.Expr, len(exprs))
		for i, expr := range exprs {
			result[i] = expr
		}
		return result
	}
	ast.SetGlobalParseFunc(parseFunc)
	
	// Инициализируем встроенные математические функции
	builtin.InitializeMathFunctions(scope.GlobalScope)
	
	// Инициализируем встроенные строковые функции
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	exprs := parser.NewParser(mainFile).Parse()

	for _, expr := range exprs {
		expr.Eval()
	}
}
