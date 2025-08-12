# Foo Language Syntax Highlighting

Комплексная поддержка синтаксиса для языка программирования Foo в различных редакторах.

## Поддерживаемые редакторы

- **VS Code** - полная поддержка с темами и грамматикой
- **Neovim/Emacs** - через Tree-sitter
- **Другие редакторы** - через TextMate grammar

## Установка

### VS Code

1. Скопируйте папку `vscode/` в директорию расширений VS Code:
   ```bash
   cp -r syntax/vscode ~/.vscode/extensions/foo-lang-1.0.0
   ```

2. Перезапустите VS Code

3. Откройте `.foo` файл - подсветка синтаксиса активируется автоматически

4. Выберите тему: `Ctrl+K Ctrl+T` → "Foo Dark" или "Foo Light"

### Neovim (с Tree-sitter)

1. Установите Tree-sitter CLI:
   ```bash
   npm install -g tree-sitter-cli
   ```

2. Скомпилируйте грамматику:
   ```bash
   cd syntax/tree-sitter
   tree-sitter generate
   ```

3. Добавьте в конфигурацию Neovim:
   ```lua
   require'nvim-treesitter.configs'.setup {
     ensure_installed = { "foo" },
     highlight = {
       enable = true,
     },
   }
   ```

### Emacs

1. Установите Tree-sitter для Emacs
2. Скомпилируйте грамматику из `tree-sitter/` директории
3. Добавьте в конфигурацию Emacs поддержку `.foo` файлов

## Возможности

### ✅ Полная поддержка синтаксиса Foo Language:

#### Ключевые слова
- **Объявления**: `let`, `const`, `fn`, `struct`, `enum`, `interface`, `impl`, `extension`
- **Управление потоком**: `if`, `else`, `for`, `match`, `return`, `yield`, `break`
- **Асинхронность**: `async`, `await`, `sleep`, `Promise`
- **Модули**: `import`, `export`, `from`, `as`
- **Макросы**: `macro`, `quote`, `unquote`, `typeof`, `type`

#### Типы данных
- **Примитивные**: `int`, `float`, `string`, `bool`
- **Union типы**: `string | int | float`
- **Optional типы**: `string?`
- **Result типы**: `Result`, `Ok`, `Err`
- **Коллекции**: массивы `[1, 2, 3]`, объекты `{key: value}`

#### Операторы
- **Арифметические**: `+`, `-`, `*`, `/`, `%`, `++`, `--`
- **Сравнения**: `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Логические**: `&&`, `||`, `!`
- **Присваивания**: `=`, `+=`, `-=`, `*=`, `/=`, `%=`
- **Стрелочные**: `=>`, `->`, `<-`
- **Тернарный**: `? :`

#### Специальные конструкции
- **Строковая интерполяция**: `"Hello ${name}"`
- **Generic функции**: `fn process<T: Interface>(item: T)`
- **Extension methods**: `extension string { fn reverse() }`
- **Макро вызовы**: `@macroName(args)`
- **Комментарии**: `//` и `/* */`

#### Встроенные функции
- **Вывод**: `print`, `println`
- **Математика**: `sin`, `cos`, `sqrt`, `abs`, `min`, `max`
- **Строки**: `strlen`, `charAt`, `substring`
- **Файлы**: `readFile`, `writeFile`, `exists`
- **HTTP**: `httpGet`, `httpPost`
- **Крипто**: `sha256Hash`, `base64Encode`
- **Regex**: `regexMatch`, `regexReplace`
- **Каналы**: `newChannel`, `send`, `receive`

### 🎨 Цветовые схемы

#### Foo Dark Theme
- **Фон**: Темно-серый (#1e1e1e)
- **Ключевые слова**: Синий (#569CD6)
- **Строки**: Оранжевый (#CE9178)
- **Комментарии**: Зеленый (#6A9955)
- **Функции**: Желтый (#DCDCAA)
- **Макросы**: Розовый (#FF6B9D)
- **Типы**: Бирюзовый (#4EC9B0)

#### Foo Light Theme  
- **Фон**: Белый (#ffffff)
- **Ключевые слова**: Синий (#0000FF)
- **Строки**: Красный (#A31515)
- **Комментарии**: Зеленый (#008000)
- **Функции**: Коричневый (#795E26)
- **Макросы**: Фиолетовый (#AF00DB)
- **Типы**: Бирюзовый (#008080)

## Примеры

Посмотрите `examples/sample.foo` для демонстрации всех возможностей подсветки синтаксиса.

## Разработка

### Структура проекта:
```
syntax/
├── vscode/                    # VS Code расширение
│   ├── package.json          # Манифест расширения
│   ├── language-configuration.json  # Настройки языка
│   ├── syntaxes/
│   │   └── foo.tmLanguage.json      # TextMate грамматика
│   └── themes/               # Цветовые темы
│       ├── foo-dark-color-theme.json
│       └── foo-light-color-theme.json
├── tree-sitter/             # Tree-sitter грамматика
│   ├── grammar.js           # Основная грамматика
│   └── package.json         # Настройки Tree-sitter
├── examples/                # Примеры кода
│   └── sample.foo           # Демонстрационный файл
└── README.md               # Документация
```

### Добавление новых фич:

1. **Новое ключевое слово**:
   - Добавьте в `vscode/syntaxes/foo.tmLanguage.json` в секцию `keywords`
   - Добавьте в `tree-sitter/grammar.js` как новое правило
   - Обновите примеры в `examples/sample.foo`

2. **Новый тип данных**:
   - Добавьте в секцию `types` в TextMate grammar
   - Добавьте правило в Tree-sitter `_type` choice
   - Добавьте цвета в темы

3. **Новый оператор**:
   - Добавьте в секцию `operators` 
   - Обновите `binary_expression` в Tree-sitter
   - Протестируйте приоритет операторов

## Лицензия

MIT License - см. основной проект foo_lang_v2.

## Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой фичи
3. Внесите изменения в соответствующие грамматики
4. Протестируйте на примерах
5. Создайте Pull Request

## Поддержка

- GitHub Issues: [foo_lang_v2 repository](https://github.com/your-username/foo_lang_v2)
- Документация: см. основной README.md проекта