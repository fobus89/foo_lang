package main

import "foo_lang/parser"

func main() {
	p, err := parser.NewParserFromFile("examples/main.foo")
	if err != nil {
		panic(err)
	}

	exprs := p.ParseWithModules()
	for _, expr := range exprs {
		expr.Eval()
	}
}
