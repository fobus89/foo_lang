package builtin

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"time"
)

// InitializeChannelFunctions инициализирует встроенные функции для работы с каналами
func InitializeChannelFunctions(globalScope *scope.ScopeStack) {
	// newChannel - создание канала
	newChannelFunc := func(args []*value.Value) *value.Value {
		if len(args) > 1 {
			return value.NewString("Error: newChannel() requires 0 or 1 argument ([bufferSize])")
		}
		
		// Аргумент - размер буфера (опциональный)
		bufferSize := 0
		if len(args) == 1 {
			if bufVal, ok := args[0].Any().(int64); ok {
				bufferSize = int(bufVal)
			}
		}
		
		ch := value.NewChannel(bufferSize)
		return value.NewChannelValue(ch)
	}
	globalScope.Set("newChannel", value.NewValue(newChannelFunc))
	
	// send - отправка в канал (синтаксический сахар для ch <- value)
	sendFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: send() requires 2 arguments (channel, value)")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: first argument must be a channel")
		}
		
		err := chVal.SendBlocking(args[1])
		if err != nil {
			return value.NewString(fmt.Sprintf("Error: %v", err))
		}
		
		return value.NewString("sent")
	}
	globalScope.Set("send", value.NewValue(sendFunc))
	
	// receive - получение из канала (синтаксический сахар для <-ch)
	receiveFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: receive() requires 1 argument (channel)")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: argument must be a channel")
		}
		
		result, err := chVal.ReceiveBlocking()
		if err != nil {
			return value.NewString(fmt.Sprintf("Error: %v", err))
		}
		
		return result
	}
	globalScope.Set("receive", value.NewValue(receiveFunc))
	
	// close - закрытие канала
	closeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: close() requires 1 argument (channel)")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: argument must be a channel")
		}
		
		chVal.Close()
		return value.NewString("closed")
	}
	globalScope.Set("close", value.NewValue(closeFunc))
	
	// len - количество элементов в канале
	lenFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: len() requires 1 argument")
		}
		
		switch val := args[0].Any().(type) {
		case *value.Channel:
			return value.NewInt64(int64(val.Len()))
		case []interface{}:
			return value.NewInt64(int64(len(val)))
		case string:
			return value.NewInt64(int64(len(val)))
		default:
			return value.NewString("Error: len() requires array, string, or channel")
		}
	}
	globalScope.Set("len", value.NewValue(lenFunc))
	
	// cap - емкость канала
	capFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: cap() requires 1 argument")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: cap() requires a channel")
		}
		
		return value.NewInt64(int64(chVal.Cap()))
	}
	globalScope.Set("cap", value.NewValue(capFunc))
	
	// tryReceive - неблокирующее получение из канала
	tryReceiveFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: tryReceive() requires 1 argument (channel)")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: argument must be a channel")
		}
		
		result, success := chVal.TryReceive()
		if !success {
			return value.NewString("no_value")
		}
		
		return result
	}
	globalScope.Set("tryReceive", value.NewValue(tryReceiveFunc))
	
	// channelInfo - информация о канале
	channelInfoFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: channelInfo() requires 1 argument (channel)")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: argument must be a channel")
		}
		
		return value.NewString(chVal.String())
	}
	globalScope.Set("channelInfo", value.NewValue(channelInfoFunc))
	
	// trySend - неблокирующая отправка в канал
	trySendFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: trySend() requires 2 arguments (channel, value)")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: first argument must be a channel")
		}
		
		success := chVal.TrySend(args[1])
		return value.NewBool(success)
	}
	globalScope.Set("trySend", value.NewValue(trySendFunc))
	
	// channelSelect - select для множественного выбора каналов
	channelSelectFunc := func(args []*value.Value) *value.Value {
		if len(args) < 1 {
			return value.NewString("Error: channelSelect() requires at least 1 argument (array of channels)")
		}
		
		channels, ok := args[0].Any().([]interface{})
		if !ok {
			return value.NewString("Error: first argument must be array of channels")
		}
		
		// Создаем select операцию
		sel := value.NewSelect()
		for _, ch := range channels {
			if channel, ok := ch.(*value.Value); ok {
				if chVal, ok := channel.Any().(*value.Channel); ok {
					sel.AddReceiveCase(chVal)
				}
			}
		}
		
		// Выполняем select
		index, result, err := sel.Execute()
		if err != nil {
			return value.NewString(fmt.Sprintf("Error: %v", err))
		}
		
		// Возвращаем объект с результатом
		resultObj := map[string]interface{}{
			"index": int64(index),
			"value": result.Any(),
		}
		return value.NewValue(resultObj)
	}
	globalScope.Set("channelSelect", value.NewValue(channelSelectFunc))
	
	// channelTimeout - операции с настраиваемым таймаутом
	channelTimeoutFunc := func(args []*value.Value) *value.Value {
		if len(args) != 3 {
			return value.NewString("Error: channelTimeout() requires 3 arguments (channel, operation, timeoutMs)")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: first argument must be a channel")
		}
		
		operation, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: second argument must be operation string ('send' or 'receive')")
		}
		
		timeoutMs, ok := args[2].Any().(int64)
		if !ok {
			return value.NewString("Error: third argument must be timeout in milliseconds")
		}
		
		switch operation {
		case "receive":
			// Создаем канал для результата
			resultChan := make(chan *value.Value, 1)
			errorChan := make(chan error, 1)
			
			// Запускаем получение в горутине
			go func() {
				val, err := chVal.ReceiveBlocking()
				if err != nil {
					errorChan <- err
				} else {
					resultChan <- val
				}
			}()
			
			// Ждем результат или таймаут
			select {
			case result := <-resultChan:
				return result
			case err := <-errorChan:
				return value.NewString(fmt.Sprintf("Error: %v", err))
			case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
				return value.NewString("timeout")
			}
		default:
			return value.NewString("Error: unsupported operation. Use 'receive'")
		}
	}
	globalScope.Set("channelTimeout", value.NewValue(channelTimeoutFunc))
	
	// channelRange - итерация по каналу до закрытия
	channelRangeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: channelRange() requires 1 argument (channel)")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: argument must be a channel")
		}
		
		var results []interface{}
		for {
			val, ok := chVal.TryReceive()
			if !ok {
				// Канал пуст, проверяем закрыт ли он
				if chVal.IsClosed() {
					break
				}
				// Канал не закрыт, ждем немного
				time.Sleep(10 * time.Millisecond)
				continue
			}
			results = append(results, val.Any())
		}
		
		return value.NewValue(results)
	}
	globalScope.Set("channelRange", value.NewValue(channelRangeFunc))
	
	// channelDrain - очистка канала
	channelDrainFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: channelDrain() requires 1 argument (channel)")
		}
		
		chVal, ok := args[0].Any().(*value.Channel)
		if !ok {
			return value.NewString("Error: argument must be a channel")
		}
		
		var drained []interface{}
		for {
			val, ok := chVal.TryReceive()
			if !ok {
				break // Канал пуст
			}
			drained = append(drained, val.Any())
		}
		
		return value.NewValue(drained)
	}
	globalScope.Set("channelDrain", value.NewValue(channelDrainFunc))
}

// Worker представляет worker для обработки задач из канала
type Worker struct {
	ID      int
	JobChan <-chan *value.Value
	Quit    chan bool
}

// NewWorker создает нового worker
func NewWorker(id int, jobChan <-chan *value.Value) *Worker {
	return &Worker{
		ID:      id,
		JobChan: jobChan,
		Quit:    make(chan bool),
	}
}

// Start запускает worker
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.JobChan:
				// Обрабатываем задачу
				fmt.Printf("Worker %d processing job: %v\n", w.ID, job.Any())
			case <-w.Quit:
				fmt.Printf("Worker %d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop останавливает worker
func (w *Worker) Stop() {
	w.Quit <- true
}

// WorkerPool представляет пул workers
type WorkerPool struct {
	Workers  []*Worker
	JobQueue chan *value.Value
	QuitAll  chan bool
}

// NewWorkerPool создает новый пул workers
func NewWorkerPool(numWorkers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		Workers:  make([]*Worker, 0, numWorkers),
		JobQueue: make(chan *value.Value, queueSize),
		QuitAll:  make(chan bool),
	}
}

// Start запускает пул workers
func (wp *WorkerPool) Start() {
	for i := 0; i < cap(wp.Workers); i++ {
		worker := NewWorker(i+1, wp.JobQueue)
		wp.Workers = append(wp.Workers, worker)
		worker.Start()
	}
}

// AddJob добавляет задачу в очередь
func (wp *WorkerPool) AddJob(job *value.Value) {
	wp.JobQueue <- job
}

// Stop останавливает все workers
func (wp *WorkerPool) Stop() {
	for _, worker := range wp.Workers {
		worker.Stop()
	}
	close(wp.JobQueue)
}