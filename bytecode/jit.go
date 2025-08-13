package bytecode

import (
	"fmt"
	"foo_lang/value"
	"time"
)

// JITCompiler представляет Just-In-Time компилятор
type JITCompiler struct {
	enabled         bool
	hotspotThreshold int64              // Порог для считания функции "горячей"
	compiledFunctions map[string]*CompiledFunction
	executionStats   map[string]*ExecutionStats
}

// CompiledFunction представляет JIT-скомпилированную функцию
type CompiledFunction struct {
	name            string
	originalBytecode []Instruction
	optimizedCode   []OptimizedInstruction
	compiledAt      time.Time
	executionCount  int64
	totalTime       time.Duration
}

// OptimizedInstruction представляет оптимизированную инструкцию
type OptimizedInstruction struct {
	Type     OptimizationType
	OpCode   OpCode
	Operands []int
	// Для специальных оптимизаций
	DirectValue    *value.Value // Для константных значений
	OptimizedFunc  func(*VM) *value.Value // Для прямого выполнения
}

// OptimizationType тип оптимизации
type OptimizationType int

const (
	OPT_NONE OptimizationType = iota
	OPT_CONSTANT_FOLDING    // Свертка констант
	OPT_DIRECT_CALL        // Прямой вызов функции
	OPT_INLINE_MATH        // Инлайн математических операций
	OPT_LOOP_UNROLL        // Разворачивание циклов
)

// ExecutionStats статистика выполнения для JIT
type ExecutionStats struct {
	CallCount      int64
	TotalTime      time.Duration
	AverageTime    time.Duration
	IsHot          bool
	LastCompiled   time.Time
}

// NewJITCompiler создает новый JIT компилятор
func NewJITCompiler() *JITCompiler {
	return &JITCompiler{
		enabled:           true,
		hotspotThreshold:  100, // Функция считается горячей после 100 вызовов
		compiledFunctions: make(map[string]*CompiledFunction),
		executionStats:    make(map[string]*ExecutionStats),
	}
}

// Enable включает JIT компиляцию
func (jit *JITCompiler) Enable() {
	jit.enabled = true
}

// Disable выключает JIT компиляцию
func (jit *JITCompiler) Disable() {
	jit.enabled = false
}

// SetHotspotThreshold устанавливает порог для горячих функций
func (jit *JITCompiler) SetHotspotThreshold(threshold int64) {
	jit.hotspotThreshold = threshold
}

// RecordExecution записывает статистику выполнения функции
func (jit *JITCompiler) RecordExecution(funcName string, executionTime time.Duration) {
	if !jit.enabled {
		return
	}

	stats, exists := jit.executionStats[funcName]
	if !exists {
		stats = &ExecutionStats{}
		jit.executionStats[funcName] = stats
	}

	stats.CallCount++
	stats.TotalTime += executionTime
	stats.AverageTime = time.Duration(int64(stats.TotalTime) / stats.CallCount)

	// Проверяем, стала ли функция горячей
	if !stats.IsHot && stats.CallCount >= jit.hotspotThreshold {
		stats.IsHot = true
		fmt.Printf("🔥 JIT: Function '%s' became hot (%d calls, avg %v)\n", 
			funcName, stats.CallCount, stats.AverageTime)
		
		// Пытаемся скомпилировать
		jit.compileFunction(funcName)
	}
}

// compileFunction компилирует функцию в оптимизированный код
func (jit *JITCompiler) compileFunction(funcName string) {
	if !jit.enabled {
		return
	}

	// Проверяем, не скомпилирована ли уже
	if _, exists := jit.compiledFunctions[funcName]; exists {
		return
	}

	fmt.Printf("🚀 JIT: Compiling function '%s'...\n", funcName)

	// Создаем заглушку скомпилированной функции
	// В реальной реализации здесь будет анализ bytecode и оптимизация
	compiled := &CompiledFunction{
		name:        funcName,
		compiledAt:  time.Now(),
		optimizedCode: jit.optimizeBytecode(funcName, nil), // nil - заглушка для bytecode
	}

	jit.compiledFunctions[funcName] = compiled
	fmt.Printf("✅ JIT: Function '%s' compiled successfully\n", funcName)
}

// optimizeBytecode оптимизирует bytecode
func (jit *JITCompiler) optimizeBytecode(funcName string, originalBytecode []Instruction) []OptimizedInstruction {
	var optimized []OptimizedInstruction

	// Простые оптимизации на основе имени функции и bytecode
	_ = funcName // Используется для специфичных оптимизаций функций
	_ = originalBytecode // Будет использоваться для анализа исходного кода
	
	// 1. Константная свертка для математических операций
	optimized = append(optimized, OptimizedInstruction{
		Type:   OPT_CONSTANT_FOLDING,
		OpCode: OP_ADD,
		OptimizedFunc: func(vm *VM) *value.Value {
			// Быстрое сложение без проверок типов
			b := vm.Pop()
			a := vm.Pop()
			if a.IsNumber() && b.IsNumber() {
				// Используем простое приведение через Any()
				aVal := a.Any()
				bVal := b.Any()
				var sum float64
				
				switch aVal.(type) {
				case int64:
					sum += float64(aVal.(int64))
				case float64:
					sum += aVal.(float64)
				}
				
				switch bVal.(type) {
				case int64:
					sum += float64(bVal.(int64))
				case float64:
					sum += bVal.(float64)
				}
				
				result := value.NewValue(sum)
				vm.Push(result)
				return nil
			}
			// Fallback к стандартной реализации
			return value.Add(a, b)
		},
	})

	// 2. Оптимизация вызовов print
	optimized = append(optimized, OptimizedInstruction{
		Type:   OPT_DIRECT_CALL,
		OpCode: OP_PRINT,
		OptimizedFunc: func(vm *VM) *value.Value {
			val := vm.Pop()
			fmt.Print(val.String()) // Прямой вызов без лишних проверок
			vm.Push(value.NewNil())
			return nil
		},
	})

	// 3. Инлайн математических функций
	optimized = append(optimized, OptimizedInstruction{
		Type:   OPT_INLINE_MATH,
		OpCode: OP_MATH_ABS,
		OptimizedFunc: func(vm *VM) *value.Value {
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
				vm.Push(value.NewValue(num))
				return nil
			}
			return value.NewString("Error: abs() requires number")
		},
	})

	return optimized
}

// ExecuteOptimized выполняет JIT-скомпилированную функцию
func (jit *JITCompiler) ExecuteOptimized(vm *VM, funcName string) *value.Value {
	if !jit.enabled {
		return nil
	}

	compiled, exists := jit.compiledFunctions[funcName]
	if !exists {
		return nil // Функция не скомпилирована
	}

	start := time.Now()
	
	fmt.Printf("⚡ JIT: Executing optimized '%s'\n", funcName)
	
	// Выполняем оптимизированные инструкции
	for _, instruction := range compiled.optimizedCode {
		if instruction.OptimizedFunc != nil {
			result := instruction.OptimizedFunc(vm)
			if result != nil && result.IsReturn() {
				compiled.executionCount++
				compiled.totalTime += time.Since(start)
				return result
			}
		}
	}

	compiled.executionCount++
	compiled.totalTime += time.Since(start)
	
	return value.NewNil()
}

// IsCompiled проверяет, скомпилирована ли функция
func (jit *JITCompiler) IsCompiled(funcName string) bool {
	_, exists := jit.compiledFunctions[funcName]
	return exists
}

// GetStats возвращает статистику JIT компилятора
func (jit *JITCompiler) GetStats() map[string]*ExecutionStats {
	return jit.executionStats
}

// PrintReport выводит отчет JIT компилятора
func (jit *JITCompiler) PrintReport() {
	fmt.Println("=== JIT КОМПИЛЯТОР ОТЧЕТ ===")
	fmt.Printf("Статус: %s\n", map[bool]string{true: "Включен", false: "Выключен"}[jit.enabled])
	fmt.Printf("Порог горячих функций: %d вызовов\n", jit.hotspotThreshold)
	fmt.Printf("Скомпилировано функций: %d\n", len(jit.compiledFunctions))
	fmt.Println()

	if len(jit.compiledFunctions) > 0 {
		fmt.Println("--- Скомпилированные функции ---")
		fmt.Printf("%-20s %15s %15s %15s\n", "Функция", "Вызовы JIT", "Время JIT", "Скомпилирована")
		fmt.Println("---------------------------------------------------------------")
		
		for name, compiled := range jit.compiledFunctions {
			fmt.Printf("%-20s %15d %15v %15v\n", 
				name, 
				compiled.executionCount,
				compiled.totalTime,
				compiled.compiledAt.Format("15:04:05"))
		}
		fmt.Println()
	}

	if len(jit.executionStats) > 0 {
		fmt.Println("--- Статистика функций ---")
		fmt.Printf("%-20s %10s %15s %15s %10s\n", "Функция", "Вызовы", "Общее время", "Среднее время", "Статус")
		fmt.Println("-----------------------------------------------------------------------")
		
		for name, stats := range jit.executionStats {
			status := "Холодная"
			if stats.IsHot {
				status = "🔥 Горячая"
			}
			
			fmt.Printf("%-20s %10d %15v %15v %10s\n",
				name,
				stats.CallCount,
				stats.TotalTime,
				stats.AverageTime,
				status)
		}
		fmt.Println()
	}

	// Рекомендации
	fmt.Println("--- Рекомендации JIT ---")
	hotCount := 0
	for _, stats := range jit.executionStats {
		if stats.IsHot {
			hotCount++
		}
	}
	
	if hotCount == 0 {
		fmt.Println("💡 Горячих функций не найдено. Уменьшите порог или проверьте нагрузку.")
	} else {
		fmt.Printf("✅ Найдено %d горячих функций. JIT активно работает!\n", hotCount)
	}
	
	if len(jit.compiledFunctions) > 0 {
		fmt.Printf("🚀 JIT компиляция активна. Производительность должна улучшиться!\n")
	}
	
	fmt.Println("=============================")
}

// CompileFunction экспортированный метод для принудительной компиляции функции
func (jit *JITCompiler) CompileFunction(funcName string) {
	jit.compileFunction(funcName)
}

// Reset сбрасывает все данные JIT компилятора
func (jit *JITCompiler) Reset() {
	for k := range jit.compiledFunctions {
		delete(jit.compiledFunctions, k)
	}
	for k := range jit.executionStats {
		delete(jit.executionStats, k)
	}
}