package value

import (
	"fmt"
	"sync"
	"time"
)

// ChannelState представляет состояние канала
type ChannelState int

const (
	ChannelOpen ChannelState = iota
	ChannelClosed
)

// Channel представляет канал для коммуникации между горутинами
type Channel struct {
	buffer   chan *Value  // Буферизованный канал
	state    ChannelState // Состояние канала
	mu       sync.RWMutex // Мьютекс для защиты состояния
	capacity int          // Размер буфера
	closed   bool         // Флаг закрытия
}

// NewChannel создает новый канал с указанным размером буфера
func NewChannel(bufferSize int) *Channel {
	if bufferSize < 0 {
		bufferSize = 0
	}
	
	return &Channel{
		buffer:   make(chan *Value, bufferSize),
		state:    ChannelOpen,
		capacity: bufferSize,
		closed:   false,
	}
}

// Send отправляет значение в канал (неблокирующая операция если канал буферизован)
func (ch *Channel) Send(value *Value) error {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	
	if ch.closed {
		return fmt.Errorf("cannot send to closed channel")
	}
	
	select {
	case ch.buffer <- value:
		return nil
	case <-time.After(5 * time.Second): // Таймаут 5 секунд
		return fmt.Errorf("channel send timeout")
	}
}

// SendBlocking отправляет значение в канал (блокирующая операция)
func (ch *Channel) SendBlocking(value *Value) error {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	
	if ch.closed {
		return fmt.Errorf("cannot send to closed channel")
	}
	
	ch.buffer <- value
	return nil
}

// Receive получает значение из канала (блокирующая операция)
func (ch *Channel) Receive() (*Value, error) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	
	select {
	case value, ok := <-ch.buffer:
		if !ok {
			return NewString(""), fmt.Errorf("channel is closed")
		}
		return value, nil
	case <-time.After(5 * time.Second): // Таймаут 5 секунд
		return NewString(""), fmt.Errorf("channel receive timeout")
	}
}

// ReceiveBlocking получает значение из канала без таймаута
func (ch *Channel) ReceiveBlocking() (*Value, error) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	
	value, ok := <-ch.buffer
	if !ok {
		return NewString(""), fmt.Errorf("channel is closed")
	}
	return value, nil
}

// TryReceive пытается получить значение из канала (неблокирующая операция)
func (ch *Channel) TryReceive() (*Value, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	
	if ch.closed && len(ch.buffer) == 0 {
		return NewString(""), false
	}
	
	select {
	case value, ok := <-ch.buffer:
		return value, ok
	default:
		return NewString(""), false
	}
}

// Close закрывает канал
func (ch *Channel) Close() {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	
	if !ch.closed {
		ch.closed = true
		ch.state = ChannelClosed
		close(ch.buffer)
	}
}

// IsClosed проверяет, закрыт ли канал
func (ch *Channel) IsClosed() bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.closed
}

// Len возвращает количество элементов в буфере канала
func (ch *Channel) Len() int {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return len(ch.buffer)
}

// Cap возвращает емкость буфера канала
func (ch *Channel) Cap() int {
	return ch.capacity
}

// String возвращает строковое представление канала
func (ch *Channel) String() string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	
	status := "open"
	if ch.closed {
		status = "closed"
	}
	
	return fmt.Sprintf("chan(cap:%d, len:%d, %s)", ch.capacity, len(ch.buffer), status)
}

// IsEmpty проверяет, пуст ли канал
func (ch *Channel) IsEmpty() bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return len(ch.buffer) == 0
}

// IsFull проверяет, полон ли канал
func (ch *Channel) IsFull() bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return len(ch.buffer) == ch.capacity
}

// Select представляет операцию select для множественного выбора каналов
type Select struct {
	cases   []SelectCase
	timeout time.Duration
}

// SelectCase представляет один case в select операции
type SelectCase struct {
	Channel   *Channel
	Value     *Value  // Для send операций
	IsSend    bool    // true для send, false для receive
	IsDefault bool    // true для default case
}

// NewSelect создает новую select операцию
func NewSelect() *Select {
	return &Select{
		cases:   make([]SelectCase, 0),
		timeout: 5 * time.Second, // Дефолтный таймаут
	}
}

// AddReceiveCase добавляет receive case
func (s *Select) AddReceiveCase(ch *Channel) {
	s.cases = append(s.cases, SelectCase{
		Channel: ch,
		IsSend:  false,
	})
}

// AddSendCase добавляет send case
func (s *Select) AddSendCase(ch *Channel, value *Value) {
	s.cases = append(s.cases, SelectCase{
		Channel: ch,
		Value:   value,
		IsSend:  true,
	})
}

// AddDefaultCase добавляет default case
func (s *Select) AddDefaultCase() {
	s.cases = append(s.cases, SelectCase{
		IsDefault: true,
	})
}

// Execute выполняет select операцию
func (s *Select) Execute() (int, *Value, error) {
	// Простая реализация select - проверяем все cases по очереди
	for i, selectCase := range s.cases {
		if selectCase.IsDefault {
			continue // Default обрабатываем в конце
		}
		
		if selectCase.IsSend {
			// Пытаемся отправить
			if !selectCase.Channel.IsFull() && !selectCase.Channel.IsClosed() {
				err := selectCase.Channel.Send(selectCase.Value)
				if err == nil {
					return i, NewString("sent"), nil
				}
			}
		} else {
			// Пытаемся получить
			if !selectCase.Channel.IsEmpty() {
				value, ok := selectCase.Channel.TryReceive()
				if ok {
					return i, value, nil
				}
			}
		}
	}
	
	// Если ничего не сработало, проверяем default
	for i, selectCase := range s.cases {
		if selectCase.IsDefault {
			return i, NewString("default"), nil
		}
	}
	
	// Если нет default, возвращаем ошибку
	return -1, NewString(""), fmt.Errorf("all cases blocked and no default")
}