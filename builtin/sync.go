package builtin

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"sync"
	"sync/atomic"
	"time"
)

// Глобальные карты для хранения примитивов синхронизации
var (
	mutexes      = make(map[string]*sync.Mutex)
	rwMutexes    = make(map[string]*sync.RWMutex)
	waitGroups   = make(map[string]*sync.WaitGroup)
	semaphores   = make(map[string]chan struct{})
	atomicInts   = make(map[string]*int64)
	conditions   = make(map[string]*sync.Cond)
	barriers     = make(map[string]*barrier)
	syncMapsMu   sync.Mutex
	mutexCounter int
)

// barrier реализует барьер синхронизации
type barrier struct {
	n       int
	count   int
	mu      sync.Mutex
	cond    *sync.Cond
	broken  bool
}

func newBarrier(n int) *barrier {
	b := &barrier{n: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *barrier) wait() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if b.broken {
		return fmt.Errorf("barrier is broken")
	}
	
	b.count++
	if b.count == b.n {
		// Последний поток достиг барьера - освобождаем всех
		b.count = 0
		b.cond.Broadcast()
	} else {
		// Ждем остальных
		b.cond.Wait()
	}
	
	return nil
}

// InitializeSyncFunctions инициализирует встроенные функции синхронизации
func InitializeSyncFunctions(globalScope *scope.ScopeStack) {

	// ============ МЬЮТЕКСЫ ============

	// newMutex - создает новый мьютекс
	newMutexFunc := func(args []*value.Value) *value.Value {
		name := ""
		if len(args) == 1 {
			if n, ok := args[0].Any().(string); ok {
				name = n
			} else {
				return value.NewString("Error: newMutex() optional argument must be string (name)")
			}
		} else if len(args) > 1 {
			return value.NewString("Error: newMutex() requires 0-1 arguments ([name])")
		}
		
		// Генерируем имя если не предоставлено
		if name == "" {
			syncMapsMu.Lock()
			mutexCounter++
			name = fmt.Sprintf("mutex_%d", mutexCounter)
			syncMapsMu.Unlock()
		}
		
		syncMapsMu.Lock()
		if _, exists := mutexes[name]; exists {
			syncMapsMu.Unlock()
			return value.NewString(fmt.Sprintf("Error: mutex '%s' already exists", name))
		}
		mutexes[name] = &sync.Mutex{}
		syncMapsMu.Unlock()
		
		return value.NewString(name)
	}
	globalScope.Set("newMutex", value.NewValue(newMutexFunc))

	// mutexLock - блокирует мьютекс
	mutexLockFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: mutexLock() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: mutexLock() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		mutex, exists := mutexes[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: mutex '%s' does not exist", name))
		}
		
		mutex.Lock()
		return value.NewBool(true)
	}
	globalScope.Set("mutexLock", value.NewValue(mutexLockFunc))

	// mutexUnlock - разблокирует мьютекс
	mutexUnlockFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: mutexUnlock() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: mutexUnlock() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		mutex, exists := mutexes[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: mutex '%s' does not exist", name))
		}
		
		mutex.Unlock()
		return value.NewBool(true)
	}
	globalScope.Set("mutexUnlock", value.NewValue(mutexUnlockFunc))

	// mutexTryLock - пытается заблокировать мьютекс без ожидания
	mutexTryLockFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: mutexTryLock() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: mutexTryLock() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		mutex, exists := mutexes[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: mutex '%s' does not exist", name))
		}
		
		// Go мьютексы не имеют TryLock в стандартной библиотеке до Go 1.18
		// Эмулируем через канал
		locked := make(chan bool, 1)
		go func() {
			mutex.Lock()
			locked <- true
		}()
		
		select {
		case <-locked:
			return value.NewBool(true)
		case <-time.After(1 * time.Millisecond):
			return value.NewBool(false)
		}
	}
	globalScope.Set("mutexTryLock", value.NewValue(mutexTryLockFunc))

	// ============ READ-WRITE МЬЮТЕКСЫ ============

	// newRWMutex - создает новый read-write мьютекс
	newRWMutexFunc := func(args []*value.Value) *value.Value {
		name := ""
		if len(args) == 1 {
			if n, ok := args[0].Any().(string); ok {
				name = n
			} else {
				return value.NewString("Error: newRWMutex() optional argument must be string (name)")
			}
		} else if len(args) > 1 {
			return value.NewString("Error: newRWMutex() requires 0-1 arguments ([name])")
		}
		
		// Генерируем имя если не предоставлено
		if name == "" {
			syncMapsMu.Lock()
			mutexCounter++
			name = fmt.Sprintf("rwmutex_%d", mutexCounter)
			syncMapsMu.Unlock()
		}
		
		syncMapsMu.Lock()
		if _, exists := rwMutexes[name]; exists {
			syncMapsMu.Unlock()
			return value.NewString(fmt.Sprintf("Error: rwmutex '%s' already exists", name))
		}
		rwMutexes[name] = &sync.RWMutex{}
		syncMapsMu.Unlock()
		
		return value.NewString(name)
	}
	globalScope.Set("newRWMutex", value.NewValue(newRWMutexFunc))

	// rwMutexRLock - блокирует мьютекс для чтения
	rwMutexRLockFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: rwMutexRLock() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: rwMutexRLock() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		mutex, exists := rwMutexes[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: rwmutex '%s' does not exist", name))
		}
		
		mutex.RLock()
		return value.NewBool(true)
	}
	globalScope.Set("rwMutexRLock", value.NewValue(rwMutexRLockFunc))

	// rwMutexRUnlock - разблокирует мьютекс для чтения
	rwMutexRUnlockFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: rwMutexRUnlock() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: rwMutexRUnlock() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		mutex, exists := rwMutexes[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: rwmutex '%s' does not exist", name))
		}
		
		mutex.RUnlock()
		return value.NewBool(true)
	}
	globalScope.Set("rwMutexRUnlock", value.NewValue(rwMutexRUnlockFunc))

	// rwMutexLock - блокирует мьютекс для записи
	rwMutexLockFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: rwMutexLock() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: rwMutexLock() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		mutex, exists := rwMutexes[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: rwmutex '%s' does not exist", name))
		}
		
		mutex.Lock()
		return value.NewBool(true)
	}
	globalScope.Set("rwMutexLock", value.NewValue(rwMutexLockFunc))

	// rwMutexUnlock - разблокирует мьютекс для записи
	rwMutexUnlockFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: rwMutexUnlock() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: rwMutexUnlock() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		mutex, exists := rwMutexes[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: rwmutex '%s' does not exist", name))
		}
		
		mutex.Unlock()
		return value.NewBool(true)
	}
	globalScope.Set("rwMutexUnlock", value.NewValue(rwMutexUnlockFunc))

	// ============ СЕМАФОРЫ ============

	// newSemaphore - создает новый семафор
	newSemaphoreFunc := func(args []*value.Value) *value.Value {
		if len(args) < 1 || len(args) > 2 {
			return value.NewString("Error: newSemaphore() requires 1-2 arguments (capacity, [name])")
		}
		
		capacity, ok := args[0].Any().(int64)
		if !ok {
			if floatVal, ok := args[0].Any().(float64); ok {
				capacity = int64(floatVal)
			} else {
				return value.NewString("Error: newSemaphore() first argument must be numeric (capacity)")
			}
		}
		
		if capacity <= 0 {
			return value.NewString("Error: semaphore capacity must be positive")
		}
		
		name := ""
		if len(args) == 2 {
			if n, ok := args[1].Any().(string); ok {
				name = n
			} else {
				return value.NewString("Error: newSemaphore() second argument must be string (name)")
			}
		}
		
		// Генерируем имя если не предоставлено
		if name == "" {
			syncMapsMu.Lock()
			mutexCounter++
			name = fmt.Sprintf("semaphore_%d", mutexCounter)
			syncMapsMu.Unlock()
		}
		
		syncMapsMu.Lock()
		if _, exists := semaphores[name]; exists {
			syncMapsMu.Unlock()
			return value.NewString(fmt.Sprintf("Error: semaphore '%s' already exists", name))
		}
		semaphores[name] = make(chan struct{}, capacity)
		syncMapsMu.Unlock()
		
		return value.NewString(name)
	}
	globalScope.Set("newSemaphore", value.NewValue(newSemaphoreFunc))

	// semaphoreAcquire - захватывает семафор
	semaphoreAcquireFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: semaphoreAcquire() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: semaphoreAcquire() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		sem, exists := semaphores[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: semaphore '%s' does not exist", name))
		}
		
		sem <- struct{}{}
		return value.NewBool(true)
	}
	globalScope.Set("semaphoreAcquire", value.NewValue(semaphoreAcquireFunc))

	// semaphoreRelease - освобождает семафор
	semaphoreReleaseFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: semaphoreRelease() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: semaphoreRelease() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		sem, exists := semaphores[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: semaphore '%s' does not exist", name))
		}
		
		select {
		case <-sem:
			return value.NewBool(true)
		default:
			return value.NewString("Error: semaphore release without acquire")
		}
	}
	globalScope.Set("semaphoreRelease", value.NewValue(semaphoreReleaseFunc))

	// semaphoreTryAcquire - пытается захватить семафор без ожидания
	semaphoreTryAcquireFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: semaphoreTryAcquire() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: semaphoreTryAcquire() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		sem, exists := semaphores[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: semaphore '%s' does not exist", name))
		}
		
		select {
		case sem <- struct{}{}:
			return value.NewBool(true)
		default:
			return value.NewBool(false)
		}
	}
	globalScope.Set("semaphoreTryAcquire", value.NewValue(semaphoreTryAcquireFunc))

	// ============ WAITGROUP ============

	// newWaitGroup - создает новую WaitGroup
	newWaitGroupFunc := func(args []*value.Value) *value.Value {
		name := ""
		if len(args) == 1 {
			if n, ok := args[0].Any().(string); ok {
				name = n
			} else {
				return value.NewString("Error: newWaitGroup() optional argument must be string (name)")
			}
		} else if len(args) > 1 {
			return value.NewString("Error: newWaitGroup() requires 0-1 arguments ([name])")
		}
		
		// Генерируем имя если не предоставлено
		if name == "" {
			syncMapsMu.Lock()
			mutexCounter++
			name = fmt.Sprintf("waitgroup_%d", mutexCounter)
			syncMapsMu.Unlock()
		}
		
		syncMapsMu.Lock()
		if _, exists := waitGroups[name]; exists {
			syncMapsMu.Unlock()
			return value.NewString(fmt.Sprintf("Error: waitgroup '%s' already exists", name))
		}
		waitGroups[name] = &sync.WaitGroup{}
		syncMapsMu.Unlock()
		
		return value.NewString(name)
	}
	globalScope.Set("newWaitGroup", value.NewValue(newWaitGroupFunc))

	// waitGroupAdd - добавляет счетчик к WaitGroup
	waitGroupAddFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: waitGroupAdd() requires 2 arguments (name, delta)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: waitGroupAdd() first argument must be string (name)")
		}
		
		delta, ok := args[1].Any().(int64)
		if !ok {
			if floatVal, ok := args[1].Any().(float64); ok {
				delta = int64(floatVal)
			} else {
				return value.NewString("Error: waitGroupAdd() second argument must be numeric (delta)")
			}
		}
		
		syncMapsMu.Lock()
		wg, exists := waitGroups[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: waitgroup '%s' does not exist", name))
		}
		
		wg.Add(int(delta))
		return value.NewBool(true)
	}
	globalScope.Set("waitGroupAdd", value.NewValue(waitGroupAddFunc))

	// waitGroupDone - уменьшает счетчик WaitGroup на 1
	waitGroupDoneFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: waitGroupDone() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: waitGroupDone() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		wg, exists := waitGroups[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: waitgroup '%s' does not exist", name))
		}
		
		wg.Done()
		return value.NewBool(true)
	}
	globalScope.Set("waitGroupDone", value.NewValue(waitGroupDoneFunc))

	// waitGroupWait - ждет пока счетчик WaitGroup не станет 0
	waitGroupWaitFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: waitGroupWait() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: waitGroupWait() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		wg, exists := waitGroups[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: waitgroup '%s' does not exist", name))
		}
		
		wg.Wait()
		return value.NewBool(true)
	}
	globalScope.Set("waitGroupWait", value.NewValue(waitGroupWaitFunc))

	// ============ АТОМАРНЫЕ ОПЕРАЦИИ ============

	// newAtomic - создает новую атомарную переменную
	newAtomicFunc := func(args []*value.Value) *value.Value {
		if len(args) < 1 || len(args) > 2 {
			return value.NewString("Error: newAtomic() requires 1-2 arguments (initialValue, [name])")
		}
		
		initialValue, ok := args[0].Any().(int64)
		if !ok {
			if floatVal, ok := args[0].Any().(float64); ok {
				initialValue = int64(floatVal)
			} else {
				return value.NewString("Error: newAtomic() first argument must be numeric (initialValue)")
			}
		}
		
		name := ""
		if len(args) == 2 {
			if n, ok := args[1].Any().(string); ok {
				name = n
			} else {
				return value.NewString("Error: newAtomic() second argument must be string (name)")
			}
		}
		
		// Генерируем имя если не предоставлено
		if name == "" {
			syncMapsMu.Lock()
			mutexCounter++
			name = fmt.Sprintf("atomic_%d", mutexCounter)
			syncMapsMu.Unlock()
		}
		
		syncMapsMu.Lock()
		if _, exists := atomicInts[name]; exists {
			syncMapsMu.Unlock()
			return value.NewString(fmt.Sprintf("Error: atomic '%s' already exists", name))
		}
		atomicInts[name] = &initialValue
		syncMapsMu.Unlock()
		
		return value.NewString(name)
	}
	globalScope.Set("newAtomic", value.NewValue(newAtomicFunc))

	// atomicGet - получает значение атомарной переменной
	atomicGetFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: atomicGet() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: atomicGet() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		atomicVar, exists := atomicInts[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: atomic '%s' does not exist", name))
		}
		
		val := atomic.LoadInt64(atomicVar)
		return value.NewInt64(val)
	}
	globalScope.Set("atomicGet", value.NewValue(atomicGetFunc))

	// atomicSet - устанавливает значение атомарной переменной
	atomicSetFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: atomicSet() requires 2 arguments (name, value)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: atomicSet() first argument must be string (name)")
		}
		
		newVal, ok := args[1].Any().(int64)
		if !ok {
			if floatVal, ok := args[1].Any().(float64); ok {
				newVal = int64(floatVal)
			} else {
				return value.NewString("Error: atomicSet() second argument must be numeric (value)")
			}
		}
		
		syncMapsMu.Lock()
		atomicVar, exists := atomicInts[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: atomic '%s' does not exist", name))
		}
		
		atomic.StoreInt64(atomicVar, newVal)
		return value.NewBool(true)
	}
	globalScope.Set("atomicSet", value.NewValue(atomicSetFunc))

	// atomicAdd - атомарно добавляет значение
	atomicAddFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: atomicAdd() requires 2 arguments (name, delta)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: atomicAdd() first argument must be string (name)")
		}
		
		delta, ok := args[1].Any().(int64)
		if !ok {
			if floatVal, ok := args[1].Any().(float64); ok {
				delta = int64(floatVal)
			} else {
				return value.NewString("Error: atomicAdd() second argument must be numeric (delta)")
			}
		}
		
		syncMapsMu.Lock()
		atomicVar, exists := atomicInts[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: atomic '%s' does not exist", name))
		}
		
		newVal := atomic.AddInt64(atomicVar, delta)
		return value.NewInt64(newVal)
	}
	globalScope.Set("atomicAdd", value.NewValue(atomicAddFunc))

	// atomicCompareAndSwap - атомарное сравнение и замена
	atomicCompareAndSwapFunc := func(args []*value.Value) *value.Value {
		if len(args) != 3 {
			return value.NewString("Error: atomicCompareAndSwap() requires 3 arguments (name, expected, new)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: atomicCompareAndSwap() first argument must be string (name)")
		}
		
		expected, ok := args[1].Any().(int64)
		if !ok {
			if floatVal, ok := args[1].Any().(float64); ok {
				expected = int64(floatVal)
			} else {
				return value.NewString("Error: atomicCompareAndSwap() second argument must be numeric (expected)")
			}
		}
		
		newVal, ok := args[2].Any().(int64)
		if !ok {
			if floatVal, ok := args[2].Any().(float64); ok {
				newVal = int64(floatVal)
			} else {
				return value.NewString("Error: atomicCompareAndSwap() third argument must be numeric (new)")
			}
		}
		
		syncMapsMu.Lock()
		atomicVar, exists := atomicInts[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: atomic '%s' does not exist", name))
		}
		
		swapped := atomic.CompareAndSwapInt64(atomicVar, expected, newVal)
		return value.NewBool(swapped)
	}
	globalScope.Set("atomicCompareAndSwap", value.NewValue(atomicCompareAndSwapFunc))

	// ============ БАРЬЕРЫ ============

	// newBarrier - создает новый барьер
	newBarrierFunc := func(args []*value.Value) *value.Value {
		if len(args) < 1 || len(args) > 2 {
			return value.NewString("Error: newBarrier() requires 1-2 arguments (n, [name])")
		}
		
		n, ok := args[0].Any().(int64)
		if !ok {
			if floatVal, ok := args[0].Any().(float64); ok {
				n = int64(floatVal)
			} else {
				return value.NewString("Error: newBarrier() first argument must be numeric (n)")
			}
		}
		
		if n <= 0 {
			return value.NewString("Error: barrier n must be positive")
		}
		
		name := ""
		if len(args) == 2 {
			if nm, ok := args[1].Any().(string); ok {
				name = nm
			} else {
				return value.NewString("Error: newBarrier() second argument must be string (name)")
			}
		}
		
		// Генерируем имя если не предоставлено
		if name == "" {
			syncMapsMu.Lock()
			mutexCounter++
			name = fmt.Sprintf("barrier_%d", mutexCounter)
			syncMapsMu.Unlock()
		}
		
		syncMapsMu.Lock()
		if _, exists := barriers[name]; exists {
			syncMapsMu.Unlock()
			return value.NewString(fmt.Sprintf("Error: barrier '%s' already exists", name))
		}
		barriers[name] = newBarrier(int(n))
		syncMapsMu.Unlock()
		
		return value.NewString(name)
	}
	globalScope.Set("newBarrier", value.NewValue(newBarrierFunc))

	// barrierWait - ждет на барьере
	barrierWaitFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: barrierWait() requires 1 argument (name)")
		}
		
		name, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: barrierWait() argument must be string (name)")
		}
		
		syncMapsMu.Lock()
		barrier, exists := barriers[name]
		syncMapsMu.Unlock()
		
		if !exists {
			return value.NewString(fmt.Sprintf("Error: barrier '%s' does not exist", name))
		}
		
		err := barrier.wait()
		if err != nil {
			return value.NewString(fmt.Sprintf("Error: %v", err))
		}
		
		return value.NewBool(true)
	}
	globalScope.Set("barrierWait", value.NewValue(barrierWaitFunc))

	// ============ ОЧИСТКА ============

	// syncCleanup - очищает все примитивы синхронизации
	syncCleanupFunc := func(args []*value.Value) *value.Value {
		syncMapsMu.Lock()
		defer syncMapsMu.Unlock()
		
		// Очищаем все карты
		mutexes = make(map[string]*sync.Mutex)
		rwMutexes = make(map[string]*sync.RWMutex)
		waitGroups = make(map[string]*sync.WaitGroup)
		semaphores = make(map[string]chan struct{})
		atomicInts = make(map[string]*int64)
		conditions = make(map[string]*sync.Cond)
		barriers = make(map[string]*barrier)
		mutexCounter = 0
		
		return value.NewBool(true)
	}
	globalScope.Set("syncCleanup", value.NewValue(syncCleanupFunc))
}