# Foo Language Syntax Highlighting - Installation Guide

Инструкции по установке syntax highlighting для Foo Language в различных редакторах.

## 🚀 Быстрая установка

### VS Code (рекомендуется)

**Метод 1: Через VSIX пакет**
```bash
# Соберите расширение
cd syntax/vscode
npm install -g vsce
vsce package

# Установите
code --install-extension foo-lang-1.0.0.vsix
```

**Метод 2: Локально**
```bash
# Windows
xcopy /E /I syntax\vscode %USERPROFILE%\.vscode\extensions\foo-lang-1.0.0

# macOS/Linux  
cp -r syntax/vscode ~/.vscode/extensions/foo-lang-1.0.0
```

### Neovim с nvim-treesitter

```bash
# 1. Установите tree-sitter
npm install -g tree-sitter-cli

# 2. Соберите грамматику
cd syntax/tree-sitter
tree-sitter generate

# 3. Добавьте в конфигурацию Neovim
```

```lua
-- ~/.config/nvim/init.lua
require'nvim-treesitter.configs'.setup {
  ensure_installed = { "foo" },
  highlight = { enable = true },
}
```

### Sublime Text

```bash
# Скопируйте TextMate grammar
cp syntax/vscode/syntaxes/foo.tmLanguage.json \
   "~/Library/Application Support/Sublime Text 3/Packages/User/"
```

### Atom

```bash
# Создайте пакет Atom
apm init foo-language-support
# Скопируйте grammars и themes
```

## 📋 Подробные инструкции

### VS Code - Пошаговая установка

1. **Подготовка**:
   ```bash
   cd syntax/vscode
   npm install
   ```

2. **Сборка расширения**:
   ```bash
   vsce package
   ```
   
3. **Установка**:
   - Откройте VS Code
   - Command Palette (Ctrl+Shift+P)
   - "Extensions: Install from VSIX..."
   - Выберите `foo-lang-1.0.0.vsix`

4. **Проверка**:
   - Откройте файл `.foo`
   - Syntax highlighting должен работать автоматически
   - Смените тему: Ctrl+K Ctrl+T → "Foo Dark" или "Foo Light"

### Neovim - Детальная настройка

1. **Установка зависимостей**:
   ```bash
   # Tree-sitter CLI
   npm install -g tree-sitter-cli
   
   # nvim-treesitter плагин (через packer)
   use 'nvim-treesitter/nvim-treesitter'
   ```

2. **Сборка грамматики**:
   ```bash
   cd syntax/tree-sitter
   tree-sitter generate
   tree-sitter build-wasm  # опционально
   ```

3. **Установка parser**:
   ```bash
   mkdir -p ~/.local/share/nvim/site/pack/packer/start/nvim-treesitter/parsers
   cp libtree-sitter-foo.so ~/.local/share/nvim/site/pack/packer/start/nvim-treesitter/parsers/
   ```

4. **Установка queries**:
   ```bash
   mkdir -p ~/.local/share/nvim/site/pack/packer/start/nvim-treesitter/queries/foo
   cp queries/highlights.scm ~/.local/share/nvim/site/pack/packer/start/nvim-treesitter/queries/foo/
   ```

5. **Конфигурация Neovim**:
   ```lua
   -- ~/.config/nvim/lua/treesitter-config.lua
   require'nvim-treesitter.configs'.setup {
     ensure_installed = { "foo", "lua", "vim" },
     
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
   
   -- Автоопределение типа файла
   vim.cmd([[
     augroup FooFiletype
       autocmd!
       autocmd BufNewFile,BufRead *.foo set filetype=foo
     augroup END
   ]])
   ```

### Emacs с tree-sitter

1. **Установка tree-sitter для Emacs**:
   ```elisp
   ;; ~/.emacs.d/init.el
   (use-package tree-sitter
     :ensure t
     :config
     (global-tree-sitter-mode)
     (add-hook 'tree-sitter-after-on-hook #'tree-sitter-hl-mode))

   (use-package tree-sitter-langs
     :ensure t
     :after tree-sitter)
   ```

2. **Настройка foo-mode**:
   ```elisp
   ;; Создайте ~/.emacs.d/foo-mode.el
   (defvar foo-mode-syntax-table
     (let ((table (make-syntax-table)))
       (modify-syntax-entry ?/ ". 124b" table)
       (modify-syntax-entry ?* ". 23" table)
       (modify-syntax-entry ?\n "> b" table)
       table))

   (define-derived-mode foo-mode prog-mode "Foo"
     "Major mode for editing Foo files."
     :syntax-table foo-mode-syntax-table
     (setq-local comment-start "// ")
     (setq-local comment-end ""))

   (add-to-list 'auto-mode-alist '("\\.foo\\'" . foo-mode))
   (add-to-list 'tree-sitter-major-mode-language-alist '(foo-mode . foo))
   ```

### Sublime Text

1. **Создание пакета**:
   ```bash
   mkdir -p "~/Library/Application Support/Sublime Text 3/Packages/Foo"
   ```

2. **Конвертация grammar**:
   ```bash
   # Используйте online конвертер или sublime-syntax-converter
   # Конвертируйте foo.tmLanguage.json в .sublime-syntax
   ```

3. **Настройка подсветки**:
   Создайте `Foo.sublime-syntax` на основе TextMate grammar

### Vim (без tree-sitter)

1. **Создание syntax файла**:
   ```vim
   " ~/.vim/syntax/foo.vim
   if exists("b:current_syntax")
     finish
   endif

   syn keyword fooKeyword let const fn struct enum interface impl extension
   syn keyword fooControl if else for match return yield break async await
   syn keyword fooBoolean true false
   syn keyword fooNull null

   syn match fooNumber '\v\d+'
   syn match fooFloat '\v\d+\.\d+'
   syn region fooString start='"' end='"' 
   syn region fooComment start='//' end='$'

   hi def link fooKeyword Keyword
   hi def link fooControl Statement  
   hi def link fooBoolean Boolean
   hi def link fooNull Constant
   hi def link fooNumber Number
   hi def link fooFloat Float
   hi def link fooString String
   hi def link fooComment Comment

   let b:current_syntax = "foo"
   ```

2. **Автодетекция типа файла**:
   ```vim
   " ~/.vim/ftdetect/foo.vim
   au BufRead,BufNewFile *.foo set filetype=foo
   ```

## 🔧 Устранение неполадок

### VS Code

**Проблема**: Подсветка не работает
```bash
# Решение 1: Перезагрузить окно
Ctrl+Shift+P → "Developer: Reload Window"

# Решение 2: Проверить ассоциацию файлов
Ctrl+Shift+P → "Change Language Mode" → "Foo"

# Решение 3: Переустановить расширение
```

**Проблема**: Темы не отображаются
```bash
# Убедитесь что themes зарегистрированы в package.json
# Перезапустите VS Code
```

### Neovim

**Проблема**: Parser не найден
```bash
# Проверьте установку
:TSInstallInfo foo

# Переустановите
:TSUninstall foo
:TSInstall foo
```

**Проблема**: Highlight queries не работают
```bash
# Проверьте путь к queries
ls ~/.local/share/nvim/site/pack/packer/start/nvim-treesitter/queries/foo/

# Перезагрузите конфигурацию
:luafile ~/.config/nvim/init.lua
```

### Tree-sitter сборка

**Проблема**: Ошибки компиляции
```bash
# Убедитесь что установлен C компилятор
gcc --version

# Очистите и пересоберите
tree-sitter clean
tree-sitter generate
```

## ✅ Проверка установки

Создайте файл `test.foo`:
```foo
// Тест syntax highlighting
let message: string = "Hello, Foo!"
fn greet(name: string) -> string {
    return "Hello, " + name
}

@macro_call(argument)
```

Проверьте что:
- ✅ Ключевые слова выделены синим
- ✅ Строки оранжевым/красным  
- ✅ Комментарии зеленым
- ✅ Функции желтым
- ✅ Макросы розовым/фиолетовым
- ✅ Auto-completion работает
- ✅ Bracket matching активен

## 📚 Дополнительные ресурсы

- [VS Code Extension API](https://code.visualstudio.com/api)
- [TextMate Grammar Guide](https://macromates.com/manual/en/language_grammars)
- [Tree-sitter Documentation](https://tree-sitter.github.io/tree-sitter/)
- [Nvim-treesitter Wiki](https://github.com/nvim-treesitter/nvim-treesitter/wiki)

## 🆘 Поддержка

Если у вас проблемы с установкой:
1. Проверьте требования к системе
2. Следуйте инструкциям точно
3. Создайте Issue в репозитории foo_lang_v2
4. Приложите логи ошибок и информацию о системе