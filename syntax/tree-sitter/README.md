# Tree-sitter Grammar for Foo Language

Tree-sitter грамматика для языка программирования Foo, обеспечивающая быстрый и точный парсинг для редакторов Neovim, Emacs, и других.

## Установка

### Требования
- Node.js (>= 14.0)
- Tree-sitter CLI
- C компилятор (GCC, Clang, или MSVC)

### Установка Tree-sitter CLI
```bash
npm install -g tree-sitter-cli
```

### Сборка грамматики
```bash
cd syntax/tree-sitter
tree-sitter generate
tree-sitter build-wasm  # Опционально для веб-использования
```

### Тестирование
```bash
tree-sitter test
tree-sitter parse ../examples/sample.foo
```

## Интеграция с редакторами

### Neovim с nvim-treesitter

1. Добавьте в конфигурацию:
```lua
-- ~/.config/nvim/init.lua или ~/.config/nvim/lua/config.lua
require'nvim-treesitter.configs'.setup {
  ensure_installed = { "foo" },
  
  highlight = {
    enable = true,
    additional_vim_regex_highlighting = false,
  },
  
  indent = {
    enable = true
  },
  
  incremental_selection = {
    enable = true,
    keymaps = {
      init_selection = "gnn",
      node_incremental = "grn", 
      scope_incremental = "grc",
      node_decremental = "grm",
    },
  },
  
  textobjects = {
    select = {
      enable = true,
      lookahead = true,
      keymaps = {
        ["af"] = "@function.outer",
        ["if"] = "@function.inner",
        ["ac"] = "@class.outer",
        ["ic"] = "@class.inner",
      },
    },
  },
}
```

2. Установите highlight queries:
```bash
mkdir -p ~/.local/share/nvim/site/pack/packer/start/nvim-treesitter/queries/foo
cp queries/highlights.scm ~/.local/share/nvim/site/pack/packer/start/nvim-treesitter/queries/foo/
```

### Emacs с tree-sitter

1. Установите tree-sitter для Emacs:
```elisp
(use-package tree-sitter
  :ensure t
  :config
  (global-tree-sitter-mode)
  (add-hook 'tree-sitter-after-on-hook #'tree-sitter-hl-mode))

(use-package tree-sitter-langs
  :ensure t
  :after tree-sitter)
```

2. Зарегистрируйте foo-mode:
```elisp
(add-to-list 'tree-sitter-major-mode-language-alist '(foo-mode . foo))
(add-to-list 'auto-mode-alist '("\\.foo\\'" . foo-mode))
```

## Возможности грамматики

### ✅ Полная поддержка синтаксиса Foo

#### Основные конструкции
- Объявления переменных: `let`, `const`
- Функции: `fn name() {}`, анонимные функции
- Структуры: `struct Name { fields }`
- Интерфейсы: `interface Name { methods }`
- Enum'ы: `enum Status { ACTIVE, INACTIVE }`

#### Generic типы
- Generic функции: `fn process<T: Constraint>()`
- Ограничения типов: `T: Interface + Trait`
- Union типы: `string | int | float`
- Optional типы: `string?`

#### Async/await
- Асинхронные функции: `async fn()`
- Ожидание результата: `await promise`
- Promise API: `Promise.all()`, `Promise.any()`

#### Макросы и метапрограммирование
- Определение макросов: `macro name() {}`
- Вызов макросов: `@macroName(args)`
- Type introspection: `typeof()`, `type()`

#### Модульная система
- Импорты: `import { item } from "module"`
- Экспорты: `export fn name()`
- Алиасы: `import * as Alias from "module"`

### 🎯 Оптимизация производительности

Tree-sitter грамматика оптимизирована для:
- **Быстрый парсинг** - инкрементальное обновление AST
- **Восстановление от ошибок** - продолжение парсинга после ошибок  
- **Память** - эффективное использование памяти для больших файлов
- **Потоковый парсинг** - поддержка больших файлов

### 📊 Highlight Queries

Файл `queries/highlights.scm` обеспечивает:
- Семантическое выделение цветом
- Контекстно-зависимые цвета
- Поддержка вложенных конструкций
- Совместимость с темами редакторов

## Разработка

### Структура грамматики

```javascript
// Основные правила
source_file: $ => repeat($._statement)

_statement: $ => choice(
  $.variable_declaration,
  $.function_declaration, 
  $.struct_declaration,
  // ... другие конструкции
)
```

### Приоритеты операторов

Правильно настроенные приоритеты для математических и логических операторов:
```javascript
binary_expression: $ => choice(
  prec.left(11, seq($._expression, '*', $._expression)),  // Высший
  prec.left(10, seq($._expression, '+', $._expression)),
  prec.left(8, seq($._expression, '<', $._expression)),
  prec.left(4, seq($._expression, '&&', $._expression)),
  prec.left(3, seq($._expression, '||', $._expression))   // Низший
)
```

### Добавление новых конструкций

1. Обновите `grammar.js`:
```javascript
new_feature: $ => seq(
  'keyword',
  field('name', $.identifier),
  // ... правила
)
```

2. Добавьте в `_statement` или `_expression`:
```javascript
_statement: $ => choice(
  // ... существующие
  $.new_feature
)
```

3. Обновите highlight queries:
```scheme
"keyword" @keyword.new_feature
(new_feature name: (identifier) @entity.name)
```

4. Протестируйте:
```bash
tree-sitter generate
tree-sitter test
```

### Отладка грамматики

```bash
# Просмотр AST
tree-sitter parse file.foo

# Отладка конкретного правила  
tree-sitter parse file.foo --debug

# Проверка query
tree-sitter query queries/highlights.scm file.foo
```

## Тестирование

### Unit тесты
```bash
tree-sitter test
```

### Корпус тестов
Используйте файлы из `../examples/` для тестирования:
```bash
tree-sitter parse ../examples/sample.foo --quiet
```

### Производительность
```bash
time tree-sitter parse large_file.foo
```

## Contributing

1. Форкните репозиторий
2. Внесите изменения в `grammar.js`
3. Обновите `queries/highlights.scm` если нужно
4. Добавьте тесты в `test/corpus/`
5. Запустите `tree-sitter test`
6. Создайте Pull Request

## License

MIT - см. основной проект foo_lang_v2

## Ресурсы

- [Tree-sitter документация](https://tree-sitter.github.io/tree-sitter/)
- [Creating parsers](https://tree-sitter.github.io/tree-sitter/creating-parsers)
- [Grammar DSL](https://tree-sitter.github.io/tree-sitter/using-parsers#pattern-matching-with-queries)