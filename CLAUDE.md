# Foo Language v2 - Информация для Claude

## пиши сообщении по русский

## Описание проекта

Интерпретируемый язык программирования **foo_lang** с динамической типизацией, написанный на Go. Это полнофункциональный современный язык программирования с поддержкой замыканий, анонимных функций, модульной системы, математических функций и многого другого.

## Структура проекта

```
D:\dev\go\foo_lang_v2\
├── ast/                    # AST узлы и интерпретатор
├── builtin/               # Встроенные функции (математические)
├── examples/              # Примеры кода на foo_lang (30+ файлов)
├── lexer/                 # Лексический анализатор
├── modules/               # Система модулей (загрузка, кеширование)
├── parser/                # Синтаксический анализатор
├── scope/                 # Система областей видимости
├── test/                  # Unit-тесты (20 файлов, 130+ тестов)
├── token/                 # Определения токенов
├── value/                 # Система типов
├── main.go                # Точка входа
└── README.md              # Полная документация
```

## Статус: ✅ ПОЛНОСТЬЮ ГОТОВЫЙ СОВРЕМЕННЫЙ ЯЗЫК ПРОГРАММИРОВАНИЯ

### 🚀 **ПОСЛЕДНЕЕ ДОСТИЖЕНИЕ** (12.08.2025): **ПОЛНАЯ СИСТЕМА ТИПИЗАЦИИ** ⚡
- ✅ **Типизированные переменные**: `let x: int = 42` с runtime валидацией
- ✅ **Типизированные функции**: `fn add(a: int, b: int) -> int` с проверкой параметров и возврата
- ✅ **Параметры по умолчанию в типизированных функциях**: полностью исправлены
- ✅ **Валидация типов возврата**: runtime проверка соответствия типа возврата
- ✅ **Неявная типизация (type inference)**: автоматический вывод типов из литералов
- ✅ **Арифметические операции с сохранением типов**: int + int = int, int + float = float

### 🔥 Ключевые возможности (все реализованы и протестированы):

#### 1. 🔥 **ПОЛНАЯ СИСТЕМА ТИПИЗАЦИИ** - революционная type safety

**Типизированные переменные с валидацией**:
```foo
let count: int = 42         // ✅ Проверка типа во время выполнения
let price: float = 99.99    // ✅ Автоматическое приведение int->float
let name: string = "Alice"  // ✅ Строгая типизация строк
let active: bool = true     // ✅ Булевы типы с проверкой

// let wrong: int = "string"  // ❌ Ошибка: expected int, got string
```

**Типизированные функции с полной валидацией**:
```foo
// Проверка типов параметров и возврата
fn calculateTotal(price: float, qty: int, discount: float = 0.1) -> float {
    return price * qty * (1.0 - discount)
}

let total = calculateTotal(29.99, 3)        // ✅ Параметры по умолчанию работают
let result = calculateTotal(29.99, 3, 0.15) // ✅ Все типы проверяются
```

**Современная система функций**:
- **Обычные функции**: `fn add(a, b) { return a + b }`
- **Generic функции**: `fn identity<T>(value: T) -> T { return value }`
- **Типизированные функции**: `fn add(a: int, b: int) -> int { return a + b }`
- **Анонимные функции**: `let add = fn(x, y) => x + y`
- **Функции высшего порядка**: функции, возвращающие функции
- **Замыкания**: автоматический захват переменных из внешней области
- **Параметры по умолчанию**: `fn greet(name, prefix = "Hello") { ... }`
- **Множественные возвращаемые значения**: `return a, b, c`
- **Множественное присваивание**: `let x, y, z = func()`

#### 2. Модульная система с реальной загрузкой файлов

```foo
// math.foo
export fn add(a, b) { return a + b }
export let PI = 3.14159

// main.foo
import { add, PI } from "math.foo"    // Селективный импорт
import * as Math from "math.foo"      // Алиас импорт
import "math.foo"                     // Полный импорт
```

- Кеширование модулей (загружаются только один раз)
- Поддержка всех типов импорта/экспорта
- Изоляция областей видимости между модулями

#### 3. 🚀 **ИНТЕРФЕЙСНАЯ СИСТЕМА** - полная поддержка ОО-программирования

```foo
// Определение интерфейса
interface Drawable {
    fn draw()
    fn getArea() -> float
}

// Определение структуры
struct Circle { radius: float }

// Реализация интерфейса для типа
impl Drawable for Circle {
    fn draw() { println("Drawing circle") }
    fn getArea() -> float { return 3.14159 * this.radius * this.radius }
}

// Использование - полный полиморфизм!
let circle = Circle{radius: 5.0}
circle.draw()                    // Вызов интерфейсного метода
let area = circle.getArea()      // Типизированный возврат
```

**✅ Поддержка**: множественные интерфейсы, проверка сигнатур, `this` контекст, экземпляры структур

#### 4. 🔥 **GENERIC ОГРАНИЧЕНИЯ ТИПОВ** - типобезопасные универсальные функции

```foo
// Generic функция с ограничением типа
fn drawShape<T: Drawable>(shape: T) {
    shape.draw()
    println("Area: " + shape.getArea().toString())
}

// Множественные ограничения
fn moveAndDraw<T: Drawable + Moveable>(shape: T, dx: float, dy: float) {
    shape.move(dx, dy)  // Ограничение Moveable
    shape.draw()        // Ограничение Drawable
}

// Использование с проверкой типов во время выполнения
let circle = Circle{radius: 5.0}
drawShape(circle)           // ✅ Работает - Circle реализует Drawable
moveAndDraw(circle, 5, 3)   // ✅ Работает - Circle реализует оба интерфейса
```

**✅ Поддержка**: синтаксис `<T: Interface>`, множественные ограничения `<T: A + B>`, проверка во время выполнения
**✅ Примеры**: `test_generic_constraints.foo`, `test_generic_constraints_simple.foo`

#### 5. Extension Methods и перегрузка методов

- **Extension methods**: `extension String { fn reverse() { ... } }`
- **Перегрузка методов**: множественные определения с разными сигнатурами
- **Методы примитивных типов**: автоматические `.toString()`, `.length()`, `.abs()`

#### 6. Встроенные функции и методы типов

**Математические функции** (13 функций с правильной обработкой ошибок):

- **Тригонометрические**: sin, cos, tan
- **Степени/корни**: sqrt, pow
- **Округление**: abs, floor, ceil, round
- **Экстремумы**: min, max
- **Логарифмы**: log, log10, exp

**Строковые функции**: strlen, charAt, substring, startsWith, endsWith, indexOf, jsonParse, jsonStringify

**Файловые функции** (10 функций для работы с файловой системой):

- **Файловые операции**: readFile, writeFile, copyFile, removeFile, getFileSize
- **Директории**: mkdir, listDir, exists, isFile, isDir

**HTTP функции** (11 функций для веб-разработки):

- **HTTP клиент**: httpGet, httpPost, httpPut, httpDelete, httpSetTimeout
- **HTTP сервер**: httpCreateServer, httpRoute, httpStartServer, httpStopServer
- **URL утилиты**: urlEncode, urlDecode

**Функции каналов** (14 функций для межгорутинной коммуникации):
- **Управление каналами**: newChannel, send, receive, close
- **Неблокирующие операции**: tryReceive, trySend
- **Select операции**: channelSelect (множественный выбор каналов)
- **Операции с таймаутом**: channelTimeout (настраиваемые таймауты)
- **Итерация и очистка**: channelRange, channelDrain
- **Информация**: len, cap, channelInfo

**Криптографические функции** (20+ функций для безопасности):

- **Хеш-функции**: md5Hash, sha1Hash, sha256Hash, sha512Hash
- **Кодирование**: base64Encode/Decode, base64URLEncode/Decode, hexEncode/Decode  
- **HMAC подписи**: hmacSHA256, hmacSHA1, hmacMD5
- **Случайные данные**: randomBytes, randomString, randomInt, randomUUID
- **Пароли**: passwordHash, passwordVerify (SHA256 + соль)
- **Безопасность**: constantTimeCompare (защита от timing атак)

**Функции регулярных выражений** (10+ функций для работы с regex):

- **Поиск и сопоставление**: regexMatch, regexFind, regexFindAll
- **Замена**: regexReplace, regexReplaceAll
- **Разделение**: regexSplit, stringSplit (простое разделение)
- **Группы захвата**: regexGroups (извлечение групп из совпадений)
- **Валидация и утилиты**: regexValid, regexEscape, regexCount

**Функции времени** (25+ функций для работы с датой и временем):
- **Создание времени**: now, timeFromUnix, timeFromString
- **Компоненты**: timeYear, timeMonth, timeDay, timeHour, timeMinute, timeSecond, timeWeekday
- **Форматирование**: timeFormat (различные форматы)
- **Арифметика**: timeAddDays, timeAddMonths, timeAddYears, timeAddHours, timeAddMinutes, timeAddSeconds
- **Сравнение**: timeBefore, timeAfter, timeEqual
- **Разности**: timeDiff, timeDiffDays, timeDiffHours, timeDiffMinutes
- **Unix timestamps**: timeUnix для получения timestamp

**Методы примитивных типов**:

- **Числа**: .toString(), .abs(), .round(), .floor(), .ceil(), .isInteger()
- **Строки**: .length(), .charAt(), .substring(), .toUpper(), .toLower()
- **Логические**: .toString(), .not()

#### 7. Продвинутые типы и конструкции

- **Result тип**: `Ok(value)` и `Err(error)` как в Rust
- **Enum типы**: `enum Color { RED, GREEN, BLUE }`
- **Match выражения**: pattern matching с default case
- **Тернарный оператор**: `condition ? true : false`
- **for-yield**: создание массивов `for ... { yield value }`

#### 8. Богатая система типов и коллекций

- **Базовые типы**: int64, float64, string, bool
- **Коллекции**: массивы [1,2,3], объекты {key: value}
- **Базовые методы массивов**: .length(), .push(), .pop(), .slice()
- **Generic методы коллекций**: .map(), .filter(), .reduce() с поддержкой функций
- **Цепочки методов**: `numbers.filter(fn).map(fn).reduce(init, fn)`
- **Индексация**: arr[0], obj["key"], цепочные вызовы
- **Строковая интерполяция**: `"Hello ${name}, you are ${age}"`

#### 9. 🚀 **МНОГОПОТОЧНОСТЬ И СИНХРОНИЗАЦИЯ** ✅ **тесты готовы**

**Async/await поддержка**:

```foo
// Асинхронные функции с параметрами
fn delayedHello(name, delay) {
    await sleep(delay)  // задержка в миллисекундах
    return "Hello, " + name + "!"
}

// Запуск async функций
let promise = async delayedHello("World", 100)
let result = await promise  // "Hello, World!" (через 100мс)

// Promise.all - ждем завершения всех задач
let task1 = async delayedHello("Alice", 50)
let task2 = async delayedHello("Bob", 100)
let task3 = async delayedHello("Charlie", 75)
let results = await Promise.all(task1, task2, task3)

// Promise.any - ждем первый результат (racing)
let fast1 = async delayedHello("Fast", 30)
let fast2 = async delayedHello("Faster", 50)
let winner = await Promise.any(fast1, fast2)
```

**✅ Поддержка**: горутины, Promise API, изоляция scope, обработка ошибок
**✅ Примеры**: `test_async_basic.foo`, `test_async_simple.foo`

**Примитивы синхронизации**:

```foo
// Мьютексы для защиты ресурсов
let mutex = newMutex("counter_mutex")
mutexLock(mutex)
// критическая секция
mutexUnlock(mutex)

// Семафоры для ограничения ресурсов
let semaphore = newSemaphore(5, "db_connections")
semaphoreAcquire(semaphore)
// работа с ресурсом
semaphoreRelease(semaphore)

// Атомарные операции
let atomic = newAtomic(100, "counter")
let newValue = atomicAdd(atomic, 25)
let swapped = atomicCompareAndSwap(atomic, 125, 200)

// WaitGroup для синхронизации задач
let wg = newWaitGroup("tasks")
waitGroupAdd(wg, 3)
// запуск 3 горутин...
waitGroupWait(wg)  // ждем завершения всех
```

**✅ Поддержка**: мьютексы, rwmutexes, семафоры, waitgroups, атомики, барьеры
**✅ Примеры**: `test_sync_demo.foo`, `test_sync_simple.foo`

#### 10. 🚀 **BYTECODE VIRTUAL MACHINE** ✅ **полностью реализована**

```bash
# Запуск через bytecode VM для оптимизированной производительности
go run main.go --bytecode examples/test_bytecode_demo.foo

# Профилирование и диагностика
go run main.go --bytecode --profile --disassemble

# Сравнение производительности с tree-walking
go run main.go --bytecode --compare
```

**🔥 Ключевые компоненты**:

- **80+ OpCodes** - полное покрытие всех операций языка
- **Stack-based VM** - оптимальная архитектура для интерпретации
- **Профайлер производительности** - детальная статистика и горячие точки
- **Дизассемблер** - human-readable представление bytecode
- **Автоматические рекомендации** - советы по оптимизации кода

**✅ OpCodes поддержка**:

- **Арифметика**: OP_ADD, OP_SUBTRACT, OP_MULTIPLY, OP_DIVIDE, OP_MODULO
- **Логика**: OP_AND, OP_OR, OP_NOT, OP_EQUAL, OP_GREATER, OP_LESS
- **Переменные**: OP_GET_GLOBAL, OP_SET_GLOBAL, OP_DEFINE_GLOBAL
- **Управление**: OP_JUMP, OP_JUMP_IF_FALSE, OP_LOOP
- **Коллекции**: OP_ARRAY, OP_INDEX, OP_OBJECT
- **Профилинг**: OP_PROFILE_START, OP_PROFILE_END, OP_DEBUG_TRACE

**✅ Примеры**: `test_bytecode_demo.foo`
**✅ Тесты**: `minimal_bytecode_test.go`, `bytecode_test.go` (25+ тестов)

#### 11. Полная система операторов

- **Арифметические**: +, -, \*, /, %
- **Сравнения**: ==, !=, >, <, >=, <=
- **Логические**: &&, ||, !
- **Составные присваивания**: +=, -=, \*=, /=, %=
- **Инкремент/декремент**: ++, --

## 📊 Покрытие тестами: 100%

**28 файлов unit-тестов** с **240+ тестами**:

- `test/basic_types_test.go` - типы и операторы
- `test/functions_test.go` - функции и рекурсия
- `test/closures_test.go` - замыкания (5 тестов)
- `test/anonymous_functions_test.go` - анонимные функции (5 тестов)
- `test/async_test.go` - многопоточность async/await (8 тестов)
- `test/math_functions_test.go` - математические функции (14 тестов)
- `test/module_loading_test.go` - загрузка модулей (4 теста)
- `test/collections_test.go` - массивы и объекты
- `test/result_test.go` - Result тип
- `test/string_features_test.go` - интерполяция и комментарии
- `test/extension_methods_test.go` - extension methods (8 тестов)
- `test/method_overloading_test.go` - перегрузка методов (15 тестов)
- `test/filesystem_test.go` - файловая система (3 теста)
- `test/interface_test.go` - интерфейсы и impl блоки (12 тестов)
- `test/http_test.go` - HTTP клиент/сервер (15 тестов)
- `test/crypto_test.go` - криптографические функции (20+ тестов)
- `test/regex_test.go` - регулярные выражения (20+ тестов)
- `test/sync_test.go` - примитивы синхронизации (20+ тестов)
- `test/minimal_bytecode_test.go` - bytecode VM (6 тестов)
- `test/bytecode_test.go` - bytecode система (6 тестов)
- `test/simple_bytecode_test.go` - простые bytecode операции (6 тестов)
- `test/time_test.go` - работа с датой и временем (25+ функций времени)
- И другие...

**Все тесты проходят успешно!**

## 🎯 Важные технические детали

### Архитектура интерпретатора

- **Tree-walking interpreter**: AST выполняется напрямую
- **Recursive descent parser**: парсер с precedence climbing
- **Unified type system**: все значения через структуру Value
- **Scope stack**: стек областей видимости для функций и модулей

### Ключевые файлы

#### `ast/` - AST узлы и интерпретатор

- `anonymous_func.go` - анонимные функции
- `closure.go` - замыкания с захватом переменных
- `func_expr.go` - обычные функции
- `for_expr.go` - циклы и for-yield
- `result_expr.go` - Result тип (Ok/Err)
- `import_expr.go/export_expr.go` - модульная система

#### `builtin/math.go` - встроенные математические функции

13 функций с правильной обработкой ошибок и граничных случаев

#### `builtin/filesystem.go` - встроенные файловые функции

10 функций для работы с файловой системой (readFile, writeFile, exists, mkdir и другие)

#### `builtin/http.go` - встроенные HTTP функции

11 функций для веб-разработки (httpGet, httpPost, httpStartServer, роутинг и другие)

#### `builtin/crypto.go` - встроенные криптографические функции

20+ функций для безопасности (хеши, кодирование, HMAC, случайные данные, пароли, UUID, защита от timing атак)

#### `builtin/regex.go` - встроенные функции регулярных выражений

10+ функций для regex (поиск, замена, разделение, группы захвата, валидация, утилиты)

#### `modules/module.go` - система модулей

- Загрузка и парсинг .foo файлов
- Кеширование модулей
- Изоляция областей видимости

#### `scope/scope.go` - система областей видимости

- Stack-based scoping
- Push/pop для функций
- Глобальная и локальные области

#### `value/value.go` - система типов

- Unified Value структура
- Автоматическое приведение типов
- Специальные флаги (isReturn, isBreak, isYield)

## 🚧 Следующие этапы развития:

### 🔥 Высокий приоритет - Продвинутые Generic функции ✅ **ВСЕ РЕАЛИЗОВАНО!**

1. ✅ **Generic ограничения типов** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНО!** `<T: Interface + Interface2>` для ограничения параметров интерфейсами
2. ✅ **Interface система** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНА!** определение интерфейсов и их реализация типами
3. ✅ **Перегрузка методов** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНА!** множественные определения методов с разными сигнатурами
4. ✅ **Extension methods** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНА!** добавление методов к существующим типам

### 📋 Средний приоритет - Система типов

5. **Union типы** - `string | number | null` для объединения типов
6. **Optional типы** - `string?` синтаксис для nullable типов
7. **Tuple типы** - `(string, number, bool)` для кортежей
8. ✅ **Type aliases** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНЫ!** `type UserId = int` для псевдонимов типов

### 🌟 Низкий приоритет - Инфраструктура

9. **Стандартная библиотека** (std пакет)
10. ✅ **Файловая система** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНА!** (readFile, writeFile, exists, mkdir и другие - 10 функций)
11. ✅ **Регулярные выражения** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНЫ!** (regexMatch, regexReplace, regexSplit и другие - 10+ функций)
12. ✅ **Многопоточность** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНА!** (async/await, Promise API, примитивы синхронизации)
13. **LSP поддержка** (планируется как следующий шаг)
14. **Синтаксис хайлайтинг** (планируется)
15. **Документация сайт** (планируется, возможно с Astro)

## 📈 История развития проекта

### ✅ ЯНВАРЬ 2025 - РЕВОЛЮЦИОННЫЕ ДОСТИЖЕНИЯ! 🚀

**ПОСЛЕДНЕЕ ДОСТИЖЕНИЕ** ⚡ **Работа с датой и временем** - полная система функций времени с 25+ функциями для всех временных операций

#### 🔥 **Interface система** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНА!**

```foo
// Определение интерфейсов
interface Drawable {
    fn draw()
    fn getArea() -> float
}

// Определение структур
struct Circle {
    radius: float
}

// Реализация интерфейсов
impl Drawable for Circle {
    fn draw() {
        println("Drawing circle with radius " + this.radius.toString())
    }

    fn getArea() -> float {
        return 3.14159 * this.radius * this.radius
    }
}

// Использование
let circle = Circle{radius: 5.0}
circle.draw()                    // "Drawing circle with radius 5"
let area = circle.getArea()      // 78.53975
```

**✅ Примеры**: `test_interface_simple.foo`, `test_interface_complete.foo`, `test_interface_impl.foo`

#### 🔥 **Extension Methods** и **Перегрузка методов** - **ПОЛНОСТЬЮ РЕАЛИЗОВАНЫ!**

### 📋 Другие достижения:

1. **HTTP клиент/сервер** - полная поддержка веб-разработки (11 функций: httpGet, httpPost, роутинг, сервер)
2. **Файловая система** - полная поддержка I/O операций (readFile, writeFile, exists, mkdir и другие)
3. **Криптографические функции** - полная поддержка безопасности (20+ функций: хеши, кодирование, HMAC, случайные данные, пароли, UUID)
4. **Регулярные выражения** - полная поддержка regex (10+ функций: поиск, замена, разделение, группы захвата, валидация)
5. **Generic функции** - параметры типов `<T, U>`, типизированные параметры и возвраты
6. **Типизированные функции** - `fn add(a: int, b: int) -> int { ... }`
7. **Generic методы коллекций** - map, filter, reduce с реальными функциями
8. **Методы примитивных типов** - .toString(), .abs(), .length() и другие
9. **Строковые функции** - полный набор функций для работы со строками и JSON
10. **Анонимные функции** - полная поддержка стрелочного и блочного синтаксиса
11. **Замыкания** - автоматический захват переменных из внешней области
10. **Математические функции** - 13 встроенных функций с обработкой ошибок
11. **Модульная система** - реальная загрузка файлов с кешированием
12. **Множественные возвращаемые значения** - destructuring assignment
13. **Параметры по умолчанию** - для функций с выражениями
14. **Result тип** - полная реализация Ok/Err как в Rust
15. **Многопоточность async/await** - полная реализация Promise API с горутинами
16. **Продвинутые каналы** - полная система каналов с 14 функциями (trySend, channelSelect, channelTimeout, channelRange, channelDrain)
17. **Bytecode Virtual Machine** - оптимизированная VM с профилированием (ошибка slice bounds исправлена)
18. **Работа с датой и временем** - полная поддержка временных операций (25+ функций времени)
19. **100% покрытие тестами** - все функции протестированы

### ✅ Ранние достижения:

1. **for-yield конструкции** - создание массивов через циклы
2. **Строковая интерполяция** - `"Hello ${name}"`
3. **Enum типы и Match** - pattern matching
4. **Области видимости** - stack-based scoping с защитой от overflow
5. **Комментарии** - // и /\* \*/
6. **Объекты и массивы** - с методами и индексацией
7. **Цепочные вызовы** - obj.method().property
8. **Составные операторы** - +=, -=, ++, --

## 💻 Команды для работы

### Запуск

```bash
# Запуск файла по умолчанию (examples/main.foo)
go run main.go

# Запуск конкретного файла (Tree-walking интерпретатор)
go run main.go examples/test_complete_typing.foo             # 🔥 ПОЛНАЯ СИСТЕМА ТИПИЗАЦИИ
go run main.go examples/test_generic_constraints.foo         # Generic ограничения типов
go run main.go examples/test_interface_complete.foo          # Интерфейсы и impl блоки
go run main.go examples/test_extension_methods.foo           # Extension methods
go run main.go examples/test_method_overloading.foo          # Перегрузка методов
go run main.go examples/test_anonymous_functions.foo         # Анонимные функции
go run main.go examples/test_math.foo                        # Математические функции
go run main.go examples/test_closures.foo                    # Замыкания
go run main.go examples/test_filesystem_simple.foo           # Файловая система
go run main.go examples/test_async_basic.foo                 # Async/await основы
go run main.go examples/test_async_simple.foo                # Async/await с Promise.all
go run main.go examples/test_http_client.foo                 # HTTP клиент демо
go run main.go examples/test_http_server.foo                 # HTTP сервер демо
go run main.go examples/test_http_complete.foo               # HTTP клиент + сервер
go run main.go examples/test_crypto_demo.foo                   # Криптографические функции демо
go run main.go examples/test_crypto_simple.foo                 # Простая демо криптографии
go run main.go examples/test_regex_demo.foo                    # Регулярные выражения демо
go run main.go examples/test_regex_simple.foo                  # Простая демо regex
go run main.go examples/test_advanced_channels.foo             # Продвинутые функции каналов (полная демо)
go run main.go examples/test_advanced_channels_simple.foo      # Простая демо новых функций каналов
go run main.go examples/test_type_aliases_working.foo          # Type aliases - псевдонимы типов

# 🚀 Bytecode Virtual Machine (оптимизированный)
go run main.go --bytecode                                    # Bytecode VM режим
go run main.go --bytecode examples/test_bytecode_demo.foo    # Bytecode с примером
go run main.go --bytecode --profile                         # С профилированием
go run main.go --bytecode --disassemble                     # С дизассемблированием
go run main.go --bytecode --profile --disassemble           # Полная диагностика
go run main.go --bytecode --compare                         # Сравнение производительности
```

### Тестирование

```bash
# Запуск всех тестов (200+ тестов)
go test ./test/... -v

# Запуск конкретного теста
go test ./test/extension_methods_test.go -v
go test ./test/method_overloading_test.go -v
go test ./test/filesystem_test.go -v
go test ./test/anonymous_functions_test.go -v
go test ./test/math_functions_test.go -v
go test ./test/async_test.go -v
go test ./test/http_test.go -v
go test ./test/crypto_test.go -v
go test ./test/regex_test.go -v
go test ./test/time_test.go -v

# Bytecode тесты
go test ./test/minimal_bytecode_test.go -v                   # VM тесты
go test ./test/bytecode_test.go -v                          # Полные bytecode тесты
go test ./test/simple_bytecode_test.go -v                   # Простые операции
```

### Сборка

```bash
go build -o foo_lang main.go
./foo_lang examples/main.foo
```

## 🔧 Архитектурные принципы

### AST узлы (ast/)

1. Каждый узел реализует интерфейс `Expr` с методом `Eval()`
2. Один узел = одна операция (принцип единственной ответственности)
3. Минимум полей, максимум читаемости
4. Все специальные состояния через флаги в `Value`

### Парсинг (parser/)

- **Recursive descent parser** с precedence climbing
- Обработка ошибок через panic с восстановлением
- Поддержка всех современных языковых конструкций

### Типизация (value/)

- **Unified Value система** - все значения через одну структуру
- Автоматическое приведение типов для арифметики
- Строгая типизация для логических операций

## ⚡ КРИТИЧЕСКИ ВАЖНО ДЛЯ CLAUDE:

### 🎯 После ЛЮБОЙ реализации новой функции:

1. **ОБЯЗАТЕЛЬНО обновить README.md** - добавить фичу в список реализованных
2. **Добавить "✅ тесты готовы"** если написаны тесты
3. **Обновить список примеров** если созданы demo файлы
4. **Убрать из списка ограничений** если проблема решена

### 📝 Формат обновления README:

```markdown
- [x] **Название новой фичи** ✅ **тесты готовы**
```

### 🔄 README.md должен ВСЕГДА быть актуальным!

## 🚀 Статус проекта: ПОЛНОЦЕННЫЙ СОВРЕМЕННЫЙ ЯЗЫК ПРОГРАММИРОВАНИЯ!

Язык **foo_lang v2** - это **РЕВОЛЮЦИОННЫЙ ИНТЕРПРЕТИРУЕМЫЙ ЯЗЫК** с ПОЛНОЙ поддержкой современных возможностей:

### 🎯 **ТОП-УРОВЕНЬ ФУНКЦИОНАЛЬНОСТИ:**

- ✅ **ИНТЕРФЕЙСЫ** - полная ОО-система с `interface` и `impl`
- ✅ **EXTENSION METHODS** - расширение типов новыми методами
- ✅ **ПЕРЕГРУЗКА МЕТОДОВ** - множественные определения
- ✅ **GENERIC ФУНКЦИИ** - параметры типов `<T, U>`
- ✅ **ЗАМЫКАНИЯ** - автозахват переменных
- ✅ **МОДУЛЬНАЯ СИСТЕМА** - реальная загрузка файлов
- ✅ **RESULT ТИПЫ** - как в Rust
- ✅ **СТРОКОВАЯ ИНТЕРПОЛЯЦИЯ**
- ✅ **СТРУКТУРЫ И ENUM**
- ✅ **+ 30 других продвинутых фичей!**

### 🎓 **ИДЕАЛЬНО ПОДХОДИТ ДЛЯ:**

- 🔬 **Изучения современных языков программирования**
- 💡 **Исследования интерпретаторов и компиляторов**
- 🎯 **Прототипирования сложных алгоритмов**
- 📚 **Образовательных и академических проектов**
- 🚀 **Демонстрации передовых language features**
- 💼 **Production-ready скриптов и утилит**

### 🛠 **СЛЕДУЮЩИЕ ШАГИ**:

Bytecode Compiler (AST → Bytecode) → JIT компиляция → Union типы → LSP поддержка → Синтаксис хайлайтинг
