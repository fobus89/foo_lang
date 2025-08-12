package test

import (
	"testing"
	"foo_lang/bytecode"
	"foo_lang/scope"
	"foo_lang/builtin"
)

func TestSimpleBytecodeBasic(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeMathFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	// Создаем chunk напрямую (без компилятора)
	chunk := bytecode.NewChunk()
	
	// Эмулируем простую программу: 10 + 5
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(10))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(5))}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 1)
	
	// Выполняем в VM для проверки
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
	
	t.Logf("✅ Simple bytecode test passed: %d", result.Any())
}

func TestSimpleBytecodeVM(t *testing.T) {
	// Создаем простой chunk вручную для тестирования VM
	chunk := bytecode.NewChunk()
	
	// Добавляем: push 15, push 5, add
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(15))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(5))}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 1)
	
	// Создаем VM и выполняем
	vm := bytecode.NewVM(chunk, scope.NewScopeStack())
	result := vm.Run()
	
	// Проверяем результат
	if result == nil {
		t.Error("VM returned nil result")
		return
	}
	
	if resultVal, ok := result.Any().(int64); !ok || resultVal != 20 {
		t.Errorf("Expected 20, got %v (type %T)", result.Any(), result.Any())
	} else {
		t.Logf("✅ VM executed successfully, result: %v", resultVal)
	}
}

func TestSimpleBytecodeArithmetic(t *testing.T) {
	// Тест арифметических операций
	chunk := bytecode.NewChunk()
	
	// Тест: 10 * 3 + 2 = 32
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(10))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(3))}, 1)
	chunk.WriteInstruction(bytecode.OP_MULTIPLY, nil, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(2))}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 1)
	
	vm := bytecode.NewVM(chunk, scope.NewScopeStack())
	result := vm.Run()
	
	if result == nil {
		t.Error("VM returned nil result")
		return
	}
	
	if resultVal, ok := result.Any().(int64); !ok || resultVal != 32 {
		t.Errorf("Expected 32, got %v", result.Any())
	}
	
	t.Logf("✅ Arithmetic test passed: %v", result.Any())
}

func TestSimpleBytecodeLogical(t *testing.T) {
	// Тест логических операций
	chunk := bytecode.NewChunk()
	
	// Тест: true && false = false
	chunk.WriteInstruction(bytecode.OP_TRUE, nil, 1)
	chunk.WriteInstruction(bytecode.OP_FALSE, nil, 1)
	chunk.WriteInstruction(bytecode.OP_AND, nil, 1)
	
	vm := bytecode.NewVM(chunk, scope.NewScopeStack())
	result := vm.Run()
	
	if result == nil {
		t.Error("VM returned nil result")
		return
	}
	
	// Логические операции могут возвращать оригинальные значения, проверяем truthiness
	if !result.IsTruthy() {
		t.Logf("✅ Logical test passed: %v (falsy)", result.Any())
	} else {
		t.Errorf("Expected falsy result, got %v", result.Any())
	}
}

func TestSimpleBytecodeComparison(t *testing.T) {
	// Тест операций сравнения
	chunk := bytecode.NewChunk()
	
	// Тест: 5 > 3 = true
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(5))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(3))}, 1)
	chunk.WriteInstruction(bytecode.OP_GREATER, nil, 1)
	
	vm := bytecode.NewVM(chunk, scope.NewScopeStack())
	result := vm.Run()
	
	if result == nil {
		t.Error("VM returned nil result")
		return
	}
	
	if resultVal, ok := result.Any().(bool); !ok || resultVal != true {
		t.Errorf("Expected true, got %v", result.Any())
	}
	
	t.Logf("✅ Comparison test passed: %v", result.Any())
}

func TestSimpleBytecodeProfiler(t *testing.T) {
	// Тест профайлера с простыми операциями
	chunk := bytecode.NewChunk()
	
	// Добавляем несколько операций для профилирования
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(1))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(2))}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(3))}, 1)
	chunk.WriteInstruction(bytecode.OP_MULTIPLY, nil, 1)
	
	vm := bytecode.NewVM(chunk, scope.NewScopeStack())
	result := vm.Run()
	
	// Проверяем профайлер
	profiler := vm.GetProfiler()
	if profiler.GetInstructionCount(bytecode.OP_CONSTANT) != 3 {
		t.Errorf("Expected 3 OP_CONSTANT instructions, got %d", 
			profiler.GetInstructionCount(bytecode.OP_CONSTANT))
	}
	
	if profiler.GetInstructionCount(bytecode.OP_ADD) != 1 {
		t.Errorf("Expected 1 OP_ADD instruction, got %d", 
			profiler.GetInstructionCount(bytecode.OP_ADD))
	}
	
	if profiler.GetInstructionCount(bytecode.OP_MULTIPLY) != 1 {
		t.Errorf("Expected 1 OP_MULTIPLY instruction, got %d", 
			profiler.GetInstructionCount(bytecode.OP_MULTIPLY))
	}
	
	// Проверяем результат: (1 + 2) * 3 = 9
	if resultVal, ok := result.Any().(int64); !ok || resultVal != 9 {
		t.Errorf("Expected 9, got %v", result.Any())
	}
	
	t.Logf("✅ Profiler test passed, total time: %v", profiler.GetTotalTime())
}