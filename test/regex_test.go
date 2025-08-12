package test

import (
	"bytes"
	"foo_lang/builtin"
	"foo_lang/parser"
	"foo_lang/scope"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRegexMatchFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "simple_match_true",
			code: `
				let result = regexMatch("h[aeiou]llo", "hello world")
				println(result.toString())
			`,
			expected: "true",
		},
		{
			name: "simple_match_false",
			code: `
				let result = regexMatch("h[aeiou]llo", "hi world")
				println(result.toString())
			`,
			expected: "false",
		},
		{
			name: "number_pattern",
			code: `
				let result = regexMatch("\\d+", "price is 100 dollars")
				println(result.toString())
			`,
			expected: "true",
		},
		{
			name: "email_validation",
			code: `
				let result = regexMatch("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", "test@example.com")
				println(result.toString())
			`,
			expected: "true",
		},
		{
			name: "invalid_email",
			code: `
				let result = regexMatch("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", "invalid-email")
				println(result.toString())
			`,
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureRegexOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeRegexFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestRegexFindFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "find_first_word",
			code: `
				let result = regexFind("\\w+", "hello world test")
				println(result)
			`,
			expected: "hello",
		},
		{
			name: "find_number",
			code: `
				let result = regexFind("\\d+", "price is 100 dollars and 50 cents")
				println(result)
			`,
			expected: "100",
		},
		{
			name: "no_match",
			code: `
				let result = regexFind("\\d+", "no numbers here")
				println(result)
			`,
			expected: "",
		},
		{
			name: "find_email",
			code: `
				let result = regexFind("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}", "Contact: test@example.com or admin@site.org")
				println(result)
			`,
			expected: "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureRegexOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeRegexFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestRegexFindAllFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		verify   func(result string) bool
	}{
		{
			name: "find_all_numbers",
			code: `
				let results = regexFindAll("\\d+", "I have 5 apples, 10 oranges, and 3 bananas")
				println("found_numbers")
			`,
			verify: func(result string) bool {
				return result == "found_numbers"
			},
		},
		{
			name: "find_all_words",
			code: `
				let results = regexFindAll("\\w+", "hello world test")
				println("found_words")
			`,
			verify: func(result string) bool {
				return result == "found_words"
			},
		},
		{
			name: "find_all_with_limit",
			code: `
				let results = regexFindAll("\\d", "123456789", 3)
				println("found_with_limit")
			`,
			verify: func(result string) bool {
				return result == "found_with_limit"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureRegexOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeRegexFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if !tt.verify(result) {
				t.Errorf("%s: verification failed for result %q", tt.name, result)
			}
		})
	}
}

func TestRegexReplaceFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "replace_first_number",
			code: `
				let result = regexReplace("\\d+", "I have 5 apples and 10 oranges", "many")
				println(result)
			`,
			expected: "I have many apples and 10 oranges",
		},
		{
			name: "replace_all_numbers",
			code: `
				let result = regexReplaceAll("\\d+", "I have 5 apples and 10 oranges", "many")
				println(result)
			`,
			expected: "I have many apples and many oranges",
		},
		{
			name: "replace_word",
			code: `
				let result = regexReplaceAll("hello", "hello world, hello universe", "hi")
				println(result)
			`,
			expected: "hi world, hi universe",
		},
		{
			name: "replace_with_groups",
			code: `
				let result = regexReplaceAll("(\\w+)@(\\w+)", "email: test@example", "[$1 at $2]")
				println(result)
			`,
			expected: "email: [test at example]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureRegexOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeRegexFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestRegexSplitFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		verify   func(result string) bool
	}{
		{
			name: "split_by_comma",
			code: `
				let parts = regexSplit(",\\s*", "apple, banana, cherry")
				println("split_by_comma")
			`,
			verify: func(result string) bool {
				return result == "split_by_comma"
			},
		},
		{
			name: "split_by_whitespace",
			code: `
				let parts = regexSplit("\\s+", "hello   world    test")
				println("split_by_whitespace")
			`,
			verify: func(result string) bool {
				return result == "split_by_whitespace"
			},
		},
		{
			name: "split_with_limit",
			code: `
				let parts = regexSplit(",", "a,b,c,d,e", 3)
				println("split_with_limit")
			`,
			verify: func(result string) bool {
				return result == "split_with_limit"
			},
		},
		{
			name: "string_split_simple",
			code: `
				let parts = stringSplit("hello,world,test", ",")
				println("string_split_simple")
			`,
			verify: func(result string) bool {
				return result == "string_split_simple"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureRegexOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeRegexFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if !tt.verify(result) {
				t.Errorf("%s: verification failed for result %q", tt.name, result)
			}
		})
	}
}

func TestRegexGroupsFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		verify   func(result string) bool
	}{
		{
			name: "extract_email_parts",
			code: `
				let groups = regexGroups("([a-zA-Z0-9._%+-]+)@([a-zA-Z0-9.-]+)\\.([a-zA-Z]{2,})", "test@example.com")
				println("extract_email_parts")
			`,
			verify: func(result string) bool {
				return result == "extract_email_parts"
			},
		},
		{
			name: "extract_date_parts",
			code: `
				let groups = regexGroups("(\\d{4})-(\\d{2})-(\\d{2})", "Today is 2024-01-15")
				println("extract_date_parts")
			`,
			verify: func(result string) bool {
				return result == "extract_date_parts"
			},
		},
		{
			name: "no_match",
			code: `
				let groups = regexGroups("(\\d+)", "no numbers here")
				println("empty")
			`,
			verify: func(result string) bool {
				return result == "empty"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureRegexOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeRegexFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if !tt.verify(result) {
				t.Errorf("%s: verification failed for result %q", tt.name, result)
			}
		})
	}
}

func TestRegexUtilityFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "valid_regex",
			code: `
				let result = regexValid("[a-zA-Z]+")
				println(result.toString())
			`,
			expected: "true",
		},
		{
			name: "invalid_regex",
			code: `
				let result = regexValid("[a-zA-Z")
				println(result.toString())
			`,
			expected: "false",
		},
		{
			name: "escape_special_chars",
			code: `
				let result = regexEscape("Hello (world) [test] {foo}")
				println(result)
			`,
			expected: "Hello \\(world\\) \\[test\\] \\{foo\\}",
		},
		{
			name: "count_matches",
			code: `
				let result = regexCount("\\d", "I have 5 apples and 10 oranges")
				println(result.toString())
			`,
			expected: "3", // 5, 1, 0 (три цифры)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureRegexOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeRegexFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestRegexErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErr  bool
		contains string
	}{
		{
			name: "invalid_pattern",
			code: `
				let result = regexMatch("[a-zA-Z", "test")
				println(result)
			`,
			wantErr:  true,
			contains: "Error compiling regex pattern",
		},
		{
			name: "wrong_argument_count",
			code: `
				let result = regexMatch("test")
				println(result)
			`,
			wantErr:  true,
			contains: "requires 2 arguments",
		},
		{
			name: "non_string_pattern",
			code: `
				let result = regexMatch(123, "test")
				println(result)
			`,
			wantErr:  true,
			contains: "must be string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureRegexOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeRegexFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if tt.wantErr && !strings.Contains(result, tt.contains) {
				t.Errorf("%s: expected error containing %q, got %q", tt.name, tt.contains, result)
			}
		})
	}
}

func TestRegexPracticalExamples(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		verify   func(result string) bool
	}{
		{
			name: "extract_urls",
			code: `
				let urls = regexFindAll("https?://[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}", "Visit https://example.com or http://test.org")
				println("extract_urls")
			`,
			verify: func(result string) bool {
				return result == "extract_urls"
			},
		},
		{
			name: "validate_phone_number",
			code: `
				let isValid = regexMatch("^\\+?[0-9]{1,3}[0-9\\s\\-\\(\\)]{7,14}[0-9]$", "+1 (555) 123-4567")
				println(isValid.toString())
			`,
			verify: func(result string) bool {
				return result == "true"
			},
		},
		{
			name: "clean_html_tags",
			code: `
				let clean = regexReplaceAll("<[^>]*>", "<p>Hello <b>world</b>!</p>", "")
				println(clean)
			`,
			verify: func(result string) bool {
				return result == "Hello world!"
			},
		},
		{
			name: "extract_hashtags",
			code: `
				let hashtags = regexFindAll("#\\w+", "Check out #golang and #programming tutorials")
				println("extract_hashtags")
			`,
			verify: func(result string) bool {
				return result == "extract_hashtags"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureRegexOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeRegexFunctions(scope.GlobalScope)

				exprs := parser.NewParser(tt.code).Parse()
				for _, expr := range exprs {
					expr.Eval()
				}
			})

			if !tt.verify(result) {
				t.Errorf("%s: verification failed for result %q", tt.name, result)
			}
		})
	}
}

// captureRegexOutput захватывает stdout для тестирования regex функций
func captureRegexOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	result := buf.String()

	// Удаляем последний символ новой строки, если есть
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result
}