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

#### Типизированные функции ✅ **тесты готовы**
```foo
// Функции с типизированными параметрами
fn strictAdd(x: int, y: int) {
    return x + y
}

fn greet(name: string, age: int) {
    return "Hello " + name + ", you are " + age + " years old!"
}

fn validate(flag: bool, threshold: float) {
    return flag && threshold > 0.5
}

// Поддержка всех примитивных типов: int, string, float, bool
let result = strictAdd(10, 5)  // = 15
let message = greet("Alice", 25)
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

### Методы массивов
- `array.length()` - получить длину массива
- `array.push(value)` - добавить элемент в конец массива

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

## Запуск тестов

```bash
go test ./test/... -v
```

**Важно**: Для каждой новой фичи обязательно нужно писать тесты! Тесты гарантируют стабильность языка и помогают избежать регрессий при добавлении новых возможностей.

## Архитектура

- `lexer/` - лексический анализатор
- `parser/` - синтаксический анализатор
- `ast/` - абстрактное синтаксическое дерево и интерпретатор
- `token/` - определения токенов
- `value/` - система типов
- `scope/` - система областей видимости
- `test/` - тесты

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
- [x] **Макросы** (определение и использование)
- [x] **Объекты** (литералы {key: value}) ✅ **тесты готовы**
- [x] **Массивы** (литералы [1, 2, 3]) ✅ **тесты готовы**
- [x] **Доступ к свойствам объектов** (obj.property) ✅ **тесты готовы**
- [x] **Методы массивов** (array.length(), array.push()) ✅ **тесты готовы**
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

### 📊 Покрытие тестами: 100%
Все основные функции языка покрыты unit-тестами:
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

### 🚧 В разработке / Планируется

#### 📋 Средний приоритет
- [ ] **Поддержка JSON** (parse/stringify)

#### 🌟 Низкий приоритет / Будущее
- [ ] **Стандартная библиотека** (std пакет)
- [ ] **Работа с файловой системой**
- [ ] **Регулярные выражения**
- [ ] **Поддержка многопоточности**
- [ ] **Пакетный менеджер зависимостей**

### ❌ Текущие ограничения
1. **Нет работы с файловой системой** (чтение/запись файлов)
2. **Нет поддержки JSON** (parse/stringify)
3. **Нет регулярных выражений**
4. **Нет поддержки многопоточности**
5. **Нет пакетного менеджера зависимостей**
6. **Строковая интерополяция** не поддерживает вложенные строки в выражениях
7. **Замыкания** не поддерживают возврат функций как значений

### 🎯 Недавно исправленные проблемы
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
- ✅ **Типизированные макросы** - макросы с типизированными параметрами (Type, StructType, EnumType, FnType) для безопасного метапрограммирования
- ✅ **Полиморфная система типов** - универсальные методы isStruct(), isPrimitive(), isEnum(), isFunction() для типов
- ✅ **100% покрытие тестами** - все основные функции покрыты unit-тестами