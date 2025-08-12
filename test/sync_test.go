package test

import (
	"bytes"
	"foo_lang/builtin"
	"foo_lang/parser"
	"foo_lang/scope"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMutexFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "create_and_use_mutex",
			code: `
				let mutex = newMutex("test_mutex")
				println("Created: " + mutex)
				let locked = mutexLock(mutex)
				println("Locked: " + locked.toString())
				let unlocked = mutexUnlock(mutex)
				println("Unlocked: " + unlocked.toString())
			`,
			expected: "Created: test_mutex\nLocked: true\nUnlocked: true",
		},
		{
			name: "auto_generated_mutex_name",
			code: `
				let mutex = newMutex()
				println("Created with auto name")
				let locked = mutexLock(mutex)
				println("Locked: " + locked.toString())
				let unlocked = mutexUnlock(mutex)
				println("Unlocked: " + unlocked.toString())
			`,
			expected: "Created with auto name\nLocked: true\nUnlocked: true",
		},
		{
			name: "mutex_already_exists_error",
			code: `
				let mutex1 = newMutex("duplicate")
				println("First: " + mutex1)
				let mutex2 = newMutex("duplicate")
				println("Second: " + mutex2)
			`,
			expected: "First: duplicate\nSecond: Error: mutex 'duplicate' already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureSyncOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeSyncFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestRWMutexFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "create_and_use_rwmutex",
			code: `
				let rwmutex = newRWMutex("test_rwmutex")
				println("Created: " + rwmutex)
				let rlocked = rwMutexRLock(rwmutex)
				println("RLocked: " + rlocked.toString())
				let runlocked = rwMutexRUnlock(rwmutex)
				println("RUnlocked: " + runlocked.toString())
			`,
			expected: "Created: test_rwmutex\nRLocked: true\nRUnlocked: true",
		},
		{
			name: "rwmutex_write_operations",
			code: `
				let rwmutex = newRWMutex("test_rw")
				let locked = rwMutexLock(rwmutex)
				println("Write locked: " + locked.toString())
				let unlocked = rwMutexUnlock(rwmutex)
				println("Write unlocked: " + unlocked.toString())
			`,
			expected: "Write locked: true\nWrite unlocked: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureSyncOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeSyncFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestSemaphoreFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "create_and_use_semaphore",
			code: `
				let sem = newSemaphore(2, "test_sem")
				println("Created: " + sem)
				let acquired = semaphoreAcquire(sem)
				println("Acquired: " + acquired.toString())
				let released = semaphoreRelease(sem)
				println("Released: " + released.toString())
			`,
			expected: "Created: test_sem\nAcquired: true\nReleased: true",
		},
		{
			name: "semaphore_try_acquire",
			code: `
				let sem = newSemaphore(1)
				let acquired1 = semaphoreTryAcquire(sem)
				println("First try: " + acquired1.toString())
				let acquired2 = semaphoreTryAcquire(sem)
				println("Second try: " + acquired2.toString())
			`,
			expected: "First try: true\nSecond try: false",
		},
		{
			name: "semaphore_invalid_capacity",
			code: `
				let sem = newSemaphore(0)
				println("Result: " + sem)
			`,
			expected: "Result: Error: semaphore capacity must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureSyncOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeSyncFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestWaitGroupFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "create_and_use_waitgroup",
			code: `
				let wg = newWaitGroup("test_wg")
				println("Created: " + wg)
				let added = waitGroupAdd(wg, 1)
				println("Added: " + added.toString())
				let done = waitGroupDone(wg)
				println("Done: " + done.toString())
			`,
			expected: "Created: test_wg\nAdded: true\nDone: true",
		},
		{
			name: "waitgroup_multiple_operations",
			code: `
				let wg = newWaitGroup()
				let added = waitGroupAdd(wg, 3)
				println("Added 3: " + added.toString())
				waitGroupDone(wg)
				waitGroupDone(wg)
				let done = waitGroupDone(wg)
				println("Final done: " + done.toString())
			`,
			expected: "Added 3: true\nFinal done: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureSyncOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeSyncFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestAtomicFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "create_and_use_atomic",
			code: `
				let atomic = newAtomic(42, "test_atomic")
				println("Created: " + atomic)
				let value = atomicGet(atomic)
				println("Value: " + value.toString())
				let set = atomicSet(atomic, 100)
				println("Set: " + set.toString())
				let newValue = atomicGet(atomic)
				println("New value: " + newValue.toString())
			`,
			expected: "Created: test_atomic\nValue: 42\nSet: true\nNew value: 100",
		},
		{
			name: "atomic_add_operations",
			code: `
				let atomic = newAtomic(10)
				let added = atomicAdd(atomic, 5)
				println("Added result: " + added.toString())
				let value = atomicGet(atomic)
				println("Final value: " + value.toString())
			`,
			expected: "Added result: 15\nFinal value: 15",
		},
		{
			name: "atomic_compare_and_swap",
			code: `
				let atomic = newAtomic(50)
				let swapped1 = atomicCompareAndSwap(atomic, 50, 75)
				println("Swapped (50->75): " + swapped1.toString())
				let swapped2 = atomicCompareAndSwap(atomic, 50, 100)
				println("Swapped (50->100): " + swapped2.toString())
				let value = atomicGet(atomic)
				println("Final value: " + value.toString())
			`,
			expected: "Swapped (50->75): true\nSwapped (50->100): false\nFinal value: 75",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureSyncOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeSyncFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestBarrierFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "create_barrier",
			code: `
				let barrier = newBarrier(2, "test_barrier")
				println("Created: " + barrier)
			`,
			expected: "Created: test_barrier",
		},
		{
			name: "barrier_invalid_size",
			code: `
				let barrier = newBarrier(0)
				println("Result: " + barrier)
			`,
			expected: "Result: Error: barrier n must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureSyncOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeSyncFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestSyncErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErr  bool
		contains string
	}{
		{
			name: "mutex_not_found",
			code: `
				let result = mutexLock("nonexistent")
				println(result)
			`,
			wantErr:  true,
			contains: "does not exist",
		},
		{
			name: "semaphore_release_without_acquire",
			code: `
				let sem = newSemaphore(1)
				let result = semaphoreRelease(sem)
				println(result)
			`,
			wantErr:  true,
			contains: "release without acquire",
		},
		{
			name: "wrong_argument_type",
			code: `
				let result = newMutex(123)
				println(result)
			`,
			wantErr:  true,
			contains: "must be string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureSyncOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeSyncFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if tt.wantErr && !strings.Contains(result, tt.contains) {
				t.Errorf("%s: expected error containing %q, got %q", tt.name, tt.contains, result)
			}
		})
	}
}

func TestSyncCleanupFunction(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "cleanup_all_sync_primitives",
			code: `
				let mutex = newMutex("test")
				let sem = newSemaphore(1, "test_sem")
				let atomic = newAtomic(42, "test_atomic")
				println("Created primitives")
				let cleaned = syncCleanup()
				println("Cleanup: " + cleaned.toString())
				// После очистки примитивы должны быть недоступны
				let result = mutexLock("test")
				println("After cleanup: " + result)
			`,
			expected: "Created primitives\nCleanup: true\nAfter cleanup: Error: mutex 'test' does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureSyncOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeSyncFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestSyncConcurrencyBasic(t *testing.T) {
	// Простые тесты синхронизации без реальной многопоточности
	// (так как foo_lang использует async/await для горутин)
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "multiple_mutex_operations",
			code: `
				let mutex = newMutex("concurrent_test")
				mutexLock(mutex)
				println("Locked")
				mutexUnlock(mutex)
				println("Unlocked")
				mutexLock(mutex)
				println("Locked again")
				mutexUnlock(mutex)
				println("Unlocked again")
			`,
			expected: "Locked\nUnlocked\nLocked again\nUnlocked again",
		},
		{
			name: "semaphore_capacity_management",
			code: `
				let sem = newSemaphore(2)
				semaphoreAcquire(sem)
				println("Acquired 1")
				semaphoreAcquire(sem)
				println("Acquired 2")
				let canAcquire = semaphoreTryAcquire(sem)
				println("Try acquire 3: " + canAcquire.toString())
				semaphoreRelease(sem)
				println("Released 1")
				let canAcquireNow = semaphoreTryAcquire(sem)
				println("Try acquire after release: " + canAcquireNow.toString())
			`,
			expected: "Acquired 1\nAcquired 2\nTry acquire 3: false\nReleased 1\nTry acquire after release: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureSyncOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeSyncFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

// captureSyncOutput захватывает stdout для тестирования sync функций
func captureSyncOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию с таймаутом для предотвращения зависания
	done := make(chan struct{})
	go func() {
		defer func() {
			if err := recover(); err != nil {
				// Игнорируем панику в тестах
			}
			close(done)
		}()
		f()
	}()

	select {
	case <-done:
		// Функция завершилась нормально
	case <-time.After(5 * time.Second):
		// Таймаут - функция зависла
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	result := buf.String()

	// Удаляем последний символ новой строки, если есть
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result
}