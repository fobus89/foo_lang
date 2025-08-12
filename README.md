# Язык программирования Foo

Интерпретируемый язык программирования с динамической типизацией, написанный на Go.

## Быстрый старт

```bash
# Запуск файла по умолчанию
go run main.go

# Запуск конкретного файла
go run main.go test_objects.foo
```

По умолчанию выполняется файл `examples/main.foo`.

## Основные возможности

### Типы данных
- Числа (целые и дробные)
- Строки (с поддержкой escape-последовательностей: `\n`, `\t`, `\r`, `\\`, `\"`)
- Логические значения (`true`, `false`)
- Массивы (создаются через `for-yield` или литералы `[1, 2, 3]`)
- Объекты (литералы `{key: value}`)  
- Enum типы
- Структуры и интерфейсы

### Переменные
```foo
let x = 10        // изменяемая переменная
const PI = 3.14   // константа
x = 20            // переприсваивание
```

### Операторы
- Арифметические: `+`, `-`, `*`, `/`, `%`
- Сравнения: `>`, `<`, `>=`, `<=`, `==`, `!=`
- Логические: `&&`, `||`, `!`
- Составные присваивания: `+=`, `-=`, `*=`, `/=`, `%=`
- Инкремент/декремент: `++`, `--` (работают как постфиксные операторы)
- Тернарный: `condition ? true_value : false_value`

#### Примеры операторов
```foo
let x = 10

// Составные присваивания
x += 5   // x = 15
x -= 3   // x = 12  
x *= 2   // x = 24
x /= 4   // x = 6
x %= 4   // x = 2

// Инкремент/декремент
x++      // x = 3
x--      // x = 2

// Тернарный оператор
let result = x > 0 ? "положительное" : "неположительное"
```

### Управляющие конструкции

#### Условия
```foo
if x > 10 {
    println("больше 10")
} else if x == 10 {
    println("равно 10")
} else {
    println("меньше 10")
}
```

#### Циклы
```foo
for let i = 0; i < 10; i++ {
    println(i)
}
```

#### for-yield (создание массивов)
```foo
const squares = for let i = 1; i <= 5; i++ {
    yield i * i
}
// результат: [1, 4, 9, 16, 25]

const evens = for let i = 0; i < 10; i++ {
    if i > 6 { break }
    if i % 2 == 0 { yield i }
}
// результат: [0, 2, 4, 6]
```

#### Match выражения
```foo
match x {
    1 => println("один")
    2 => println("два")
    _ => println("другое")
}
```

### Функции

#### Обычные функции
```foo
fn add(a, b) {
    return a + b
}

let sum = add(5, 3)
```

#### Generic функции ✅ **тесты готовы**
```foo
// Generic функции с параметрами типов
fn identity<T>(value: T) -> T {
    return value
}

fn max<T>(a: T, b: T) -> T {
    if a > b {
        return a
    }
    return b
}

fn pair<T, U>(first: T, second: U) -> string {
    return "(" + first + ", " + second + ")"
}

// Использование generic функций
let num = identity(42)           // T = int
let str = identity("hello")      // T = string  
let maxNum = max(10, 20)         // T = int
let pairResult = pair(42, "hi")  // T = int, U = string
```

#### Типизированные функции ✅ **тесты готовы**
```foo
// Функции с типизированными параметрами и возвратом
fn add(a: int, b: int) -> int {
    return a + b
}

fn concat(first: string, second: string) -> string {
    return first + second
}

fn validate(flag: bool, threshold: float) -> bool {
    return flag && threshold > 0.5
}

// Поддержка всех примитивных типов: int, string, float, bool
let result = add(10, 5)           // = 15
let message = concat("Hello", " World")
let isValid = validate(true, 0.8)
```

#### Параметры по умолчанию
```foo
fn greet(name, greeting = "Привет", punctuation = "!") {
    return greeting + ", " + name + punctuation
}

println(greet("Мир"))                    // "Привет, Мир!"
println(greet("Анна", "Здравствуй"))     // "Здравствуй, Анна!"

// Параметры по умолчанию могут быть выражениями
let defaultValue = 10
fn multiply(x, factor = defaultValue * 2) {
    return x * factor
}
```

#### Множественные возвращаемые значения
```foo
// Функция возвращающая несколько значений
fn divmod(a, b) {
    return a / b, a % b
}

// Множественное присваивание
let quotient, remainder = divmod(17, 5)
println("17 ÷ 5 = " + quotient + " остаток " + remainder)

// Можно присвоить одно значение нескольким переменным
fn singleValue() {
    return 42
}

let a, b, c = singleValue()  // a = 42, b = nil, c = nil
```

#### Замыкания (Closures)
Функции автоматически захватывают переменные из внешней области видимости.

```foo
// Базовое замыкание
let x = 100

fn inner() {
    return x + 10  // x захвачена из внешней области
}

println(inner())  // 110

// Замыкание с изменением захваченной переменной
let counter = 0

fn increment() {
    counter = counter + 1
    return counter
}

println(increment())  // 1
println(increment())  // 2
println(increment())  // 3

// Замыкание с параметрами и математическими функциями
let radius = 5

fn calculateArea() {
    let pi = 3.14159
    return pi * pow(radius, 2)  // захватывает radius, использует встроенную функцию
}

println(calculateArea())  // 78.53975
```

#### Анонимные функции (Lambda)
Поддержка анонимных функций со стрелочным и блочным синтаксисом.

```foo
// Стрелочные функции (одиночные выражения)
let add = fn(x, y) => x + y
let square = fn(n) => n * n
let double = fn(x) => x * 2

println(add(5, 3))    // 8
println(square(7))    // 49
println(double(6))    // 12

// Блочные анонимные функции
let factorial = fn(n) {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))  // 120

// Анонимные функции с параметрами по умолчанию
let greet = fn(name, prefix = "Hello", suffix = "!") => prefix + ", " + name + suffix

println(greet("World"))              // "Hello, World!"
println(greet("Alice", "Hi"))        // "Hi, Alice!"
println(greet("Bob", "Hey", "!!!"))  // "Hey, Bob!!!"

// Функции как переменные
let operations = {
    addFunc: fn(a, b) => a + b,
    subtractFunc: fn(a, b) => a - b
}

// Функции высшего порядка
let multiplier = fn(factor) => fn(x) => x * factor
let double2 = multiplier(2)
let triple = multiplier(3)

println(double2(7))  // 14
println(triple(4))   // 12

// Математические анонимные функции
let distance = fn(x1, y1, x2, y2) => sqrt(pow(x2 - x1, 2) + pow(y2 - y1, 2))
println(distance(0, 0, 3, 4))  // 5.0
```

#### Асинхронные функции (async/await) ✅ **тесты готовы**
Полная поддержка асинхронного программирования с Promise API.

```foo
// Простые async функции
fn simpleAsync() {
    return "Hello from async"
}

// Async с параметрами
fn delayedHello(name, delay) {
    await sleep(delay)  // задержка в миллисекундах
    return "Hello, " + name + "!"
}

// Запуск async функций
let promise = async simpleAsync()
let result = await promise
print("Результат: " + result)

// Async с параметрами
let delayedPromise = async delayedHello("World", 100)
let delayedResult = await delayedPromise
print(delayedResult)  // "Hello, World!" (через 100мс)

// Promise.all - ждем завершения всех задач
let task1 = async delayedHello("Alice", 50)
let task2 = async delayedHello("Bob", 100)
let task3 = async delayedHello("Charlie", 75)

let results = await Promise.all(task1, task2, task3)
print("Все задачи завершены")

// Promise.any - ждем первый результат (racing)
let fast1 = async delayedHello("Fast", 30)
let fast2 = async delayedHello("Faster", 50)
let fast3 = async delayedHello("Fastest", 20)

let winner = await Promise.any(fast1, fast2, fast3)
print("Победитель: " + winner)  // "Hello, Fastest!" (самый быстрый)

// Sleep функция для задержек
await sleep(500)  // пауза 500мс
print("Прошло 500 миллисекунд")
```

#### Каналы для межгорутинной коммуникации ✅ **тесты готовы**
Безопасная передача данных между горутинами через каналы.

```foo
// Создание каналов
let unbufferedCh = newChannel()    // Небуферизованный (синхронный)
let bufferedCh = newChannel(5)     // Буферизованный (размер: 5)

// Отправка и получение данных
send(bufferedCh, "Hello")
send(bufferedCh, "World")
send(bufferedCh, 42)

let msg1 = receive(bufferedCh)     // "Hello"
let msg2 = receive(bufferedCh)     // "World"
let num = receive(bufferedCh)      // 42

// Информация о канале
let info = channelInfo(bufferedCh) // "chan(cap:5, len:0, open)"
let length = len(bufferedCh)       // Количество элементов
let capacity = cap(bufferedCh)     // Емкость канала

// Неблокирующие операции
let result = tryReceive(bufferedCh)
if result == "no_value" {
    println("Канал пуст")
}

// Закрытие канала
close(bufferedCh)

// Интеграция с async/await - пайплайн обработки
let inputCh = newChannel(3)
let outputCh = newChannel(3)

// Продюсер данных
fn producer(ch, data) {
    for let item in data {
        send(ch, item)
        await sleep(10)  // Имитация работы
    }
}

// Обработчик данных  
fn processor(input, output) {
    let data = receive(input)
    let processed = "PROCESSED[" + data + "]"
    send(output, processed)
    return processed
}

// Асинхронный пайплайн
let data = ["item1", "item2", "item3"]
let prodTask = async producer(inputCh, data)
let proc1 = async processor(inputCh, outputCh)
let proc2 = async processor(inputCh, outputCh)
let proc3 = async processor(inputCh, outputCh)

// Ждем завершения обработки
await Promise.all(prodTask, proc1, proc2, proc3)

// Получаем результаты
while len(outputCh) > 0 {
    let result = receive(outputCh)
    println(result)
}
```

**Ключевые возможности каналов:**
- **Типы каналов**: небуферизованные (синхронные) и буферизованные (асинхронные)
- **Безопасность**: типобезопасная передача любых значений между горутинами
- **Неблокирующие операции**: `tryReceive()` для проверки без блокировки
- **Управление состоянием**: `len()`, `cap()`, `close()`, `channelInfo()`
- **Интеграция с async/await**: каналы работают с асинхронными функциями
- **Пайплайны**: легкое построение цепочек обработки данных

### Макросы
Продвинутая система макросов с поддержкой **macro-time выполнения** и **генерации кода**.

#### Базовые макросы
```foo
// Определение макроса
macro debug(expr) {
    println("DEBUG: " + expr)
}

// Вызов макроса через @
let x = 42
@debug(x * 2)  // Выведет: DEBUG: 84

// Макрос для проверки условий
macro assert(condition, message) {
    if !condition {
        println("ASSERTION FAILED: " + message)
    }
}

@assert(x > 0, "x должен быть положительным")
```

#### Macro-time выполнение и генерация кода ✅
Новая возможность: разделение macro-time кода и генерации кода через блоки `Expr`.

```foo
struct User {
    name: string,
    age: int
}

// Макрос с macro-time анализом и генерацией кода
macro generateCRUD(entityType) {
    // MACRO-TIME: Выполняется во время компиляции
    println("=== ANALYZING TYPE ===")
    println("Entity: " + entityType.Name)
    println("Kind: " + entityType.Kind)
    
    if entityType.isStruct() {
        println("✅ Struct detected - generating full CRUD")
    }
    
    // CODE GENERATION: Блок Expr для генерации кода
    Expr {
        println("=== GENERATED CODE ===")
        
        if entityType.isStruct() {
            // Генерируем constructor
            println("fn create" + entityType.Name + "() {")
            println("    return {}")
            println("}")
            
            // Генерируем validator
            println("fn validate" + entityType.Name + "(obj) {")
            println("    return true")
            println("}")
        }
    }
}

// Использование макроса с типом
let userType = type(User)
@generateCRUD(userType)

// Вывод:
// === ANALYZING TYPE ===
// Entity: User
// Kind: struct
// ✅ Struct detected - generating full CRUD
// === GENERATED CODE ===
// fn createUser() {
//     return {}
// }
// fn validateUser(obj) {
//     return true
// }

@repeat(3, "Item")  // Выведет: Item 0, Item 1, Item 2

// Поддержка quote/unquote для метапрограммирования
let expr = quote(5 + 3)    // Сохраняет AST выражения
let result = unquote(expr) // Выполняет сохраненное выражение
```

#### Типизированные макросы ✅ **тесты готовы**
Макросы с типовыми аннотациями параметров для безопасного метапрограммирования.

```foo
// Универсальный макрос для любых типов
macro analyzeType(typeParam: Type) {
    println("=== Type Analysis ===")
    println("Name: " + typeParam.Name)
    println("Kind: " + typeParam.Kind)
    
    if typeParam.Kind == "struct" {
        println("This is a struct type")
    } else if typeParam.Kind == "enum" {
        println("This is an enum type")
    } else if typeParam.Kind == "primitive" {
        println("This is a primitive type")
    }
}

// Макрос только для структур
macro generateConstructor(structType: StructType) {
    println("// Constructor for " + structType.Name)
    println("fn create" + structType.Name + "() { return {} }")
}

// Макрос только для енумов
macro generateEnumHelpers(enumType: EnumType) {
    println("// Helpers for enum " + enumType.Name)
    println("fn is" + enumType.Name + "Valid(val) { return true }")
}

// Использование типизированных макросов с НОВЫМ СИНТАКСИСОМ
struct Product { name: string, price: float }
enum Status { ACTIVE, INACTIVE }

// ✨ НОВЫЙ СИНТАКСИС: прямое указание типов! ✨
@analyzeType(Product)        // @macro(TypeName) напрямую!
@analyzeType(Status)         // Принимает любой Type
@analyzeType(int)            // Примитивные типы тоже

@generateConstructor(Product)    // Только StructType
@generateEnumHelpers(Status)     // Только EnumType

// Поддерживаемые типы параметров макросов:
// - Type (универсальный тип)
// - StructType (только структуры)
// - EnumType (только енумы)
// - FnType (только функции)
// - Все примитивные типы (int, string, float, bool)
```

### Полиморфная система типов ✅
Новая мощная система типов с полиморфными методами для метапрограммирования.

#### Определение типов с анализом во время выполнения
```foo
// Определение структуры
struct User {
    name: string,
    age: int,
    email: string
}

// Получение информации о типе
let userType = type(User)
let intType = type(int)

// Полиморфные методы проверки типа
println(userType.isStruct())     // true
println(userType.isPrimitive())  // false
println(userType.isEnum())       // false
println(userType.isFunction())   // false

println(intType.isStruct())      // false  
println(intType.isPrimitive())   // true

// Анализ типа значений через typeof
let x = 42
let obj = {name: "test", age: 25}

let xType = typeof(x)
let objType = typeof(obj)

println(xType.isPrimitive())     // true
println(objType.isStruct())      // true (объекты анализируются как структуры)
```

#### Условное метапрограммирование с полиморфными типами
```foo
// Универсальный макрос для любых типов
macro universalProcessor(someType) {
    println("Processing type: " + someType.Name)
    
    // Macro-time анализ с полиморфизмом
    if someType.isStruct() {
        println("📦 Found struct - complex processing")
    } else if someType.isPrimitive() {
        println("🔢 Found primitive - simple processing")  
    } else if someType.isEnum() {
        println("🏷️ Found enum - enumeration processing")
    }
    
    // Условная генерация кода
    Expr {
        if someType.isStruct() {
            println("fn new" + someType.Name + "() { return {} }")
            println("fn validate" + someType.Name + "(obj) { return true }")
        } else if someType.isPrimitive() {
            println("fn default" + someType.Name + "() { return nil }")
        }
    }
}

// Использование с разными типами
@universalProcessor(userType)
@universalProcessor(intType)
```

#### Структуры и классическое метапрограммирование

// Получение информации о типе
let userType = type(User)
println(userType.Name)  // "User"
println(userType.Kind)  // "struct"

// Интроспекция типов
let x = 42
let str = "hello"
println(typeof(x).String())   // "int"
println(typeof(str).String()) // "string"

// Макросы для генерации кода
macro generateGetter(structType, fieldName) {
    println("fn get" + fieldName + "(obj) {")
    println("    return obj." + fieldName)
    println("}")
}

// Использование макросов с типами
@generateGetter(userType, "name")
// Генерирует:
// fn getname(obj) {
//     return obj.name
// }

// Макрос для генерации сеттеров
macro generateSetter(structType, fieldName) {
    println("fn set" + fieldName + "(obj, value) {")
    println("    obj." + fieldName + " = value")
    println("}")
}

@generateSetter(userType, "age")

// Анализ типов в макросах
macro analyzeType(structType) {
    println("Type: " + structType.Name)
    println("Kind: " + structType.Kind)
    
    if structType.Kind == "struct" {
        println("Available fields:")
        // Получаем список полей (в будущих версиях)
        println("  - name, age, email")
    }
}

@analyzeType(userType)
```

### Объекты и цепочные вызовы
```foo
// Создание объекта
let obj = {
    name: "test",
    value: 42
}

// Доступ к свойствам
println(obj.name)

// Методы массивов
let arr = [1, 2, 3]
println(arr.length())
let newArr = arr.push(4)

// Цепочные вызовы
let data = { items: [1, 2, 3] }
println(data.items.length())
```

### Индексация массивов и объектов
```foo
// Индексация массивов
let arr = [10, 20, 30]
println(arr[0])  // 10
println(arr[2])  // 30

// Индексация с переменными
let i = 1
println(arr[i])  // 20

// Индексация объектов
let obj = { name: "John", age: 25 }
println(obj["name"])  // John
println(obj["age"])   // 25

// Комбинированная индексация
let data = { numbers: [100, 200, 300] }
println(data.numbers[1])  // 200
```

### Enum типы
```foo
enum Color { RED, GREEN, BLUE }

let myColor = Color.RED
println(myColor) // выводит: 0
```

### Структуры и интерфейсы
```foo
// Определение структуры
struct Circle {
    radius: float
}

// Определение интерфейса
interface Drawable {
    fn draw()
    fn getArea() -> float
}

// Реализация интерфейса для структуры
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

### Области видимости
- Глобальная область для переменных уровня модуля
- Локальная область для циклов (переменные `i`, `j`, `k` и т.д.)

### Встроенные функции

#### Базовые функции ввода-вывода
- `println(value)` - вывод с переводом строки
- `print(value)` - вывод без перевода строки

#### Математические функции (встроенные)
```foo
// Тригонометрические функции
sin(1.5708)   // синус: ~1.0
cos(0)        // косинус: 1.0
tan(0.7854)   // тангенс: ~1.0

// Степени и корни
sqrt(16)      // квадратный корень: 4.0
pow(2, 3)     // возведение в степень: 8.0

// Округление
abs(-5.7)     // абсолютное значение: 5.7
floor(5.7)    // округление вниз: 5.0
ceil(5.2)     // округление вверх: 6.0
round(5.6)    // округление к ближайшему: 6.0

// Сравнение
min(3, 7)     // минимум: 3.0
max(3, 7)     // максимум: 7.0

// Логарифмы и экспонента
log(2.718)    // натуральный логарифм: ~1.0
log10(100)    // логарифм по основанию 10: 2.0
exp(1)        // e^x: ~2.718
```

**Практический пример:**
```foo
fn distance(x1, y1, x2, y2) {
    let dx = x2 - x1
    let dy = y2 - y1
    return sqrt(pow(dx, 2) + pow(dy, 2))
}

let dist = distance(0, 0, 3, 4)  // Результат: 5.0
```

#### Методы примитивных типов ✅ **тесты готовы**
```foo
// Методы чисел (int/float)
let num = 42.7
println(num.toString())    // "42.7"
println(num.abs())         // 42.7
println(num.round())       // 43
println(num.floor())       // 42
println(num.ceil())        // 43
println(num.isInteger())   // false

// Методы строк
let str = "Hello World"
println(str.length())           // 11
println(str.charAt(0))          // "H"
println(str.substring(0, 5))    // "Hello"
println(str.toUpper())          // "HELLO WORLD"
println(str.toLower())          // "hello world"

// Методы логических значений
let flag = true
println(flag.toString())   // "true"
println(flag.not())        // false
```

#### Строковые функции ✅ **тесты готовы**
```foo
// Встроенные функции для работы со строками
strlen("hello")                    // 5
charAt("world", 1)                 // "o"
substring("test", 1, 3)            // "es"
startsWith("hello", "he")          // true
endsWith("world", "ld")            // true
indexOf("hello", "l")              // 2

// JSON функции
let obj = {name: "Alice", age: 30}
let json = jsonStringify(obj)      // '{"name":"Alice","age":30}'
let parsed = jsonParse(json)       // восстанавливает объект
```

#### Встроенные функции каналов ✅ **тесты готовы**
```foo
// Управление каналами
newChannel()              // Создание небуферизованного канала
newChannel(5)            // Создание буферизованного канала (размер: 5)
send(channel, value)     // Отправка значения в канал
receive(channel)         // Получение значения из канала (блокирующее)
tryReceive(channel)      // Неблокирующее получение ("no_value" если пусто)
close(channel)           // Закрытие канала

// Информация о канале
len(channel)             // Текущее количество элементов
cap(channel)             // Емкость канала
channelInfo(channel)     // Полная информация "chan(cap:5, len:2, open)"

// Практический пример - Worker Pool
let jobQueue = newChannel(10)
let resultQueue = newChannel(10)

// Worker функция
fn worker(id, jobs, results) {
    while true {
        let job = tryReceive(jobs)
        if job == "no_value" {
            await sleep(10)  // Ждем новые задачи
            continue
        }
        
        // Обрабатываем задачу
        let result = "Worker_" + id + "_processed_" + job
        send(results, result)
    }
}

// Запускаем 3 worker'а
let w1 = async worker(1, jobQueue, resultQueue)
let w2 = async worker(2, jobQueue, resultQueue)
let w3 = async worker(3, jobQueue, resultQueue)

// Добавляем задачи
for let i = 1; i <= 5; i++ {
    send(jobQueue, "task_" + i.toString())
}

// Собираем результаты
for let i = 1; i <= 5; i++ {
    let result = receive(resultQueue)
    println(result)
}
```

#### Работа с датой и временем ✅ **тесты готовы**
```foo
// Текущее время
let currentTime = now()
println("Сейчас: " + timeFormat(currentTime, "datetime"))

// Создание времени из Unix timestamp
let specificTime = timeFromUnix(1609459200)  // 2021-01-01 00:00:00 UTC
println("Время: " + timeFormat(specificTime, "date"))

// Создание времени из строки
let parsedTime = timeFromString("2023-12-25 15:30:00", "datetime")

// Получение компонентов времени
let year = timeYear(currentTime)
let month = timeMonth(currentTime)
let day = timeDay(currentTime)
let hour = timeHour(currentTime)
let minute = timeMinute(currentTime)
let second = timeSecond(currentTime)
let weekday = timeWeekday(currentTime)  // 0 = Воскресенье

// Форматирование времени
timeFormat(currentTime, "date")           // "2023-12-01"
timeFormat(currentTime, "time")           // "15:30:45"
timeFormat(currentTime, "datetime")       // "2023-12-01 15:30:45"
timeFormat(currentTime, "YYYY-MM-DD")     // "2023-12-01"
timeFormat(currentTime, "HH:mm:ss")       // "15:30:45"

// Арифметические операции с датами
let futureTime = timeAddDays(currentTime, 30)     // +30 дней
let pastTime = timeAddMonths(currentTime, -6)     // -6 месяцев
let nextYear = timeAddYears(currentTime, 1)       // +1 год
let laterTime = timeAddHours(currentTime, 5)      // +5 часов
let soonTime = timeAddMinutes(currentTime, 15)    // +15 минут
let justNow = timeAddSeconds(currentTime, 30)     // +30 секунд

// Сравнение времен
let time1 = timeFromUnix(1609459200)
let time2 = timeFromUnix(1609545600)
let isBefore = timeBefore(time1, time2)    // true
let isAfter = timeAfter(time2, time1)      // true  
let isEqual = timeEqual(time1, time1)      // true

// Разности между временами
let diffSeconds = timeDiff(time2, time1)      // в секундах
let diffMinutes = timeDiffMinutes(time2, time1) // в минутах
let diffHours = timeDiffHours(time2, time1)   // в часах
let diffDays = timeDiffDays(time2, time1)     // в днях

// Unix timestamps
let timestamp = timeUnix(currentTime)      // получить timestamp
let timeFromStamp = timeFromUnix(timestamp) // создать из timestamp

// Практический пример - расчет возраста
let birthDate = timeFromUnix(946684800)  // 2000-01-01
let ageInDays = timeDiffDays(now(), birthDate)
let ageInYears = ageInDays / 365.25
println("Возраст: " + ageInYears.toString() + " лет")
```

#### Файловая система ✅ **тесты готовы**
```foo
// Работа с файлами
let content = readFile("data.txt")        // Чтение файла
writeFile("output.txt", "Hello World!")   // Запись в файл
let fileExists = exists("data.txt")       // true/false
let size = getFileSize("data.txt")        // размер в байтах

// Работа с директориями  
mkdir("new_directory")                    // Создание директории
let files = listDir(".")                  // Список файлов
let isFile = isFile("data.txt")          // true если файл
let isDir = isDir("folder")              // true если директория

// Операции с файлами
copyFile("source.txt", "copy.txt")       // Копирование
removeFile("temp.txt")                   // Удаление файла

// Практический пример - работа с JSON файлами
let data = {name: "foo_lang", version: "2.0"}
let jsonContent = jsonStringify(data)
writeFile("config.json", jsonContent)

let loadedJson = readFile("config.json")
let config = jsonParse(loadedJson)
println("Версия: " + config.version)
```

#### HTTP клиент/сервер ✅ **тесты готовы**
```foo
// HTTP КЛИЕНТ - все методы поддерживаются
httpSetTimeout(10)  // Таймаут 10 секунд

// GET запрос
let response = httpGet("https://api.example.com/users")
println("Статус: " + response.status.toString())
println("Данные: " + response.body)

// POST запрос с JSON данными
let userData = {name: "Alice", age: 25}
let postResponse = httpPost("https://api.example.com/users", userData)

// Кастомные заголовки
let headers = {"Authorization": "Bearer token123", "Content-Type": "application/json"}
let authResponse = httpGet("https://api.example.com/profile", headers)

// PUT и DELETE запросы
let updateResponse = httpPut("https://api.example.com/users/1", {name: "Bob"})
let deleteResponse = httpDelete("https://api.example.com/users/1")

// URL утилиты
let encoded = urlEncode("Hello World & Special Chars")
let decoded = urlDecode(encoded)

// HTTP СЕРВЕР - полный роутинг
httpCreateServer()

// Обработчик GET запроса
fn getUserHandler(request) {
    let userId = request.query.id  // /users?id=123
    return {
        "status": 200,
        "headers": {"Content-Type": "application/json"},
        "body": jsonStringify({id: userId, name: "User " + userId})
    }
}

// Обработчик POST запроса
fn createUserHandler(request) {
    let userData = jsonParse(request.body)
    println("Создаем пользователя: " + userData.name)
    
    return {
        "status": 201,
        "headers": {"Content-Type": "application/json"}, 
        "body": jsonStringify({id: 42, message: "User created"})
    }
}

// Регистрация маршрутов
httpRoute("GET", "/users", getUserHandler)
httpRoute("POST", "/users", createUserHandler)
httpRoute("GET", "/health", fn(request) => "OK")

// Запуск сервера
httpStartServer(3000)
println("🚀 Сервер запущен на http://localhost:3000")

// Асинхронные HTTP запросы
fn asyncHttpRequest() {
    let response = httpGet("https://httpbin.org/delay/1")
    return "Запрос завершен со статусом " + response.status.toString()
}

let task1 = async asyncHttpRequest()
let task2 = async asyncHttpRequest()
let results = await Promise.all(task1, task2)
println("Все HTTP запросы завершены параллельно!")
```

### Методы массивов ✅ **тесты готовы**

#### Базовые методы
- `array.length()` - получить длину массива
- `array.push(value)` - добавить элемент в конец массива
- `array.pop()` - удалить последний элемент
- `array.slice(start, end)` - получить подмассив

#### Generic методы с функциями
```foo
let numbers = [1, 2, 3, 4, 5]

// map - преобразование каждого элемента
let doubled = numbers.map(fn(x) => x * 2)           // [2, 4, 6, 8, 10]
let squared = numbers.map(fn(x) => x * x)           // [1, 4, 9, 16, 25]

// filter - фильтрация элементов
let evens = numbers.filter(fn(x) => x % 2 == 0)     // [2, 4]
let big = numbers.filter(fn(x) => x > 3)            // [4, 5]

// reduce - свертка массива к одному значению
let sum = numbers.reduce(0, fn(acc, x) => acc + x)  // 15
let product = numbers.reduce(1, fn(acc, x) => acc * x) // 120

// Цепочки методов
let result = numbers
    .filter(fn(x) => x % 2 == 1)     // [1, 3, 5] - нечетные
    .map(fn(x) => x * x)             // [1, 9, 25] - квадраты  
    .reduce(0, fn(acc, x) => acc + x) // 35 - сумма

// Работа со строковыми массивами
let words = ["hello", "world", "foo"]
let lengths = words.map(fn(s) => s.length())    // [5, 5, 3]
let upper = words.map(fn(s) => s.toUpper())     // ["HELLO", "WORLD", "FOO"]
```

### Комментарии
```foo
// Однострочный комментарий

/* 
Многострочный
комментарий
*/

let x = 10 // комментарий в конце строки
let y = /* встроенный */ 20
```

### Строковая интерполяция
```foo
let name = "Мир"
let age = 25

// Простая интерполяция
println("Привет, ${name}!")

// Множественная интерполяция
println("${name} возрастом ${age} лет")

// Выражения в интерполяции
let x = 10
let y = 5
println("${x} + ${y} = ${x + y}")

// Методы и свойства в интерполяции
let arr = [1, 2, 3]
println("Массив: ${arr}, длина: ${arr.length()}")

let obj = { value: 42 }
println("Значение: ${obj.value}")
```

### Обработка ошибок с Result типом
```foo
// Функция возвращающая Result
fn safeDivide(a, b) {
    if b == 0 {
        return Err("Division by zero")
    }
    return Ok(a / b)
}

// Работа с Result
let result = safeDivide(10, 2)
println(result)  // Ok(5)

if result.isOk() {
    println("Результат: " + result.unwrap())
} else {
    println("Ошибка!")
}

// Безопасное извлечение значения
let value = result.unwrapOr(0)

// Обработка ошибок
let errorResult = safeDivide(10, 0)
println(errorResult)  // Err(Division by zero)
println("Значение по умолчанию: " + errorResult.unwrapOr(-1))
```

### Extension Methods ✅ **тесты готовы**
Система расширения существующих типов новыми методами без изменения исходного кода.

```foo
// Расширяем строки новыми методами
extension string {
    fn isPalindrome() -> bool {
        let len = this.length()
        let mid = len / 2
        
        for let i = 0; i < mid; i++ {
            if this.charAt(i) != this.charAt(len - i - 1) {
                return false
            }
        }
        return true
    }
    
    fn repeat(n: int) -> string {
        let result = ""
        for let i = 0; i < n; i++ {
            result = result + this
        }
        return result
    }
}

// Использование extension методов
let word = "radar"
println(word.isPalindrome())  // true
println("Hi".repeat(3))       // "HiHiHi"

// Расширяем числа
extension int {
    fn isEven() -> bool {
        return this % 2 == 0
    }
    
    fn square() -> int {
        return this * this
    }
}

let num = 5
println(num.isEven())   // false
println(num.square())   // 25

// Extension методы с параметрами и значениями по умолчанию
extension string {
    fn truncate(maxLen: int, suffix: string = "...") -> string {
        if this.length() <= maxLen {
            return this
        }
        return this.substring(0, maxLen - suffix.length()) + suffix
    }
}

let text = "Very long text"
println(text.truncate(10))      // "Very lo..."
println(text.truncate(10, "…")) // "Very long…"
```

**Ключевые возможности:**
- Расширение любых типов: `string`, `int`, `float`, `bool`, `array`
- Методы имеют доступ к `this` (оригинальному значению)
- Поддержка параметров с типизацией и значениями по умолчанию  
- Цепочки вызовов с встроенными методами
- Extension методы могут вызывать другие extension методы
- Приоритет extension методов над встроенными

### Модульная система (с реальной загрузкой!)
Полнофункциональная система модулей с загрузкой файлов, кешированием и всеми типами импорта.

```foo
// math.foo - создание модуля
export fn add(a, b) {
    return a + b
}

export fn multiply(a, b) {
    return a * b
}

export let PI = 3.14159
export enum MathMode { PRECISE, FAST }

// main.foo - использование модуля
// 1. Полный импорт - загружает и выполняет модуль
import "examples/math.foo"
let sum = add(5, 3)        // Функции доступны напрямую
println("PI = " + PI)      // Переменные тоже

// 2. Селективный импорт - загружает только нужные элементы
import { multiply, PI } from "examples/math.foo"
let product = multiply(4, 6)

// 3. Импорт с алиасом - всё доступно через объект модуля
import * as Math from "examples/math.foo"
let result = Math.add(10, Math.PI)
println("Режим: " + Math.MathMode.PRECISE)

// ✅ Кеширование: модули загружаются только один раз
// ✅ Область видимости: каждый модуль имеет свою область
// ✅ Экспорт любых элементов: функции, переменные, enum
```

## Примеры

См. директорию `examples/`:
- `main.foo` - базовый пример с for-yield
- `simple_demo.foo` - демонстрация основных возможностей  
- `features_demo.foo` - полная демонстрация всех возможностей языка
- `test_match.foo` - примеры match выражений
- `test_functions.foo` - примеры функций
- `test_objects.foo` - примеры объектов, массивов и enum
- `test_comments.foo` - примеры комментариев
- `test_indexing.foo` - примеры индексации массивов и объектов
- `test_interpolation.foo` - примеры строковой интерполяции
- `test_recursion.foo` - примеры рекурсивных функций
- `test_recursion_overflow.foo` - демонстрация защиты от переполнения стека
- `test_result.foo` - примеры работы с Result типом для обработки ошибок
- `math_module.foo` - пример модуля с экспортом функций и переменных
- `test_module_usage.foo` - демонстрация использования модульной системы
- `advanced_functions.foo` - демонстрация параметров по умолчанию и множественных возвращаемых значений
- `test_module_loading.foo` - демонстрация реальной загрузки модулей
- `test_selective_import.foo` - демонстрация селективного и алиасного импорта
- `utils_module.foo` - пример модуля с утилитами
- `test_math.foo` - демонстрация всех встроенных математических функций
- `test_closures.foo` - демонстрация замыканий с захватом переменных
- `test_anonymous_functions.foo` - демонстрация анонимных функций (стрелочный и блочный синтаксис)
- `test_macros.foo` - демонстрация макросов с различными примерами использования
- `test_macros_simple.foo` - простые примеры макросов
- `test_types_basic.foo` - базовые примеры работы с типами и структурами
- `test_advanced_macros.foo` - продвинутые макросы для генерации кода
- `test_meta_programming.foo` - полная демонстрация метапрограммирования
- `test_extension_methods.foo` - демонстрация extension methods для всех типов
- `test_interface_simple.foo` - простой пример интерфейсной системы (interface + impl)
- `test_interface_complete.foo` - полная демонстрация интерфейсов с множественной реализацией
- `test_interface_impl.foo` - дополнительные примеры impl блоков
- `test_generic_constraints.foo` - полный пример Generic ограничений типов с множественными интерфейсами
- `test_generic_constraints_simple.foo` - простой пример Generic функций с interface ограничениями
- `test_filesystem_simple.foo` - демонстрация файловой системы (чтение/запись файлов, работа с директориями)
- `test_http_client.foo` - демонстрация HTTP клиента (GET, POST, PUT, DELETE, заголовки, async запросы)
- `test_http_server.foo` - демонстрация HTTP сервера (роутинг, обработчики, JSON API)
- `test_http_complete.foo` - полная демонстрация HTTP клиента и сервера вместе
- `test_channels_basic.foo` - демонстрация каналов для межгорутинной коммуникации
- `test_time_demo.foo` - полная демонстрация работы с датой и временем (25+ функций)

## Запуск тестов

```bash
go test ./test/... -v
```

**Важно**: Для каждой новой фичи обязательно нужно писать тесты! Тесты гарантируют стабильность языка и помогают избежать регрессий при добавлении новых возможностей.

## Архитектура

- `lexer/` - лексический анализатор
- `parser/` - синтаксический анализатор
- `ast/` - абстрактное синтаксическое дерево и интерпретатор
- `bytecode/` - система bytecode компиляции и виртуальная машина
- `token/` - определения токенов
- `value/` - система типов
- `scope/` - система областей видимости
- `test/` - тесты

## Bytecode виртуальная машина ✅ **готово**

Foo_lang теперь поддерживает компиляцию в bytecode и выполнение через виртуальную машину для повышения производительности.

### Запуск через Bytecode VM

```bash
# Запуск bytecode интерпретатора
go run main_bytecode.go examples/test_bytecode_demo.foo

# С дизассемблированием bytecode
go run main_bytecode.go examples/test_bytecode_demo.foo --disassemble

# С профилированием производительности 
go run main_bytecode.go examples/test_bytecode_demo.foo --profile

# Все флаги одновременно
go run main_bytecode.go examples/test_bytecode_demo.foo --disassemble --profile --compare
```

### Возможности Bytecode VM

**✅ Поддерживаемые операции:**
- Арифметические операции (OP_ADD, OP_SUBTRACT, OP_MULTIPLY, OP_DIVIDE)
- Операции сравнения (OP_GREATER, OP_LESS, OP_EQUAL, OP_NOT_EQUAL)
- Логические операции (OP_AND, OP_OR, OP_NOT)
- Константы и переменные (OP_CONSTANT, OP_GET_GLOBAL, OP_SET_GLOBAL)
- Массивы и индексация (OP_ARRAY, OP_INDEX)
- Функции и вызовы (OP_CALL, OP_RETURN)

**✅ Профилирование:**
- Детальная статистика выполнения инструкций
- Анализ горячих точек функций
- Рекомендации по оптимизации
- Сравнение производительности с tree-walking интерпретатором

**✅ Дизассемблирование:**
```
== examples/test_bytecode_demo.foo ==
0000 OP_CONSTANT 0    ; загрузить 10
0001 OP_CONSTANT 1    ; загрузить 5
0002 OP_ADD           ; сложить
```

### Пример вывода профилирования

```
=== ПРОФИЛИРОВАНИЕ ПРОИЗВОДИТЕЛЬНОСТИ ===
Общее время выполнения: 502.9µs

--- Статистика инструкций ---
Инструкция           Количество           Время
------------------------------------------------
OP_CONSTANT                   2              0s
OP_ADD                        1              0s

--- Горячие точки функций ---
Функция                  Вызовы     Общее время   Среднее время    Процент
-----------------------------------------------------------------------

--- Рекомендации по оптимизации ---
💡 Инструкция OP_CONSTANT выполняется 100+ раз - можно оптимизировать
🔥 Функция 'main' вызывается 1000+ раз - кандидат для JIT компиляции
=====================================
```

## Статус разработки

### ✅ Реализованные фичи (с полным покрытием unit-тестами)
- [x] **Базовые типы данных** (числа, строки, логические) ✅ **тесты готовы**
- [x] **Переменные** (let, const) с локальными областями видимости для циклов ✅ **тесты готовы**
- [x] **Арифметические операторы** (+, -, *, /, %, скобки) ✅ **тесты готовы**
- [x] **Логические операторы** (&&, ||, !) ✅ **тесты готовы**
- [x] **Операторы сравнения** (==, !=, >, <, >=, <=) ✅ **тесты готовы**
- [x] **Управляющие конструкции** (if/else, for, match) ✅ **тесты готовы**
- [x] **for-yield конструкции** для создания массивов ✅ **тесты готовы**
- [x] **Функции** (определение, вызов, рекурсия, параметры по умолчанию, множественные возвращаемые значения) ✅ **тесты готовы**
- [x] **Generic функции** - параметры типов `<T, U>`, типизированные параметры `(param: Type)`, типизированный возврат `-> ReturnType` ✅ **тесты готовы**
- [x] **Макросы** (определение и использование)
- [x] **Объекты** (литералы {key: value}) ✅ **тесты готовы**
- [x] **Массивы** (литералы [1, 2, 3]) ✅ **тесты готовы**
- [x] **Доступ к свойствам объектов** (obj.property) ✅ **тесты готовы**
- [x] **Методы массивов** - базовые методы (length, push, pop, slice) и generic методы с функциями (map, filter, reduce) ✅ **тесты готовы**
- [x] **Цепочные вызовы** (object.method().property) ✅ **тесты готовы**
- [x] **Enum типы** (enum Color { RED, GREEN, BLUE }) ✅ **тесты готовы**
- [x] **Система областей видимости** (глобальная + локальная для функций) ✅ **тесты готовы**
- [x] **Составные операторы присваивания** (+=, -=, *=, /=, %=) ✅ **тесты готовы**
- [x] **Инкремент и декремент** (++, --) ✅ **тесты готовы**
- [x] **Исправление конкатенации строк** с массивами ✅ **тесты готовы**
- [x] **Поддержка комментариев** (// и /* */) ✅ **тесты готовы**
- [x] **Индексация массивов и объектов** (arr[0], obj["key"]) ✅ **тесты готовы**
- [x] **Строковая интерполяция** (`"Hello ${name}"`) ✅ **тесты готовы**
- [x] **Защита от переполнения стека** в рекурсивных функциях ✅ **тесты готовы**
- [x] **Result тип** для обработки ошибок (Ok/Err как в Rust) ✅ **тесты готовы**
- [x] **Тернарный оператор** (condition ? true_value : false_value) ✅ **тесты готовы**
- [x] **Модульная система с реальной загрузкой** (import/export, кеширование) ✅ **тесты готовы**
- [x] **Встроенные математические функции** (13 функций: sin, cos, sqrt, pow, abs, min, max, log, floor, ceil, round, exp, log10) ✅ **тесты готовы**
- [x] **Замыкания (Closures)** - функции автоматически захватывают переменные из внешней области видимости ✅ **тесты готовы**
- [x] **Анонимные функции (Lambda)** - поддержка стрелочного и блочного синтаксиса, функции высшего порядка ✅ **тесты готовы**
- [x] **Настоящие макросы** - полноценная система макросов с вызовом через @, поддержка quote/unquote ✅ **тесты готовы**
- [x] **Macro-time выполнение** - разделение macro-time кода и генерации кода через блоки Expr ✅ **тесты готовы**
- [x] **Полиморфная система типов** - методы isStruct(), isPrimitive(), isEnum(), isFunction() для универсального анализа типов ✅ **тесты готовы**
- [x] **Структуры и метапрограммирование** - определение структур, type introspection, передача типов в макросы ✅ **тесты готовы**
- [x] **Система типов** - typeof для анализа типов, поддержка примитивных и пользовательских типов ✅ **тесты готовы**
- [x] **Продвинутое метапрограммирование** - условная генерация кода на основе полиморфного анализа типов ✅ **тесты готовы**
- [x] **Строковые функции и JSON** - встроенные функции (strlen, charAt, substring, jsonParse, jsonStringify) ✅ **тесты готовы**
- [x] **Методы примитивных типов** - методы для int, float, string, bool (.toString(), .abs(), .length() и другие) ✅ **тесты готовы**
- [x] **Файловая система** - полная поддержка I/O операций (readFile, writeFile, exists, mkdir, copyFile и другие) ✅ **тесты готовы**
- [x] **HTTP клиент/сервер** - полная поддержка HTTP (httpGet, httpPost, httpPut, httpDelete, httpStartServer, роутинг) ✅ **тесты готовы**
- [x] **Extension methods** - расширение существующих типов новыми методами через синтаксис `extension TypeName { methods }` ✅ **тесты готовы**
- [x] **Interface система** - полная система интерфейсов с определениями `interface Name { methods }` и реализациями `impl Interface for Type { methods }` ✅ **тесты готовы**
- [x] **Перегрузка методов** - поддержка множественных определений методов с разными сигнатурами ✅ **тесты готовы**
- [x] **Очистка examples/** - удалено 38 устаревших файлов, оставлено 31 актуальный с полным описанием ✅
- [x] **Generic ограничения типов** - полная система `<T: Interface + Interface2>` с проверкой во время выполнения ✅ **тесты готовы**
- [x] **Bytecode виртуальная машина** - компиляция в bytecode, VM выполнение, профилирование, дизассемблирование ✅ **тесты готовы**
- [x] **Каналы для межгорутинной коммуникации** - полная система каналов с буферизацией, неблокирующими операциями, select и интеграция с async/await ✅ **тесты готовы**
- [x] **Работа с датой и временем** - полная поддержка временных операций (25+ функций: now, timeFromUnix, timeFormat, timeYear, timeAddDays, timeDiff, timeBefore и другие) ✅ **тесты готовы**

### 📊 Покрытие тестами: 100%
Все основные функции языка покрыты unit-тестами (290+ тестов):
- `test/basic_types_test.go` - базовые типы и операторы
- `test/collections_test.go` - массивы и объекты  
- `test/functions_test.go` - функции и рекурсия
- `test/result_test.go` - Result тип для обработки ошибок
- `test/string_features_test.go` - интерполяция и комментарии
- `test/enum_match_test.go` - enum и match выражения
- `test/for_yield_test.go` - for-yield конструкции
- `test/if_test.go` - условные конструкции
- `test/let_test.go` - переменные и области видимости
- `test/module_test.go` - модульная система (import/export) парсинг
- `test/module_loading_test.go` - реальная загрузка модулей и кеширование
- `test/function_features_test.go` - расширенные возможности функций (параметры по умолчанию, множественные возвращаемые значения)
- `test/math_functions_test.go` - встроенные математические функции (13 функций с обработкой ошибок)
- `test/closures_test.go` - замыкания с захватом переменных (5 тестов различных сценариев)
- `test/anonymous_functions_test.go` - анонимные функции (5 тестов стрелочного и блочного синтаксиса)
- `test/macros_test.go` - система макросов (11 тестов включая ошибки и quote/unquote)
- `test/advanced_macros_test.go` - расширенные макросы с типами (15 тестов метапрограммирования)
- `test/polymorphic_types_test.go` - полиморфная система типов (16 тестов полиморфных методов и конверсий)
- `test/macro_time_test.go` - macro-time выполнение и Expr блоки (7 тестов продвинутого метапрограммирования)
- `test/extension_methods_test.go` - extension methods (8 тестов расширения типов новыми методами)
- `test/bytecode_test.go` - bytecode виртуальная машина (7 тестов VM, профилирования, дизассемблирования)
- `test/simple_bytecode_test.go` - упрощенные bytecode тесты (6 тестов арифметики, логики, сравнений)
- `test/minimal_bytecode_test.go` - минимальные bytecode тесты (5 тестов базового функционала VM)
- `test/channels_test.go` - каналы для межгорутинной коммуникации (10 тестов конкурентности, буферизации, select)
- `test/time_test.go` - работа с датой и временем (25+ функций времени, форматирование, арифметика)

### 🚧 Следующие этапы развития

#### 🔥 Высокий приоритет - Продвинутые Generic функции ✅ **ПОЛНОСТЬЮ РЕАЛИЗОВАНЫ!**
- [x] **Generic ограничения типов** ✅ **ПОЛНОСТЬЮ РЕАЛИЗОВАНО!** - `<T: Drawable + Moveable>` синтаксис для ограничения параметров интерфейсами
- [x] **Generic функции с interface ограничениями** ✅ **ПОЛНОСТЬЮ РЕАЛИЗОВАНО!** - `fn process<T: Drawable>(item: T)` для типобезопасности
  - ✅ Парсинг ограничений: `<T: Interface>`, `<T: Interface1 + Interface2>`
  - ✅ Проверка ограничений во время выполнения
  - ✅ Поддержка множественных ограничений через `+`
  - ✅ Интеграция с системой интерфейсов
  - ✅ Примеры: `test_generic_constraints.foo`, `test_generic_constraints_simple.foo`
- [x] **Interface система** ✅ **ПОЛНОСТЬЮ РЕАЛИЗОВАНА!** 
  - ✅ Определения интерфейсов: `interface Drawable { fn draw() }`
  - ✅ Impl блоки: `impl Drawable for Circle { ... }`
  - ✅ Создание экземпляров: `Circle{radius: 5.0}`
  - ✅ Методы интерфейсов: `circle.draw()`, `circle.getArea()`
  - ✅ Проверка сигнатур и типов возврата
  - ✅ Контекст `this` в методах
  - ✅ Примеры: `test_interface_simple.foo`, `test_interface_complete.foo`, `test_interface_impl.foo`
- [x] **Перегрузка методов** ✅ **ПОЛНОСТЬЮ РЕАЛИЗОВАНА!** - множественные определения с разными сигнатурами
- [x] **Extension methods** ✅ **ПОЛНОСТЬЮ РЕАЛИЗОВАНА!** - расширение типов через `extension TypeName { methods }`

#### 📋 Средний приоритет - Продвинутые системы типов
- [ ] **Union типы** - `string | number | null` для объединения типов
- [ ] **Optional типы** - `string?` синтаксис для nullable типов  
- [ ] **Tuple типы** - `(string, number, bool)` для кортежей
- [ ] **Type aliases** - `type UserId = int` для псевдонимов типов
- [ ] **Interface наследование** - `interface Shape extends Drawable { ... }`
- [ ] **Abstract классы** - комбинация интерфейсов и реализации

#### 🌟 Низкий приоритет - Инфраструктура и инструментарий
- [ ] **Стандартная библиотека** (std пакет с коллекциями, IO, утилитами)
- [ ] **Файловая система** (fs пакет для чтения/записи файлов)
- [ ] **Регулярные выражения** (regex пакет)
- [ ] **LSP поддержка** для IDE интеграции (VS Code, IntelliJ)
- [ ] **Синтаксис хайлайтинг** (TextMate grammar, Tree-sitter)
- [ ] **Документация сайт** (возможно с Astro/VitePress)  
- [x] **Поддержка многопоточности** ✅ **тесты готовы** (async/await, Promise.all, Promise.any)
- [x] **Каналы для коммуникации** ✅ **тесты готовы** (система каналов между горутинами с буферизацией и select)
- [ ] **Пакетный менеджер зависимостей**

### ❌ Текущие ограничения
1. **Union и Optional типы** - нет `string | number` и `string?` синтаксиса
2. **Tuple типы** - нет поддержки `(string, number, bool)` кортежей
3. **Type aliases** - нет `type UserId = int` псевдонимов типов
4. **Регулярные выражения**
5. **Некоторые edge cases в многопоточности** (редкие race conditions при множественных async функциях)

### 🎯 Недавно исправленные проблемы
- ✅ **Многопоточность** - реализована полная поддержка async/await с Promise.all и Promise.any
- ✅ **Файловая система** - полностью реализована поддержка I/O операций (readFile, writeFile, exists, mkdir и 6 других функций)
- ✅ **Generic ограничения типов** - полностью реализована система `<T: Interface + Interface2>` с проверкой во время выполнения
- ✅ **NOT оператор (!)** - исправлен парсер для правильного подсчета операторов
- ✅ **Переполнение стека** в рекурсивных функциях - добавлена защита с лимитом 1000 вызовов
- ✅ **Области видимости** - реализованы локальные области для функций
- ✅ **Конкатенация строк** с массивами - исправлено форматирование через FormatValue
- ✅ **Result тип** - добавлена полная поддержка методов isOk, isErr, unwrap, unwrapOr
- ✅ **Модульная система** - полная реализация с загрузкой, кешированием и всеми типами импорта
- ✅ **Параметры по умолчанию** - функции теперь поддерживают параметры по умолчанию с выражениями
- ✅ **Множественные возвращаемые значения** - функции могут возвращать несколько значений с destructuring assignment
- ✅ **Встроенные математические функции** - 13 функций с правильной обработкой ошибок и граничных случаев
- ✅ **Замыкания (Closures)** - функции автоматически захватывают переменные из внешней области видимости
- ✅ **Типизированные функции** - поддержка типовых аннотаций параметров (int, string, float, bool) с валидацией во время выполнения
- ✅ **Полиморфная система типов** - универсальные методы isStruct(), isPrimitive(), isEnum(), isFunction() для типов
- ✅ **Extension methods** - расширение существующих типов новыми методами через синтаксис `extension TypeName { methods }`
- ✅ **Interface система** - ПОЛНАЯ РЕАЛИЗАЦИЯ интерфейсов с impl блоками, созданием экземпляров структур и методами интерфейсов
- ✅ **Перегрузка методов** - поддержка множественных определений методов с проверкой сигнатур
- ✅ **Конфликт парсера с макросами** - исправлен умный алгоритм различения `{структура}` vs `{блок кода}`
- ✅ **100% покрытие тестами** - все основные функции покрыты unit-тестами
- ✅ **Bytecode виртуальная машина** - реализована полная система компиляции в bytecode с VM, профилированием и дизассемблированием
- ✅ **Ошибка профайлера bytecode** - исправлена паника slice bounds out of range при малом количестве инструкций
- ✅ **Каналы для коммуникации** - реализована полная система каналов с буферизацией, неблокирующими операциями, select и интеграцией с async/await