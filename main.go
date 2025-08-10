package main

import (
	"fmt"
	"foo_lang/parser"
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

	mainFile, _ := os.ReadFile("examples/main.foo")

	exprs := parser.NewParser(mainFile).Parse()

	for _, expr := range exprs {
		expr.Eval()
	}
}
