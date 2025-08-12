package scope

import (
	"fmt"
	"foo_lang/value"
	"sync"
)

// Scope представляет область видимости переменных
type Scope struct {
	parent *Scope
	vars   map[string]*value.Value
	mu     sync.RWMutex
}

// NewScope создает новую область видимости
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent: parent,
		vars:   make(map[string]*value.Value),
	}
}

// Set устанавливает переменную в текущей области
func (s *Scope) Set(name string, val *value.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vars[name] = val
}

// Get получает переменную, ищет в текущей области и родительских
func (s *Scope) Get(name string) (*value.Value, bool) {
	s.mu.RLock()
	val, exists := s.vars[name]
	s.mu.RUnlock()
	
	if exists {
		return val, true
	}
	if s.parent != nil {
		return s.parent.Get(name)
	}
	return nil, false
}

// Has проверяет, существует ли переменная в текущей области (не в родительских)
func (s *Scope) Has(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.vars[name]
	return exists
}

// Update обновляет существующую переменную (ищет в текущей и родительских областях)
func (s *Scope) Update(name string, val *value.Value) bool {
	s.mu.Lock()
	_, exists := s.vars[name]
	if exists {
		s.vars[name] = val
		s.mu.Unlock()
		return true
	}
	s.mu.Unlock()
	
	if s.parent != nil {
		return s.parent.Update(name, val)
	}
	return false
}

// ScopeStack - стек областей видимости
type ScopeStack struct {
	current        *Scope
	recursionDepth int
	maxRecursion   int
}

// NewScopeStack создает новый стек с глобальной областью
func NewScopeStack() *ScopeStack {
	return &ScopeStack{
		current:        NewScope(nil), // глобальная область
		recursionDepth: 0,
		maxRecursion:   1000, // Максимальная глубина рекурсии
	}
}

// Push создает новую локальную область
func (ss *ScopeStack) Push() {
	ss.current = NewScope(ss.current)
}

// Pop удаляет текущую локальную область
func (ss *ScopeStack) Pop() {
	if ss.current.parent != nil {
		ss.current = ss.current.parent
	}
}

// Set устанавливает переменную в текущей области
func (ss *ScopeStack) Set(name string, val *value.Value) {
	ss.current.Set(name, val)
}

// Get получает переменную
func (ss *ScopeStack) Get(name string) (*value.Value, bool) {
	return ss.current.Get(name)
}

// Has проверяет существование в текущей области
func (ss *ScopeStack) Has(name string) bool {
	return ss.current.Has(name)
}

// Update обновляет существующую переменную
func (ss *ScopeStack) Update(name string, val *value.Value) bool {
	return ss.current.Update(name, val)
}

// PushFunction увеличивает счетчик рекурсии и создает новую область
func (ss *ScopeStack) PushFunction() error {
	ss.recursionDepth++
	if ss.recursionDepth > ss.maxRecursion {
		return fmt.Errorf("maximum recursion depth exceeded (%d)", ss.maxRecursion)
	}
	ss.Push()
	return nil
}

// PopFunction уменьшает счетчик рекурсии и удаляет область
func (ss *ScopeStack) PopFunction() {
	ss.Pop()
	if ss.recursionDepth > 0 {
		ss.recursionDepth--
	}
}

// SetMaxRecursion устанавливает максимальную глубину рекурсии
func (ss *ScopeStack) SetMaxRecursion(max int) {
	ss.maxRecursion = max
}

// GetRecursionDepth возвращает текущую глубину рекурсии
func (ss *ScopeStack) GetRecursionDepth() int {
	return ss.recursionDepth
}

// GetAll возвращает все переменные из всех областей видимости (для экспорта модулей)
func (ss *ScopeStack) GetAll() map[string]*value.Value {
	result := make(map[string]*value.Value)
	
	// Собираем все переменные, начиная с глобальной области
	scope := ss.current
	for scope != nil {
		scope.mu.RLock()
		for name, val := range scope.vars {
			// Переменные из более локальных областей имеют приоритет
			if _, exists := result[name]; !exists {
				result[name] = val
			}
		}
		scope.mu.RUnlock()
		scope = scope.parent
	}
	
	return result
}

// GlobalScope - глобальный стек областей видимости
// ВАЖНО: Эта переменная инициализируется парсером при создании
var GlobalScope *ScopeStack