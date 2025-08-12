package test

import (
	"bytes"
	"foo_lang/builtin"
	"foo_lang/parser"
	"foo_lang/scope"
	"foo_lang/value"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestChannelBasicOperations(t *testing.T) {
	// Инициализируем scope с функциями каналов
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeMathFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	builtin.InitializeChannelFunctions(scope.GlobalScope)

	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "create_unbuffered_channel",
			code: `
				let ch = newChannel()
				let info = channelInfo(ch)
				println(info)
			`,
			want: "chan(cap:0, len:0, open)",
		},
		{
			name: "create_buffered_channel",
			code: `
				let ch = newChannel(3)
				let info = channelInfo(ch)
				println(info)
			`,
			want: "chan(cap:3, len:0, open)",
		},
		{
			name: "send_and_receive",
			code: `
				let ch = newChannel(2)
				send(ch, "hello")
				send(ch, "world")
				let msg1 = receive(ch)
				let msg2 = receive(ch)
				println(msg1 + " " + msg2)
			`,
			want: "hello world",
		},
		{
			name: "channel_length_and_capacity",
			code: `
				let ch = newChannel(3)
				send(ch, 1)
				send(ch, 2)
				println("len: " + len(ch).toString() + ", cap: " + cap(ch).toString())
			`,
			want: "len: 2, cap: 3",
		},
		{
			name: "close_channel",
			code: `
				let ch = newChannel(1)
				send(ch, "test")
				close(ch)
				let info = channelInfo(ch)
				println(info)
			`,
			want: "chan(cap:1, len:1, closed)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureChannelOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				// Регистрируем все встроенные функции
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeChannelFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.want {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.want, result)
			}
		})
	}
}

func TestChannelNonBlockingOperations(t *testing.T) {
	// Инициализируем scope
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeChannelFunctions(scope.GlobalScope)

	t.Run("try_receive_empty_channel", func(t *testing.T) {
		result := captureChannelOutput(func() {
			scope.GlobalScope = scope.NewScopeStack()
			builtin.InitializeChannelFunctions(scope.GlobalScope)

			code := `
				let ch = newChannel(1)
				let result = tryReceive(ch)
				println(result)
			`
			exprs := parser.NewParser(code).Parse()
			for _, expr := range exprs {
				expr.Eval()
			}
		})

		if result != "no_value" {
			t.Errorf("expected 'no_value', got %q", result)
		}
	})

	t.Run("try_receive_with_data", func(t *testing.T) {
		result := captureChannelOutput(func() {
			scope.GlobalScope = scope.NewScopeStack()
			builtin.InitializeChannelFunctions(scope.GlobalScope)

			code := `
				let ch = newChannel(1)
				send(ch, "available")
				let result = tryReceive(ch)
				println(result)
			`
			exprs := parser.NewParser(code).Parse()
			for _, expr := range exprs {
				expr.Eval()
			}
		})

		if result != "available" {
			t.Errorf("expected 'available', got %q", result)
		}
	})
}

func TestChannelWithNumbers(t *testing.T) {
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeChannelFunctions(scope.GlobalScope)

	result := captureChannelOutput(func() {
		scope.GlobalScope = scope.NewScopeStack()
		builtin.InitializeChannelFunctions(scope.GlobalScope)

		code := `
			let ch = newChannel(3)
			send(ch, 42)
			send(ch, 3.14)
			send(ch, 100)
			
			let num1 = receive(ch)
			let num2 = receive(ch)
			let num3 = receive(ch)
			
			println(num1.toString() + "," + num2.toString() + "," + num3.toString())
		`
		exprs := parser.NewParser(code).Parse()
		for _, expr := range exprs {
			expr.Eval()
		}
	})

	expected := "42,3.14,100"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestChannelDirectAPI(t *testing.T) {
	// Тест прямого API каналов (без foo_lang кода)

	t.Run("channel_creation", func(t *testing.T) {
		ch := value.NewChannel(2)
		if ch.Cap() != 2 {
			t.Errorf("expected capacity 2, got %d", ch.Cap())
		}
		if ch.Len() != 0 {
			t.Errorf("expected length 0, got %d", ch.Len())
		}
		if ch.IsClosed() {
			t.Error("expected channel to be open")
		}
	})

	t.Run("send_and_receive", func(t *testing.T) {
		ch := value.NewChannel(1)
		testValue := value.NewString("test")

		err := ch.Send(testValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		received, err := ch.Receive()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if received.Any().(string) != "test" {
			t.Errorf("expected 'test', got %v", received.Any())
		}
	})

	t.Run("channel_closure", func(t *testing.T) {
		ch := value.NewChannel(1)
		testValue := value.NewString("before_close")

		ch.Send(testValue)
		ch.Close()

		if !ch.IsClosed() {
			t.Error("expected channel to be closed")
		}

		// Должны все еще суметь получить данные, отправленные до закрытия
		received, ok := ch.TryReceive()
		if ok && received.Any().(string) != "before_close" {
			t.Errorf("expected 'before_close', got %v", received.Any())
		}
	})

	t.Run("buffer_overflow", func(t *testing.T) {
		ch := value.NewChannel(1)

		// Заполняем буфер
		err := ch.Send(value.NewString("first"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Проверяем, что канал полон
		if !ch.IsFull() {
			t.Error("expected channel to be full")
		}
	})

	t.Run("empty_channel_check", func(t *testing.T) {
		ch := value.NewChannel(1)

		if !ch.IsEmpty() {
			t.Error("expected channel to be empty")
		}

		ch.Send(value.NewString("data"))

		if ch.IsEmpty() {
			t.Error("expected channel to not be empty")
		}
	})
}

func TestChannelConcurrency(t *testing.T) {
	// Тест конкурентности каналов
	ch := value.NewChannel(20) // Увеличиваем буфер, чтобы вместить все сообщения

	// Запускаем несколько горутин для отправки данных
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				msg := value.NewString("msg_" + string(rune(n)) + "_" + string(rune(j)))
				ch.Send(msg)
			}
		}(i)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// Проверяем, что данные получены
	if ch.Len() != 15 {
		t.Errorf("expected 15 messages, got %d", ch.Len())
	}

	// Получаем все сообщения
	count := 0
	for !ch.IsEmpty() {
		_, ok := ch.TryReceive()
		if ok {
			count++
		}
	}

	if count != 15 {
		t.Errorf("expected to receive 15 messages, got %d", count)
	}
}

func TestChannelSelect(t *testing.T) {
	// Тест select операций
	sel := value.NewSelect()
	ch1 := value.NewChannel(1)
	ch2 := value.NewChannel(1)

	// Добавляем receive cases
	sel.AddReceiveCase(ch1)
	sel.AddReceiveCase(ch2)
	sel.AddDefaultCase()

	// Отправляем в первый канал
	ch1.Send(value.NewString("from_ch1"))

	// Выполняем select
	caseIndex, result, err := sel.Execute()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if caseIndex != 0 {
		t.Errorf("expected case 0, got %d", caseIndex)
	}

	if result.Any().(string) != "from_ch1" {
		t.Errorf("expected 'from_ch1', got %v", result.Any())
	}
}

func TestChannelErrorHandling(t *testing.T) {
	t.Run("send_to_closed_channel", func(t *testing.T) {
		ch := value.NewChannel(1)
		ch.Close()

		err := ch.Send(value.NewString("test"))
		if err == nil {
			t.Error("expected error when sending to closed channel")
		}
	})

	t.Run("invalid_buffer_size", func(t *testing.T) {
		ch := value.NewChannel(-5)
		if ch.Cap() != 0 {
			t.Errorf("expected capacity 0 for negative buffer size, got %d", ch.Cap())
		}
	})
}

// captureOutput захватывает stdout для тестирования
func captureChannelOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	result := buf.String()
	return strings.TrimSpace(result)
}
