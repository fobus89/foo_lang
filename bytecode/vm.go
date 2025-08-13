package bytecode

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
)

// VM представляет виртуальную машину для выполнения bytecode
type VM struct {
	chunk       *Chunk
	ip          int                 // instruction pointer
	stack       []*value.Value      // стек операций
	sp          int                 // stack pointer
	globals     map[string]*value.Value
	scope       *scope.ScopeStack
	callFrames  []CallFrame         // стек вызовов функций
	profiler    *Profiler           // профайлер производительности
	breakpoints map[int]bool        // точки останова для debugger'а
	debugMode   bool               // режим отладки
	jit         *JITCompiler       // JIT компилятор
}

// CallFrame представляет кадр вызова функции
type CallFrame struct {
	function *value.Value
	ip       int
	basePtr  int // указатель на начало локальных переменных в стеке
}

// NewVM создает новую виртуальную машину
func NewVM(chunk *Chunk, scopeStack *scope.ScopeStack) *VM {
	return &VM{
		chunk:       chunk,
		ip:          0,
		stack:       make([]*value.Value, 0, 256),
		sp:          0,
		globals:     make(map[string]*value.Value),
		scope:       scopeStack,
		callFrames:  make([]CallFrame, 0, 64),
		profiler:    NewProfiler(),
		breakpoints: make(map[int]bool),
		debugMode:   false,
		jit:         NewJITCompiler(),
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

		// 🔥 Проверка breakpoint'ов для debugger'а
		if vm.debugMode && vm.isBreakpoint(instruction.Line) {
			fmt.Printf("\n🔴 BREAKPOINT: Line %d, IP=%d\n", instruction.Line, vm.ip)
			vm.PrintStack()
			vm.printLocalVariables()
			
			// В реальном debugger'е здесь будет интерактивная сессия
			fmt.Printf("Press Enter to continue...")
			fmt.Scanln() // Простое ожидание ввода
		}

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

	// Локальные переменные
	case OP_GET_LOCAL:
		slotIndex := instruction.Operands[0]
		if vm.sp <= slotIndex {
			return value.NewString("Error: local variable index out of bounds")
		}
		localVal := vm.stack[slotIndex]
		vm.Push(localVal)

	case OP_SET_LOCAL:
		slotIndex := instruction.Operands[0]
		val := vm.Peek(0) // не извлекаем из стека
		if vm.sp <= slotIndex {
			return value.NewString("Error: local variable index out of bounds")
		}
		vm.stack[slotIndex] = val

	// Объекты
	case OP_OBJECT:
		size := instruction.Operands[0]
		obj := make(map[string]*value.Value)
		
		// Извлекаем пары ключ-значение из стека
		for i := 0; i < size; i++ {
			val := vm.Pop()
			key := vm.Pop()
			obj[key.String()] = val
		}
		vm.Push(value.NewValue(obj))

	case OP_SET_INDEX:
		val := vm.Pop()
		index := vm.Pop()  
		obj := vm.Pop()
		// Простая заглушка для установки индекса
		fmt.Printf("Setting index %s of %s to %s\n", index.String(), obj.String(), val.String())
		vm.Push(value.NewNil())

	// Вызовы функций
	case OP_CALL:
		argCount := instruction.Operands[0]
		
		// Извлекаем аргументы из стека
		args := make([]*value.Value, argCount)
		for i := argCount - 1; i >= 0; i-- {
			args[i] = vm.Pop()
		}
		
		// Извлекаем функцию
		function := vm.Pop()
		
		// Вызываем функцию (пока заглушка)
		result := vm.callFunction(function, args)
		vm.Push(result)

	case OP_RETURN:
		// Возвращаем значение из функции
		result := vm.Pop()
		result.SetReturn(true)
		return result

	// Замыкания (заглушка)
	case OP_CLOSURE:
		functionIndex := instruction.Operands[0]
		function := vm.chunk.Constants[functionIndex]
		vm.Push(value.FromInterface(function))

	case OP_CALL_FUNCTION:
		argCount := instruction.Operands[0]
		// Аналогично OP_CALL, но для именованных функций
		args := make([]*value.Value, argCount)
		for i := argCount - 1; i >= 0; i-- {
			args[i] = vm.Pop()
		}
		function := vm.Pop()
		result := vm.callFunction(function, args)
		vm.Push(result)

	// Async/await (заглушки)
	case OP_ASYNC:
		_ = vm.Pop() // val не используется
		// Создаем асинхронную задачу (заглушка)
		vm.Push(value.NewValue("async_task"))

	case OP_AWAIT:
		_ = vm.Pop() // promise не используется
		// Ожидаем выполнения (заглушка)
		result := value.NewValue("awaited_result")
		vm.Push(result)

	case OP_SLEEP:
		duration := vm.Pop()
		// Спим указанное количество миллисекунд (заглушка)
		fmt.Printf("Sleep for %s ms\n", duration.String())
		vm.Push(value.NewNil())

	// HTTP операции (заглушки)
	case OP_HTTP_GET:
		url := vm.Pop()
		fmt.Printf("HTTP GET: %s\n", url.String())
		result := value.NewValue(map[string]*value.Value{
			"status": value.NewValue(200),
			"body":   value.NewValue("response body"),
		})
		vm.Push(result)

	// Математические функции
	case OP_MATH_SIN:
		val := vm.Pop()
		if val.IsNumber() {
			import_math := 3.14159 // заглушка для math.Sin
			result := value.NewValue(import_math) 
			vm.Push(result)
		} else {
			vm.Push(value.NewString("Error: sin() requires number"))
		}

	case OP_MATH_COS:
		val := vm.Pop()
		if val.IsNumber() {
			result := value.NewValue(1.0) // заглушка
			vm.Push(result)
		} else {
			vm.Push(value.NewString("Error: cos() requires number"))
		}

	case OP_MATH_SQRT:
		val := vm.Pop()
		if val.IsNumber() {
			// Простая заглушка для sqrt
			valAny := val.Any()
			var num float64
			
			switch valAny.(type) {
			case int64:
				num = float64(valAny.(int64))
			case float64:
				num = valAny.(float64)
			}
			
			if num >= 0 {
				result := value.NewValue(num / 2) // приблизительная заглушка
				vm.Push(result)
			} else {
				vm.Push(value.NewString("Error: sqrt of negative number"))
			}
		} else {
			vm.Push(value.NewString("Error: sqrt() requires number"))
		}

	case OP_MATH_POW:
		exponent := vm.Pop()
		base := vm.Pop()
		if base.IsNumber() && exponent.IsNumber() {
			// Заглушка для pow
			baseAny := base.Any()
			expAny := exponent.Any()
			var baseNum, expNum float64
			
			switch baseAny.(type) {
			case int64:
				baseNum = float64(baseAny.(int64))
			case float64:
				baseNum = baseAny.(float64)
			}
			
			switch expAny.(type) {
			case int64:
				expNum = float64(expAny.(int64))
			case float64:
				expNum = expAny.(float64)
			}
			
			result := value.NewValue(baseNum * expNum) // простая заглушка
			vm.Push(result)
		} else {
			vm.Push(value.NewString("Error: pow() requires numbers"))
		}

	case OP_MATH_ABS:
		val := vm.Pop()
		if val.IsNumber() {
			valAny := val.Any()
			var num float64
			
			switch valAny.(type) {
			case int64:
				num = float64(valAny.(int64))
			case float64:
				num = valAny.(float64)
			}
			
			if num < 0 {
				num = -num
			}
			result := value.NewValue(num)
			vm.Push(result)
		} else {
			vm.Push(value.NewString("Error: abs() requires number"))
		}

	// Строковые операции
	case OP_STRING_LEN:
		str := vm.Pop()
		length := len(str.String())
		vm.Push(value.NewValue(int64(length)))

	case OP_STRING_CONCAT:
		b := vm.Pop()
		a := vm.Pop()
		result := value.NewValue(a.String() + b.String())
		vm.Push(result)

	case OP_STRING_CHAR_AT:
		index := vm.Pop()
		str := vm.Pop()
		s := str.String()
		i := int(index.Int())
		if i >= 0 && i < len(s) {
			char := string(s[i])
			vm.Push(value.NewValue(char))
		} else {
			vm.Push(value.NewString("Error: string index out of bounds"))
		}

	case OP_STRING_SUBSTRING:
		end := vm.Pop()
		start := vm.Pop()
		str := vm.Pop()
		s := str.String()
		startIdx := int(start.Int())
		endIdx := int(end.Int())
		if startIdx >= 0 && endIdx <= len(s) && startIdx <= endIdx {
			result := s[startIdx:endIdx]
			vm.Push(value.NewValue(result))
		} else {
			vm.Push(value.NewString("Error: invalid substring indices"))
		}

	// Методы и свойства (заглушки)
	case OP_METHOD_CALL:
		methodIndex := instruction.Operands[0]
		methodName := vm.chunk.Constants[methodIndex].(string)
		obj := vm.Pop()
		fmt.Printf("Method call: %s.%s()\n", obj.String(), methodName)
		vm.Push(value.NewNil())

	case OP_PROPERTY_ACCESS:
		propIndex := instruction.Operands[0]
		propName := vm.chunk.Constants[propIndex].(string)
		obj := vm.Pop()
		fmt.Printf("Property access: %s.%s\n", obj.String(), propName)
		vm.Push(value.NewNil())

	// Типы и метапрограммирование
	case OP_TYPE_OF:
		_ = vm.Pop() // val не используется в заглушке
		typeName := "unknown" // В реальной реализации нужно добавить метод TypeString в value.Value
		vm.Push(value.NewValue(typeName))

	case OP_TYPE_CHECK:
		expectedType := vm.Pop().String()
		_ = vm.Pop() // val не используется в заглушке
		actualType := "unknown" // В реальной реализации нужно добавить метод TypeString в value.Value
		isMatch := actualType == expectedType
		vm.Push(value.NewValue(isMatch))

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

// callFunction вызывает функцию с заданными аргументами
func (vm *VM) callFunction(function *value.Value, args []*value.Value) *value.Value {
	// Пока это заглушка - в полной реализации здесь будет создание CallFrame,
	// установка локальных переменных и выполнение тела функции
	
	functionName := function.String()
	
	// Простые встроенные функции
	switch functionName {
	case "print":
		if len(args) > 0 {
			fmt.Print(args[0].String())
		}
		return value.NewNil()
		
	case "println":
		if len(args) > 0 {
			fmt.Println(args[0].String())
		}
		return value.NewNil()
		
	case "len":
		if len(args) > 0 {
			arg := args[0]
			if arg.IsString() {
				return value.NewValue(int64(len(arg.String())))
			} else {
				// Для других типов пробуем привести к массиву
				if arr, ok := arg.Any().([]interface{}); ok {
					return value.NewValue(int64(len(arr)))
				}
			}
		}
		return value.NewValue(int64(0))
		
	case "type":
		if len(args) > 0 {
			return value.NewValue("unknown") // В реальной реализации нужно добавить метод TypeString в value.Value
		}
		return value.NewValue("unknown")
		
	default:
		// Для пользовательских функций - заглушка
		fmt.Printf("Calling user function: %s with %d args\n", functionName, len(args))
		return value.NewValue("function_result")
	}
}

// SetBreakpoint устанавливает точку останова для debugger'а
func (vm *VM) SetBreakpoint(line int) {
	if vm.breakpoints == nil {
		vm.breakpoints = make(map[int]bool)
	}
	vm.breakpoints[line] = true
}

// RemoveBreakpoint удаляет точку останова
func (vm *VM) RemoveBreakpoint(line int) {
	if vm.breakpoints != nil {
		delete(vm.breakpoints, line)
	}
}

// isBreakpoint проверяет, есть ли breakpoint на данной строке
func (vm *VM) isBreakpoint(line int) bool {
	if vm.breakpoints == nil {
		return false
	}
	return vm.breakpoints[line]
}

// printLocalVariables выводит локальные переменные для debugger'а
func (vm *VM) printLocalVariables() {
	fmt.Printf("Local variables:\n")
	
	// Выводим глобальные переменные
	if len(vm.globals) > 0 {
		fmt.Printf("  Globals:\n")
		for name, val := range vm.globals {
			fmt.Printf("    %s = %s\n", name, val.String())
		}
	}
	
	// Если есть CallFrames, выводим информацию о них
	if len(vm.callFrames) > 0 {
		frame := vm.callFrames[len(vm.callFrames)-1]
		fmt.Printf("  Current frame: %s (base=%d)\n", frame.function.String(), frame.basePtr)
	}
}

// EnableDebugMode включает режим отладки
func (vm *VM) EnableDebugMode() {
	vm.debugMode = true
}

// DisableDebugMode отключает режим отладки
func (vm *VM) DisableDebugMode() {
	vm.debugMode = false
}

// IsDebugMode возвращает статус режима отладки
func (vm *VM) IsDebugMode() bool {
	return vm.debugMode
}

// GetJIT возвращает JIT компилятор
func (vm *VM) GetJIT() *JITCompiler {
	return vm.jit
}

// EnableJIT включает JIT компиляцию
func (vm *VM) EnableJIT() {
	vm.jit.Enable()
}

// DisableJIT выключает JIT компиляцию
func (vm *VM) DisableJIT() {
	vm.jit.Disable()
}

// IsBreakpoint проверяет, есть ли breakpoint на данной строке (экспортированный метод для тестов)
func (vm *VM) IsBreakpoint(line int) bool {
	return vm.isBreakpoint(line)
}
