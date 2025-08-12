package test

import (
	"testing"
	"foo_lang/bytecode"
	"foo_lang/scope"
)

func TestMinimalBytecodeVM(t *testing.T) {
	// Создаем chunk напрямую без компилятора
	chunk := bytecode.NewChunk()
	
	// Простой тест: 10 + 5 = 15
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(10))}, 1)
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
	
	if resultVal, ok := result.Any().(int64); !ok || resultVal != 15 {
		t.Errorf("Expected 15, got %v (type %T)", result.Any(), result.Any())
	} else {
		t.Logf("✅ VM test passed: %d", resultVal)
	}
}

func TestMinimalBytecodeArithmetic(t *testing.T) {
	// Тест различных арифметических операций
	tests := []struct {
		name     string
		a, b     int64
		op       bytecode.OpCode
		expected int64
	}{
		{"Addition", 7, 3, bytecode.OP_ADD, 10},
		{"Subtraction", 10, 4, bytecode.OP_SUBTRACT, 6},
		{"Multiplication", 6, 7, bytecode.OP_MULTIPLY, 42},
		{"Division", 20, 4, bytecode.OP_DIVIDE, 5},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chunk := bytecode.NewChunk()
			
			chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(test.a)}, 1)
			chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(test.b)}, 1)
			chunk.WriteInstruction(test.op, nil, 1)
			
			vm := bytecode.NewVM(chunk, scope.NewScopeStack())
			result := vm.Run()
			
			if result == nil {
				t.Errorf("%s: VM returned nil", test.name)
				return
			}
			
			// Для деления может быть float64
			if test.op == bytecode.OP_DIVIDE {
				if resultVal, ok := result.Any().(float64); !ok || int64(resultVal) != test.expected {
					t.Errorf("%s: expected %d, got %v", test.name, test.expected, result.Any())
				}
			} else {
				if resultVal, ok := result.Any().(int64); !ok || resultVal != test.expected {
					t.Errorf("%s: expected %d, got %v", test.name, test.expected, result.Any())
				}
			}
			
			t.Logf("✅ %s: %d %s %d = %v", test.name, test.a, 
				map[bytecode.OpCode]string{
					bytecode.OP_ADD: "+", bytecode.OP_SUBTRACT: "-",
					bytecode.OP_MULTIPLY: "*", bytecode.OP_DIVIDE: "/",
				}[test.op], test.b, result.Any())
		})
	}
}

func TestMinimalBytecodeComparison(t *testing.T) {
	// Тест операций сравнения
	tests := []struct {
		name     string
		a, b     int64
		op       bytecode.OpCode
		expected bool
	}{
		{"Greater", 5, 3, bytecode.OP_GREATER, true},
		{"Greater False", 2, 5, bytecode.OP_GREATER, false},
		{"Less", 2, 8, bytecode.OP_LESS, true},
		{"Less False", 8, 2, bytecode.OP_LESS, false},
		{"Equal", 7, 7, bytecode.OP_EQUAL, true},
		{"Equal False", 7, 8, bytecode.OP_EQUAL, false},
		{"Not Equal", 5, 3, bytecode.OP_NOT_EQUAL, true},
		{"Not Equal False", 5, 5, bytecode.OP_NOT_EQUAL, false},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chunk := bytecode.NewChunk()
			
			chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(test.a)}, 1)
			chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(test.b)}, 1)
			chunk.WriteInstruction(test.op, nil, 1)
			
			vm := bytecode.NewVM(chunk, scope.NewScopeStack())
			result := vm.Run()
			
			if result == nil {
				t.Errorf("%s: VM returned nil", test.name)
				return
			}
			
			if resultVal, ok := result.Any().(bool); !ok || resultVal != test.expected {
				t.Errorf("%s: expected %v, got %v", test.name, test.expected, result.Any())
			} else {
				t.Logf("✅ %s: %d %s %d = %v", test.name, test.a,
					map[bytecode.OpCode]string{
						bytecode.OP_GREATER: ">", bytecode.OP_LESS: "<",
						bytecode.OP_EQUAL: "==", bytecode.OP_NOT_EQUAL: "!=",
					}[test.op], test.b, resultVal)
			}
		})
	}
}

func TestMinimalBytecodeLogical(t *testing.T) {
	// Тест логических операций
	tests := []struct {
		name     string
		a, b     bool
		op       bytecode.OpCode
		expected bool
	}{
		{"AND True", true, true, bytecode.OP_AND, true},
		{"AND False", true, false, bytecode.OP_AND, false},
		{"OR True", true, false, bytecode.OP_OR, true},
		{"OR False", false, false, bytecode.OP_OR, false},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chunk := bytecode.NewChunk()
			
			// Добавляем boolean значения
			if test.a {
				chunk.WriteInstruction(bytecode.OP_TRUE, nil, 1)
			} else {
				chunk.WriteInstruction(bytecode.OP_FALSE, nil, 1)
			}
			
			if test.b {
				chunk.WriteInstruction(bytecode.OP_TRUE, nil, 1)
			} else {
				chunk.WriteInstruction(bytecode.OP_FALSE, nil, 1)
			}
			
			chunk.WriteInstruction(test.op, nil, 1)
			
			vm := bytecode.NewVM(chunk, scope.NewScopeStack())
			result := vm.Run()
			
			if result == nil {
				t.Errorf("%s: VM returned nil", test.name)
				return
			}
			
			// Логические операции могут возвращать оригинальные значения
			resultBool := false
			if val, ok := result.Any().(bool); ok {
				resultBool = val
			} else {
				// Проверяем truthiness для других типов
				resultBool = result.IsTruthy()
			}
			
			if resultBool != test.expected {
				t.Errorf("%s: expected %v, got %v (result: %v)", test.name, test.expected, resultBool, result.Any())
			} else {
				t.Logf("✅ %s: %v %s %v = %v", test.name, test.a,
					map[bytecode.OpCode]string{
						bytecode.OP_AND: "&&", bytecode.OP_OR: "||",
					}[test.op], test.b, resultBool)
			}
		})
	}
}

func TestMinimalBytecodeProfiler(t *testing.T) {
	// Тест профайлера
	chunk := bytecode.NewChunk()
	
	// Создаем последовательность операций для профилирования
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(1))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(2))}, 1)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(3))}, 1)
	chunk.WriteInstruction(bytecode.OP_MULTIPLY, nil, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(1))}, 1)
	chunk.WriteInstruction(bytecode.OP_SUBTRACT, nil, 1)
	
	vm := bytecode.NewVM(chunk, scope.NewScopeStack())
	result := vm.Run()
	
	// Проверяем результат: ((1 + 2) * 3) - 1 = 8
	if resultVal, ok := result.Any().(int64); !ok || resultVal != 8 {
		t.Errorf("Expected 8, got %v", result.Any())
	}
	
	// Проверяем профайлер
	profiler := vm.GetProfiler()
	
	expectedCounts := map[bytecode.OpCode]int64{
		bytecode.OP_CONSTANT: 4,
		bytecode.OP_ADD:      1,
		bytecode.OP_MULTIPLY: 1,
		bytecode.OP_SUBTRACT: 1,
	}
	
	for opcode, expectedCount := range expectedCounts {
		actualCount := profiler.GetInstructionCount(opcode)
		if actualCount != expectedCount {
			t.Errorf("Expected %d %s instructions, got %d", 
				expectedCount, 
				map[bytecode.OpCode]string{
					bytecode.OP_CONSTANT: "OP_CONSTANT",
					bytecode.OP_ADD: "OP_ADD", 
					bytecode.OP_MULTIPLY: "OP_MULTIPLY",
					bytecode.OP_SUBTRACT: "OP_SUBTRACT",
				}[opcode], 
				actualCount)
		}
	}
	
	t.Logf("✅ Profiler test passed, result: %v, total time: %v", 
		result.Any(), profiler.GetTotalTime())
}

func TestMinimalBytecodeDisassembly(t *testing.T) {
	// Тест дизассемблирования
	chunk := bytecode.NewChunk()
	
	// Простая программа
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant(int64(42))}, 1)
	chunk.WriteInstruction(bytecode.OP_CONSTANT, []int{chunk.AddConstant("Hello")}, 2)
	chunk.WriteInstruction(bytecode.OP_ADD, nil, 3)
	
	// Проверяем, что можем дизассемблировать без паники
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Disassembly caused panic: %v", r)
		}
	}()
	
	// Дизассемблируем chunk (вывод идет в stdout)
	bytecode.DisassembleChunk(chunk, "test_program")
	
	// Дизассемблируем отдельные инструкции
	for i, instruction := range chunk.Code {
		bytecode.DisassembleInstruction(&instruction, i)
	}
	
	t.Log("✅ Disassembly completed without errors")
}