package builtin

import (
	"foo_lang/scope"
	"foo_lang/value"
)

// InitializeGlobalObjects создает глобальные объекты System, IO, Console, Process, Debug, Memory
func InitializeGlobalObjects(globalScope *scope.ScopeStack) {
	// Создаем глобальный IO объект
	ioObject := createIOObject()
	globalScope.Set("IO", value.NewValue(ioObject))
	
	// Создаем глобальный System объект
	systemObject := createSystemObject()
	globalScope.Set("System", value.NewValue(systemObject))
	
	// Создаем глобальный Console объект
	consoleObject := createConsoleObject()
	globalScope.Set("Console", value.NewValue(consoleObject))
	
	// Создаем глобальный Process объект
	processObject := createProcessObject()
	globalScope.Set("Process", value.NewValue(processObject))
	
	// Создаем глобальный Debug объект
	debugObject := createDebugObject()
	globalScope.Set("Debug", value.NewValue(debugObject))
	
	// Создаем глобальный Memory объект
	memoryObject := createMemoryObject()
	globalScope.Set("Memory", value.NewValue(memoryObject))
	
	// Создаем глобальный CLI объект
	cliObject := createCLIObject()
	globalScope.Set("CLI", value.NewValue(cliObject))
}

// createIOObject создает объект IO с методами ввода/вывода
func createIOObject() map[string]*value.Value {
	ioObject := map[string]*value.Value{
		// Методы ввода
		"input": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Input(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"readLine": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := ReadLine(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"inputNumber": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := InputNumber(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getChar": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetChar(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		// Методы вывода
		"write": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Write(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"writeLn": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := WriteLn(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"printf": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Printf(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"putChar": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := PutChar(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
	}
	
	return ioObject
}

// createSystemObject создает объект System с системными функциями
func createSystemObject() map[string]*value.Value {
	systemObject := map[string]*value.Value{
		"getOS": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetOS(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getEnv": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetEnv(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"setEnv": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := SetEnv(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getAllEnv": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetAllEnv(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getWorkingDir": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetWorkingDir(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"changeDir": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := ChangeDir(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getHostname": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetHostname(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"exit": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Exit(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
	}
	
	return systemObject
}

// createConsoleObject создает объект Console для консольных операций
func createConsoleObject() map[string]*value.Value {
	consoleObject := map[string]*value.Value{
		"printf": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Printf(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"writeLn": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := WriteLn(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"write": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Write(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"input": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Input(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"readLine": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := ReadLine(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
	}
	
	return consoleObject
}

// createProcessObject создает объект Process для управления процессами
func createProcessObject() map[string]*value.Value {
	processObject := map[string]*value.Value{
		"exec": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Exec(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"spawn": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Spawn(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"kill": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Kill(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getPid": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetPid(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
	}
	
	return processObject
}

// createDebugObject создает объект Debug для отладки
func createDebugObject() map[string]*value.Value {
	debugObject := map[string]*value.Value{
		"debug": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Debug(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"trace": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Trace(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"typeOf": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := TypeOf(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"sizeOf": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := SizeOf(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"assert": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Assert(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"profile": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Profile(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"benchmark": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := Benchmark(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
	}
	
	return debugObject
}

// createMemoryObject создает объект Memory для работы с памятью
func createMemoryObject() map[string]*value.Value {
	memoryObject := map[string]*value.Value{
		"stats": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := MemStats(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"gc": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GC(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
	}
	
	return memoryObject
}

// createCLIObject создает объект CLI для работы с аргументами командной строки
func createCLIObject() map[string]*value.Value {
	cliObject := map[string]*value.Value{
		"getArgs": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetArgs(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getArg": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetArg(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getArgCount": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetArgCount(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getScriptName": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetScriptName(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getScriptPath": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetScriptPath(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getScriptDir": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetScriptDir(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"hasArg": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := HasArg(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"getFlag": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := GetFlag(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
		
		"parseArgs": value.NewValue(func(args []*value.Value) *value.Value {
			result, err := ParseArgs(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		}),
	}
	
	return cliObject
}