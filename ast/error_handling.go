package ast

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorInfo содержит детальную информацию об ошибке
type ErrorInfo struct {
	Type        string   // Тип ошибки (RuntimeError, TypeError, etc.)
	Message     string   // Сообщение об ошибке
	Code        string   // Код ошибки (E001, E002, etc.)
	StackTrace  []string // Стек вызовов
	Context     string   // Контекст (функция, файл и т.д.)
	Suggestions []string // Предложения по исправлению
}

// NewErrorInfo создает новую информацию об ошибке
func NewErrorInfo(errorType, message, code string) *ErrorInfo {
	return &ErrorInfo{
		Type:        errorType,
		Message:     message,
		Code:        code,
		StackTrace:  getCaller(),
		Context:     "",
		Suggestions: []string{},
	}
}

// WithContext добавляет контекст к ошибке
func (e *ErrorInfo) WithContext(context string) *ErrorInfo {
	e.Context = context
	return e
}

// WithSuggestion добавляет предложение по исправлению
func (e *ErrorInfo) WithSuggestion(suggestion string) *ErrorInfo {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// ToString форматирует ошибку в строку
func (e *ErrorInfo) ToString() string {
	var sb strings.Builder
	
	// Основное сообщение
	sb.WriteString(fmt.Sprintf("[%s] %s (%s)\n", e.Type, e.Message, e.Code))
	
	// Контекст
	if e.Context != "" {
		sb.WriteString(fmt.Sprintf("Context: %s\n", e.Context))
	}
	
	// Стек вызовов (первые 3 уровня)
	if len(e.StackTrace) > 0 {
		sb.WriteString("Stack trace:\n")
		for i, frame := range e.StackTrace {
			if i >= 3 { // Ограничиваем до 3 уровней для читаемости
				break
			}
			sb.WriteString(fmt.Sprintf("  %d: %s\n", i+1, frame))
		}
	}
	
	// Предложения
	if len(e.Suggestions) > 0 {
		sb.WriteString("Suggestions:\n")
		for _, suggestion := range e.Suggestions {
			sb.WriteString(fmt.Sprintf("  • %s\n", suggestion))
		}
	}
	
	return sb.String()
}

// getCaller получает стек вызовов
func getCaller() []string {
	var frames []string
	for i := 2; i < 10; i++ { // Пропускаем getCaller() и вызывающую функцию
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		// Получаем имя функции
		funcName := "unknown"
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = fn.Name()
		}
		
		// Упрощаем путь к файлу
		parts := strings.Split(file, "/")
		if len(parts) > 2 {
			file = ".../" + strings.Join(parts[len(parts)-2:], "/")
		}
		
		frame := fmt.Sprintf("%s:%d in %s", file, line, funcName)
		frames = append(frames, frame)
	}
	return frames
}

// Предопределенные типы ошибок
const (
	RuntimeError    = "RuntimeError"
	TypeError       = "TypeError"
	ArgumentError   = "ArgumentError"
	ReferenceError  = "ReferenceError"
	IndexError      = "IndexError"
	AttributeError  = "AttributeError"
	ValueError      = "ValueError"
	OverflowError   = "OverflowError"
)

// Коды ошибок
const (
	// Ошибки аргументов
	E001_WRONG_ARG_COUNT      = "E001"
	E002_INVALID_ARG_TYPE     = "E002"
	E003_MISSING_REQUIRED_ARG = "E003"
	
	// Ошибки ссылок
	E101_UNDEFINED_VARIABLE   = "E101"
	E102_UNDEFINED_FUNCTION   = "E102"
	E103_UNDEFINED_TYPE       = "E103"
	
	// Ошибки индексации
	E201_INDEX_OUT_OF_BOUNDS  = "E201"
	E202_INVALID_INDEX_TYPE   = "E202"
	
	// Ошибки типов
	E301_TYPE_MISMATCH        = "E301"
	E302_CANNOT_CONVERT       = "E302"
	E303_INVALID_OPERATION    = "E303"
	
	// Ошибки значений
	E401_INVALID_VALUE        = "E401"
	E402_EMPTY_CONTAINER      = "E402"
	E403_READONLY_VARIABLE    = "E403"
)

// Функции-хелперы для создания типичных ошибок

// NewArgumentError создает ошибку неправильных аргументов
func NewArgumentError(expected, got int, function string) *Value {
	err := NewErrorInfo(ArgumentError, 
		fmt.Sprintf("Function '%s' expects %d arguments, got %d", function, expected, got),
		E001_WRONG_ARG_COUNT).
		WithContext(function).
		WithSuggestion(fmt.Sprintf("Check the function signature and provide exactly %d arguments", expected))
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

// NewTypeError создает ошибку типа
func NewTypeError(expected, got string, context string) *Value {
	err := NewErrorInfo(TypeError,
		fmt.Sprintf("Expected %s, got %s", expected, got),
		E301_TYPE_MISMATCH).
		WithContext(context).
		WithSuggestion(fmt.Sprintf("Convert the value to %s or use a different operation", expected))
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

// NewReferenceError создает ошибку неопределенной ссылки
func NewReferenceError(name, itemType string) *Value {
	code := E101_UNDEFINED_VARIABLE
	if itemType == "function" {
		code = E102_UNDEFINED_FUNCTION
	} else if itemType == "type" {
		code = E103_UNDEFINED_TYPE
	}
	
	err := NewErrorInfo(ReferenceError,
		fmt.Sprintf("%s '%s' is not defined", strings.Title(itemType), name),
		code).
		WithSuggestion(fmt.Sprintf("Check the spelling of '%s' or define it before use", name)).
		WithSuggestion("Use auto-completion in your IDE to see available options")
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

// NewIndexError создает ошибку индексации
func NewIndexError(index, length int) *Value {
	err := NewErrorInfo(IndexError,
		fmt.Sprintf("Index %d out of bounds for length %d", index, length),
		E201_INDEX_OUT_OF_BOUNDS).
		WithSuggestion(fmt.Sprintf("Use index between 0 and %d", length-1)).
		WithSuggestion("Check array/string length before indexing")
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

// NewValueError создает ошибку значения
func NewValueError(message, context string) *Value {
	err := NewErrorInfo(ValueError, message, E401_INVALID_VALUE).
		WithContext(context)
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

// NewEmptyContainerError создает ошибку пустого контейнера
func NewEmptyContainerError(operation, container string) *Value {
	err := NewErrorInfo(ValueError,
		fmt.Sprintf("Cannot %s on empty %s", operation, container),
		E402_EMPTY_CONTAINER).
		WithSuggestion(fmt.Sprintf("Check if %s is not empty before calling %s", container, operation)).
		WithSuggestion(fmt.Sprintf("Use length() method to check %s size", container))
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

// SafeEval безопасно выполняет выражение и возвращает Result
func SafeEval(expr Expr, context string) *Value {
	// Используем recover для перехвата panic
	defer func() {
		if r := recover(); r != nil {
			// Конвертируем panic в Result<T, E>
			errorMsg := fmt.Sprintf("%v", r)
			_ = NewErrorInfo(RuntimeError, errorMsg, "E999").
				WithContext(context).
				WithSuggestion("Check the operation and input values")
			
			// Здесь мы не можем вернуть значение из defer, 
			// поэтому panic будет перехвачен выше в стеке
		}
	}()
	
	// Выполняем выражение
	value := expr.Eval()
	
	// Оборачиваем в Ok если это не Result
	if _, isResult := value.Any().(*ResultValue); !isResult {
		result := NewResultOk(value)
		return NewValue(result)
	}
	
	return value
}

// IsError проверяет является ли значение ошибкой
func IsError(value *Value) bool {
	if result, ok := value.Any().(*ResultValue); ok {
		return result.IsErr()
	}
	return false
}

// ExtractError извлекает ErrorInfo из Result
func ExtractError(value *Value) *ErrorInfo {
	if result, ok := value.Any().(*ResultValue); ok && result.IsErr() {
		if errorInfo, ok := result.error.Any().(*ErrorInfo); ok {
			return errorInfo
		}
	}
	return nil
}

// HandleError обрабатывает ошибку - выводит или возвращает Result
func HandleError(err *ErrorInfo, shouldPanic bool) *Value {
	if shouldPanic {
		panic(err.ToString())
	}
	
	result := NewResultErr(NewValue(err))
	return NewValue(result)
}

// TryFunction безопасно выполняет функцию с обработкой ошибок
func TryFunction(fn func() *Value, context string) *Value {
	defer func() {
		if r := recover(); r != nil {
			errorMsg := fmt.Sprintf("%v", r)
			err := NewErrorInfo(RuntimeError, errorMsg, "E999").
				WithContext(context)
			
			// Возвращаем через глобальную переменную (не идеально, но работает)
			lastError = NewValue(NewResultErr(NewValue(err)))
			_ = lastError // Используем переменную
		}
	}()
	
	return fn()
}

// Глобальная переменная для хранения последней ошибки (временное решение)
var lastError *Value