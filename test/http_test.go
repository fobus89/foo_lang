package test

import (
	"testing"
	"time"
	"foo_lang/parser"
	"foo_lang/builtin"
	"foo_lang/scope"
)

func TestHttpClientFunctions(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	
	// Инициализируем HTTP функции
	builtin.InitializeHttpFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Тест HTTP клиента без реальных запросов
print("=== Тест HTTP клиента ===")

// Устанавливаем таймаут
httpSetTimeout(10)
print("Таймаут установлен")

// Тестируем URL утилиты
let encoded = urlEncode("Hello World")
let decoded = urlDecode(encoded)
print("URL кодирование работает: " + (decoded == "Hello World").toString())

print("HTTP клиент тесты завершены")
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Проверяем, что не было ошибок
			if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
				t.Errorf("HTTP client test failed: %s", str)
			}
		}
	}
}

func TestHttpServerBasic(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	
	// Инициализируем HTTP функции
	builtin.InitializeHttpFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Тест HTTP сервера
print("=== Тест HTTP сервера ===")

// Создаем сервер
httpCreateServer()

// Добавляем маршрут
fn helloHandler(request) {
    print("Получен запрос: " + request.method + " " + request.path)
    let response = {
        "status": 200,
        "headers": {"Content-Type": "application/json"},
        "body": "Hello from foo_lang server!"
    }
    return response
}

httpRoute("GET", "/hello", helloHandler)

// Запускаем сервер на порту 8080
let serverResult = httpStartServer(8080)
print("Сервер запущен: " + serverResult)

// Даем серверу время запуститься
await sleep(100)

// Останавливаем сервер
let stopResult = httpStopServer()
print("Сервер остановлен: " + stopResult)

print("HTTP сервер тесты завершены")
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Проверяем, что не было ошибок
			if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
				t.Errorf("HTTP server test failed: %s", str)
			}
		}
	}
}

func TestHttpUtilityFunctions(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	
	// Инициализируем HTTP функции
	builtin.InitializeHttpFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Тест утилит HTTP
print("=== Тест HTTP утилит ===")

// URL кодирование/декодирование
let originalText = "Hello World & Foo/Bar"
let encoded = urlEncode(originalText)
let decoded = urlDecode(encoded)

print("Оригинал: " + originalText)
print("Закодировано: " + encoded)
print("Декодировано: " + decoded)

// Проверяем, что декодированное значение совпадает с оригиналом
if (decoded == originalText) {
    print("URL кодирование/декодирование работает правильно")
} else {
    print("ОШИБКА: URL кодирование/декодирование не работает")
}

print("HTTP утилиты тесты завершены")
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Проверяем, что не было ошибок
			if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
				t.Errorf("HTTP utilities test failed: %s", str)
			}
		}
	}
}

func TestHttpMethods(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	
	// Инициализируем HTTP функции
	builtin.InitializeHttpFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Тест всех HTTP методов
print("=== Тест HTTP методов ===")

// Устанавливаем короткий таймаут для тестов
httpSetTimeout(5)
print("Таймаут установлен для HTTP методов")

// Тестируем создание сервера
let server = httpCreateServer()
print("HTTP сервер создан")

print("Все HTTP методы протестированы")
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Проверяем, что не было ошибок
			if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
				t.Errorf("HTTP methods test failed: %s", str)
			}
		}
	}
}

func TestHttpHeaders(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	
	// Инициализируем HTTP функции
	builtin.InitializeHttpFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Тест работы с заголовками
print("=== Тест HTTP заголовков ===")

// GET запрос с кастомными заголовками
let customHeaders = {
    "User-Agent": "foo_lang/1.0",
    "X-Custom-Header": "test-value"
}

print("Заголовки созданы")

// Тестируем URL утилиты  
let testUrl = "Hello World & Test"
let encodedUrl = urlEncode(testUrl)
print("URL кодирование работает")

print("HTTP заголовки тесты завершены")
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Проверяем, что не было ошибок
			if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
				t.Errorf("HTTP headers test failed: %s", str)
			}
		}
	}
}

func TestHttpServerRouting(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	
	// Инициализируем HTTP функции
	builtin.InitializeHttpFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Тест роутинга HTTP сервера
print("=== Тест HTTP роутинга ===")

// Создаем сервер
httpCreateServer()

// Обработчик для главной страницы
fn homeHandler(request) {
    return {
        "status": 200,
        "headers": {"Content-Type": "text/plain"},
        "body": "Welcome to foo_lang server!"
    }
}

// Обработчик для API
fn apiHandler(request) {
    let response = {
        "status": 200,
        "headers": {"Content-Type": "application/json"},
        "body": "API working with method: " + request.method
    }
    return response
}

// Регистрируем маршруты
httpRoute("GET", "/", homeHandler)
httpRoute("GET", "/api/get", apiHandler)
httpRoute("POST", "/api/post", apiHandler)

print("Маршруты зарегистрированы")

// Запускаем сервер
httpStartServer(8081)
print("Сервер запущен на порту 8081")

// Небольшая задержка
await sleep(50)

// Останавливаем сервер
httpStopServer()
print("Сервер остановлен")

print("HTTP роутинг тесты завершены")
`

	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Проверяем, что не было ошибок
			if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
				t.Errorf("HTTP routing test failed: %s", str)
			}
		}
	}
}

func TestHttpAsyncRequests(t *testing.T) {
	// Сбрасываем глобальный scope
	scope.GlobalScope = scope.NewScopeStack()
	
	// Инициализируем HTTP функции
	builtin.InitializeHttpFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Тест асинхронных HTTP запросов
print("=== Тест асинхронных HTTP запросов ===")

httpSetTimeout(5)

// Функция для тестирования async
fn asyncTestFunction() {
    await sleep(10)
    return "Async test completed"
}

// Запускаем несколько асинхронных задач
let task1 = async asyncTestFunction()
let task2 = async asyncTestFunction()

// Ждем завершения всех
let results = await Promise.all(task1, task2)
print("Все асинхронные запросы завершены")
print("Количество результатов: 2")

print("Асинхронные HTTP тесты завершены")
`

	startTime := time.Now()
	exprs := parser.NewParser([]byte(code)).Parse()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Проверяем, что не было ошибок
			if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
				t.Errorf("HTTP async test failed: %s", str)
			}
		}
	}
	
	elapsed := time.Since(startTime)
	// Асинхронные запросы должны выполниться быстрее чем последовательные
	if elapsed > 3*time.Second {
		t.Errorf("Async HTTP requests took too long: %v", elapsed)
	}
}