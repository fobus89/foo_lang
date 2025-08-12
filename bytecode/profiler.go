package bytecode

import (
	"fmt"
	"time"
	"sort"
)

// Profiler собирает статистику производительности выполнения bytecode
type Profiler struct {
	enabled           bool
	startTime         time.Time
	totalTime         time.Duration
	instructionCounts map[OpCode]int64
	instructionTimes  map[OpCode]time.Duration
	functionTimes     map[string]time.Duration
	functionCalls     map[string]int64
	functionStack     []FunctionFrame
	hotspots          []Hotspot
}

// FunctionFrame представляет кадр функции в профайлере
type FunctionFrame struct {
	name      string
	startTime time.Time
}

// Hotspot представляет горячую точку в коде
type Hotspot struct {
	Name       string
	CallCount  int64
	TotalTime  time.Duration
	AvgTime    time.Duration
	Percentage float64
}

// NewProfiler создает новый профайлер
func NewProfiler() *Profiler {
	return &Profiler{
		enabled:           true,
		instructionCounts: make(map[OpCode]int64),
		instructionTimes:  make(map[OpCode]time.Duration),
		functionTimes:     make(map[string]time.Duration),
		functionCalls:     make(map[string]int64),
		functionStack:     make([]FunctionFrame, 0, 64),
		hotspots:          make([]Hotspot, 0, 100),
	}
}

// Enable включает профилирование
func (p *Profiler) Enable() {
	p.enabled = true
}

// Disable выключает профилирование
func (p *Profiler) Disable() {
	p.enabled = false
}

// StartExecution начинает профилирование выполнения
func (p *Profiler) StartExecution() {
	if !p.enabled {
		return
	}
	p.startTime = time.Now()
}

// EndExecution завершает профилирование выполнения
func (p *Profiler) EndExecution() {
	if !p.enabled {
		return
	}
	p.totalTime = time.Since(p.startTime)
	p.calculateHotspots()
}

// RecordInstruction записывает выполнение инструкции
func (p *Profiler) RecordInstruction(opcode OpCode) {
	if !p.enabled {
		return
	}
	
	start := time.Now()
	p.instructionCounts[opcode]++
	
	// Простая симуляция времени выполнения инструкции
	// В реальной реализации здесь будет измерение реального времени
	elapsed := time.Since(start)
	p.instructionTimes[opcode] += elapsed
}

// StartFunction начинает профилирование функции
func (p *Profiler) StartFunction(name string) {
	if !p.enabled {
		return
	}
	
	frame := FunctionFrame{
		name:      name,
		startTime: time.Now(),
	}
	p.functionStack = append(p.functionStack, frame)
	p.functionCalls[name]++
}

// EndFunction завершает профилирование функции
func (p *Profiler) EndFunction(name string) {
	if !p.enabled || len(p.functionStack) == 0 {
		return
	}
	
	// Ищем кадр функции в стеке
	for i := len(p.functionStack) - 1; i >= 0; i-- {
		if p.functionStack[i].name == name {
			elapsed := time.Since(p.functionStack[i].startTime)
			p.functionTimes[name] += elapsed
			
			// Удаляем кадр из стека
			p.functionStack = append(p.functionStack[:i], p.functionStack[i+1:]...)
			break
		}
	}
}

// calculateHotspots вычисляет горячие точки
func (p *Profiler) calculateHotspots() {
	p.hotspots = p.hotspots[:0] // очищаем предыдущие результаты
	
	for name, totalTime := range p.functionTimes {
		callCount := p.functionCalls[name]
		avgTime := time.Duration(0)
		if callCount > 0 {
			avgTime = time.Duration(int64(totalTime) / callCount)
		}
		
		percentage := float64(totalTime) / float64(p.totalTime) * 100
		
		hotspot := Hotspot{
			Name:       name,
			CallCount:  callCount,
			TotalTime:  totalTime,
			AvgTime:    avgTime,
			Percentage: percentage,
		}
		
		p.hotspots = append(p.hotspots, hotspot)
	}
	
	// Сортируем по времени выполнения (убыванию)
	sort.Slice(p.hotspots, func(i, j int) bool {
		return p.hotspots[i].TotalTime > p.hotspots[j].TotalTime
	})
}

// PrintReport выводит отчет профилирования
func (p *Profiler) PrintReport() {
	fmt.Println("=== ПРОФИЛИРОВАНИЕ ПРОИЗВОДИТЕЛЬНОСТИ ===")
	fmt.Printf("Общее время выполнения: %v\n", p.totalTime)
	fmt.Println()
	
	// Статистика инструкций
	fmt.Println("--- Статистика инструкций ---")
	
	type instructionStat struct {
		opcode OpCode
		count  int64
		time   time.Duration
	}
	
	var stats []instructionStat
	for opcode, count := range p.instructionCounts {
		stats = append(stats, instructionStat{
			opcode: opcode,
			count:  count,
			time:   p.instructionTimes[opcode],
		})
	}
	
	// Сортируем по количеству выполнений
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].count > stats[j].count
	})
	
	fmt.Printf("%-20s %10s %15s\n", "Инструкция", "Количество", "Время")
	fmt.Println("------------------------------------------------")
	for _, stat := range stats {
		opcodeStr := p.getOpCodeString(stat.opcode)
		fmt.Printf("%-20s %10d %15v\n", opcodeStr, stat.count, stat.time)
	}
	fmt.Println()
	
	// Горячие точки функций
	fmt.Println("--- Горячие точки функций ---")
	fmt.Printf("%-20s %10s %15s %15s %10s\n", "Функция", "Вызовы", "Общее время", "Среднее время", "Процент")
	fmt.Println("-----------------------------------------------------------------------")
	
	for _, hotspot := range p.hotspots {
		if hotspot.Percentage >= 1.0 { // показываем только функции, занимающие >1% времени
			fmt.Printf("%-20s %10d %15v %15v %9.2f%%\n",
				hotspot.Name,
				hotspot.CallCount,
				hotspot.TotalTime,
				hotspot.AvgTime,
				hotspot.Percentage)
		}
	}
	fmt.Println()
	
	// Рекомендации по оптимизации
	fmt.Println("--- Рекомендации по оптимизации ---")
	if len(p.hotspots) > 0 {
		topHotspot := p.hotspots[0]
		if topHotspot.Percentage > 50 {
			fmt.Printf("⚠️  Функция '%s' занимает %.1f%% времени выполнения - требует оптимизации\n", 
				topHotspot.Name, topHotspot.Percentage)
		}
		
		if topHotspot.CallCount > 1000 {
			fmt.Printf("🔥 Функция '%s' вызывается %d раз - кандидат для JIT компиляции\n", 
				topHotspot.Name, topHotspot.CallCount)
		}
	}
	
	// Проверяем часто используемые инструкции
	for _, stat := range stats[:3] { // топ-3 инструкции
		if stat.count > 100 {
			opcodeStr := p.getOpCodeString(stat.opcode)
			fmt.Printf("💡 Инструкция %s выполняется %d раз - можно оптимизировать\n", 
				opcodeStr, stat.count)
		}
	}
	
	fmt.Println("=====================================")
}

// getOpCodeString возвращает строковое представление опкода
func (p *Profiler) getOpCodeString(opcode OpCode) string {
	switch opcode {
	case OP_CONSTANT:
		return "OP_CONSTANT"
	case OP_ADD:
		return "OP_ADD"
	case OP_SUBTRACT:
		return "OP_SUBTRACT"
	case OP_MULTIPLY:
		return "OP_MULTIPLY"
	case OP_DIVIDE:
		return "OP_DIVIDE"
	case OP_GET_GLOBAL:
		return "OP_GET_GLOBAL"
	case OP_SET_GLOBAL:
		return "OP_SET_GLOBAL"
	case OP_CALL:
		return "OP_CALL"
	case OP_PRINT:
		return "OP_PRINT"
	default:
		return fmt.Sprintf("OP_%d", int(opcode))
	}
}

// GetHotspots возвращает список горячих точек
func (p *Profiler) GetHotspots() []Hotspot {
	return p.hotspots
}

// GetTotalTime возвращает общее время выполнения
func (p *Profiler) GetTotalTime() time.Duration {
	return p.totalTime
}

// GetInstructionCount возвращает количество выполнений инструкции
func (p *Profiler) GetInstructionCount(opcode OpCode) int64 {
	return p.instructionCounts[opcode]
}

// Reset сбрасывает все собранные данные профилирования
func (p *Profiler) Reset() {
	p.startTime = time.Time{}
	p.totalTime = 0
	
	for k := range p.instructionCounts {
		delete(p.instructionCounts, k)
	}
	for k := range p.instructionTimes {
		delete(p.instructionTimes, k)
	}
	for k := range p.functionTimes {
		delete(p.functionTimes, k)
	}
	for k := range p.functionCalls {
		delete(p.functionCalls, k)
	}
	
	p.functionStack = p.functionStack[:0]
	p.hotspots = p.hotspots[:0]
}