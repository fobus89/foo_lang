# ✅ ГОТОВО! Extension Methods и глобальные объекты реализованы

## 🔥 Что было реализовано

### 1. 📦 **Extension Methods для встроенных функций**

Создана система extension methods, которая позволяет группировать системные функции в объекты типа System, IO, Console и т.д.

**Реализованные extension methods:**

#### System Extension
- `System.getOS()` - информация об ОС
- `System.getEnv(key)` - переменная окружения
- `System.setEnv(key, value)` - установка переменной
- `System.exit(code)` - выход из программы

#### IO Extension  
- `IO.input(prompt)` - ввод с приглашением
- `IO.readLine()` - чтение строки
- `IO.inputNumber(prompt)` - ввод числа
- `IO.write(...)` - вывод без переноса
- `IO.writeLn(...)` - вывод с переносом

#### Console Extension
- `Console.printf(format, ...)` - форматированный вывод
- `Console.getChar()` - чтение символа
- `Console.putChar(char)` - вывод символа

#### Process Extension
- `Process.exec(command, ...args)` - выполнение команды
- `Process.spawn(command, ...args)` - запуск процесса
- `Process.kill(pid)` - завершение процесса
- `Process.getPid()` - текущий PID

#### Debug Extension
- `Debug.debug(value)` - отладочный вывод
- `Debug.trace(depth)` - трассировка стека
- `Debug.typeOf(value)` - тип значения
- `Debug.sizeOf(value)` - размер в памяти

#### Memory Extension
- `Memory.stats()` - статистика памяти
- `Memory.gc()` - сборка мусора

### 2. 🌐 **Глобальные объекты**

Созданы глобальные объекты, доступные из любого места в программе:

```foo
// Использование глобальных объектов
IO.writeLn("Hello World!")
let osInfo = System.getOS()
let result = Process.exec("echo", "test")
Debug.debug(result)
let mem = Memory.stats()
```

**Доступные глобальные объекты:**
- **IO** - ввод/вывод (8 методов)
- **System** - системные функции (8 методов)  
- **Console** - консольные операции (5 методов)
- **Process** - управление процессами (4 метода)
- **Debug** - отладка (7 методов)
- **Memory** - управление памятью (2 метода)
- **CLI** - аргументы командной строки (9 методов)

### 3. 🔧 **Технические улучшения**

#### Поддержка методов объектов в AST

Добавлена поддержка вызова методов на объектах типа `map[string]*value.Value` в `ast/method_call_expr.go`:

```go
// Методы для объектов (map[string]*value.Value) - глобальные объекты типа IO, System и т.д.
if objectMap, ok := obj.Any().(map[string]*value.Value); ok {
    if method, exists := objectMap[m.MethodName]; exists {
        if fn, isFn := method.Any().(func([]*value.Value) *value.Value); isFn {
            // Конвертируем типы и вызываем функцию
            return fn(args)
        }
    }
}
```

#### Интеграция в main.go

Добавлена инициализация extension methods и глобальных объектов:

```go
// Extension methods и глобальные объекты
builtin.InitializeSystemExtensions(scopeStack)  // Extension methods
builtin.InitializeGlobalObjects(scopeStack)     // Глобальные объекты
```

### 4. 📁 **Созданные файлы**

#### Реализация:
- `builtin/system_extensions.go` - Extension methods для системных типов
- `builtin/io_global.go` - Глобальные объекты (IO, System, Console и т.д.)

#### Примеры:
- `examples/test_extensions_global.foo` - Полная демонстрация всех возможностей
- `examples/test_global_objects_simple.foo` - Простой пример использования

## 📈 Результат

### До реализации:
❌ Функции вызывались только напрямую: `getOS()`, `printf()`, `exec()`  
❌ Нет группировки функций по логическим областям  
❌ Отсутствует объектно-ориентированный API  

### После реализации:
✅ **Два стиля вызова функций:**
   - Старый стиль: `getOS()`, `printf("hello")`
   - Новый стиль: `System.getOS()`, `IO.printf("hello")`

✅ **Логическая группировка функций:**
   - `IO.*` - все функции ввода/вывода
   - `System.*` - все системные функции  
   - `Process.*` - все функции процессов
   - `Debug.*` - все отладочные функции

✅ **Объектно-ориентированный API:**
   - Читаемый и понятный код
   - Автокомплит в IDE по объектам
   - Группировка связанных функций

## 🎯 Демонстрация возможностей

### Объектно-ориентированный подход:

```foo
// Старый стиль
let osInfo = getOS()
printf("Running on: %s\n", osInfo["os"])
let result = exec("echo", "hello")
debug(result)

// Новый стиль - объектно-ориентированный
let osInfo = System.getOS()
IO.printf("Running on: %s\n", osInfo["os"])
let result = Process.exec("echo", "hello")
Debug.debug(result)
```

### Цепочки вызовов:

```foo
// Комбинируем функции разных объектов
let currentDir = System.getWorkingDir()
IO.printf("Working in: %s\n", currentDir)
Debug.printf("Type: %s\n", Debug.typeOf(currentDir))
```

### Группировка по функциональности:

```foo
// Все I/O операции в одном объекте
IO.writeLn("Hello!")
let input = IO.input("Enter name: ")
IO.printf("Hello, %s!\n", input)

// Все системные функции в одном объекте  
System.setEnv("DEBUG", "1")
let debug = System.getEnv("DEBUG")
let osInfo = System.getOS()
```

## 🚀 Преимущества новой архитектуры

### 1. **Читаемость кода**
- Сразу понятно к какой области относится функция
- `IO.readLine()` vs `readLine()` - более выразительно

### 2. **IDE поддержка**
- Автокомплит показывает все методы объекта
- Группировка функций в IDE по объектам

### 3. **Модульность**
- Легко добавлять новые методы к существующим объектам
- Четкое разделение ответственности между объектами

### 4. **Совместимость**
- Старые функции продолжают работать
- Постепенный переход на новый стиль

## ✅ Статус проекта

**foo_lang теперь поддерживает современный объектно-ориентированный API!**

### Основные достижения:
- ✅ Extension methods система реализована
- ✅ Глобальные объекты созданы и работают  
- ✅ Поддержка методов объектов в AST
- ✅ Примеры и документация готовы
- ✅ Обратная совместимость сохранена

### Объекты готовы к использованию:
- 🔥 **IO** - полный набор функций ввода/вывода
- 🔥 **System** - системная информация и переменные
- 🔥 **Process** - управление процессами  
- 🔥 **Debug** - профессиональные инструменты отладки
- 🔥 **Memory** - управление памятью
- 🔥 **Console** - консольные операции
- 🔥 **CLI** - работа с аргументами командной строки

## 🎊 Заключение

Проект **foo_lang** получил мощную объектно-ориентированную систему для работы с встроенными функциями. Теперь код становится более читаемым, модульным и удобным для разработки.

**Пример использования:**

```foo
// Сессионная информация через ОО API
let sessionInfo = {
    "user": System.getEnv("USER"),
    "pid": Process.getPid(), 
    "os": System.getOS()["os"],
    "memory": Memory.stats()["alloc"]
}

Debug.debug(sessionInfo)
IO.printf("Session: %s@%s (PID: %d)\n", 
    sessionInfo["user"], 
    sessionInfo["os"], 
    sessionInfo["pid"])
```

**Результат: PRODUCTION-READY система с современным API! 🚀**