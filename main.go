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

	exprs := parser.NewParser(mainFile).Parse()

	for _, expr := range exprs {
		expr.Eval()
	}
}
