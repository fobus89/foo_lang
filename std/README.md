# 📚 foo_lang Standard Libraries

Стандартная библиотека foo_lang v3 предоставляет мощные модули для работы с файловой системой, HTTP, криптографией и временем.

## 🚀 Быстрый старт

```foo
// Импорт модулей
import * as fs from "std/fs.foo"
import * as http from "std/http.foo" 
import * as crypto from "std/crypto.foo"
import * as time from "std/time.foo"
```

## 📦 Модули

### 🗂️ std.fs - Файловая система (40+ функций)

Полнофункциональная работа с файлами и директориями.

```foo
import * as fs from "std/fs.foo"

// Основные операции
fs.writeFile("demo.txt", "Hello, World!")
let content = fs.readFile("demo.txt").unwrap()

// Работа с директориями
fs.mkdir("testdir")
let files = fs.listDir(".").unwrap()

// Утилиты
let info = fs.getFileInfo("demo.txt").unwrap()
let found = fs.findFiles(".", ".foo").unwrap()
```

**Основные функции:**
- `readFile()`, `writeFile()`, `appendFile()`
- `mkdir()`, `mkdirAll()`, `listDir()`, `removeDir()`
- `copyFile()`, `copyDir()`, `getFileInfo()`
- `findFiles()`, `walkDir()`, `getDirSize()`
- `basename()`, `dirname()`, `extname()`, `joinPath()`

### 🌐 std.http - HTTP клиент и сервер (50+ функций)

Современный HTTP клиент и сервер с middleware поддержкой.

```foo
import * as http from "std/http.foo"

// HTTP клиент
let response = http.get("https://api.example.com/data")
let postResult = http.post("https://api.example.com/users", {
    "name": "Alice",
    "email": "alice@example.com"
})

// HTTP сервер
http.createServer()
http.get_route("/", fn(req, res) {
    res.send("Hello, World!")
})
http.startServer(8080)
```

**Основные функции:**
- **Клиент**: `get()`, `post()`, `put()`, `delete()`, `request()`
- **Сервер**: `createServer()`, `route()`, `startServer()`, `stopServer()`
- **Утилиты**: `buildURL()`, `parseURL()`, `parseCookies()`, `buildCookie()`
- **Middleware**: `use()`, `logger()`, `cors()`, `json()`
- **Auth**: `basicAuth()`, `parseBasicAuth()`

### 🔐 std.crypto - Криптография и безопасность (30+ функций)

Полный набор криптографических функций для безопасности.

```foo
import * as crypto from "std/crypto.foo"

// Хеширование
let hash = crypto.sha256("sensitive data")
let hmac = crypto.hmacSHA256("message", "secret_key")

// Кодирование
let encoded = crypto.base64Encode("data")
let hex = crypto.hexEncode("binary data")

// Пароли и токены
let passwordData = crypto.hashPassword("myPassword123").unwrap()
let isValid = crypto.verifyPassword("myPassword123", passwordData["hash"], passwordData["salt"])

// JWT токены
let token = crypto.generateJWT({"user": "alice"}, "secret").unwrap()
let payload = crypto.verifyJWT(token, "secret").unwrap()
```

**Основные функции:**
- **Хеши**: `md5()`, `sha1()`, `sha256()`, `sha512()`
- **HMAC**: `hmacSHA256()`, `hmacSHA1()`, `hmacMD5()`
- **Кодирование**: `base64Encode()`, `hexEncode()`, `base64URLEncode()`
- **Пароли**: `hashPassword()`, `verifyPassword()`, `checkPasswordStrength()`
- **Случайные данные**: `randomBytes()`, `randomString()`, `randomUUID()`
- **JWT**: `generateJWT()`, `verifyJWT()`, `generateSecretKey()`

### ⏰ std.time - Дата и время (60+ функций)

Комплексная работа с датами, временем и временными зонами.

```foo
import * as time from "std/time.foo"

// Текущее время и создание
let now = time.now()
let specificTime = time.fromString("2023-01-15 14:30:00")
let unixTime = time.fromUnix(1673789400)

// Форматирование
let formatted = time.format(now, "2006-01-02 15:04:05")
let iso = time.formatISO(now)
let human = time.formatHuman(now)

// Арифметика
let tomorrow = time.addDays(now, 1)
let nextWeek = time.addDays(now, 7)
let hourLater = time.addHours(now, 1)

// Сравнения и разности  
let diff = time.diffDays(nextWeek, now).unwrap()
let isBefore = time.before(now, tomorrow)
```

**Основные функции:**
- **Создание**: `now()`, `fromUnix()`, `fromString()`, `fromISO()`
- **Компоненты**: `year()`, `month()`, `day()`, `hour()`, `minute()`, `second()`
- **Форматирование**: `format()`, `formatISO()`, `formatHuman()`, `toUnix()`
- **Арифметика**: `addDays()`, `addHours()`, `addMonths()`, `addYears()`
- **Разности**: `diffDays()`, `diffHours()`, `diffMinutes()`, `diffSeconds()`
- **Сравнения**: `before()`, `after()`, `equal()`, `isSameDay()`
- **Утилиты**: `isLeapYear()`, `daysInMonth()`, `startOfDay()`, `endOfDay()`

## 🎯 Примеры использования

### Создание защищенного веб-API

```foo
import * as http from "std/http.foo"
import * as crypto from "std/crypto.foo"
import * as time from "std/time.foo"
import * as fs from "std/fs.foo"

// Middleware для аутентификации
fn authMiddleware(req, res, next) {
    let token = req["headers"]["Authorization"]
    if token != null {
        let payload = crypto.verifyJWT(token, "secret_key")
        if payload.isOk() {
            req["user"] = payload.unwrap()
            return next()
        }
    }
    res.status(401).send("Unauthorized")
}

// Создаем API сервер
http.createServer()
http.use(http.logger())
http.use(http.cors())
http.use(authMiddleware)

http.get_route("/api/time", fn(req, res) {
    res.json({
        "current_time": time.formatISO(time.now()),
        "user": req["user"]["username"]
    })
})

http.post_route("/api/upload", fn(req, res) {
    let filename = "uploads/" + crypto.randomUUID() + ".txt"
    let result = fs.writeFile(filename, req["body"])
    
    if result.isOk() {
        res.json({"status": "success", "file": filename})
    } else {
        res.status(500).json({"error": "Upload failed"})
    }
})

http.startServer(3000)
```

### Система логирования с шифрованием

```foo
import * as fs from "std/fs.foo"
import * as crypto from "std/crypto.foo"
import * as time from "std/time.foo"

fn secureLogger(message, level = "INFO") {
    let timestamp = time.formatISO(time.now())
    let logEntry = "[" + timestamp + "] " + level + ": " + message
    
    // Шифруем лог
    let encrypted = crypto.encryptText(logEntry, "log_encryption_key")
    if encrypted.isOk() {
        let logFile = "logs/" + time.format(time.now(), "2006-01-02") + ".log"
        fs.appendFile(logFile, encrypted.unwrap() + "\n")
    }
}

// Использование
secureLogger("User authentication successful", "INFO")
secureLogger("Failed login attempt", "WARNING")
secureLogger("System error occurred", "ERROR")
```

### Система бэкапов с расписанием

```foo
import * as fs from "std/fs.foo"
import * as crypto from "std/crypto.foo"
import * as time from "std/time.foo"

fn createBackup(sourceDir, backupDir) {
    let timestamp = time.format(time.now(), "20060102_150405")
    let backupName = "backup_" + timestamp + ".tar.gz"
    let backupPath = fs.joinPath(backupDir, backupName)
    
    // Создаем архив (в реальности нужна поддержка tar/gzip)
    let result = fs.copyDir(sourceDir, backupPath)
    
    if result.isOk() {
        // Создаем контрольную сумму
        let files = fs.listDir(backupPath).unwrap()
        let manifest = ""
        
        for let file in files {
            let content = fs.readFile(fs.joinPath(backupPath, file)).unwrap()
            let hash = crypto.sha256(content)
            manifest = manifest + file + ":" + hash + "\n"
        }
        
        fs.writeFile(fs.joinPath(backupPath, "MANIFEST.txt"), manifest)
        
        return Ok({
            "backup_name": backupName,
            "created_at": time.now(),
            "files_count": files.length()
        })
    }
    
    return result
}

// Автоматические бэкапы
fn scheduleBackups() {
    let lastBackup = time.now()
    
    while true {
        let currentTime = time.now()
        let hoursSinceBackup = time.diffHours(currentTime, lastBackup).unwrapOr(0)
        
        if hoursSinceBackup >= 24 {
            println("Creating daily backup...")
            let result = createBackup("./data", "./backups")
            
            if result.isOk() {
                println("Backup created successfully!")
                lastBackup = currentTime
            }
        }
        
        // Ждем час перед следующей проверкой
        time.sleep(3600000) // 1 час в миллисекундах
    }
}
```

## 🔧 Установка и настройка

1. Убедитесь, что у вас установлен foo_lang v3
2. Стандартные библиотеки находятся в директории `std/`
3. Используйте импорт: `import * as модуль from "std/модуль.foo"`

## 📖 Документация

Каждый модуль содержит подробные комментарии и примеры использования. Для получения полной документации изучите исходный код модулей:

- `std/fs.foo` - файловая система
- `std/http.foo` - HTTP клиент/сервер
- `std/crypto.foo` - криптография
- `std/time.foo` - дата и время

## 🚀 Производительность

Все стандартные библиотеки оптимизированы для высокой производительности:

- **Кеширование**: модули загружаются единожды
- **Ленивая инициализация**: функции создаются по требованию
- **Встроенные функции**: критический код использует Go реализации
- **JIT оптимизация**: горячие пути автоматически оптимизируются

## 🤝 Вклад в проект

Стандартные библиотеки foo_lang открыты для улучшений:

1. Добавьте новые функции в соответствующие модули
2. Создайте тесты для новой функциональности  
3. Обновите документацию
4. Отправьте pull request

## 📄 Лицензия

Стандартные библиотеки foo_lang распространяются под той же лицензией, что и основной интерпретатор.

---

**🎉 Наслаждайтесь разработкой с foo_lang standard libraries!**