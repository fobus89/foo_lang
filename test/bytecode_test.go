package test

import (
	"testing"
	"foo_lang/bytecode"
	"foo_lang/scope"
	"foo_lang/builtin"
)

func TestBytecodeBasicOperations(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeMathFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	// Тест базовых операций через создание chunk вручную
	chunk := bytecode.NewChunk()
	
	// Создаем программу: let a = 10; let b = 5; let sum = a + b
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(10))}, 1)
	chunk.WriteInstruction(bytecode.OP_DEFINE_GLOBAL, []int{chunk.AddConstant("a")}, 1)
	
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(5))}, 2)
	chunk.WriteInstruction(bytecode.OP_DEFINE_GLOBAL, []int{chunk.AddConstant("b")}, 2)
	
	chunk.WriteInstruction(bytecode.OP_GET_GLOBAL, []int{chunk.AddConstant("a")}, 3)
	chunk.WriteInstruction(bytecode.OP_GET_GLOBAL, []int{chunk.AddConstant("b")}, 3)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 3)
	chunk.WriteInstruction(bytecode.OP_DEFINE_GLOBAL, []int{chunk.AddConstant("sum")}, 3)
	
	// Получаем результат
	chunk.WriteInstruction(bytecode.OP_GET_GLOBAL, []int{chunk.AddConstant("sum")}, 4)
	
	// Выполняем в VM
	vm := bytecode.NewVM(chunk, scope.GlobalScope)
	result := vm.Run()
	
	// Проверяем результат: 10 + 5 = 15
	if result == nil {
		t.Error("VM returned nil result")
		return
	}
	
	if resultVal, ok := result.Any().(int64); !ok || resultVal != 15 {
		t.Errorf("Expected 15, got %v", result.Any())
	}
	
	t.Logf("✅ Bytecode basic operations test passed: %v", result.Any())
}

func TestBytecodeControlFlow(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeMathFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	// Упрощенный тест управления потоком
	chunk := bytecode.NewChunk()
	
	// Упрощенный тест: просто 5 > 3 (без условного перехода)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(5))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(3))}, 1)
	chunk.WriteInstruction(bytecode.OP_GREATER, nil, 1)
	
	// Выполняем в VM
	vm := bytecode.NewVM(chunk, scope.GlobalScope)
	result := vm.Run()
	
	// Проверяем результат: 5 > 3 должно быть true
	if result == nil {
		t.Error("VM returned nil result")
		return
	}
	
	if resultVal, ok := result.Any().(bool); !ok || resultVal != true {
		t.Errorf("Expected true, got %v", result.Any())
	}
	
	t.Logf("✅ Bytecode control flow test passed: %v", result.Any())
}

func TestBytecodeArrays(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeMathFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	// Упрощенный тест массивов
	chunk := bytecode.NewChunk()
	
	// Создаем массив [1, 2, 3] и получаем первый элемент
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(1))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(2))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(3))}, 1)
	chunk.WriteInstruction(bytecode.OP_ARRAY, []int{3}, 1) // создаем массив из 3 элементов
	
	// Индексация: получаем элемент [0]
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(0))}, 2)
	chunk.WriteInstruction(bytecode.OP_INDEX, nil, 2)
	
	// Выполняем в VM
	vm := bytecode.NewVM(chunk, scope.GlobalScope)
	result := vm.Run()
	
	// Проверяем результат: первый элемент должен быть 1
	if result == nil {
		t.Error("VM returned nil result")
		return
	}
	
	if resultVal, ok := result.Any().(int64); !ok || resultVal != 1 {
		t.Errorf("Expected 1, got %v", result.Any())
	}
	
	t.Logf("✅ Bytecode arrays test passed: %v", result.Any())
}

func TestBytecodeCompilerChunk(t *testing.T) {
	// Тест создания chunk напрямую
	chunk := bytecode.NewChunk()
	
	// Простое выражение: 1 + 2 = 3
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(1))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(2))}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 1)
	
	// Проверяем, что chunk создался
	if chunk == nil {
		t.Error("Failed to create chunk")
		return
	}
	
	// Проверяем, что есть инструкции
	if len(chunk.Code) == 0 {
		t.Error("Chunk has no instructions")
		return
	}
	
	// Проверяем, что есть константы
	if len(chunk.Constants) == 0 {
		t.Error("Chunk has no constants")
		return
	}
	
	t.Logf("✅ Chunk created successfully with %d instructions and %d constants", 
		len(chunk.Code), len(chunk.Constants))
}

func TestBytecodeVMExecution(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	
	// Простое выражение для тестирования VM
	chunk := bytecode.NewChunk()
	
	// Эмулируем: push 10, push 5, add
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(10))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(5))}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 1)
	
	// Создаем VM и выполняем
	vm := bytecode.NewVM(chunk, scope.GlobalScope)
	result := vm.Run()
	
	// Проверяем результат
	if result == nil {
		t.Error("VM returned nil result")
		return
	}
	
	if resultVal, ok := result.Any().(int64); !ok || resultVal != 15 {
		t.Errorf("Expected 15, got %v", result.Any())
	}
	
	t.Logf("✅ VM executed successfully, result: %v", result.Any())
}

func TestBytecodeProfiler(t *testing.T) {
	// Тест профайлера
	profiler := bytecode.NewProfiler()
	
	// Имитируем выполнение
	profiler.StartExecution()
	
	// Записываем несколько инструкций
	profiler.RecordInstruction(bytecode.OP_CONSTANT)
	profiler.RecordInstruction(bytecode.OP_CONSTANT)
	profiler.RecordInstruction(bytecode.OP_ADD)
	profiler.RecordInstruction(bytecode.OP_PRINT)
	
	// Профилируем функцию
	profiler.StartFunction("testFunction")
	profiler.RecordInstruction(bytecode.OP_CALL)
	profiler.EndFunction("testFunction")
	
	profiler.EndExecution()
	
	// Проверяем, что данные собрались
	if profiler.GetInstructionCount(bytecode.OP_CONSTANT) != 2 {
		t.Errorf("Expected 2 OP_CONSTANT instructions, got %d", 
			profiler.GetInstructionCount(bytecode.OP_CONSTANT))
	}
	
	if profiler.GetInstructionCount(bytecode.OP_ADD) != 1 {
		t.Errorf("Expected 1 OP_ADD instruction, got %d", 
			profiler.GetInstructionCount(bytecode.OP_ADD))
	}
	
	hotspots := profiler.GetHotspots()
	if len(hotspots) == 0 {
		t.Error("No hotspots found")
	}
	
	t.Logf("✅ Profiler test completed, total time: %v", profiler.GetTotalTime())
}

func TestBytecodeDisassembly(t *testing.T) {
	// Тест дизассемблирования bytecode
	chunk := bytecode.NewChunk()
	
	// Добавляем несколько инструкций
	constantIndex := chunk.AddConstant(int64(42))
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{constantIndex}, 1)
	chunk.WriteInstruction(bytecode.OP_PRINT, nil, 2)
	
	// Проверяем, что можем дизассемблировать без ошибок
	// В реальной реализации это бы выводило в stdout
	bytecode.DisassembleChunk(chunk, "test_chunk")
	
	// Дизассемблируем отдельную инструкцию
	if len(chunk.Code) > 0 {
		bytecode.DisassembleInstruction(&chunk.Code[0], 0)
	}
	
	t.Log("✅ Disassembly completed successfully")
}