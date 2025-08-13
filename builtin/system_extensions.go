package builtin

import (
	"foo_lang/scope"
	"foo_lang/value"
)

// InitializeSystemExtensions создает extension methods для системных типов
func InitializeSystemExtensions(globalScope *scope.ScopeStack) {
	// Регистрируем System как расширение для глобального объекта
	RegisterSystemExtension()
	RegisterIOExtension()
	RegisterConsoleExtension()
	RegisterProcessExtension()
	RegisterDebugExtension()
	RegisterMemoryExtension()
}

// System extension - системные функции как методы
func RegisterSystemExtension() {
	// System.getOS()
	systemGetOS := &SystemExtensionMethod{
		Name: "getOS",
		Func: func(args []*value.Value) *value.Value {
			result, err := GetOS(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("System", "getOS", systemGetOS)
	
	// System.getEnv(key)
	systemGetEnv := &SystemExtensionMethod{
		Name: "getEnv",
		Func: func(args []*value.Value) *value.Value {
			result, err := GetEnv(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("System", "getEnv", systemGetEnv)
	
	// System.setEnv(key, value)
	systemSetEnv := &SystemExtensionMethod{
		Name: "setEnv",
		Func: func(args []*value.Value) *value.Value {
			result, err := SetEnv(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("System", "setEnv", systemSetEnv)
	
	// System.exit(code)
	systemExit := &SystemExtensionMethod{
		Name: "exit",
		Func: func(args []*value.Value) *value.Value {
			result, err := Exit(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("System", "exit", systemExit)
}

// IO extension - функции ввода/вывода
func RegisterIOExtension() {
	// IO.input(prompt)
	ioInput := &SystemExtensionMethod{
		Name: "input",
		Func: func(args []*value.Value) *value.Value {
			result, err := Input(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("IO", "input", ioInput)
	
	// IO.readLine()
	ioReadLine := &SystemExtensionMethod{
		Name: "readLine",
		Func: func(args []*value.Value) *value.Value {
			result, err := ReadLine(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("IO", "readLine", ioReadLine)
	
	// IO.inputNumber(prompt)
	ioInputNumber := &SystemExtensionMethod{
		Name: "inputNumber",
		Func: func(args []*value.Value) *value.Value {
			result, err := InputNumber(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("IO", "inputNumber", ioInputNumber)
	
	// IO.write(...)
	ioWrite := &SystemExtensionMethod{
		Name: "write",
		Func: func(args []*value.Value) *value.Value {
			result, err := Write(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("IO", "write", ioWrite)
	
	// IO.writeLn(...)
	ioWriteLn := &SystemExtensionMethod{
		Name: "writeLn",
		Func: func(args []*value.Value) *value.Value {
			result, err := WriteLn(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("IO", "writeLn", ioWriteLn)
}

// Console extension - консольные функции
func RegisterConsoleExtension() {
	// Console.printf(format, ...)
	consolePrintf := &SystemExtensionMethod{
		Name: "printf",
		Func: func(args []*value.Value) *value.Value {
			result, err := Printf(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Console", "printf", consolePrintf)
	
	// Console.getChar()
	consoleGetChar := &SystemExtensionMethod{
		Name: "getChar",
		Func: func(args []*value.Value) *value.Value {
			result, err := GetChar(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Console", "getChar", consoleGetChar)
	
	// Console.putChar(char)
	consolePutChar := &SystemExtensionMethod{
		Name: "putChar",
		Func: func(args []*value.Value) *value.Value {
			result, err := PutChar(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Console", "putChar", consolePutChar)
}

// Process extension - функции процессов
func RegisterProcessExtension() {
	// Process.exec(command, ...args)
	processExec := &SystemExtensionMethod{
		Name: "exec",
		Func: func(args []*value.Value) *value.Value {
			result, err := Exec(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Process", "exec", processExec)
	
	// Process.spawn(command, ...args)
	processSpawn := &SystemExtensionMethod{
		Name: "spawn",
		Func: func(args []*value.Value) *value.Value {
			result, err := Spawn(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Process", "spawn", processSpawn)
	
	// Process.kill(pid)
	processKill := &SystemExtensionMethod{
		Name: "kill",
		Func: func(args []*value.Value) *value.Value {
			result, err := Kill(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Process", "kill", processKill)
	
	// Process.getPid()
	processGetPid := &SystemExtensionMethod{
		Name: "getPid",
		Func: func(args []*value.Value) *value.Value {
			result, err := GetPid(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Process", "getPid", processGetPid)
}

// Debug extension - функции отладки
func RegisterDebugExtension() {
	// Debug.debug(value)
	debugDebug := &SystemExtensionMethod{
		Name: "debug",
		Func: func(args []*value.Value) *value.Value {
			result, err := Debug(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Debug", "debug", debugDebug)
	
	// Debug.trace(depth)
	debugTrace := &SystemExtensionMethod{
		Name: "trace",
		Func: func(args []*value.Value) *value.Value {
			result, err := Trace(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Debug", "trace", debugTrace)
	
	// Debug.typeOf(value)
	debugTypeOf := &SystemExtensionMethod{
		Name: "typeOf",
		Func: func(args []*value.Value) *value.Value {
			result, err := TypeOf(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Debug", "typeOf", debugTypeOf)
	
	// Debug.sizeOf(value)
	debugSizeOf := &SystemExtensionMethod{
		Name: "sizeOf",
		Func: func(args []*value.Value) *value.Value {
			result, err := SizeOf(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Debug", "sizeOf", debugSizeOf)
}

// Memory extension - функции памяти
func RegisterMemoryExtension() {
	// Memory.stats()
	memoryStats := &SystemExtensionMethod{
		Name: "stats",
		Func: func(args []*value.Value) *value.Value {
			result, err := MemStats(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Memory", "stats", memoryStats)
	
	// Memory.gc()
	memoryGC := &SystemExtensionMethod{
		Name: "gc",
		Func: func(args []*value.Value) *value.Value {
			result, err := GC(args)
			if err != nil {
				return value.NewValue("Error: " + err.Error())
			}
			return result
		},
	}
	value.RegisterExtensionMethod("Memory", "gc", memoryGC)
}

// SystemExtensionMethod реализует интерфейс extension method для системных функций
type SystemExtensionMethod struct {
	Name string
	Func func(args []*value.Value) *value.Value
}

func (s *SystemExtensionMethod) Call(receiver *value.Value, args []*value.Value) *value.Value {
	// Для системных extension methods receiver не используется
	return s.Func(args)
}