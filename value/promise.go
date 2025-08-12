package value

import (
	"sync"
	"time"
)

// PromiseState представляет состояние промиса
type PromiseState int

const (
	PromisePending PromiseState = iota
	PromiseFulfilled
	PromiseRejected
)

// Promise представляет асинхронную операцию
type Promise struct {
	state    PromiseState
	value    *Value
	error    *Value
	mu       sync.Mutex
	done     chan struct{}
	handlers []func()
}

// NewPromise создает новый промис
func NewPromise() *Promise {
	return &Promise{
		state:    PromisePending,
		done:     make(chan struct{}),
		handlers: make([]func(), 0),
	}
}

// Resolve выполняет промис с результатом
func (p *Promise) Resolve(value *Value) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.state != PromisePending {
		return
	}
	
	p.state = PromiseFulfilled
	p.value = value
	close(p.done)
	
	// Выполняем все обработчики
	for _, handler := range p.handlers {
		go handler()
	}
}

// Reject отклоняет промис с ошибкой
func (p *Promise) Reject(err *Value) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.state != PromisePending {
		return
	}
	
	p.state = PromiseRejected
	p.error = err
	close(p.done)
	
	// Выполняем все обработчики
	for _, handler := range p.handlers {
		go handler()
	}
}

// Then добавляет обработчик для результата
func (p *Promise) Then(handler func()) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.state != PromisePending {
		// Если промис уже выполнен, запускаем обработчик сразу
		go handler()
	} else {
		p.handlers = append(p.handlers, handler)
	}
}

// Wait ожидает завершения промиса
func (p *Promise) Wait() {
	<-p.done
}

// WaitWithTimeout ожидает завершения с таймаутом
func (p *Promise) WaitWithTimeout(timeout time.Duration) bool {
	select {
	case <-p.done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// GetState возвращает текущее состояние
func (p *Promise) GetState() PromiseState {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.state
}

// GetValue возвращает значение (если fulfilled)
func (p *Promise) GetValue() *Value {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value
}

// GetError возвращает ошибку (если rejected)
func (p *Promise) GetError() *Value {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.error
}

// PromiseAll ожидает завершения всех промисов
func PromiseAll(promises []*Promise) *Promise {
	result := NewPromise()
	
	if len(promises) == 0 {
		result.Resolve(NewValue([]interface{}{}))
		return result
	}
	
	go func() {
		values := make([]*Value, len(promises))
		var wg sync.WaitGroup
		wg.Add(len(promises))
		
		hasError := false
		var firstError *Value
		var errorMu sync.Mutex
		
		for i, promise := range promises {
			go func(idx int, p *Promise) {
				defer wg.Done()
				p.Wait()
				
				errorMu.Lock()
				defer errorMu.Unlock()
				
				if p.GetState() == PromiseRejected {
					if !hasError {
						hasError = true
						firstError = p.GetError()
					}
				} else {
					values[idx] = p.GetValue()
				}
			}(i, promise)
		}
		
		wg.Wait()
		
		if hasError {
			result.Reject(firstError)
		} else {
			// Конвертируем в массив interface{}
			interfaceValues := make([]interface{}, len(values))
			for i, v := range values {
				if v != nil {
					interfaceValues[i] = v.Any()
				}
			}
			result.Resolve(NewValue(interfaceValues))
		}
	}()
	
	return result
}

// PromiseAny ожидает завершения любого промиса
func PromiseAny(promises []*Promise) *Promise {
	result := NewPromise()
	
	if len(promises) == 0 {
		result.Reject(NewValue("No promises provided"))
		return result
	}
	
	go func() {
		done := make(chan struct{})
		var once sync.Once
		
		errorCount := 0
		var errorMu sync.Mutex
		errors := make([]*Value, len(promises))
		
		for i, promise := range promises {
			go func(idx int, p *Promise) {
				p.Wait()
				
				if p.GetState() == PromiseFulfilled {
					once.Do(func() {
						result.Resolve(p.GetValue())
						close(done)
					})
				} else {
					errorMu.Lock()
					errors[idx] = p.GetError()
					errorCount++
					if errorCount == len(promises) {
						// Все промисы отклонены
						once.Do(func() {
							result.Reject(NewValue("All promises were rejected"))
							close(done)
						})
					}
					errorMu.Unlock()
				}
			}(i, promise)
		}
		
		<-done
	}()
	
	return result
}