package test

import (
	"fmt"
	"foo_lang/bytecode"
	"foo_lang/scope"
	"foo_lang/value"
	"testing"
	"time"
)

func TestVMBasicOperations(t *testing.T) {
	// Создаем простой bytecode chunk
	chunk := bytecode.NewChunk()
	
	// Добавляем константы
	constIndex1 := chunk.AddConstant(int64(10))
	constIndex2 := chunk.AddConstant(int64(20))
	
	// Добавляем инструкции: 10 + 20
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{constIndex1}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{constIndex2}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, []int{}, 1)
	
	// Создаем VM и выполняем
	scopeStack := scope.NewScopeStack()
	vm := bytecode.NewVM(chunk, scopeStack)
	
	result := vm.Run()
	
	if result == nil {
		t.Error("Expected result, got nil")
		return
	}
	
	// Результат должен быть 30 (10 + 20) 
	// Проверяем тип результата и приводим к числу
	resultAny := result.Any()
	var resultNum float64
	
	switch resultAny.(type) {
	case int64:
		resultNum = float64(resultAny.(int64))
	case float64:
		resultNum = resultAny.(float64)
	default:
		t.Errorf("Expected numeric result, got %T", resultAny)
		return
	}
	
	if resultNum != 30.0 {
		t.Errorf("Expected 30, got %v", resultNum)
	}
}

func TestJITCompiler(t *testing.T) {
	scopeStack := scope.NewScopeStack()
	chunk := bytecode.NewChunk()
	vm := bytecode.NewVM(chunk, scopeStack)
	
	jit := vm.GetJIT()
	
	// Тестируем, что JIT включен по умолчанию
	if jit == nil {
		t.Error("JIT compiler should be initialized")
		return
	}
	
	// Симулируем множественные вызовы функции
	for i := 0; i < 150; i++ {
		jit.RecordExecution("testFunction", time.Microsecond)
	}
	
	// Проверяем, что функция стала горячей
	stats := jit.GetStats()
	if stat, exists := stats["testFunction"]; exists {
		if !stat.IsHot {
			t.Error("Function should be hot after 150 calls")
		}
	} else {
		t.Error("Function statistics not found")
	}
	
	// Проверяем, что функция скомпилирована
	if !jit.IsCompiled("testFunction") {
		t.Error("Function should be compiled after becoming hot")
	}
}

func TestDebugger(t *testing.T) {
	scopeStack := scope.NewScopeStack()
	chunk := bytecode.NewChunk()
	vm := bytecode.NewVM(chunk, scopeStack)
	
	// Проверяем, что debugger выключен по умолчанию
	if vm.IsDebugMode() {
		t.Error("Debug mode should be disabled by default")
	}
	
	// Включаем debug режим
	vm.EnableDebugMode()
	if !vm.IsDebugMode() {
		t.Error("Debug mode should be enabled")
	}
	
	// Устанавливаем breakpoint
	vm.SetBreakpoint(5)
	if !vm.IsBreakpoint(5) {
		t.Error("Breakpoint should be set on line 5")
	}
	
	// Удаляем breakpoint
	vm.RemoveBreakpoint(5)
	if vm.IsBreakpoint(5) {
		t.Error("Breakpoint should be removed from line 5")
	}
}

func TestProfiler(t *testing.T) {
	scopeStack := scope.NewScopeStack()
	chunk := bytecode.NewChunk()
	vm := bytecode.NewVM(chunk, scopeStack)
	
	profiler := vm.GetProfiler()
	if profiler == nil {
		t.Error("Profiler should be initialized")
		return
	}
	
	// Записываем выполнение инструкций
	profiler.RecordInstruction(bytecode.OP_ADD)
	profiler.RecordInstruction(bytecode.OP_ADD)
	profiler.RecordInstruction(bytecode.OP_MULTIPLY)
	
	// Проверяем счетчики
	if profiler.GetInstructionCount(bytecode.OP_ADD) != 2 {
		t.Errorf("Expected 2 ADD instructions, got %d", profiler.GetInstructionCount(bytecode.OP_ADD))
	}
	
	if profiler.GetInstructionCount(bytecode.OP_MULTIPLY) != 1 {
		t.Errorf("Expected 1 MULTIPLY instruction, got %d", profiler.GetInstructionCount(bytecode.OP_MULTIPLY))
	}
}

func TestStringOperations(t *testing.T) {
	chunk := bytecode.NewChunk()
	scopeStack := scope.NewScopeStack()
	
	// Тестируем OP_STRING_LEN
	strIndex := chunk.AddConstant("Hello")
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{strIndex}, 1)
	chunk.WriteInstruction(bytecode.OP_STRING_LEN, []int{}, 1)
	
	vm := bytecode.NewVM(chunk, scopeStack)
	result := vm.Run()
	
	if result == nil {
		t.Error("Expected result, got nil")
		return
	}
	
	// Длина "Hello" должна быть 5
	if result.Any().(int64) != 5 {
		t.Errorf("Expected length 5, got %v", result.Any())
	}
}

func TestMathOperations(t *testing.T) {
	chunk := bytecode.NewChunk()
	scopeStack := scope.NewScopeStack()
	
	// Тестируем OP_MATH_ABS с отрицательным числом
	numIndex := chunk.AddConstant(int64(-42))
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{numIndex}, 1)
	chunk.WriteInstruction(bytecode.OP_MATH_ABS, []int{}, 1)
	
	vm := bytecode.NewVM(chunk, scopeStack)
	result := vm.Run()
	
	if result == nil {
		t.Error("Expected result, got nil")
		return
	}
	
	// abs(-42) должен быть 42
	resultAny := result.Any()
	var resultNum float64
	
	switch resultAny.(type) {
	case int64:
		resultNum = float64(resultAny.(int64))
	case float64:
		resultNum = resultAny.(float64)
	}
	
	if resultNum != 42.0 {
		t.Errorf("Expected 42, got %v", resultNum)
	}
}

func TestArrayOperations(t *testing.T) {
	chunk := bytecode.NewChunk()
	scopeStack := scope.NewScopeStack()
	
	// Создаем массив из 3 элементов: [1, 2, 3]
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(1))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(2))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(3))}, 1)
	chunk.WriteInstruction(bytecode.OP_ARRAY, []int{3}, 1)
	
	vm := bytecode.NewVM(chunk, scopeStack)
	result := vm.Run()
	
	if result == nil {
		t.Error("Expected result, got nil")
		return
	}
	
	// Проверяем, что результат - массив
	if arr, ok := result.Any().([]*value.Value); ok {
		if len(arr) != 3 {
			t.Errorf("Expected array length 3, got %d", len(arr))
		}
	} else {
		t.Error("Expected array result")
	}
}

func TestJITOptimizations(t *testing.T) {
	scopeStack := scope.NewScopeStack()
	chunk := bytecode.NewChunk()
	vm := bytecode.NewVM(chunk, scopeStack)
	
	jit := vm.GetJIT()
	
	// Тестируем, что JIT может выполнить оптимизированную функцию
	// (пока это заглушка, но структура готова)
	result := jit.ExecuteOptimized(vm, "nonExistentFunction")
	if result != nil {
		t.Error("Should return nil for non-compiled function")
	}
	
	// Компилируем функцию вручную для теста
	jit.CompileFunction("testOptimizedFunction")
	
	if !jit.IsCompiled("testOptimizedFunction") {
		t.Error("Function should be compiled")
	}
}

// Benchmarking test для измерения производительности JIT
func BenchmarkVMExecution(b *testing.B) {
	chunk := bytecode.NewChunk()
	scopeStack := scope.NewScopeStack()
	
	// Простые арифметические операции
	for i := 0; i < 100; i++ {
		chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(i))}, 1)
		chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(2))}, 1)
		chunk.WriteInstruction(bytecode.OP_ADD, []int{}, 1)
		chunk.WriteInstruction(bytecode.OP_POP, []int{}, 1) // Убираем результат со стека
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		vm := bytecode.NewVM(chunk, scopeStack)
		vm.Run()
	}
}

// Интеграционный тест всей системы
func TestIntegrationVMEnhanced(t *testing.T) {
	fmt.Println("=== Integration Test: Enhanced VM ===")
	
	scopeStack := scope.NewScopeStack()
	chunk := bytecode.NewChunk()
	
	// Создаем сложную программу
	// let x = 10 + 20; print(x)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(10))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(20))}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, []int{}, 1)
	chunk.WriteInstruction(bytecode.OP_DEFINE_GLOBAL, []int{chunk.AddConstant("x")}, 1)
	
	chunk.WriteInstruction(bytecode.OP_GET_GLOBAL, []int{chunk.AddConstant("x")}, 2)
	chunk.WriteInstruction(bytecode.OP_PRINT, []int{}, 2)
	
	vm := bytecode.NewVM(chunk, scopeStack)
	
	// Включаем все улучшения
	vm.EnableDebugMode()
	vm.EnableJIT()
	vm.SetBreakpoint(2) // Breakpoint на печати
	
	fmt.Println("Executing enhanced VM with all features enabled...")
	
	// Выполняем (в реальном тесте breakpoint будет требовать ввода)
	vm.DisableDebugMode() // Отключаем для автоматического выполнения
	result := vm.Run()
	
	// Выводим отчеты
	fmt.Println("\n--- Profiler Report ---")
	vm.GetProfiler().PrintReport()
	
	fmt.Println("\n--- JIT Report ---")
	vm.GetJIT().PrintReport()
	
	if result == nil {
		t.Error("Expected successful execution")
	}
	
	fmt.Println("=== Integration Test Completed ===")
}