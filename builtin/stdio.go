package builtin

import (
	"bufio"
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"os"
	"strconv"
	"strings"
)

var stdinReader = bufio.NewReader(os.Stdin)

// ReadLine читает строку из stdin
func ReadLine(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("readLine expects 0 arguments, got %d", len(args))
	}

	line, err := stdinReader.ReadString('\n')
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: reading line: %v", err)), nil
	}

	// Удаляем символ новой строки
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")

	return value.NewValue(line), nil
}

// Input читает строку из stdin с опциональным prompt
func Input(args []*value.Value) (*value.Value, error) {
	prompt := ""
	if len(args) == 1 {
		if !args[0].IsString() {
			return nil, fmt.Errorf("input prompt must be a string")
		}
		prompt = args[0].String()
	} else if len(args) > 1 {
		return nil, fmt.Errorf("input expects 0 or 1 arguments, got %d", len(args))
	}

	// Выводим prompt если есть
	if prompt != "" {
		fmt.Print(prompt)
	}

	line, err := stdinReader.ReadString('\n')
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: reading input: %v", err)), nil
	}

	// Удаляем символ новой строки
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")

	return value.NewValue(line), nil
}

// InputNumber читает число из stdin
func InputNumber(args []*value.Value) (*value.Value, error) {
	prompt := ""
	if len(args) == 1 {
		if !args[0].IsString() {
			return nil, fmt.Errorf("inputNumber prompt must be a string")
		}
		prompt = args[0].String()
	} else if len(args) > 1 {
		return nil, fmt.Errorf("inputNumber expects 0 or 1 arguments, got %d", len(args))
	}

	// Выводим prompt если есть
	if prompt != "" {
		fmt.Print(prompt)
	}

	line, err := stdinReader.ReadString('\n')
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: reading input: %v", err)), nil
	}

	// Удаляем символ новой строки
	line = strings.TrimSpace(line)

	// Пробуем парсить как float
	if floatVal, err := strconv.ParseFloat(line, 64); err == nil {
		// Проверяем, является ли это целым числом
		if floatVal == float64(int64(floatVal)) {
			return value.NewValue(int64(floatVal)), nil
		}
		return value.NewValue(floatVal), nil
	}

	// Пробуем парсить как int
	if intVal, err := strconv.ParseInt(line, 10, 64); err == nil {
		return value.NewValue(intVal), nil
	}

	return value.NewValue(fmt.Sprintf("Error: invalid number: %s", line)), nil
}

// Printf форматированный вывод в stdout
func Printf(args []*value.Value) (*value.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("printf expects at least 1 argument")
	}

	if !args[0].IsString() {
		return nil, fmt.Errorf("printf format must be a string")
	}

	format := args[0].String()
	printArgs := make([]interface{}, 0, len(args)-1)

	// Конвертируем аргументы в interface{} для fmt.Printf
	for i := 1; i < len(args); i++ {
		if args[i].IsNumber() {
			if args[i].IsFloat64() {
				printArgs = append(printArgs, args[i].Float64())
			} else {
				printArgs = append(printArgs, args[i].Int64())
			}
		} else if args[i].IsBool() {
			printArgs = append(printArgs, args[i].Bool())
		} else if args[i].IsString() {
			printArgs = append(printArgs, args[i].String())
		} else if args[i].Any() == nil {
			printArgs = append(printArgs, "null")
		} else {
			printArgs = append(printArgs, args[i].String())
		}
	}

	fmt.Printf(format, printArgs...)
	return value.NewValue(nil), nil
}

// GetChar читает один символ из stdin
func GetChar(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getChar expects 0 arguments, got %d", len(args))
	}

	char, err := stdinReader.ReadByte()
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: reading char: %v", err)), nil
	}

	return value.NewString(string(char)), nil
}

// PutChar выводит один символ в stdout
func PutChar(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("putChar expects 1 argument, got %d", len(args))
	}

	var char string
	if args[0].IsString() {
		str := args[0].String()
		if len(str) > 0 {
			char = string(str[0])
		} else {
			return nil, fmt.Errorf("putChar: empty string")
		}
	} else if args[0].IsNumber() {
		// Принимаем ASCII код
		code := args[0].Int64()
		if code >= 0 && code <= 127 {
			char = string(rune(code))
		} else {
			return nil, fmt.Errorf("putChar: invalid ASCII code %d", code)
		}
	} else {
		return nil, fmt.Errorf("putChar expects a string or int")
	}

	fmt.Print(char)
	return value.NewValue(nil), nil
}

// WriteLn выводит строку с переносом строки
func WriteLn(args []*value.Value) (*value.Value, error) {
	if len(args) == 0 {
		fmt.Println()
		return value.NewValue(nil), nil
	}

	values := make([]interface{}, len(args))
	for i, arg := range args {
		if arg.IsNumber() {
			if arg.IsFloat64() {
				values[i] = arg.Float64()
			} else {
				values[i] = arg.Int64()
			}
		} else if arg.IsBool() {
			values[i] = arg.Bool()
		} else if arg.IsString() {
			values[i] = arg.String()
		} else if arg.Any() == nil {
			values[i] = "null"
		} else {
			values[i] = arg.String()
		}
	}

	fmt.Println(values...)
	return value.NewValue(nil), nil
}

// Write выводит без переноса строки
func Write(args []*value.Value) (*value.Value, error) {
	for _, arg := range args {
		if arg.IsNumber() {
			if arg.IsFloat64() {
				fmt.Print(arg.Float64())
			} else {
				fmt.Print(arg.Int64())
			}
		} else if arg.IsBool() {
			fmt.Print(arg.Bool())
		} else if arg.IsString() {
			fmt.Print(arg.String())
		} else if arg.Any() == nil {
			fmt.Print("null")
		} else {
			fmt.Print(arg.String())
		}
	}
	return value.NewValue(nil), nil
}

// InitializeStdioFunctions инициализирует встроенные STDIO функции
func InitializeStdioFunctions(globalScope *scope.ScopeStack) {
	// printf функция (адаптированная)
	printfFunc := func(args []*value.Value) *value.Value {
		result, err := Printf(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("printf", value.NewValue(printfFunc))
	
	// writeLn функция
	writeLnFunc := func(args []*value.Value) *value.Value {
		result, err := WriteLn(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("writeLn", value.NewValue(writeLnFunc))
	
	// input функция
	inputFunc := func(args []*value.Value) *value.Value {
		result, err := Input(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("input", value.NewValue(inputFunc))
	
	// readLine функция
	readLineFunc := func(args []*value.Value) *value.Value {
		result, err := ReadLine(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("readLine", value.NewValue(readLineFunc))
}