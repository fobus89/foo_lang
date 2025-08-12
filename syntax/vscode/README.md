# Foo Language VS Code Extension

VS Code расширение для поддержки языка программирования Foo.

## Установка для разработки

1. Клонируйте репозиторий foo_lang_v2
2. Откройте папку `syntax/vscode` в VS Code
3. Нажмите F5 для запуска Extension Development Host
4. В новом окне откройте .foo файл для тестирования

## Установка как расширения

### Метод 1: Локальная установка
```bash
# Скопируйте папку в директорию расширений
cp -r syntax/vscode ~/.vscode/extensions/foo-lang-1.0.0
# Перезапустите VS Code
```

### Метод 2: Через VSIX (рекомендуется)
```bash
# Установите vsce
npm install -g vsce

# Соберите расширение
cd syntax/vscode
vsce package

# Установите .vsix файл
code --install-extension foo-lang-1.0.0.vsix
```

## Возможности

### ✅ Syntax Highlighting
- Полная поддержка всех конструкций Foo
- Ключевые слова, операторы, типы
- Строки с интерполяцией
- Комментарии (// и /* */)
- Макросы и их вызовы

### ✅ Language Features
- Auto-closing pairs для скобок и кавычек
- Comment toggling (Ctrl+/)
- Block comment toggle (Shift+Alt+A)
- Bracket matching
- Indentation rules

### ✅ Themes
- **Foo Dark** - темная тема, оптимизированная для Foo
- **Foo Light** - светлая тема
- Совместимость с другими темами VS Code

### ✅ File Association
- Автоматическое определение .foo файлов
- Корректная настройка синтаксиса

## Конфигурация

### Language Configuration
Файл `language-configuration.json` содержит:
- Настройки комментариев
- Auto-closing pairs для скобок
- Правила отступов
- Word pattern для выделения слов

### TextMate Grammar
Файл `syntaxes/foo.tmLanguage.json` содержит полную грамматику:
- Scope patterns для всех элементов языка
- Вложенные паттерны для строковой интерполяции
- Приоритеты для операторов
- Контекстные правила

## Customization

### Добавление новых ключевых слов
1. Откройте `syntaxes/foo.tmLanguage.json`
2. Найдите секцию `keywords` 
3. Добавьте новое слово в нужную группу
4. Перезапустите VS Code

### Изменение цветов темы
1. Откройте `themes/foo-dark-color-theme.json` или `foo-light-color-theme.json`
2. Измените значения в `tokenColors`
3. Сохраните и перезапустите VS Code

### Настройка отступов
1. Откройте `language-configuration.json`
2. Измените `indentationRules`
3. Настройте `increaseIndentPattern` и `decreaseIndentPattern`

## Разработка

### Структура файлов
```
vscode/
├── package.json                      # Манифест расширения
├── language-configuration.json       # Настройки языка  
├── syntaxes/
│   └── foo.tmLanguage.json           # TextMate грамматика
├── themes/
│   ├── foo-dark-color-theme.json     # Темная тема
│   └── foo-light-color-theme.json    # Светлая тема
└── README.md                         # Документация
```

### Тестирование

1. Откройте `examples/sample.foo` в Extension Development Host
2. Проверьте подсветку всех конструкций:
   - Ключевые слова должны быть синими
   - Строки оранжевыми
   - Комментарии зелеными
   - Функции желтыми
   - Макросы розовыми

3. Протестируйте функции:
   - Auto-closing для скобок
   - Comment toggle
   - Bracket matching
   - Go to matching bracket

### Публикация

```bash
# Войдите в Azure DevOps
vsce login <publisher>

# Опубликуйте расширение
vsce publish

# Или создайте VSIX для локального распространения
vsce package
```

## Troubleshooting

### Подсветка не работает
1. Убедитесь что файл имеет расширение `.foo`
2. Проверьте что расширение активировано: F1 → "Developer: Reload Window"
3. Проверьте консоль разработчика: Help → Toggle Developer Tools

### Темы не применяются
1. Откройте Command Palette (Ctrl+Shift+P)
2. Выберите "Preferences: Color Theme"
3. Найдите "Foo Dark" или "Foo Light"

### Скобки не закрываются автоматически
1. Проверьте настройки VS Code: `"editor.autoClosingBrackets": "always"`
2. Убедитесь что `language-configuration.json` загружен корректно

## Contributing

1. Форкните репозиторий foo_lang_v2
2. Внесите изменения в `syntax/vscode/`
3. Протестируйте в Extension Development Host
4. Создайте Pull Request

## License

MIT - см. основной проект foo_lang_v2