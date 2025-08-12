package builtin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"foo_lang/value"
)

// HTTP клиент с настройками
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// HTTP сервер экземпляр
var httpServer *http.Server
var serverMux *http.ServeMux

// HttpFunction представляет HTTP функцию
type HttpFunction struct {
	name string
	fn   func([]*value.Value) *value.Value
}

func (hf *HttpFunction) Eval() *value.Value {
	return value.NewValue(hf)
}

func (hf *HttpFunction) Call(args []*value.Value) *value.Value {
	return hf.fn(args)
}

func (hf *HttpFunction) String() string {
	return "builtin function " + hf.name
}

func (hf *HttpFunction) Name() string {
	return hf.name
}

func InitializeHttpFunctions(scopeStack ScopeStack) {
	// HTTP клиент функции
	scopeStack.Set("httpGet", value.NewValue(&HttpFunction{
		name: "httpGet",
		fn: HttpGet,
	}))
	
	scopeStack.Set("httpPost", value.NewValue(&HttpFunction{
		name: "httpPost",
		fn: HttpPost,
	}))
	
	scopeStack.Set("httpPut", value.NewValue(&HttpFunction{
		name: "httpPut",
		fn: HttpPut,
	}))
	
	scopeStack.Set("httpDelete", value.NewValue(&HttpFunction{
		name: "httpDelete",
		fn: HttpDelete,
	}))
	
	// HTTP сервер функции
	scopeStack.Set("httpCreateServer", value.NewValue(&HttpFunction{
		name: "httpCreateServer",
		fn: HttpCreateServer,
	}))
	
	scopeStack.Set("httpRoute", value.NewValue(&HttpFunction{
		name: "httpRoute",
		fn: HttpRoute,
	}))
	
	scopeStack.Set("httpStartServer", value.NewValue(&HttpFunction{
		name: "httpStartServer",
		fn: HttpStartServer,
	}))
	
	scopeStack.Set("httpStopServer", value.NewValue(&HttpFunction{
		name: "httpStopServer",
		fn: HttpStopServer,
	}))
	
	// Утилиты
	scopeStack.Set("httpSetTimeout", value.NewValue(&HttpFunction{
		name: "httpSetTimeout",
		fn: HttpSetTimeout,
	}))
	
	scopeStack.Set("urlEncode", value.NewValue(&HttpFunction{
		name: "urlEncode",
		fn: UrlEncode,
	}))
	
	scopeStack.Set("urlDecode", value.NewValue(&HttpFunction{
		name: "urlDecode",
		fn: UrlDecode,
	}))
}

// HttpGet выполняет HTTP GET запрос
// Использование: httpGet(url, [headers])
func HttpGet(args []*value.Value) *value.Value {
	if len(args) < 1 {
		return value.NewValue("Error: httpGet() requires at least 1 argument (url)")
	}
	
	url, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("Error: first argument must be a string (url)")
	}
	
	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to create request: %v", err))
	}
	
	// Добавляем заголовки если переданы
	if len(args) >= 2 {
		if headers, ok := args[1].Any().(map[string]*value.Value); ok {
			for key, val := range headers {
				if headerValue, ok := val.Any().(string); ok {
					req.Header.Set(key, headerValue)
				}
			}
		}
	}
	
	// Выполняем запрос
	resp, err := httpClient.Do(req)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: request failed: %v", err))
	}
	defer resp.Body.Close()
	
	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to read response: %v", err))
	}
	
	// Создаем объект ответа
	response := map[string]*value.Value{
		"status":     value.NewValue(int64(resp.StatusCode)),
		"statusText": value.NewValue(resp.Status),
		"body":       value.NewValue(string(body)),
		"headers":    headersToMap(resp.Header),
	}
	
	return value.NewValue(response)
}

// HttpPost выполняет HTTP POST запрос
// Использование: httpPost(url, body, [headers])
func HttpPost(args []*value.Value) *value.Value {
	if len(args) < 2 {
		return value.NewValue("Error: httpPost() requires at least 2 arguments (url, body)")
	}
	
	url, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("Error: first argument must be a string (url)")
	}
	
	// Подготавливаем тело запроса
	var bodyReader io.Reader
	var contentType string
	
	switch body := args[1].Any().(type) {
	case string:
		bodyReader = strings.NewReader(body)
		contentType = "text/plain"
	case map[string]*value.Value:
		// JSON объект
		jsonData, err := valueMapToJSON(body)
		if err != nil {
			return value.NewValue(fmt.Sprintf("Error: failed to encode JSON: %v", err))
		}
		bodyReader = bytes.NewReader(jsonData)
		contentType = "application/json"
	default:
		return value.NewValue("Error: body must be a string or object")
	}
	
	// Создаем запрос
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to create request: %v", err))
	}
	
	// Устанавливаем Content-Type
	req.Header.Set("Content-Type", contentType)
	
	// Добавляем заголовки если переданы
	if len(args) >= 3 {
		if headers, ok := args[2].Any().(map[string]*value.Value); ok {
			for key, val := range headers {
				if headerValue, ok := val.Any().(string); ok {
					req.Header.Set(key, headerValue)
				}
			}
		}
	}
	
	// Выполняем запрос
	resp, err := httpClient.Do(req)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: request failed: %v", err))
	}
	defer resp.Body.Close()
	
	// Читаем тело ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to read response: %v", err))
	}
	
	// Создаем объект ответа
	response := map[string]*value.Value{
		"status":     value.NewValue(int64(resp.StatusCode)),
		"statusText": value.NewValue(resp.Status),
		"body":       value.NewValue(string(responseBody)),
		"headers":    headersToMap(resp.Header),
	}
	
	return value.NewValue(response)
}

// HttpPut выполняет HTTP PUT запрос
func HttpPut(args []*value.Value) *value.Value {
	if len(args) < 2 {
		return value.NewValue("Error: httpPut() requires at least 2 arguments (url, body)")
	}
	
	return httpMethodWithBody("PUT", args)
}

// HttpDelete выполняет HTTP DELETE запрос  
func HttpDelete(args []*value.Value) *value.Value {
	if len(args) < 1 {
		return value.NewValue("Error: httpDelete() requires at least 1 argument (url)")
	}
	
	url, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("Error: first argument must be a string (url)")
	}
	
	// Создаем запрос
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to create request: %v", err))
	}
	
	// Добавляем заголовки если переданы
	if len(args) >= 2 {
		if headers, ok := args[1].Any().(map[string]*value.Value); ok {
			for key, val := range headers {
				if headerValue, ok := val.Any().(string); ok {
					req.Header.Set(key, headerValue)
				}
			}
		}
	}
	
	// Выполняем запрос
	resp, err := httpClient.Do(req)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: request failed: %v", err))
	}
	defer resp.Body.Close()
	
	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to read response: %v", err))
	}
	
	// Создаем объект ответа
	response := map[string]*value.Value{
		"status":     value.NewValue(int64(resp.StatusCode)),
		"statusText": value.NewValue(resp.Status),
		"body":       value.NewValue(string(body)),
		"headers":    headersToMap(resp.Header),
	}
	
	return value.NewValue(response)
}

// HttpCreateServer создает новый HTTP сервер
func HttpCreateServer(args []*value.Value) *value.Value {
	// Инициализируем сервер если еще не создан
	if serverMux == nil {
		serverMux = http.NewServeMux()
	}
	
	return value.NewValue("HTTP server created")
}

// HttpRoute добавляет маршрут к HTTP серверу
// Использование: httpRoute(method, path, handler)
func HttpRoute(args []*value.Value) *value.Value {
	if len(args) < 3 {
		return value.NewValue("Error: httpRoute() requires 3 arguments (method, path, handler)")
	}
	
	method, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("Error: first argument must be a string (method)")
	}
	
	path, ok := args[1].Any().(string)
	if !ok {
		return value.NewValue("Error: second argument must be a string (path)")
	}
	
	handler := args[2]
	
	if serverMux == nil {
		serverMux = http.NewServeMux()
	}
	
	// Создаем обработчик
	serverMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод
		if strings.ToUpper(r.Method) != strings.ToUpper(method) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// Читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		
		// Создаем объект запроса
		request := map[string]*value.Value{
			"method":  value.NewValue(r.Method),
			"path":    value.NewValue(r.URL.Path),
			"query":   queryToMap(r.URL.Query()),
			"headers": headersToMap(r.Header),
			"body":    value.NewValue(string(body)),
		}
		
		// Вызываем обработчик (это должна быть функция)
		if callable, ok := handler.Any().(*HttpFunction); ok {
			response := callable.Call([]*value.Value{value.NewValue(request)})
			
			// Обрабатываем ответ
			if responseObj, ok := response.Any().(map[string]*value.Value); ok {
				// Устанавливаем статус
				if status, exists := responseObj["status"]; exists {
					if statusCode, ok := status.Any().(int64); ok {
						w.WriteHeader(int(statusCode))
					}
				}
				
				// Устанавливаем заголовки
				if headers, exists := responseObj["headers"]; exists {
					if headerMap, ok := headers.Any().(map[string]*value.Value); ok {
						for key, val := range headerMap {
							if headerValue, ok := val.Any().(string); ok {
								w.Header().Set(key, headerValue)
							}
						}
					}
				}
				
				// Отправляем тело ответа
				if body, exists := responseObj["body"]; exists {
					if bodyStr, ok := body.Any().(string); ok {
						w.Write([]byte(bodyStr))
					}
				}
			} else {
				// Простой текстовый ответ
				if responseStr, ok := response.Any().(string); ok {
					w.Write([]byte(responseStr))
				}
			}
		} else {
			http.Error(w, "Handler is not a function", http.StatusInternalServerError)
		}
	})
	
	return value.NewValue(fmt.Sprintf("Route %s %s registered", method, path))
}

// HttpStartServer запускает HTTP сервер
// Использование: httpStartServer(port)
func HttpStartServer(args []*value.Value) *value.Value {
	if len(args) < 1 {
		return value.NewValue("Error: httpStartServer() requires 1 argument (port)")
	}
	
	port, ok := args[0].Any().(int64)
	if !ok {
		return value.NewValue("Error: port must be a number")
	}
	
	if serverMux == nil {
		return value.NewValue("Error: no server created, call httpCreateServer() first")
	}
	
	// Создаем сервер
	httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: serverMux,
	}
	
	// Запускаем сервер в горутине
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()
	
	return value.NewValue(fmt.Sprintf("HTTP server started on port %d", port))
}

// HttpStopServer останавливает HTTP сервер
func HttpStopServer(args []*value.Value) *value.Value {
	if httpServer == nil {
		return value.NewValue("Error: no server running")
	}
	
	// Даем серверу 5 секунд на graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := httpServer.Shutdown(ctx); err != nil {
		return value.NewValue(fmt.Sprintf("Error stopping server: %v", err))
	}
	
	httpServer = nil
	return value.NewValue("HTTP server stopped")
}

// HttpSetTimeout устанавливает таймаут для HTTP клиента
func HttpSetTimeout(args []*value.Value) *value.Value {
	if len(args) < 1 {
		return value.NewValue("Error: httpSetTimeout() requires 1 argument (seconds)")
	}
	
	seconds, ok := args[0].Any().(int64)
	if !ok {
		return value.NewValue("Error: timeout must be a number (seconds)")
	}
	
	httpClient.Timeout = time.Duration(seconds) * time.Second
	return value.NewValue(fmt.Sprintf("HTTP timeout set to %d seconds", seconds))
}

// UrlEncode кодирует строку для URL
func UrlEncode(args []*value.Value) *value.Value {
	if len(args) < 1 {
		return value.NewValue("Error: urlEncode() requires 1 argument (string)")
	}
	
	str, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("Error: argument must be a string")
	}
	
	encoded := url.QueryEscape(str)
	return value.NewValue(encoded)
}

// UrlDecode декодирует URL строку
func UrlDecode(args []*value.Value) *value.Value {
	if len(args) < 1 {
		return value.NewValue("Error: urlDecode() requires 1 argument (string)")
	}
	
	str, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("Error: argument must be a string")
	}
	
	decoded, err := url.QueryUnescape(str)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to decode URL: %v", err))
	}
	
	return value.NewValue(decoded)
}

// Вспомогательные функции

// httpMethodWithBody обрабатывает HTTP методы с телом запроса
func httpMethodWithBody(method string, args []*value.Value) *value.Value {
	url, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("Error: first argument must be a string (url)")
	}
	
	// Подготавливаем тело запроса
	var bodyReader io.Reader
	var contentType string
	
	switch body := args[1].Any().(type) {
	case string:
		bodyReader = strings.NewReader(body)
		contentType = "text/plain"
	case map[string]*value.Value:
		// JSON объект
		jsonData, err := valueMapToJSON(body)
		if err != nil {
			return value.NewValue(fmt.Sprintf("Error: failed to encode JSON: %v", err))
		}
		bodyReader = bytes.NewReader(jsonData)
		contentType = "application/json"
	default:
		return value.NewValue("Error: body must be a string or object")
	}
	
	// Создаем запрос
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to create request: %v", err))
	}
	
	// Устанавливаем Content-Type
	req.Header.Set("Content-Type", contentType)
	
	// Добавляем заголовки если переданы
	if len(args) >= 3 {
		if headers, ok := args[2].Any().(map[string]*value.Value); ok {
			for key, val := range headers {
				if headerValue, ok := val.Any().(string); ok {
					req.Header.Set(key, headerValue)
				}
			}
		}
	}
	
	// Выполняем запрос
	resp, err := httpClient.Do(req)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: request failed: %v", err))
	}
	defer resp.Body.Close()
	
	// Читаем тело ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to read response: %v", err))
	}
	
	// Создаем объект ответа
	response := map[string]*value.Value{
		"status":     value.NewValue(int64(resp.StatusCode)),
		"statusText": value.NewValue(resp.Status),
		"body":       value.NewValue(string(responseBody)),
		"headers":    headersToMap(resp.Header),
	}
	
	return value.NewValue(response)
}

// headersToMap конвертирует HTTP заголовки в map
func headersToMap(headers http.Header) *value.Value {
	result := make(map[string]*value.Value)
	
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = value.NewValue(values[0])
		}
	}
	
	return value.NewValue(result)
}

// queryToMap конвертирует URL query в map
func queryToMap(query url.Values) *value.Value {
	result := make(map[string]*value.Value)
	
	for key, values := range query {
		if len(values) > 0 {
			result[key] = value.NewValue(values[0])
		}
	}
	
	return value.NewValue(result)
}

// valueMapToJSON конвертирует map[string]*value.Value в JSON
func valueMapToJSON(valueMap map[string]*value.Value) ([]byte, error) {
	// Конвертируем в обычный map для JSON marshalling
	normalMap := make(map[string]interface{})
	
	for key, val := range valueMap {
		normalMap[key] = val.Any()
	}
	
	return json.Marshal(normalMap)
}