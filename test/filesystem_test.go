package test

import (
	"testing"
	"os"
	"foo_lang/parser"
	"foo_lang/builtin"
	"foo_lang/scope"
)

func TestFilesystemOperations(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment(
		builtin.InitializeFilesystemFunctions,
		builtin.InitializeStringFunctions,
	)
	
	// Очистка перед тестом
	os.RemoveAll("test_fs_dir")
	
	code := `
// Создаем директорию
let mkdirResult = mkdir("test_fs_dir")
print("mkdir result: " + mkdirResult)

// Записываем файл
let content = "Hello, filesystem!\nТестовые данные на русском языке."
let writeResult = writeFile("test_fs_dir/test.txt", content)
print("writeFile result: " + writeResult)

// Проверяем существование
let fileExists = exists("test_fs_dir/test.txt")
let dirExists = exists("test_fs_dir")
print("File exists: " + fileExists.toString())
print("Dir exists: " + dirExists.toString())

// Читаем файл обратно
let readContent = readFile("test_fs_dir/test.txt")
print("Read content: " + readContent)

// Проверяем типы путей
let isFileCheck = isFile("test_fs_dir/test.txt")
let isDirCheck = isDir("test_fs_dir")
print("isFile: " + isFileCheck.toString())
print("isDir: " + isDirCheck.toString())

// Получаем размер файла
let fileSize = getFileSize("test_fs_dir/test.txt")
print("File size: " + fileSize.toString())

// Копируем файл
let copyResult = copyFile("test_fs_dir/test.txt", "test_fs_dir/copy.txt")
print("copyFile result: " + copyResult)

// Проверяем копию
let copyExists = exists("test_fs_dir/copy.txt")
print("Copy exists: " + copyExists.toString())

// Список файлов в директории
let files = listDir("test_fs_dir")
print("Directory files count: " + files.length().toString())
`

	exprs := parser.NewParser([]byte(code)).ParseWithoutScopeInit()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Check for errors in results
			if str, ok := result.Any().(string); ok && len(str) > 0 && str[0:5] == "Error" {
				t.Errorf("Filesystem operation failed: %s", str)
			}
		}
	}
	
	// Проверяем, что файлы действительно созданы
	if _, err := os.Stat("test_fs_dir/test.txt"); os.IsNotExist(err) {
		t.Error("Test file was not created")
	}
	
	if _, err := os.Stat("test_fs_dir/copy.txt"); os.IsNotExist(err) {
		t.Error("Copy file was not created")
	}
	
	// Очистка после теста
	os.RemoveAll("test_fs_dir")
}

func TestFilesystemErrors(t *testing.T) {
	// Инициализируем встроенные функции
	builtin.InitializeFilesystemFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	code := `
// Попытка прочитать несуществующий файл
let readError = readFile("nonexistent_file.txt")
print("Read nonexistent: " + readError)

// Попытка записи с некорректным путем
let writeError = writeFile("", "content")
print("Write with empty path: " + writeError)

// Проверка несуществующего файла
let existsResult = exists("definitely_does_not_exist.txt")
print("Nonexistent file exists: " + existsResult.toString())
`

	exprs := parser.NewParser([]byte(code)).ParseWithoutScopeInit()
	
	for _, expr := range exprs {
		expr.Eval() 
		// Ошибки обрабатываются через строки, не panic - это нормально
	}
}

func TestFilesystemComplexOperations(t *testing.T) {
	// Инициализируем встроенные функции
	builtin.InitializeFilesystemFunctions(scope.GlobalScope)
	builtin.InitializeStringFunctions(scope.GlobalScope)
	
	// Очистка перед тестом
	os.RemoveAll("complex_test_dir")
	
	code := `
// Создание вложенных директорий
let deepMkdir = mkdir("complex_test_dir/level1/level2/level3")
print("Deep mkdir: " + deepMkdir)

// Запись в глубокую структуру
let deepWrite = writeFile("complex_test_dir/level1/level2/level3/deep.txt", "Deep file content")
print("Deep write: " + deepWrite)

// JSON операции с файлами
let jsonData = {
    "name": "test",
    "values": [1, 2, 3],
    "nested": {"key": "value"}
}

let jsonContent = jsonStringify(jsonData)
let jsonWriteResult = writeFile("complex_test_dir/data.json", jsonContent)
print("JSON write: " + jsonWriteResult)

// Чтение JSON
let jsonReadContent = readFile("complex_test_dir/data.json")
print("JSON file read successfully: " + (jsonReadContent != "").toString())
`

	exprs := parser.NewParser([]byte(code)).ParseWithoutScopeInit()
	
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.Any() != nil {
			// Check for errors
			if str, ok := result.Any().(string); ok && len(str) > 5 && str[0:5] == "Error" {
				t.Errorf("Complex filesystem operation failed: %s", str)
			}
		}
	}
	
	// Проверяем создание глубокой структуры
	if _, err := os.Stat("complex_test_dir/level1/level2/level3/deep.txt"); os.IsNotExist(err) {
		t.Error("Deep nested file was not created")
	}
	
	// Очистка после теста
	os.RemoveAll("complex_test_dir")
}