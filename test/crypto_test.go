package test

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"foo_lang/builtin"
	"foo_lang/parser"
	"foo_lang/scope"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCryptoHashFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "md5_hash",
			code: `
				let hash = md5Hash("hello world")
				println(hash)
			`,
			expected: "5eb63bbbe01eeed093cb22bb8f5acdc3", // MD5 of "hello world"
		},
		{
			name: "sha1_hash",
			code: `
				let hash = sha1Hash("hello world")
				println(hash)
			`,
			expected: "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", // SHA1 of "hello world"
		},
		{
			name: "sha256_hash",
			code: `
				let hash = sha256Hash("hello world")
				println(hash)
			`,
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", // SHA256 of "hello world"
		},
		{
			name: "sha512_hash",
			code: `
				let hash = sha512Hash("test")
				println(hash)
			`,
			expected: "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff", // SHA512 для "test"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureCryptoOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeCryptoFunctions(scope.GlobalScope)

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

func TestCryptoBase64Functions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "base64_encode",
			code: `
				let encoded = base64Encode("Hello, World!")
				println(encoded)
			`,
			expected: base64.StdEncoding.EncodeToString([]byte("Hello, World!")),
		},
		{
			name: "base64_decode",
			code: `
				let encoded = base64Encode("Hello, World!")
				let decoded = base64Decode(encoded)
				println(decoded)
			`,
			expected: "Hello, World!",
		},
		{
			name: "base64_url_encode",
			code: `
				let encoded = base64URLEncode("Hello+World/Test=")
				println(encoded)
			`,
			expected: base64.URLEncoding.EncodeToString([]byte("Hello+World/Test=")),
		},
		{
			name: "base64_url_decode",
			code: `
				let encoded = base64URLEncode("Hello+World/Test=")
				let decoded = base64URLDecode(encoded)
				println(decoded)
			`,
			expected: "Hello+World/Test=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureCryptoOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeCryptoFunctions(scope.GlobalScope)

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

func TestCryptoHexFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "hex_encode",
			code: `
				let encoded = hexEncode("Hello")
				println(encoded)
			`,
			expected: hex.EncodeToString([]byte("Hello")),
		},
		{
			name: "hex_decode",
			code: `
				let encoded = hexEncode("Hello")
				let decoded = hexDecode(encoded)
				println(decoded)
			`,
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureCryptoOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeCryptoFunctions(scope.GlobalScope)

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

func TestCryptoHMACFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		verify   func(result string) bool
	}{
		{
			name: "hmac_sha256",
			code: `
				let hmac = hmacSHA256("secret-key", "Hello, World!")
				println(hmac)
			`,
			verify: func(result string) bool {
				// Проверяем что результат является валидным hex строкой нужной длины
				// SHA256 HMAC должен быть 64 символа (32 байта * 2)
				if len(result) != 64 {
					return false
				}
				// Проверяем что все символы hex
				for _, c := range result {
					if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
						return false
					}
				}
				return true
			},
		},
		{
			name: "hmac_sha1",
			code: `
				let hmac = hmacSHA1("secret-key", "Hello, World!")
				println(hmac)
			`,
			verify: func(result string) bool {
				// SHA1 HMAC должен быть 40 символов (20 байт * 2)
				return len(result) == 40
			},
		},
		{
			name: "hmac_md5",
			code: `
				let hmac = hmacMD5("secret-key", "Hello, World!")
				println(hmac)
			`,
			verify: func(result string) bool {
				// MD5 HMAC должен быть 32 символа (16 байт * 2)
				return len(result) == 32
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureCryptoOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeCryptoFunctions(scope.GlobalScope)

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

func TestCryptoRandomFunctions(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		verify func(result string) bool
	}{
		{
			name: "random_bytes",
			code: `
				let bytes = randomBytes(16)
				println(bytes)
			`,
			verify: func(result string) bool {
				// 16 байт = 32 hex символа
				return len(result) == 32
			},
		},
		{
			name: "random_string_default",
			code: `
				let str = randomString(10)
				println(str)
			`,
			verify: func(result string) bool {
				// Проверяем длину
				if len(result) != 10 {
					return false
				}
				// Проверяем что содержит только допустимые символы
				alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
				for _, c := range result {
					found := false
					for _, a := range alphabet {
						if c == a {
							found = true
							break
						}
					}
					if !found {
						return false
					}
				}
				return true
			},
		},
		{
			name: "random_string_custom",
			code: `
				let str = randomString(8, "ABC123")
				println(str)
			`,
			verify: func(result string) bool {
				// Проверяем длину
				if len(result) != 8 {
					return false
				}
				// Проверяем что содержит только символы из алфавита
				alphabet := "ABC123"
				for _, c := range result {
					found := false
					for _, a := range alphabet {
						if c == a {
							found = true
							break
						}
					}
					if !found {
						return false
					}
				}
				return true
			},
		},
		{
			name: "random_int",
			code: `
				let num = randomInt(10, 100)
				println(num.toString())
			`,
			verify: func(result string) bool {
				// Просто проверяем что это число
				return len(result) > 0 && (result[0] >= '0' && result[0] <= '9')
			},
		},
		{
			name: "random_uuid",
			code: `
				let uuid = randomUUID()
				println(uuid)
			`,
			verify: func(result string) bool {
				// UUID должен быть в формате xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
				parts := strings.Split(result, "-")
				if len(parts) != 5 {
					return false
				}
				// Проверяем длины частей
				if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
					return false
				}
				// Проверяем версию (4)
				if len(parts[2]) > 0 && parts[2][0] != '4' {
					return false
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureCryptoOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeMathFunctions(scope.GlobalScope)
				builtin.InitializeCryptoFunctions(scope.GlobalScope)

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

func TestCryptoPasswordFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "password_hash_and_verify_success",
			code: `
				let password = "mysecretpassword"
				let hash = passwordHash(password)
				let isValid = passwordVerify(password, hash)
				println(isValid.toString())
			`,
			expected: "true",
		},
		{
			name: "password_hash_and_verify_failure",
			code: `
				let password = "mysecretpassword"
				let hash = passwordHash(password)
				let isValid = passwordVerify("wrongpassword", hash)
				println(isValid.toString())
			`,
			expected: "false",
		},
		{
			name: "password_hash_with_salt",
			code: `
				let password = "test123"
				let salt = "customsalt"
				let hash = passwordHash(password, salt)
				let isValid = passwordVerify(password, hash)
				println(isValid.toString())
			`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureCryptoOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeCryptoFunctions(scope.GlobalScope)

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

func TestCryptoConstantTimeCompare(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "constant_time_compare_equal",
			code: `
				let result = constantTimeCompare("hello", "hello")
				println(result.toString())
			`,
			expected: "true",
		},
		{
			name: "constant_time_compare_not_equal",
			code: `
				let result = constantTimeCompare("hello", "world")
				println(result.toString())
			`,
			expected: "false",
		},
		{
			name: "constant_time_compare_different_length",
			code: `
				let result = constantTimeCompare("hello", "hello world")
				println(result.toString())
			`,
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureCryptoOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeCryptoFunctions(scope.GlobalScope)

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

func TestCryptoErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErr  bool
		contains string
	}{
		{
			name: "invalid_base64",
			code: `
				let result = base64Decode("invalid-base64!")
				println(result)
			`,
			wantErr:  true,
			contains: "Error decoding base64",
		},
		{
			name: "invalid_hex",
			code: `
				let result = hexDecode("invalid-hex-string!")
				println(result)
			`,
			wantErr:  true,
			contains: "Error decoding hex",
		},
		{
			name: "random_bytes_too_large",
			code: `
				let result = randomBytes(2000)
				println(result)
			`,
			wantErr:  true,
			contains: "length must be between 1 and 1024",
		},
		{
			name: "random_string_empty_alphabet",
			code: `
				let result = randomString(10, "")
				println(result)
			`,
			wantErr:  true,
			contains: "alphabet cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureCryptoOutput(func() {
				scope.GlobalScope = scope.NewScopeStack()
				builtin.InitializeStringFunctions(scope.GlobalScope)
				builtin.InitializeCryptoFunctions(scope.GlobalScope)

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

func TestCryptoKnownVectors(t *testing.T) {
	// Тестируем с известными векторами для проверки корректности
	t.Run("known_md5_vector", func(t *testing.T) {
		expected := fmt.Sprintf("%x", md5.Sum([]byte("The quick brown fox jumps over the lazy dog")))
		
		result := captureCryptoOutput(func() {
			scope.GlobalScope = scope.NewScopeStack()
			builtin.InitializeStringFunctions(scope.GlobalScope)
			builtin.InitializeCryptoFunctions(scope.GlobalScope)

			exprs := parser.NewParser(`
				let hash = md5Hash("The quick brown fox jumps over the lazy dog")
				println(hash)
			`).Parse()
			for _, expr := range exprs {
				expr.Eval()
			}
		})

		if result != expected {
			t.Errorf("MD5 hash mismatch: expected %q, got %q", expected, result)
		}
	})

	t.Run("known_sha256_vector", func(t *testing.T) {
		expected := fmt.Sprintf("%x", sha256.Sum256([]byte("abc")))
		
		result := captureCryptoOutput(func() {
			scope.GlobalScope = scope.NewScopeStack()
			builtin.InitializeStringFunctions(scope.GlobalScope)
			builtin.InitializeCryptoFunctions(scope.GlobalScope)

			exprs := parser.NewParser(`
				let hash = sha256Hash("abc")
				println(hash)
			`).Parse()
			for _, expr := range exprs {
				expr.Eval()
			}
		})

		if result != expected {
			t.Errorf("SHA256 hash mismatch: expected %q, got %q", expected, result)
		}
	})
}

// captureCryptoOutput захватывает stdout для тестирования криптографических функций
func captureCryptoOutput(f func()) string {
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