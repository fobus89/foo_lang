package builtin

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"foo_lang/value"
)

// Встроенные функции для работы с файловой системой

// ReadFile читает содержимое файла и возвращает строку
func ReadFile(args []*value.Value) *value.Value {
	if len(args) != 1 {
		return value.NewValue(fmt.Sprintf("readFile() requires exactly 1 argument, got %d", len(args)))
	}
	
	filename, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("readFile() requires a string argument (filename)")
	}
	
	// Читаем файл
	data, err := os.ReadFile(filename)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error reading file '%s': %v", filename, err))
	}
	
	return value.NewValue(string(data))
}

// WriteFile записывает содержимое в файл
func WriteFile(args []*value.Value) *value.Value {
	if len(args) != 2 {
		return value.NewValue(fmt.Sprintf("writeFile() requires exactly 2 arguments, got %d", len(args)))
	}
	
	filename, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("writeFile() first argument must be a string (filename)")
	}
	
	content, ok := args[1].Any().(string)
	if !ok {
		return value.NewValue("writeFile() second argument must be a string (content)")
	}
	
	// Создаем директории если нужно
	dir := filepath.Dir(filename)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return value.NewValue(fmt.Sprintf("Error creating directories for '%s': %v", filename, err))
		}
	}
	
	// Записываем файл
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error writing file '%s': %v", filename, err))
	}
	
	return value.NewValue("ok") // Успешная запись
}

// Exists проверяет существование файла или директории
func Exists(args []*value.Value) *value.Value {
	if len(args) != 1 {
		return value.NewValue(fmt.Sprintf("exists() requires exactly 1 argument, got %d", len(args)))
	}
	
	path, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("exists() requires a string argument (path)")
	}
	
	_, err := os.Stat(path)
	if err == nil {
		return value.NewValue(true) // Файл/директория существует
	}
	
	if os.IsNotExist(err) {
		return value.NewValue(false) // Файл/директория не существует
	}
	
	// Какая-то другая ошибка
	return value.NewValue(false)
}

// Mkdir создает директорию
func Mkdir(args []*value.Value) *value.Value {
	if len(args) != 1 {
		return value.NewValue(fmt.Sprintf("mkdir() requires exactly 1 argument, got %d", len(args)))
	}
	
	dirPath, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("mkdir() requires a string argument (directory path)")
	}
	
	// Создаем директорию со всеми родительскими директориями
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error creating directory '%s': %v", dirPath, err))
	}
	
	return value.NewValue("ok") // Успешное создание
}

// RemoveFile удаляет файл
func RemoveFile(args []*value.Value) *value.Value {
	if len(args) != 1 {
		return value.NewValue(fmt.Sprintf("removeFile() requires exactly 1 argument, got %d", len(args)))
	}
	
	filename, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("removeFile() requires a string argument (filename)")
	}
	
	err := os.Remove(filename)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error removing file '%s': %v", filename, err))
	}
	
	return value.NewValue("ok") // Успешное удаление
}

// CopyFile копирует файл
func CopyFile(args []*value.Value) *value.Value {
	if len(args) != 2 {
		return value.NewValue(fmt.Sprintf("copyFile() requires exactly 2 arguments, got %d", len(args)))
	}
	
	srcPath, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("copyFile() first argument must be a string (source path)")
	}
	
	dstPath, ok := args[1].Any().(string)
	if !ok {
		return value.NewValue("copyFile() second argument must be a string (destination path)")
	}
	
	// Открываем источник
	src, err := os.Open(srcPath)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error opening source file '%s': %v", srcPath, err))
	}
	defer src.Close()
	
	// Создаем директории для назначения если нужно
	dir := filepath.Dir(dstPath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return value.NewValue(fmt.Sprintf("Error creating directories for '%s': %v", dstPath, err))
		}
	}
	
	// Создаем файл назначения
	dst, err := os.Create(dstPath)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error creating destination file '%s': %v", dstPath, err))
	}
	defer dst.Close()
	
	// Копируем содержимое
	_, err = io.Copy(dst, src)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error copying from '%s' to '%s': %v", srcPath, dstPath, err))
	}
	
	return value.NewValue("ok") // Успешное копирование
}

// ListDir возвращает список файлов в директории
func ListDir(args []*value.Value) *value.Value {
	if len(args) != 1 {
		return value.NewValue(fmt.Sprintf("listDir() requires exactly 1 argument, got %d", len(args)))
	}
	
	dirPath, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("listDir() requires a string argument (directory path)")
	}
	
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error reading directory '%s': %v", dirPath, err))
	}
	
	// Создаем массив имен файлов
	var files []interface{}
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	
	return value.NewValue(files)
}

// IsFile проверяет, является ли путь файлом
func IsFile(args []*value.Value) *value.Value {
	if len(args) != 1 {
		return value.NewValue(fmt.Sprintf("isFile() requires exactly 1 argument, got %d", len(args)))
	}
	
	path, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("isFile() requires a string argument (path)")
	}
	
	info, err := os.Stat(path)
	if err != nil {
		return value.NewValue(false) // Файл не существует или ошибка
	}
	
	return value.NewValue(!info.IsDir()) // true если это файл (не директория)
}

// IsDir проверяет, является ли путь директорией
func IsDir(args []*value.Value) *value.Value {
	if len(args) != 1 {
		return value.NewValue(fmt.Sprintf("isDir() requires exactly 1 argument, got %d", len(args)))
	}
	
	path, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("isDir() requires a string argument (path)")
	}
	
	info, err := os.Stat(path)
	if err != nil {
		return value.NewValue(false) // Директория не существует или ошибка
	}
	
	return value.NewValue(info.IsDir()) // true если это директория
}

// GetFileSize возвращает размер файла в байтах
func GetFileSize(args []*value.Value) *value.Value {
	if len(args) != 1 {
		return value.NewValue(fmt.Sprintf("getFileSize() requires exactly 1 argument, got %d", len(args)))
	}
	
	filename, ok := args[0].Any().(string)
	if !ok {
		return value.NewValue("getFileSize() requires a string argument (filename)")
	}
	
	info, err := os.Stat(filename)
	if err != nil {
		return value.NewValue(fmt.Sprintf("Error getting file info for '%s': %v", filename, err))
	}
	
	return value.NewValue(info.Size()) // Размер в байтах
}

// FilesystemFunction представляет функцию файловой системы (аналогично MathFunction)
type FilesystemFunction struct {
	name string
	fn   func([]*value.Value) *value.Value
}

func (ff *FilesystemFunction) Eval() *value.Value {
	// Функции файловой системы не вызываются напрямую через Eval
	return value.NewValue(ff)
}

func (ff *FilesystemFunction) Call(args []*value.Value) *value.Value {
	return ff.fn(args)
}

func (ff *FilesystemFunction) Name() string {
	return ff.name
}

func (ff *FilesystemFunction) String() string {
	return "builtin function " + ff.name
}

// CreateFilesystemFunctions создает все функции для работы с файловой системой
func CreateFilesystemFunctions() map[string]*value.Value {
	functions := make(map[string]*value.Value)
	
	// Базовые операции с файлами
	functions["readFile"] = value.NewValue(&FilesystemFunction{"readFile", ReadFile})
	functions["writeFile"] = value.NewValue(&FilesystemFunction{"writeFile", WriteFile})
	functions["exists"] = value.NewValue(&FilesystemFunction{"exists", Exists})
	functions["removeFile"] = value.NewValue(&FilesystemFunction{"removeFile", RemoveFile})
	functions["copyFile"] = value.NewValue(&FilesystemFunction{"copyFile", CopyFile})
	functions["getFileSize"] = value.NewValue(&FilesystemFunction{"getFileSize", GetFileSize})
	
	// Операции с директориями
	functions["mkdir"] = value.NewValue(&FilesystemFunction{"mkdir", Mkdir})
	functions["listDir"] = value.NewValue(&FilesystemFunction{"listDir", ListDir})
	functions["isFile"] = value.NewValue(&FilesystemFunction{"isFile", IsFile})
	functions["isDir"] = value.NewValue(&FilesystemFunction{"isDir", IsDir})
	
	return functions
}

// InitializeFilesystemFunctions добавляет функции файловой системы в глобальную область видимости
func InitializeFilesystemFunctions(scopeStack ScopeStack) {
	functions := CreateFilesystemFunctions()
	for name, fn := range functions {
		scopeStack.Set(name, fn)
	}
}

