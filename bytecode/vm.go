package bytecode

import (
	"fmt"
	"foo_lang/value"
	"foo_lang/scope"
)

// VM представляет виртуальную машину для выполнения bytecode
type VM struct {
	chunk      *Chunk
	ip         int                // instruction pointer
	stack      []*value.Value     // стек операций
	sp         int                // stack pointer
	globals    map[string]*value.Value
	scope      *scope.ScopeStack
	callFrames []CallFrame        // стек вызовов функций
	profiler   *Profiler         // профайлер производительности
}

// CallFrame представляет кадр вызова функции
type CallFrame struct {
	function   *value.Value
	ip         int
	basePtr    int  // указатель на начало локальных переменных в стеке
}

// NewVM создает новую виртуальную машину
func NewVM(chunk *Chunk, scopeStack *scope.ScopeStack) *VM {
	return &VM{
		chunk:      chunk,
		ip:         0,
		stack:      make([]*value.Value, 0, 256),
		sp:         0,
		globals:    make(map[string]*value.Value),
		scope:      scopeStack,
		callFrames: make([]CallFrame, 0, 64),
		profiler:   NewProfiler(),
	}
}

// Reset сбрасывает состояние VM для нового выполнения
func (vm *VM) Reset() {
	vm.ip = 0
	vm.stack = vm.stack[:0]
	vm.sp = 0
	vm.callFrames = vm.callFrames[:0]
}

// Push добавляет значение в стек
func (vm *VM) Push(val *value.Value) {
	vm.stack = append(vm.stack, val)
	vm.sp++
}

// Pop извлекает значение из стека
func (vm *VM) Pop() *value.Value {
	if vm.sp == 0 {
		return value.NewString("Error: stack underflow")
	}
	vm.sp--
	val := vm.stack[vm.sp]
	vm.stack = vm.stack[:vm.sp]
	return val
}

// Peek возвращает верхнее значение стека без извлечения
func (vm *VM) Peek(distance int) *value.Value {
	index := vm.sp - 1 - distance
	if index < 0 || index >= vm.sp {
		return value.NewString("Error: stack index out of bounds")
	}
	return vm.stack[index]
}

// Run выполняет bytecode
func (vm *VM) Run() *value.Value {
	vm.profiler.StartExecution()
	defer vm.profiler.EndExecution()

	for vm.ip < len(vm.chunk.Code) {
		instruction := &vm.chunk.Code[vm.ip]
		
		// Профилирование инструкций
		vm.profiler.RecordInstruction(instruction.OpCode)
		
		result := vm.executeInstruction(instruction)
		if result != nil && result.IsReturn() {
			return result
		}
		
		vm.ip++
	}
	
	// Возвращаем последнее значение из стека или nil
	if vm.sp > 0 {
		return vm.Pop()
	}
	return value.NewNil()
}

// executeInstruction выполняет одну инструкцию
func (vm *VM) executeInstruction(instruction *Instruction) *value.Value {
	switch instruction.OpCode {
	
	// Константы и литералы
	case OP_CONSTANT:
		constantIndex := instruction.Operands[0]
		if constantIndex >= len(vm.chunk.Constants) {
			return value.NewString("Error: constant index out of bounds")
		}
		constant := vm.chunk.Constants[constantIndex]
		vm.Push(value.FromInterface(constant))
		
	case OP_NIL:
		vm.Push(value.NewNil())
		
	case OP_TRUE:
		vm.Push(value.NewBool(true))
		
	case OP_FALSE:
		vm.Push(value.NewBool(false))
	
	// Арифметические операции
	case OP_ADD:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.Add(a, b)
		})
		
	case OP_SUBTRACT:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.Subtract(a, b)
		})
		
	case OP_MULTIPLY:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.Multiply(a, b)
		})
		
	case OP_DIVIDE:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.Divide(a, b)
		})
		
	case OP_MODULO:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.Modulo(a, b)
		})
		
	case OP_NEGATE:
		operand := vm.Pop()
		vm.Push(value.Negate(operand))
	
	// Логические операции
	case OP_NOT:
		operand := vm.Pop()
		vm.Push(value.Not(operand))
		
	case OP_AND:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.And(a, b)
		})
		
	case OP_OR:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.Or(a, b)
		})
	
	// Операции сравнения
	case OP_EQUAL:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.Equal(a, b)
		})
		
	case OP_NOT_EQUAL:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.NotEqual(a, b)
		})
		
	case OP_GREATER:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.Greater(a, b)
		})
		
	case OP_GREATER_EQUAL:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.GreaterEqual(a, b)
		})
		
	case OP_LESS:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.Less(a, b)
		})
		
	case OP_LESS_EQUAL:
		return vm.binaryOperation(func(a, b *value.Value) *value.Value {
			return value.LessEqual(a, b)
		})
	
	// Переменные
	case OP_GET_GLOBAL:
		nameIndex := instruction.Operands[0]
		name := vm.chunk.Constants[nameIndex].(string)
		if val, exists := vm.globals[name]; exists {
			vm.Push(val)
		} else {
			// Пытаемся найти в scope
			if scopeVal, exists := vm.scope.Get(name); exists {
				vm.Push(scopeVal)
			} else {
				return value.NewString(fmt.Sprintf("Error: undefined variable '%s'", name))
			}
		}
		
	case OP_SET_GLOBAL:
		nameIndex := instruction.Operands[0]
		name := vm.chunk.Constants[nameIndex].(string)
		val := vm.Peek(0) // не извлекаем из стека
		vm.globals[name] = val
		vm.scope.Set(name, val)
		
	case OP_DEFINE_GLOBAL:
		nameIndex := instruction.Operands[0]
		name := vm.chunk.Constants[nameIndex].(string)
		val := vm.Pop()
		vm.globals[name] = val
		vm.scope.Set(name, val)
	
	// Управление потоком
	case OP_JUMP:
		offset := instruction.Operands[0]
		vm.ip += offset - 1 // -1 потому что ip++ в основном цикле
		
	case OP_JUMP_IF_FALSE:
		offset := instruction.Operands[0]
		condition := vm.Pop()
		if !condition.IsTruthy() {
			vm.ip += offset - 1
		}
		
	case OP_LOOP:
		offset := instruction.Operands[0]
		vm.ip -= offset + 1
	
	// Стек операции
	case OP_POP:
		vm.Pop()
		
	case OP_DUP:
		val := vm.Peek(0)
		vm.Push(val)
	
	// Массивы
	case OP_ARRAY:
		size := instruction.Operands[0]
		elements := make([]*value.Value, size)
		for i := size - 1; i >= 0; i-- {
			elements[i] = vm.Pop()
		}
		vm.Push(value.NewArray(elements))
		
	case OP_INDEX:
		index := vm.Pop()
		obj := vm.Pop()
		result := value.Index(obj, index)
		vm.Push(result)
	
	// Встроенные функции
	case OP_PRINT:
		val := vm.Pop()
		fmt.Print(val.String())
		vm.Push(value.NewNil())
		
	case OP_PRINTLN:
		val := vm.Pop()
		fmt.Println(val.String())
		vm.Push(value.NewNil())
	
	// Профилинг
	case OP_PROFILE_START:
		nameIndex := instruction.Operands[0]
		name := vm.chunk.Constants[nameIndex].(string)
		vm.profiler.StartFunction(name)
		
	case OP_PROFILE_END:
		nameIndex := instruction.Operands[0]
		name := vm.chunk.Constants[nameIndex].(string)
		vm.profiler.EndFunction(name)
	
	// Отладка
	case OP_DEBUG_TRACE:
		fmt.Printf("DEBUG: IP=%d, SP=%d, Stack: ", vm.ip, vm.sp)
		for i := 0; i < vm.sp; i++ {
			fmt.Printf("[%s] ", vm.stack[i].String())
		}
		fmt.Println()
	
	default:
		return value.NewString(fmt.Sprintf("Error: unknown opcode %d", instruction.OpCode))
	}
	
	return nil
}

// binaryOperation выполняет бинарную операцию с двумя операндами из стека
func (vm *VM) binaryOperation(op func(a, b *value.Value) *value.Value) *value.Value {
	b := vm.Pop()
	a := vm.Pop()
	result := op(a, b)
	vm.Push(result)
	return nil
}

// GetProfiler возвращает профайлер VM
func (vm *VM) GetProfiler() *Profiler {
	return vm.profiler
}

// PrintStack выводит содержимое стека для отладки
func (vm *VM) PrintStack() {
	fmt.Printf("Stack (SP=%d): ", vm.sp)
	for i := 0; i < vm.sp; i++ {
		fmt.Printf("[%s] ", vm.stack[i].String())
	}
	fmt.Println()
}