package builtin

import (
	"bytes"
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// Exec выполняет команду и возвращает результат
func Exec(args []*value.Value) (*value.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("exec expects at least 1 argument")
	}

	if !args[0].IsString() {
		return nil, fmt.Errorf("exec command must be a string")
	}

	command := args[0].String()
	
	// Дополнительные аргументы команды
	cmdArgs := make([]string, 0)
	for i := 1; i < len(args); i++ {
		if !args[i].IsString() {
			return nil, fmt.Errorf("exec arguments must be strings")
		}
		cmdArgs = append(cmdArgs, args[i].String())
	}

	// Определяем shell в зависимости от ОС
	var cmd *exec.Cmd
	if len(cmdArgs) > 0 {
		cmd = exec.Command(command, cmdArgs...)
	} else {
		// Если нет аргументов, пробуем выполнить через shell
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	
	// Создаем объект с результатом
	resultMap := map[string]*value.Value{
		"stdout":   value.NewValue(stdout.String()),
		"stderr":   value.NewValue(stderr.String()),
		"success":  value.NewValue(err == nil),
		"exitCode": value.NewValue(int64(0)),
	}
	result := value.NewValue(resultMap)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				resultMap["exitCode"] = value.NewValue(int64(status.ExitStatus()))
			}
		}
		resultMap["error"] = value.NewValue(err.Error())
	}

	return result, nil
}

// Spawn запускает процесс в фоне и возвращает PID
func Spawn(args []*value.Value) (*value.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("spawn expects at least 1 argument")
	}

	if !args[0].IsString() {
		return nil, fmt.Errorf("spawn command must be a string")
	}

	command := args[0].String()
	
	// Дополнительные аргументы команды
	cmdArgs := make([]string, 0)
	for i := 1; i < len(args); i++ {
		if !args[i].IsString() {
			return nil, fmt.Errorf("spawn arguments must be strings")
		}
		cmdArgs = append(cmdArgs, args[i].String())
	}

	var cmd *exec.Cmd
	if len(cmdArgs) > 0 {
		cmd = exec.Command(command, cmdArgs...)
	} else {
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
	}

	// Запускаем процесс в фоне
	err := cmd.Start()
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to spawn process: %v", err)), nil
	}

	// Возвращаем PID процесса
	return value.NewValue(int64(cmd.Process.Pid)), nil
}

// Kill завершает процесс по PID
func Kill(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("kill expects 1 argument")
	}

	if !args[0].IsNumber() {
		return nil, fmt.Errorf("kill pid must be an integer")
	}

	pid := args[0].Int()
	
	process, err := os.FindProcess(pid)
	if err != nil {
		return value.NewBool(false), nil
	}

	err = process.Kill()
	return value.NewValue(err == nil), nil
}

// GetEnv получает переменную окружения
func GetEnv(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("getEnv expects 1 argument")
	}

	if !args[0].IsString() {
		return nil, fmt.Errorf("getEnv key must be a string")
	}

	key := args[0].String()
	val := os.Getenv(key)
	
	if val == "" {
		return value.NewValue(nil), nil
	}
	
	return value.NewValue(val), nil
}

// SetEnv устанавливает переменную окружения
func SetEnv(args []*value.Value) (*value.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("setEnv expects 2 arguments")
	}

	if !args[0].IsString() {
		return nil, fmt.Errorf("setEnv key must be a string")
	}
	
	if !args[1].IsString() {
		return nil, fmt.Errorf("setEnv value must be a string")
	}

	key := args[0].String()
	val := args[1].String()
	
	err := os.Setenv(key, val)
	return value.NewValue(err == nil), nil
}

// GetAllEnv возвращает все переменные окружения
func GetAllEnv(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getAllEnv expects 0 arguments")
	}

	env := os.Environ()
	result := make(map[string]*value.Value)
	
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			result[pair[0]] = value.NewValue(pair[1])
		}
	}
	
	return value.NewValue(result), nil
}

// Exit выход из программы с кодом
func Exit(args []*value.Value) (*value.Value, error) {
	exitCode := 0
	
	if len(args) == 1 {
		if !args[0].IsNumber() {
			return nil, fmt.Errorf("exit code must be an integer")
		}
		exitCode = args[0].Int()
	} else if len(args) > 1 {
		return nil, fmt.Errorf("exit expects 0 or 1 arguments")
	}
	
	os.Exit(exitCode)
	return value.NewValue(nil), nil // никогда не выполнится
}

// GetWorkingDir получает текущую рабочую директорию
func GetWorkingDir(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getWorkingDir expects 0 arguments")
	}

	dir, err := os.Getwd()
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to get working directory: %v", err)), nil
	}
	
	return value.NewValue(dir), nil
}

// ChangeDir меняет текущую рабочую директорию
func ChangeDir(args []*value.Value) (*value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("changeDir expects 1 argument")
	}

	if !args[0].IsString() {
		return nil, fmt.Errorf("changeDir path must be a string")
	}

	path := args[0].String()
	err := os.Chdir(path)
	
	return value.NewValue(err == nil), nil
}

// GetPid возвращает PID текущего процесса
func GetPid(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getPid expects 0 arguments")
	}
	
	return value.NewValue(int64(os.Getpid())), nil
}

// GetHostname возвращает имя хоста
func GetHostname(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getHostname expects 0 arguments")
	}
	
	hostname, err := os.Hostname()
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error: failed to get hostname: %v", err)), nil
	}
	
	return value.NewString(hostname), nil
}

// GetOS возвращает информацию об операционной системе
func GetOS(args []*value.Value) (*value.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("getOS expects 0 arguments")
	}
	
	resultMap := map[string]*value.Value{
		"os":      value.NewValue(runtime.GOOS),
		"arch":    value.NewValue(runtime.GOARCH),
		"cpus":    value.NewValue(int64(runtime.NumCPU())),
		"version": value.NewValue(runtime.Version()),
	}
	
	return value.NewValue(resultMap), nil
}

// InitializeProcessFunctions инициализирует встроенные функции процессов
func InitializeProcessFunctions(globalScope *scope.ScopeStack) {
	// Адаптируем функции для совместимости со старым API
	
	// getOS функция
	getOSFunc := func(args []*value.Value) *value.Value {
		result, err := GetOS(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("getOS", value.NewValue(getOSFunc))
	
	// getPid функция  
	getPidFunc := func(args []*value.Value) *value.Value {
		result, err := GetPid(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("getPid", value.NewValue(getPidFunc))
	
	// getEnv функция
	getEnvFunc := func(args []*value.Value) *value.Value {
		result, err := GetEnv(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("getEnv", value.NewValue(getEnvFunc))
	
	// exec функция
	execFunc := func(args []*value.Value) *value.Value {
		result, err := Exec(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("exec", value.NewValue(execFunc))
	
	// getWorkingDir функция
	getWorkingDirFunc := func(args []*value.Value) *value.Value {
		result, err := GetWorkingDir(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("getWorkingDir", value.NewValue(getWorkingDirFunc))
	
	// spawn функция
	spawnFunc := func(args []*value.Value) *value.Value {
		result, err := Spawn(args)
		if err != nil {
			return value.NewValue("Error: " + err.Error())
		}
		return result
	}
	globalScope.Set("spawn", value.NewValue(spawnFunc))
}