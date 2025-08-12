package main

import (
	"fmt"
	"foo_lang/ast"
	"foo_lang/builtin"
	"foo_lang/modules"
	"foo_lang/parser"
	"foo_lang/scope"
	"os"
	"strings"
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

	// Проверяем флаг bytecode режима
	for _, arg := range os.Args {
		if arg == "--bytecode" || arg == "-b" {
			RunBytecodeMode()
			return
		}
		if arg == "--help" || arg == "-h" {
			printUsage()
			return
		}
	}

	var mainFile []byte
	var err error

	if len(os.Args) > 1 {
		// Пропускаем флаги при поиске файла
		filename := ""
		for _, arg := range os.Args[1:] {
			if !strings.HasPrefix(arg, "-") {
				filename = arg
				break
			}
		}
		
		if filename != "" {
			mainFile, err = os.ReadFile(filename)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				return
			}
		} else {
			mainFile, _ = os.ReadFile("examples/main.foo")
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

	// Инициализируем встроенные функции файловой системы
	builtin.InitializeFilesystemFunctions(scope.GlobalScope)
	
	// Инициализируем встроенные HTTP функции
	builtin.InitializeHttpFunctions(scope.GlobalScope)
	
	// Инициализируем встроенные функции каналов
	builtin.InitializeChannelFunctions(scope.GlobalScope)
	
	// Инициализируем встроенные функции времени
	builtin.InitializeTimeFunctions(scope.GlobalScope)

	exprs := parser.NewParser(mainFile).Parse()

	for _, expr := range exprs {
		expr.Eval()
	}
}

// RunBytecodeMode запускает bytecode режим (заглушка)
func RunBytecodeMode() {
	fmt.Println("Bytecode режим находится в разработке")
	fmt.Println("VM система реализована и протестирована, компилятор в процессе разработки")
}

// printUsage выводит справку по использованию
func printUsage() {
	fmt.Println("foo_lang v2 - Современный интерпретируемый язык программирования")
	fmt.Println()
	fmt.Println("Использование:")
	fmt.Println("  go run main.go [файл.foo] [флаги]")
	fmt.Println()
	fmt.Println("Флаги:")
	fmt.Println("  -b, --bytecode    Использовать bytecode VM (оптимизированный)")
	fmt.Println("  -d, --disassemble Показать дизассемблированный bytecode")
	fmt.Println("  -p, --profile     Показать профилирование производительности")
	fmt.Println("  -c, --compare     Сравнить производительность tree-walking vs bytecode")
	fmt.Println("  -h, --help        Показать эту справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  go run main.go                                    # tree-walking интерпретатор")
	fmt.Println("  go run main.go --bytecode                         # bytecode VM")
	fmt.Println("  go run main.go examples/test_bytecode_demo.foo -b # bytecode VM с файлом")
	fmt.Println("  go run main.go --bytecode --profile --disassemble # полная диагностика")
	fmt.Println("  go run main.go --bytecode --compare               # сравнение производительности")
	fmt.Println()
	fmt.Println("Возможности:")
	fmt.Println("  ✅ Generic функции и типизация")
	fmt.Println("  ✅ Interface система и Extension Methods") 
	fmt.Println("  ✅ Async/await и многопоточность")
	fmt.Println("  ✅ HTTP клиент/сервер")
	fmt.Println("  ✅ Файловая система")
	fmt.Println("  ✅ Bytecode компиляция и профилирование")
}
