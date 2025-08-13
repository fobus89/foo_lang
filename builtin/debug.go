package builtin

import (
	"fmt"
	"foo_lang/value"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"
)

// Debug выводит подробную информацию о значении с типом
func Debug(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("debug expects 1 argument")
	}

	val := args[0]
	output := formatDebugValue(val, 0)
	fmt.Println("[DEBUG]", output)

	return val, nil
}

// formatDebugValue форматирует значение для отладки
func formatDebugValue(val *value.Value, indent int) string {
	indentStr := strings.Repeat("  ", indent)
	
	if val == nil {
		return fmt.Sprintf("%snil", indentStr)
	}

	data := val.Any()
	if data == nil {
		return fmt.Sprintf("%snull", indentStr)
	}

	// Определяем тип и значение
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Map:
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Object[%d] {\n", v.Len()))
		for _, key := range v.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			valItem := v.MapIndex(key).Interface().(*value.Value)
			sb.WriteString(fmt.Sprintf("%s  %s: %s\n", indentStr, keyStr, formatDebugValue(valItem, indent+1)))
		}
		sb.WriteString(fmt.Sprintf("%s}", indentStr))
		return sb.String()
		
	case reflect.Slice:
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Array[%d] [\n", v.Len()))
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface().(*value.Value)
			sb.WriteString(fmt.Sprintf("%s  [%d]: %s\n", indentStr, i, formatDebugValue(item, indent+1)))
		}
		sb.WriteString(fmt.Sprintf("%s]", indentStr))
		return sb.String()
		
	case reflect.String:
		return fmt.Sprintf("String(%d): %q", len(v.String()), v.String())
		
	case reflect.Int, reflect.Int64:
		return fmt.Sprintf("Int: %d", v.Int())
		
	case reflect.Float64:
		return fmt.Sprintf("Float: %f", v.Float())
		
	case reflect.Bool:
		return fmt.Sprintf("Bool: %v", v.Bool())
		
	case reflect.Func:
		return fmt.Sprintf("Function: %s", t.String())
		
	default:
		return fmt.Sprintf("%s: %v", t.String(), data)
	}
}

// Trace выводит трассировку стека вызовов
func Trace(args []*value.Value) (*value.Value, error) {
	maxDepth := 10
	if len(args) == 1 {
		if args[0].IsNumber() {
			maxDepth = args[0].Int()
		} else {
			return nil, fmt.Errorf("trace depth must be a number")
		}
	} else if len(args) > 1 {
		return nil, fmt.Errorf("trace expects 0 or 1 arguments")
	}

	fmt.Println("=== Stack Trace ===")
	
	// Получаем стек вызовов
	pc := make([]uintptr, maxDepth)
	n := runtime.Callers(2, pc)
	
	if n == 0 {
		fmt.Println("No stack trace available")
		return value.NewValue(nil), nil
	}

	frames := runtime.CallersFrames(pc[:n])
	
	i := 0
	for {
		frame, more := frames.Next()
		
		// Пропускаем внутренние вызовы runtime
		if strings.Contains(frame.File, "runtime/") {
			if !more {
				break
			}
			continue
		}
		
		fmt.Printf("#%d %s\n", i, frame.Function)
		fmt.Printf("   %s:%d\n", frame.File, frame.Line)
		i++
		
		if !more {
			break
		}
	}
	
	fmt.Println("==================")
	
	return value.NewValue(nil), nil
}

// GetStackTrace возвращает стек вызовов как строку
func GetStackTrace(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getStackTrace expects 0 arguments")
	}

	stack := debug.Stack()
	return value.NewValue(string(stack)), nil
}

// Profile профилирует выполнение функции
func Profile(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("profile expects 1 argument (function)")
	}

	fn := args[0]
	
	// Проверяем, что это функция
	if fn.Any() == nil {
		return nil, fmt.Errorf("profile expects a function")
	}

	// Засекаем время начала
	startTime := time.Now()
	startMem := getMemStats()

	// Выполняем функцию
	// Предполагаем, что функция может быть вызвана через интерфейс
	callable, ok := fn.Any().(func([]*value.Value) (*value.Value, error))
	if !ok {
		return nil, fmt.Errorf("profile argument must be a callable function")
	}

	result, err := callable([]*value.Value{})
	
	// Засекаем время окончания
	duration := time.Since(startTime)
	endMem := getMemStats()

	// Создаем объект с результатами профилирования
	profileResult := map[string]*value.Value{
		"result":       result,
		"duration_ms":  value.NewValue(float64(duration.Milliseconds())),
		"duration_ns":  value.NewValue(int64(duration.Nanoseconds())),
		"memory_used":  value.NewValue(int64(endMem.Alloc - startMem.Alloc)),
		"gc_runs":      value.NewValue(int64(endMem.NumGC - startMem.NumGC)),
	}

	if err != nil {
		profileResult["error"] = value.NewValue(err.Error())
	}

	return value.NewValue(profileResult), nil
}

// getMemStats получает статистику памяти
func getMemStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

// Benchmark выполняет бенчмарк функции
func Benchmark(args []*value.Value) (*value.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("benchmark expects 1 or 2 arguments (function, [iterations])")
	}

	fn := args[0]
	iterations := 100
	
	if len(args) == 2 {
		if !args[1].IsNumber() {
			return nil, fmt.Errorf("benchmark iterations must be a number")
		}
		iterations = args[1].Int()
	}

	// Проверяем, что это функция
	callable, ok := fn.Any().(func([]*value.Value) (*value.Value, error))
	if !ok {
		return nil, fmt.Errorf("benchmark first argument must be a callable function")
	}

	// Прогрев
	for i := 0; i < 10; i++ {
		callable([]*value.Value{})
	}

	// Запускаем бенчмарк
	startTime := time.Now()
	startMem := getMemStats()

	for i := 0; i < iterations; i++ {
		callable([]*value.Value{})
	}

	duration := time.Since(startTime)
	endMem := getMemStats()

	// Вычисляем среднее время
	avgDuration := duration / time.Duration(iterations)

	// Создаем результат
	result := map[string]*value.Value{
		"iterations":       value.NewValue(int64(iterations)),
		"total_time_ms":    value.NewValue(float64(duration.Milliseconds())),
		"avg_time_ms":      value.NewValue(float64(avgDuration.Milliseconds())),
		"avg_time_us":      value.NewValue(float64(avgDuration.Microseconds())),
		"avg_time_ns":      value.NewValue(int64(avgDuration.Nanoseconds())),
		"ops_per_second":   value.NewValue(float64(iterations) / duration.Seconds()),
		"memory_per_op":    value.NewValue(int64((endMem.Alloc - startMem.Alloc) / uint64(iterations))),
		"total_memory":     value.NewValue(int64(endMem.Alloc - startMem.Alloc)),
		"gc_runs":          value.NewValue(int64(endMem.NumGC - startMem.NumGC)),
	}

	return value.NewValue(result), nil
}

// MemStats возвращает статистику памяти
func MemStats(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("memStats expects 0 arguments")
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := map[string]*value.Value{
		"alloc":          value.NewValue(int64(m.Alloc)),
		"total_alloc":    value.NewValue(int64(m.TotalAlloc)),
		"sys":            value.NewValue(int64(m.Sys)),
		"heap_alloc":     value.NewValue(int64(m.HeapAlloc)),
		"heap_sys":       value.NewValue(int64(m.HeapSys)),
		"heap_idle":      value.NewValue(int64(m.HeapIdle)),
		"heap_in_use":    value.NewValue(int64(m.HeapInuse)),
		"heap_released":  value.NewValue(int64(m.HeapReleased)),
		"heap_objects":   value.NewValue(int64(m.HeapObjects)),
		"stack_in_use":   value.NewValue(int64(m.StackInuse)),
		"stack_sys":      value.NewValue(int64(m.StackSys)),
		"gc_count":       value.NewValue(int64(m.NumGC)),
		"gc_cpu_percent": value.NewValue(m.GCCPUFraction * 100),
		"goroutines":     value.NewValue(int64(runtime.NumGoroutine())),
		"cpus":           value.NewValue(int64(runtime.NumCPU())),
	}

	return value.NewValue(stats), nil
}

// GC запускает сборку мусора
func GC(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("gc expects 0 arguments")
	}

	runtime.GC()
	return value.NewValue(nil), nil
}

// TypeOf возвращает тип значения
func TypeOf(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("typeOf expects 1 argument")
	}

	val := args[0]
	if val == nil {
		return value.NewValue("nil"), nil
	}

	data := val.Any()
	if data == nil {
		return value.NewValue("null"), nil
	}

	// Определяем тип
	switch data.(type) {
	case string:
		return value.NewValue("string"), nil
	case int64:
		return value.NewValue("int"), nil
	case float64:
		return value.NewValue("float"), nil
	case bool:
		return value.NewValue("bool"), nil
	case map[string]*value.Value:
		return value.NewValue("object"), nil
	case []*value.Value:
		return value.NewValue("array"), nil
	case func([]*value.Value) (*value.Value, error):
		return value.NewValue("function"), nil
	default:
		t := reflect.TypeOf(data)
		return value.NewValue(t.String()), nil
	}
}

// SizeOf возвращает размер значения в памяти
func SizeOf(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sizeOf expects 1 argument")
	}

	val := args[0]
	if val == nil {
		return value.NewValue(int64(0)), nil
	}

	data := val.Any()
	if data == nil {
		return value.NewValue(int64(0)), nil
	}

	size := calculateSize(data)
	return value.NewValue(int64(size)), nil
}

// calculateSize рекурсивно вычисляет размер значения
func calculateSize(data interface{}) uintptr {
	if data == nil {
		return 0
	}

	v := reflect.ValueOf(data)
	size := reflect.TypeOf(data).Size()

	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			size += calculateSize(key.Interface())
			size += calculateSize(v.MapIndex(key).Interface())
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			size += calculateSize(v.Index(i).Interface())
		}
	case reflect.String:
		size += uintptr(v.Len())
	}

	return size
}

// Assert проверяет условие и паникует если оно ложно
func Assert(args []*value.Value) (*value.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("assert expects 1 or 2 arguments (condition, [message])")
	}

	condition := args[0].Bool()
	message := "Assertion failed"
	
	if len(args) == 2 {
		if args[1].IsString() {
			message = args[1].String()
		}
	}

	if !condition {
		panic(message)
	}

	return value.NewValue(true), nil
}

// StartCPUProfile начинает профилирование CPU
func StartCPUProfile(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("startCPUProfile expects 1 argument (filename)")
	}

	if !args[0].IsString() {
		return nil, fmt.Errorf("startCPUProfile filename must be a string")
	}

	filename := args[0].String()
	
	f, err := os.Create(filename)
	if err != nil {
		return value.NewValue(false), nil
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		f.Close()
		return value.NewValue(false), nil
	}

	// Сохраняем файл для последующего закрытия
	cpuProfileFile = f
	
	return value.NewValue(true), nil
}

// StopCPUProfile останавливает профилирование CPU
func StopCPUProfile(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("stopCPUProfile expects 0 arguments")
	}

	pprof.StopCPUProfile()
	
	if cpuProfileFile != nil {
		cpuProfileFile.Close()
		cpuProfileFile = nil
	}
	
	return value.NewValue(nil), nil
}

var cpuProfileFile *os.File