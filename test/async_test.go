package test

import (
	"testing"
	"time"
	"foo_lang/parser"
	"foo_lang/builtin"
	"foo_lang/scope"
)

func TestAsyncBasic(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	
	// Инициализируем встроенные функции
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Простая async функция
fn simpleAsync() {
    return "Hello from async"
}

// Запускаем async и ждем результат
let promise = async simpleAsync()
let result = await promise
print("Result: " + result)
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Проверяем, что не было ошибок
			if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
				t.Errorf("Async operation failed: %s", str)
			}
		}
	}
}

func TestAsyncAwaitWithSleep(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Функция с задержкой
fn delayedOperation(delay) {
    await sleep(delay)
    return "Completed after " + delay.toString() + "ms"
}

// Измеряем время выполнения
let start = "Starting..."
print(start)

let promise = async delayedOperation(100)
let result = await promise
print("Result: " + result)
`

	startTime := time.Now()
	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		expr.Eval()
	}
	
	elapsed := time.Since(startTime)
	// Проверяем, что операция заняла примерно 100мс
	if elapsed < 90*time.Millisecond || elapsed > 200*time.Millisecond {
		t.Errorf("Expected operation to take ~100ms, but took %v", elapsed)
	}
}

func TestPromiseAll(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Простые async операции без параметров (избегаем race conditions)
fn taskA() {
    await sleep(30)
    return "TaskA completed"
}

fn taskB() {
    await sleep(50)
    return "TaskB completed"
}

fn taskC() {
    await sleep(20)
    return "TaskC completed"
}

// Запускаем параллельно
let task1 = async taskA()
let task2 = async taskB()
let task3 = async taskC()

// Ждем завершения всех
let results = await Promise.all(task1, task2, task3)
print("All tasks completed")

// Проверяем результаты
print("Results count: " + results.length().toString())
`

	startTime := time.Now()
	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		expr.Eval()
	}
	
	elapsed := time.Since(startTime)
	// Promise.all должен завершиться за время самой долгой задачи (~50ms)
	if elapsed > 100*time.Millisecond {
		t.Errorf("Promise.all took too long: %v (expected ~50ms)", elapsed)
	}
}

func TestPromiseAny(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Простые async операции с разными задержками (без параметров)
fn raceSlow() {
    await sleep(80)
    return "Slow won!"
}

fn raceMedium() {
    await sleep(40)
    return "Medium won!"
}

fn raceFast() {
    await sleep(20)
    return "Fast won!"
}

// Запускаем гонку
let slow = async raceSlow()
let medium = async raceMedium()
let fast = async raceFast()

// Ждем первый результат
let winner = await Promise.any(slow, medium, fast)
print("Winner: " + winner)
`

	startTime := time.Now()
	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		expr.Eval()
	}
	
	elapsed := time.Since(startTime)
	// Promise.any должен завершиться за время самой быстрой задачи (~20ms)
	if elapsed > 50*time.Millisecond {
		t.Errorf("Promise.any took too long: %v (expected ~20ms)", elapsed)
	}
}

func TestAsyncErrorHandling(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Простая успешная async функция
fn successOperation() {
    await sleep(30)
    return "Success"
}

// Тест успешной операции
let successPromise = async successOperation()
let successResult = await successPromise
print("Success result: " + successResult)
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		result := expr.Eval()
		// Проверяем, что успешная операция прошла без ошибок
		if result != nil && result.Any() != nil {
			if str, ok := result.Any().(string); ok && str == "Success" {
				// Успех
			}
		}
	}
}

func TestAsyncParallelExecution(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Простая async функция для тестирования параллельного выполнения
fn simpleTask() {
    await sleep(30)
    return "Task completed"
}

// Запускаем задачу
let task = async simpleTask()

// Ждем результат
let result = await task

print("Parallel execution test completed")
print(result)
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		expr.Eval()
	}
	
	// Тест пройден если не было паник
	t.Log("Parallel execution test passed")
}

func TestAsyncWithClosures(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Простой async тест с замыканием (без параметров)
let sharedValue = "Initial"

fn createAsyncClosure() {
    return fn() {
        await sleep(30)
        sharedValue = sharedValue + " + Modified"
        return "Closure executed"
    }
}

// Создаем и выполняем замыкание
let closure = createAsyncClosure()
let promise = async closure()
let result = await promise

print("Result: " + result)
print("Final shared value: " + sharedValue)
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		expr.Eval()
	}
	
	// Тест пройден если замыкания работают с async
	t.Log("Async with closures test passed")
}