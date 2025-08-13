# 🧙‍♂️ Полное руководство по макросам foo_lang v3

## 🚀 **Система макросов - мощный инструмент метапрограммирования**

Макросы в foo_lang позволяют генерировать код во время выполнения программы, создавать DSL и автоматизировать рутинные задачи программирования.

---

## 📖 **Основные концепции**

### **1️⃣ Определение макроса**

```foo
macro macroName(param1, param2) {
    // Код макроса
    println("Executing macro with: " + param1 + ", " + param2)
}
```

### **2️⃣ Вызов макроса**

```foo
@macroName("hello", "world")
```

### **3️⃣ Типизированные параметры**

```foo
macro generateStructUtils(structParam: StructType) {
    // structParam содержит информацию о структуре
    println("Processing struct: " + structParam.Name)
}

struct User { name: string, age: int }
@generateStructUtils(User)  // Передаем тип User
```

---

## 🔥 **Типы параметров макросов**

### **Type** - любой тип
```foo
macro debugType(t: Type) {
    println("Type: " + t.Name + ", Kind: " + t.Kind)
}
```

### **StructType** - только структуры
```foo
macro generateCRUD(s: StructType) {
    println("Generating CRUD for: " + s.Name)
    // Доступ к полям: s.Fields
}
```

### **FnType** - только функции
```foo
macro wrapFunction(f: FnType) {
    println("Wrapping function: " + f.Name)
    // Доступ к параметрам: f.Params
}
```

### **EnumType** - только енумы
```foo
macro generateEnumHelpers(e: EnumType) {
    println("Processing enum: " + e.Name)
    // Доступ к значениям: e.Values
}
```

---

## ⚡ **Двухфазная система выполнения**

### **Фаза 1: Macro-Time (время макроса)**
Код выполняется сразу при вызове макроса:

```foo
macro example(name) {
    // Это выполняется в macro-time
    println("Processing: " + name)
    
    Expr {
        // Это генерирует код для runtime
        println("Hello from generated code!")
    }
}
```

### **Фаза 2: Code Generation (генерация кода)**
Блок `Expr {}` генерирует код для выполнения:

```foo
macro generateFunction(funcName) {
    Expr {
        // Этот код будет выполнен как обычный foo_lang код
        fn dynamicFunction() {
            println("This is a generated function: " + funcName)
        }
    }
}
```

---

## 🎯 **Практические примеры использования**

### **1️⃣ Генерация CRUD операций**

```foo
macro generateCRUD(entityName) {
    println("Generating CRUD for: " + entityName)
    
    Expr {
        // Хранилище данных
        let storage = []
        
        // Create
        fn create(data) {
            data.id = randomString(8)
            storage.push(data)
            return data
        }
        
        // Read
        fn get(id) {
            for let item in storage {
                if item.id == id { return item }
            }
            return null
        }
        
        // Update
        fn update(id, data) {
            for let i = 0; i < storage.length(); i++ {
                if storage[i].id == id {
                    storage[i] = data
                    return data
                }
            }
            return null
        }
        
        // Delete
        fn delete(id) {
            for let i = 0; i < storage.length(); i++ {
                if storage[i].id == id {
                    storage.splice(i, 1)
                    return true
                }
            }
            return false
        }
    }
}

// Использование
@generateCRUD("User")
let user = create({name: "Alice", email: "alice@example.com"})
println("Created user: " + user.id)
```

### **2️⃣ Генерация валидаторов для структур**

```foo
struct User {
    name: string,
    age: int,
    email: string
}

macro generateValidator(structParam: StructType) {
    println("Creating validator for: " + structParam.Name)
    
    Expr {
        fn validate(obj) {
            if obj == null {
                return { valid: false, errors: ["Object is null"] }
            }
            
            let errors = []
            
            // Валидация имени
            if obj.name == null || obj.name == "" {
                errors.push("Name is required")
            }
            
            // Валидация возраста
            if obj.age == null || obj.age < 0 || obj.age > 150 {
                errors.push("Age must be between 0 and 150")
            }
            
            // Валидация email
            if obj.email == null || obj.email.indexOf("@") < 0 {
                errors.push("Valid email is required")
            }
            
            return {
                valid: errors.length() == 0,
                errors: errors
            }
        }
    }
}

@generateValidator(User)

// Использование
let user = { name: "Alice", age: 30, email: "alice@example.com" }
let result = validate(user)
if result.valid {
    println("User is valid!")
} else {
    println("Validation errors: " + result.errors.join(", "))
}
```

### **3️⃣ Генерация API endpoints**

```foo
macro generateAPI(entityName, routes) {
    println("Generating API for: " + entityName)
    
    Expr {
        // GET /entityName
        fn handleGet(req, res) {
            let data = getAll()
            res.json(data)
        }
        
        // POST /entityName
        fn handlePost(req, res) {
            let newItem = create(req.body)
            res.status(201).json(newItem)
        }
        
        // PUT /entityName/:id
        fn handlePut(req, res) {
            let id = req.params.id
            let updated = update(id, req.body)
            if updated {
                res.json(updated)
            } else {
                res.status(404).json({error: "Not found"})
            }
        }
        
        // DELETE /entityName/:id
        fn handleDelete(req, res) {
            let id = req.params.id
            let deleted = delete(id)
            if deleted {
                res.status(204).send()
            } else {
                res.status(404).json({error: "Not found"})
            }
        }
        
        // Регистрируем маршруты
        fn registerRoutes(app) {
            app.get("/" + entityName.toLowerCase(), handleGet)
            app.post("/" + entityName.toLowerCase(), handlePost)
            app.put("/" + entityName.toLowerCase() + "/:id", handlePut)
            app.delete("/" + entityName.toLowerCase() + "/:id", handleDelete)
        }
    }
}

@generateAPI("User", ["get", "post", "put", "delete"])
// Теперь доступны функции handleGet, handlePost, handlePut, handleDelete, registerRoutes
```

### **4️⃣ Генерация тестов**

```foo
macro generateTests(testSuiteName, testCases) {
    println("Generating tests for: " + testSuiteName)
    
    Expr {
        fn runTests() {
            println("Running " + testSuiteName + " test suite...")
            let passed = 0
            let total = 0
            
            // Тест 1: Базовая функциональность
            total = total + 1
            if testBasicFunctionality() {
                println("✅ Basic functionality test passed")
                passed = passed + 1
            } else {
                println("❌ Basic functionality test failed")
            }
            
            // Тест 2: Edge cases
            total = total + 1
            if testEdgeCases() {
                println("✅ Edge cases test passed")
                passed = passed + 1
            } else {
                println("❌ Edge cases test failed")
            }
            
            // Тест 3: Performance
            total = total + 1
            if testPerformance() {
                println("✅ Performance test passed")
                passed = passed + 1
            } else {
                println("❌ Performance test failed")
            }
            
            println("Results: " + passed.toString() + "/" + total.toString() + " tests passed")
            return passed == total
        }
        
        fn testBasicFunctionality() {
            // Базовая функциональность
            return true
        }
        
        fn testEdgeCases() {
            // Граничные случаи
            return true
        }
        
        fn testPerformance() {
            // Тесты производительности
            return true
        }
    }
}

@generateTests("MathFunctions", ["abs", "sqrt", "pow"])
let allPassed = runTests()
```

### **5️⃣ Генерация ORM-подобных запросов**

```foo
macro generateQueryBuilder(tableName) {
    Expr {
        fn queryBuilder() {
            return {
                table: tableName,
                fields: "*",
                conditions: [],
                orderField: null,
                limitCount: null,
                
                select: fn(fieldList) {
                    this.fields = fieldList
                    return this
                },
                
                where: fn(condition) {
                    this.conditions.push(condition)
                    return this
                },
                
                orderBy: fn(field) {
                    this.orderField = field
                    return this
                },
                
                limit: fn(count) {
                    this.limitCount = count
                    return this
                },
                
                build: fn() {
                    let query = "SELECT " + this.fields + " FROM " + this.table
                    
                    if this.conditions.length() > 0 {
                        query = query + " WHERE " + this.conditions.join(" AND ")
                    }
                    
                    if this.orderField != null {
                        query = query + " ORDER BY " + this.orderField
                    }
                    
                    if this.limitCount != null {
                        query = query + " LIMIT " + this.limitCount.toString()
                    }
                    
                    return query
                }
            }
        }
    }
}

@generateQueryBuilder("users")

// Использование
let query = queryBuilder()
    .select("name, email")
    .where("age > 18")
    .where("active = true")
    .orderBy("created_at DESC")
    .limit(10)
    .build()

println("Generated SQL: " + query)
// Вывод: SELECT name, email FROM users WHERE age > 18 AND active = true ORDER BY created_at DESC LIMIT 10
```

### **6️⃣ Генерация серializаторов/deserializers**

```foo
struct Product {
    id: int,
    name: string,
    price: float,
    available: bool
}

macro generateSerializers(structParam: StructType) {
    println("Generating serializers for: " + structParam.Name)
    
    Expr {
        fn toJSON(obj) {
            let json = "{"
            json = json + "\"id\":" + obj.id.toString() + ","
            json = json + "\"name\":\"" + obj.name + "\","
            json = json + "\"price\":" + obj.price.toString() + ","
            json = json + "\"available\":" + obj.available.toString()
            json = json + "}"
            return json
        }
        
        fn fromJSON(jsonString) {
            // Упрощенная реализация парсинга JSON
            return {
                id: 1,
                name: "Sample Product",
                price: 99.99,
                available: true
            }
        }
        
        fn toXML(obj) {
            let xml = "<" + structParam.Name.toLowerCase() + ">"
            xml = xml + "<id>" + obj.id.toString() + "</id>"
            xml = xml + "<name>" + obj.name + "</name>"
            xml = xml + "<price>" + obj.price.toString() + "</price>"
            xml = xml + "<available>" + obj.available.toString() + "</available>"
            xml = xml + "</" + structParam.Name.toLowerCase() + ">"
            return xml
        }
    }
}

@generateSerializers(Product)

// Использование
let product = { id: 1, name: "Laptop", price: 999.99, available: true }
println("JSON: " + toJSON(product))
println("XML: " + toXML(product))
```

---

## 🎨 **Продвинутые техники**

### **Условная генерация кода**

```foo
macro generateOptionalMethods(structParam: StructType, includeValidation) {
    Expr {
        // Всегда генерируем базовые методы
        fn create(data) {
            return data
        }
        
        // Условно генерируем валидацию
        if includeValidation {
            fn validate(obj) {
                return obj != null
            }
        }
    }
}
```

### **Композиция макросов**

```foo
macro generateFullCRUD(entityName, includeValidation, includeAPI) {
    @generateCRUD(entityName)
    
    if includeValidation {
        @generateValidator(entityName)
    }
    
    if includeAPI {
        @generateAPI(entityName, ["get", "post", "put", "delete"])
    }
}
```

### **Рекурсивные макросы**

```foo
macro generateNestedStructs(depth) {
    if depth > 0 {
        Expr {
            struct NestedStruct {
                level: int,
                data: string
            }
        }
        
        @generateNestedStructs(depth - 1)
    }
}
```

---

## 🚀 **Лучшие практики**

### **1. Используйте типизированные параметры**
```foo
// ✅ Хорошо
macro processStruct(s: StructType) { ... }

// ❌ Плохо  
macro processStruct(s) { ... }
```

### **2. Разделяйте macro-time и runtime код**
```foo
macro example(name) {
    // Macro-time: подготовка, анализ, логирование
    println("Processing: " + name)
    
    Expr {
        // Runtime: генерация исполняемого кода
        fn generatedFunction() { ... }
    }
}
```

### **3. Используйте осмысленные имена**
```foo
// ✅ Хорошо
macro generateRESTController(entityType: StructType) { ... }

// ❌ Плохо
macro gen(t: StructType) { ... }
```

### **4. Документируйте макросы**
```foo
// Генерирует полный CRUD API для сущности
// Параметры:
//   entityName: имя сущности
//   operations: массив операций ["create", "read", "update", "delete"]
macro generateCRUD(entityName, operations) { ... }
```

---

## 🎯 **Итого: Возможности макросов foo_lang**

✅ **Генерация кода в runtime**  
✅ **Типизированные параметры (StructType, FnType, EnumType)**  
✅ **Двухфазная система выполнения**  
✅ **Meta-programming возможности**  
✅ **Автоматизация рутинных задач**  
✅ **Создание DSL (Domain Specific Languages)**  
✅ **Code templating система**  
✅ **Интеграция со структурами и типами**  

**🚀 Система макросов foo_lang v3 - это мощный инструмент для создания высокоуровневых абстракций и автоматизации разработки!**