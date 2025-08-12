package ast

import (
	"foo_lang/scope"
	"foo_lang/value"
	"sync"
	"time"
)

// AsyncExpr представляет async выражение для запуска горутины
type AsyncExpr struct {
	Expr Expr
}

// Карта scope для каждой горутины - нужно для избежания race conditions
var goroutineScopes = make(map[int64]*scope.ScopeStack)
var goroutineCounter int64
var scopeMutex sync.RWMutex

func (a *AsyncExpr) Eval() *value.Value {
	promise := value.NewPromise()
	
	// Захватываем текущее состояние scope (все переменные)
	capturedVars := scope.GlobalScope.GetAll()
	
	// Увеличиваем счетчик горутин
	scopeMutex.Lock()
	goroutineCounter++
	goroutineID := goroutineCounter
	scopeMutex.Unlock()
	
	// Запускаем выполнение в отдельной горутине
	go func() {
		// Восстанавливаемся от паник
		defer func() {
			// Очищаем scope для этой горутины
			scopeMutex.Lock()
			delete(goroutineScopes, goroutineID)
			scopeMutex.Unlock()
			
			if r := recover(); r != nil {
				// Преобразуем панику в ошибку
				var errMsg string
				switch v := r.(type) {
				case string:
					errMsg = v
				case error:
					errMsg = v.Error()
				default:
					errMsg = "Unknown error"
				}
				promise.Reject(value.NewValue(errMsg))
			}
		}()
		
		// Создаем изолированный scope для этой горутины
		isolatedScope := scope.NewScopeStack()
		
		// Восстанавливаем все захваченные переменные в изолированном scope
		for name, val := range capturedVars {
			isolatedScope.Set(name, val)
		}
		
		// Сохраняем scope для этой горутины
		scopeMutex.Lock()
		goroutineScopes[goroutineID] = isolatedScope
		originalScope := scope.GlobalScope
		scope.GlobalScope = isolatedScope
		scopeMutex.Unlock()
		
		// Гарантируем восстановление scope в defer
		defer func() {
			scopeMutex.Lock()
			scope.GlobalScope = originalScope
			scopeMutex.Unlock()
		}()
		
		// Выполняем выражение в изолированном scope
		result := a.Expr.Eval()
		
		// Проверяем специальные флаги
		if result != nil && (result.IsReturn() || result.IsBreak()) {
			// Эти флаги не должны распространяться через async границы
			// Создаем новое значение без флагов
			cleanResult := value.NewValue(result.Any())
			result = cleanResult
		}
		
		// Резолвим промис с результатом
		promise.Resolve(result)
	}()
	
	// Возвращаем промис
	return value.NewValue(promise)
}

// AwaitExpr представляет await выражение для ожидания промиса
type AwaitExpr struct {
	Expr Expr
}

func (a *AwaitExpr) Eval() *value.Value {
	// Вычисляем выражение
	result := a.Expr.Eval()
	
	// Проверяем, является ли результат промисом
	if promise, ok := result.Any().(*value.Promise); ok {
		// Ожидаем завершения промиса
		promise.Wait()
		
		// Возвращаем результат или ошибку
		if promise.GetState() == value.PromiseFulfilled {
			return promise.GetValue()
		} else {
			// Если промис отклонен, паникуем с ошибкой
			errValue := promise.GetError()
			if errValue != nil {
				panic(errValue.Any())
			}
			panic("Promise rejected without error")
		}
	}
	
	// Если это не промис, возвращаем как есть
	return result
}

// PromiseAllExpr представляет Promise.all() выражение
type PromiseAllExpr struct {
	Args []Expr
}

func (p *PromiseAllExpr) Eval() *value.Value {
	// Вычисляем все аргументы
	promises := make([]*value.Promise, 0)
	
	for _, arg := range p.Args {
		result := arg.Eval()
		if promise, ok := result.Any().(*value.Promise); ok {
			promises = append(promises, promise)
		} else {
			// Если аргумент не промис, создаем resolved промис
			immediatePromise := value.NewPromise()
			immediatePromise.Resolve(result)
			promises = append(promises, immediatePromise)
		}
	}
	
	// Вызываем PromiseAll
	resultPromise := value.PromiseAll(promises)
	return value.NewValue(resultPromise)
}

// PromiseAnyExpr представляет Promise.any() выражение
type PromiseAnyExpr struct {
	Args []Expr
}

func (p *PromiseAnyExpr) Eval() *value.Value {
	// Вычисляем все аргументы
	promises := make([]*value.Promise, 0)
	
	for _, arg := range p.Args {
		result := arg.Eval()
		if promise, ok := result.Any().(*value.Promise); ok {
			promises = append(promises, promise)
		} else {
			// Если аргумент не промис, создаем resolved промис
			immediatePromise := value.NewPromise()
			immediatePromise.Resolve(result)
			promises = append(promises, immediatePromise)
		}
	}
	
	// Вызываем PromiseAny
	resultPromise := value.PromiseAny(promises)
	return value.NewValue(resultPromise)
}

// SleepExpr представляет функцию sleep для задержки
type SleepExpr struct {
	Duration Expr
}

func (s *SleepExpr) Eval() *value.Value {
	// Вычисляем длительность
	durationValue := s.Duration.Eval()
	
	// Преобразуем в миллисекунды
	var ms int64
	switch v := durationValue.Any().(type) {
	case int64:
		ms = v
	case float64:
		ms = int64(v)
	default:
		panic("sleep() requires a number argument (milliseconds)")
	}
	
	// Создаем промис
	promise := value.NewPromise()
	
	// Запускаем таймер в горутине
	go func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		promise.Resolve(value.NewValue(nil))
	}()
	
	return value.NewValue(promise)
}