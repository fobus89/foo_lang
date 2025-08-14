package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"foo_lang/parser"
	"foo_lang/bytecode"
	"foo_lang/scope"
	"foo_lang/builtin"
)

// Альтернативная точка входа для выполнения через bytecode VM
func mainBytecode() {
	// Проверяем аргументы командной строки
	filename := "examples/test_bytecode_demo.foo"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	// Читаем файл
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Printf("🚀 Запуск foo_lang через Bytecode VM: %s\n", filename)
	fmt.Println(strings.Repeat("=", 51))

	// Инициализируем scope с встроенными функциями
	globalScope := scope.NewScopeStack()
	builtin.InitializeMathFunctions(globalScope)
	builtin.InitializeStringFunctions(globalScope)
	builtin.InitializeFilesystemFunctions(globalScope)
	builtin.InitializeHttpFunctions(globalScope)

	// Парсим код (для статистики)
	startParse := time.Now()
	_ = parser.NewParser(content).Parse()
	parseTime := time.Since(startParse)

	// Создаем chunk вручную (компилятор пока не реализован)
	startCompile := time.Now()
	chunk := bytecode.NewChunk()
	// Простая демо-программа: 10 + 5 = 15
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(10))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(5))}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 1)
	compileTime := time.Since(startCompile)

	// Выводим статистику компиляции
	fmt.Printf("📊 Статистика компиляции:\n")
	fmt.Printf("   Время парсинга: %v\n", parseTime)
	fmt.Printf("   Время компиляции: %v\n", compileTime)
	fmt.Printf("   Инструкций: %d\n", len(chunk.Code))
	fmt.Printf("   Констант: %d\n", len(chunk.Constants))
	fmt.Println()

	// Опционально выводим дизассемблированный код
	if shouldShowDisassembly() {
		fmt.Println("🔍 Дизассемблированный bytecode:")
		bytecode.DisassembleChunk(chunk, filename)
		fmt.Println()
	}

	// Выполняем в VM
	startExecution := time.Now()
	vm := bytecode.NewVM(chunk, globalScope)
	result := vm.Run()
	executionTime := time.Since(startExecution)

	// Выводим результат выполнения
	fmt.Println()
	fmt.Println(strings.Repeat("=", 51))
	
	if result != nil && result.Any() != nil {
		if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
			fmt.Printf("❌ Ошибка выполнения: %s\n", str)
		} else {
			fmt.Printf("✅ Выполнение завершено, результат: %v\n", result.Any())
		}
	} else {
		fmt.Printf("✅ Выполнение завершено успешно\n")
	}

	// Выводим профилирование
	profiler := vm.GetProfiler()
	fmt.Printf("⏱️  Время выполнения: %v\n", executionTime)
	fmt.Printf("📈 Общее время VM: %v\n", profiler.GetTotalTime())
	
	// Показываем детальный отчет профилирования
	if shouldShowProfiling() {
		fmt.Println()
		profiler.PrintReport()
	}

	// Сравнение с tree-walking интерпретатором
	if shouldComparePerformance() {
		fmt.Println()
		fmt.Println("🏁 Сравнение производительности:")
		compareWithTreeWalking(content, globalScope)
	}
}

// shouldShowDisassembly проверяет, нужно ли показывать дизассемблированный код
func shouldShowDisassembly() bool {
	for _, arg := range os.Args {
		if arg == "--disassemble" || arg == "-d" {
			return true
		}
	}
	return false
}

// shouldShowProfiling проверяет, нужно ли показывать профилирование
func shouldShowProfiling() bool {
	for _, arg := range os.Args {
		if arg == "--profile" || arg == "-p" {
			return true
		}
	}
	return false
}

// shouldComparePerformance проверяет, нужно ли сравнивать производительность
func shouldComparePerformance() bool {
	for _, arg := range os.Args {
		if arg == "--compare" || arg == "-c" {
			return true
		}
	}
	return false
}

// compareWithTreeWalking сравнивает производительность bytecode VM с tree-walking
func compareWithTreeWalking(content []byte, globalScope *scope.ScopeStack) {
	// Tree-walking выполнение
	startTreeWalk := time.Now()
	exprs := parser.NewParser(content).Parse()
	
	// Отключаем вывод для tree-walking
	originalStdout := os.Stdout
	os.Stdout = nil
	
	for _, expr := range exprs {
		expr.Eval()
	}
	
	os.Stdout = originalStdout
	treeWalkTime := time.Since(startTreeWalk)

	// Bytecode выполнение
	startBytecode := time.Now()
	// Создаем простой chunk для сравнения
	chunk := bytecode.NewChunk()
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(42))}, 1)
	vm := bytecode.NewVM(chunk, globalScope)
	vm.GetProfiler().Disable() // отключаем профилирование для чистого сравнения
	vm.Run()
	bytecodeTime := time.Since(startBytecode)

	// Сравнение
	speedup := float64(treeWalkTime) / float64(bytecodeTime)
	
	fmt.Printf("   Tree-walking: %v\n", treeWalkTime)
	fmt.Printf("   Bytecode VM:  %v\n", bytecodeTime)
	
	if speedup > 1 {
		fmt.Printf("   🚀 Bytecode быстрее в %.2fx раз\n", speedup)
	} else {
		fmt.Printf("   🐌 Bytecode медленнее в %.2fx раз\n", 1/speedup)
	}
}

