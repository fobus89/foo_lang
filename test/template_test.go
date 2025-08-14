package test

import (
	"foo_lang/ast"
	"foo_lang/parser"
	"foo_lang/scope"
	"foo_lang/value"
	"testing"
)

func TestTemplateExpression(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "Simple variable interpolation",
			template: "Hello ${name}!",
			expected: "Hello World!",
		},
		{
			name:     "Multiple variables",
			template: "User: ${name}, Age: ${age}",
			expected: "User: Alice, Age: 25",
		},
		{
			name:     "Property access",
			template: "Type: ${user.Name}",
			expected: "Type: User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новую область видимости для каждого теста
			scope.GlobalScope.Push()
			defer scope.GlobalScope.Pop()

			// Устанавливаем переменные для интерполяции
			switch tt.name {
			case "Simple variable interpolation":
				scope.GlobalScope.Set("name", value.NewValue("World"))
			case "Multiple variables":
				scope.GlobalScope.Set("name", value.NewValue("Alice"))
				scope.GlobalScope.Set("age", value.NewValue(int64(25)))
			case "Property access":
				// Создаем TypeInfo для User
				userType := ast.NewStructTypeInfo("User", map[string]*ast.TypeInfo{
					"name": ast.NewPrimitiveTypeInfo("string"),
					"age":  ast.NewPrimitiveTypeInfo("int"),
				})
				scope.GlobalScope.Set("user", value.NewValue(userType))
			}

			// Создаем и выполняем TemplateExpr
			templateExpr := ast.NewTemplateExpr(tt.template)
			result := templateExpr.Eval()

			if result == nil {
				t.Fatalf("Expected result, got nil")
			}

			actual := result.String()
			if actual != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestMacroWithTemplateGeneration(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	// Тестовый код с макросом, использующим generate блок
	code := `
		struct User {
			name: string,
			age: int
		}

		macro generatePrint(structParam: StructType) {
			println("=== Macro-time: Processing " + structParam.Name + " ===")
			
			generate {
				fn print${structParam.Name}(obj) {
					println("Name: " + obj.name)
					println("Age: " + obj.age)
				}
			}
		}

		@generatePrint(User)
	`

	// Парсим и выполняем код
	p := parser.NewParser(code)
	exprs := p.Parse()

	// Выполняем все выражения
	for _, expr := range exprs {
		result := expr.Eval()
		if result != nil && result.IsReturn() {
			break
		}
	}

	// Проверяем, что макрос был выполнен без ошибок
	// (в данном случае просто проверяем что нет паники)
	t.Log("Macro with template generation executed successfully")
}

func TestTemplateWithForLoop(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	// Создаем TypeInfo с полями для тестирования
	userType := ast.NewStructTypeInfo("User", map[string]*ast.TypeInfo{
		"name": ast.NewPrimitiveTypeInfo("string"),
		"age":  ast.NewPrimitiveTypeInfo("int"),
		"email": ast.NewPrimitiveTypeInfo("string"),
	})

	scope.GlobalScope.Set("User", value.NewValue(userType))

	// Шаблон с циклом for
	template := `${for field in User.Fields { 
		println("Field: " + field.Name + " (Type: " + field.Type.Name + ")");
	}}`

	// Создаем и выполняем TemplateExpr
	templateExpr := ast.NewTemplateExpr(template)
	result := templateExpr.Eval()

	if result == nil {
		t.Fatalf("Expected result, got nil")
	}

	// Проверяем что результат содержит сгенерированный код для полей
	actual := result.String()
	t.Logf("Generated template result: %s", actual)

	// Проверяем что код содержит обработку полей
	expectedSubstrings := []string{"Field:", "name", "age", "email"}
	for _, expected := range expectedSubstrings {
		if !contains(actual, expected) {
			t.Errorf("Expected result to contain %q, got: %s", expected, actual)
		}
	}
}

func TestTemplateWithConditionals(t *testing.T) {
	// Инициализируем тестовое окружение
	InitTestEnvironment()

	tests := []struct {
		name      string
		template  string
		setupVars func()
		expected  string
	}{
		{
			name:     "Simple if true",
			template: "${if isEnabled { Enabled } else { Disabled }}",
			setupVars: func() {
				scope.GlobalScope.Set("isEnabled", value.NewValue(true))
			},
			expected: "Enabled",
		},
		{
			name:     "Simple if false", 
			template: "${if isEnabled { Enabled } else { Disabled }}",
			setupVars: func() {
				scope.GlobalScope.Set("isEnabled", value.NewValue(false))
			},
			expected: "Disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новую область видимости для каждого теста
			scope.GlobalScope.Push()
			defer scope.GlobalScope.Pop()

			// Устанавливаем переменные
			tt.setupVars()

			// Создаем и выполняем TemplateExpr
			templateExpr := ast.NewTemplateExpr(tt.template)
			result := templateExpr.Eval()

			if result == nil {
				t.Fatalf("Expected result, got nil")
			}

			actual := result.String()
			if actual != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr ||
		      containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}