package test

import (
	"bytes"
	"fmt"
	"foo_lang/builtin"
	"foo_lang/parser"
	"foo_lang/scope"
	"foo_lang/value"
	"io"
	"os"
	"testing"
	"time"
)

func TestTimeBasicOperations(t *testing.T) {
	// Инициализируем scope с функциями времени
	scope.GlobalScope = scope.NewScopeStack()
	builtin.InitializeMathFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	builtin.InitializeTimeFunctions(scope.GlobalScope)

	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "current_time",
			code: `
				let currentTime = now()
				let year = timeYear(currentTime)
				println("Current year: " + year.toString())
			`,
			want: fmt.Sprintf("Current year: %d", time.Now().Year()),
		},
		{
			name: "time_from_unix",
			code: `
				let timestamp = 1609459200  // 2021-01-01 00:00:00 UTC
				let time = timeFromUnix(timestamp)
				let year = timeYear(time)
				println("Year: " + year.toString())
			`,
			want: "Year: 2021",
		},
		{
			name: "time_components",
			code: `
				let timestamp = 1609459200  // 2021-01-01 00:00:00 UTC
				let time = timeFromUnix(timestamp)
				let year = timeYear(time)
				let month = timeMonth(time)
				let day = timeDay(time)
				println(year.toString() + "-" + month.toString() + "-" + day.toString())
			`,
			want: "2021-1-1",
		},
		{
			name: "time_format",
			code: `
				let timestamp = 1609459200  // 2021-01-01 00:00:00 UTC
				let time = timeFromUnix(timestamp)
				let formatted = timeFormat(time, "date")
				println(formatted)
			`,
			want: "2021-01-01",
		},
		{
			name: "time_add_days",
			code: `
				let timestamp = 1609459200  // 2021-01-01 00:00:00 UTC
				let time = timeFromUnix(timestamp)
				let newTime = timeAddDays(time, 10)
				let day = timeDay(newTime)
				println("Day after adding 10 days: " + day.toString())
			`,
			want: "Day after adding 10 days: 11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureTimeOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				// Регистрируем все встроенные функции
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeTimeFunctions(scope.GlobalScope)
				
				exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
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

func TestTimeArithmetic(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "add_months",
			code: `
				let timestamp = 1609459200  // 2021-01-01 00:00:00 UTC
				let time = timeFromUnix(timestamp)
				let newTime = timeAddMonths(time, 2)
				let month = timeMonth(newTime)
				println("Month: " + month.toString())
			`,
			want: "Month: 3",
		},
		{
			name: "add_years",
			code: `
				let timestamp = 1609459200  // 2021-01-01 00:00:00 UTC
				let time = timeFromUnix(timestamp)
				let newTime = timeAddYears(time, 2)
				let year = timeYear(newTime)
				println("Year: " + year.toString())
			`,
			want: "Year: 2023",
		},
		{
			name: "add_hours",
			code: `
				let timestamp = 1609459200  // 2021-01-01 00:00:00 UTC
				let time = timeFromUnix(timestamp)
				let newTime = timeAddHours(time, 5)
				let hour = timeHour(newTime)
				println("Hour: " + hour.toString())
			`,
			want: "Hour: 10", // UTC+5 временная зона + 5 часов = 10
		},
		{
			name: "time_difference",
			code: `
				let timestamp1 = 1609459200  // 2021-01-01 00:00:00 UTC
				let timestamp2 = 1609459260  // 2021-01-01 00:01:00 UTC (60 seconds later)
				let time1 = timeFromUnix(timestamp1)
				let time2 = timeFromUnix(timestamp2)
				let diff = timeDiff(time2, time1)
				println("Difference in seconds: " + diff.toString())
			`,
			want: "Difference in seconds: 60",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureTimeOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeTimeFunctions(scope.GlobalScope)
				
				exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
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

func TestTimeComparison(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "time_before",
			code: `
				let timestamp1 = 1609459200  // 2021-01-01 00:00:00 UTC
				let timestamp2 = 1609459260  // 2021-01-01 00:01:00 UTC
				let time1 = timeFromUnix(timestamp1)
				let time2 = timeFromUnix(timestamp2)
				let isBefore = timeBefore(time1, time2)
				println("Time1 is before Time2: " + isBefore.toString())
			`,
			want: "Time1 is before Time2: true",
		},
		{
			name: "time_after",
			code: `
				let timestamp1 = 1609459260  // 2021-01-01 00:01:00 UTC
				let timestamp2 = 1609459200  // 2021-01-01 00:00:00 UTC
				let time1 = timeFromUnix(timestamp1)
				let time2 = timeFromUnix(timestamp2)
				let isAfter = timeAfter(time1, time2)
				println("Time1 is after Time2: " + isAfter.toString())
			`,
			want: "Time1 is after Time2: true",
		},
		{
			name: "time_equal",
			code: `
				let timestamp = 1609459200
				let time1 = timeFromUnix(timestamp)
				let time2 = timeFromUnix(timestamp)
				let isEqual = timeEqual(time1, time2)
				println("Times are equal: " + isEqual.toString())
			`,
			want: "Times are equal: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureTimeOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeTimeFunctions(scope.GlobalScope)
				
				exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
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

func TestTimeFormatting(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "custom_format",
			code: `
				let timestamp = 1609459200  // 2021-01-01 00:00:00 UTC
				let time = timeFromUnix(timestamp)
				let formatted = timeFormat(time, "YYYY-MM-DD HH:mm:ss")
				println(formatted)
			`,
			want: "2021-01-01 05:00:00", // UTC+5 временная зона
		},
		{
			name: "time_only_format",
			code: `
				let timestamp = 1609462800  // 2021-01-01 01:00:00 UTC
				let time = timeFromUnix(timestamp)
				let formatted = timeFormat(time, "time")
				println(formatted)
			`,
			want: "06:00:00", // 01:00 UTC + 5 часов = 06:00
		},
		{
			name: "datetime_format",
			code: `
				let timestamp = 1609462800  // 2021-01-01 01:00:00 UTC
				let time = timeFromUnix(timestamp)
				let formatted = timeFormat(time, "datetime")
				println(formatted)
			`,
			want: "2021-01-01 06:00:00", // 01:00 UTC + 5 часов = 06:00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureTimeOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeTimeFunctions(scope.GlobalScope)
				
				exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
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

func TestTimeDirectAPI(t *testing.T) {
	t.Run("time_value_creation", func(t *testing.T) {
		now := time.Now()
		timeValue := value.NewTime(now)
		
		if !timeValue.IsTime() {
			t.Error("expected value to be time type")
		}
		
		retrievedTime := timeValue.Time()
		if !now.Equal(retrievedTime) {
			t.Errorf("expected %v, got %v", now, retrievedTime)
		}
		
		typeName := value.GetValueTypeName(timeValue)
		if typeName != "time" {
			t.Errorf("expected type name 'time', got %s", typeName)
		}
	})
	
	t.Run("time_unix_timestamp", func(t *testing.T) {
		// Тестируем конкретное время: 2021-01-01 00:00:00 UTC
		expected := time.Unix(1609459200, 0)
		timeValue := value.NewTime(expected)
		
		if !timeValue.IsTime() {
			t.Error("expected value to be time type")
		}
		
		retrieved := timeValue.Time()
		if !expected.Equal(retrieved) {
			t.Errorf("expected %v, got %v", expected, retrieved)
		}
		
		// Проверяем компоненты
		if retrieved.Year() != 2021 {
			t.Errorf("expected year 2021, got %d", retrieved.Year())
		}
		if retrieved.Month() != time.January {
			t.Errorf("expected January, got %v", retrieved.Month())
		}
		if retrieved.Day() != 1 {
			t.Errorf("expected day 1, got %d", retrieved.Day())
		}
	})
}

func TestTimeErrorHandling(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "invalid_argument_count_now",
			code: `
				let result = now(123)
				println(result)
			`,
			want: "Error: now() requires 0 arguments",
		},
		{
			name: "invalid_timeFromUnix_arg",
			code: `
				let result = timeFromUnix("invalid")
				println(result)
			`,
			want: "Error: timeFromUnix() requires numeric argument",
		},
		{
			name: "invalid_timeFormat_arg",
			code: `
				let result = timeFormat("not_a_time", "date")
				println(result)
			`,
			want: "Error: first argument must be a time value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureTimeOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeTimeFunctions(scope.GlobalScope)
				
				exprs := parser.NewParser(tt.code).ParseWithoutScopeInit()
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

// captureTimeOutput захватывает stdout для тестирования времени
func captureTimeOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

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