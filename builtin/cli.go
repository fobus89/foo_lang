package builtin

import (
	"fmt"
	"foo_lang/value"
	"os"
	"path/filepath"
	"strings"
)

// Глобальные переменные для хранения аргументов
var (
	cliArgs      []string
	scriptName   string
	isInitialized bool
)

// InitCLI инициализирует CLI аргументы
func InitCLI(args []string) {
	if len(args) > 0 {
		scriptName = args[0]
		if len(args) > 1 {
			cliArgs = args[1:]
		} else {
			cliArgs = []string{}
		}
	} else {
		scriptName = ""
		cliArgs = []string{}
	}
	isInitialized = true
}

// GetArgs возвращает массив аргументов командной строки
func GetArgs(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getArgs expects 0 arguments")
	}

	if !isInitialized {
		// Если не инициализировано, берем из os.Args
		InitCLI(os.Args)
	}

	// Создаем массив значений
	result := make([]*value.Value, len(cliArgs))
	for i, arg := range cliArgs {
		result[i] = value.NewValue(arg)
	}

	return value.NewValue(result), nil
}

// GetArg получает аргумент по индексу
func GetArg(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("getArg expects 1 argument")
	}

	if !args[0].IsNumber() {
		return nil, fmt.Errorf("getArg index must be a number")
	}

	index := args[0].Int()

	if !isInitialized {
		InitCLI(os.Args)
	}

	if index < 0 || index >= len(cliArgs) {
		return value.NewValue(nil), nil
	}

	return value.NewValue(cliArgs[index]), nil
}

// GetArgCount возвращает количество аргументов
func GetArgCount(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getArgCount expects 0 arguments")
	}

	if !isInitialized {
		InitCLI(os.Args)
	}

	return value.NewValue(int64(len(cliArgs))), nil
}

// GetScriptName возвращает имя скрипта
func GetScriptName(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getScriptName expects 0 arguments")
	}

	if !isInitialized {
		InitCLI(os.Args)
	}

	return value.NewValue(scriptName), nil
}

// GetScriptPath возвращает полный путь к скрипту
func GetScriptPath(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getScriptPath expects 0 arguments")
	}

	if !isInitialized {
		InitCLI(os.Args)
	}

	if scriptName == "" {
		return value.NewValue(""), nil
	}

	absPath, err := filepath.Abs(scriptName)
	if err != nil {
		return value.NewValue(scriptName), nil
	}

	return value.NewValue(absPath), nil
}

// GetScriptDir возвращает директорию скрипта
func GetScriptDir(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getScriptDir expects 0 arguments")
	}

	if !isInitialized {
		InitCLI(os.Args)
	}

	if scriptName == "" {
		return value.NewValue(""), nil
	}

	absPath, err := filepath.Abs(scriptName)
	if err != nil {
		return value.NewValue(""), nil
	}

	dir := filepath.Dir(absPath)
	return value.NewValue(dir), nil
}

// HasArg проверяет наличие аргумента
func HasArg(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hasArg expects 1 argument")
	}

	if !args[0].IsString() {
		return nil, fmt.Errorf("hasArg argument must be a string")
	}

	searchArg := args[0].String()

	if !isInitialized {
		InitCLI(os.Args)
	}

	for _, arg := range cliArgs {
		if arg == searchArg {
			return value.NewValue(true), nil
		}
	}

	return value.NewValue(false), nil
}

// GetFlag получает значение флага (--key=value или --key value)
func GetFlag(args []*value.Value) (*value.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("getFlag expects 1 or 2 arguments")
	}

	if !args[0].IsString() {
		return nil, fmt.Errorf("getFlag flag name must be a string")
	}

	flagName := args[0].String()
	defaultValue := (*value.Value)(nil)
	
	if len(args) == 2 {
		defaultValue = args[1]
	}

	if !isInitialized {
		InitCLI(os.Args)
	}

	// Ищем флаг в аргументах
	for i, arg := range cliArgs {
		// Проверяем формат --key=value
		if len(arg) > len(flagName)+3 && arg[:len(flagName)+3] == "--"+flagName+"=" {
			return value.NewValue(arg[len(flagName)+3:]), nil
		}
		
		// Проверяем формат --key value
		if arg == "--"+flagName {
			if i+1 < len(cliArgs) && !isFlag(cliArgs[i+1]) {
				return value.NewValue(cliArgs[i+1]), nil
			}
			// Флаг без значения - возвращаем true
			return value.NewValue(true), nil
		}
		
		// Проверяем формат -k (короткий флаг)
		if len(flagName) == 1 && arg == "-"+flagName {
			if i+1 < len(cliArgs) && !isFlag(cliArgs[i+1]) {
				return value.NewValue(cliArgs[i+1]), nil
			}
			// Флаг без значения - возвращаем true
			return value.NewValue(true), nil
		}
	}

	// Флаг не найден, возвращаем значение по умолчанию
	if defaultValue != nil {
		return defaultValue, nil
	}
	
	return value.NewValue(nil), nil
}

// isFlag проверяет, является ли строка флагом
func isFlag(s string) bool {
	return len(s) > 0 && (s[0] == '-')
}

// ParseArgs парсит аргументы в объект с флагами и позиционными аргументами
func ParseArgs(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("parseArgs expects 0 arguments")
	}

	if !isInitialized {
		InitCLI(os.Args)
	}

	flags := make(map[string]*value.Value)
	positional := make([]*value.Value, 0)

	for i := 0; i < len(cliArgs); i++ {
		arg := cliArgs[i]
		
		// Проверяем формат --key=value
		if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			parts := strings.SplitN(arg[2:], "=", 2)
			flags[parts[0]] = value.NewValue(parts[1])
			continue
		}
		
		// Проверяем формат --key value или --key
		if strings.HasPrefix(arg, "--") {
			key := arg[2:]
			if i+1 < len(cliArgs) && !isFlag(cliArgs[i+1]) {
				flags[key] = value.NewValue(cliArgs[i+1])
				i++ // Пропускаем следующий аргумент
			} else {
				flags[key] = value.NewValue(true)
			}
			continue
		}
		
		// Проверяем формат -k value или -k
		if strings.HasPrefix(arg, "-") && len(arg) == 2 {
			key := string(arg[1])
			if i+1 < len(cliArgs) && !isFlag(cliArgs[i+1]) {
				flags[key] = value.NewValue(cliArgs[i+1])
				i++ // Пропускаем следующий аргумент
			} else {
				flags[key] = value.NewValue(true)
			}
			continue
		}
		
		// Позиционный аргумент
		positional = append(positional, value.NewValue(arg))
	}

	result := map[string]*value.Value{
		"flags":      value.NewValue(flags),
		"positional": value.NewValue(positional),
		"all":        value.NewValue(cliArgs),
	}

	return value.NewValue(result), nil
}