package builtin

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"math/big"
)

// InitializeCryptoFunctions инициализирует встроенные криптографические функции
func InitializeCryptoFunctions(globalScope *scope.ScopeStack) {
	
	// ============ ХЕШ-ФУНКЦИИ ============
	
	// md5Hash - вычисление MD5 хеша
	md5HashFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: md5Hash() requires 1 argument (data)")
		}
		
		data, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: md5Hash() requires string argument")
		}
		
		hash := md5.Sum([]byte(data))
		return value.NewString(hex.EncodeToString(hash[:]))
	}
	globalScope.Set("md5Hash", value.NewValue(md5HashFunc))

	// sha1Hash - вычисление SHA1 хеша
	sha1HashFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: sha1Hash() requires 1 argument (data)")
		}
		
		data, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: sha1Hash() requires string argument")
		}
		
		hash := sha1.Sum([]byte(data))
		return value.NewString(hex.EncodeToString(hash[:]))
	}
	globalScope.Set("sha1Hash", value.NewValue(sha1HashFunc))

	// sha256Hash - вычисление SHA256 хеша
	sha256HashFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: sha256Hash() requires 1 argument (data)")
		}
		
		data, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: sha256Hash() requires string argument")
		}
		
		hash := sha256.Sum256([]byte(data))
		return value.NewString(hex.EncodeToString(hash[:]))
	}
	globalScope.Set("sha256Hash", value.NewValue(sha256HashFunc))

	// sha512Hash - вычисление SHA512 хеша
	sha512HashFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: sha512Hash() requires 1 argument (data)")
		}
		
		data, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: sha512Hash() requires string argument")
		}
		
		hash := sha512.Sum512([]byte(data))
		return value.NewString(hex.EncodeToString(hash[:]))
	}
	globalScope.Set("sha512Hash", value.NewValue(sha512HashFunc))

	// ============ BASE64 КОДИРОВАНИЕ ============
	
	// base64Encode - кодирование в Base64
	base64EncodeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: base64Encode() requires 1 argument (data)")
		}
		
		data, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: base64Encode() requires string argument")
		}
		
		encoded := base64.StdEncoding.EncodeToString([]byte(data))
		return value.NewString(encoded)
	}
	globalScope.Set("base64Encode", value.NewValue(base64EncodeFunc))

	// base64Decode - декодирование из Base64
	base64DecodeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: base64Decode() requires 1 argument (encodedData)")
		}
		
		encodedData, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: base64Decode() requires string argument")
		}
		
		decoded, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error decoding base64: %v", err))
		}
		
		return value.NewString(string(decoded))
	}
	globalScope.Set("base64Decode", value.NewValue(base64DecodeFunc))

	// base64URLEncode - кодирование в Base64 URL-safe
	base64URLEncodeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: base64URLEncode() requires 1 argument (data)")
		}
		
		data, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: base64URLEncode() requires string argument")
		}
		
		encoded := base64.URLEncoding.EncodeToString([]byte(data))
		return value.NewString(encoded)
	}
	globalScope.Set("base64URLEncode", value.NewValue(base64URLEncodeFunc))

	// base64URLDecode - декодирование из Base64 URL-safe
	base64URLDecodeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: base64URLDecode() requires 1 argument (encodedData)")
		}
		
		encodedData, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: base64URLDecode() requires string argument")
		}
		
		decoded, err := base64.URLEncoding.DecodeString(encodedData)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error decoding base64URL: %v", err))
		}
		
		return value.NewString(string(decoded))
	}
	globalScope.Set("base64URLDecode", value.NewValue(base64URLDecodeFunc))

	// ============ HEX КОДИРОВАНИЕ ============
	
	// hexEncode - кодирование в HEX
	hexEncodeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: hexEncode() requires 1 argument (data)")
		}
		
		data, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: hexEncode() requires string argument")
		}
		
		encoded := hex.EncodeToString([]byte(data))
		return value.NewString(encoded)
	}
	globalScope.Set("hexEncode", value.NewValue(hexEncodeFunc))

	// hexDecode - декодирование из HEX
	hexDecodeFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: hexDecode() requires 1 argument (hexData)")
		}
		
		hexData, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: hexDecode() requires string argument")
		}
		
		decoded, err := hex.DecodeString(hexData)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error decoding hex: %v", err))
		}
		
		return value.NewString(string(decoded))
	}
	globalScope.Set("hexDecode", value.NewValue(hexDecodeFunc))

	// ============ HMAC ФУНКЦИИ ============
	
	// hmacSHA256 - HMAC с SHA256
	hmacSHA256Func := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: hmacSHA256() requires 2 arguments (key, data)")
		}
		
		key, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: hmacSHA256() first argument must be string (key)")
		}
		
		data, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: hmacSHA256() second argument must be string (data)")
		}
		
		mac := hmac.New(sha256.New, []byte(key))
		mac.Write([]byte(data))
		signature := mac.Sum(nil)
		
		return value.NewString(hex.EncodeToString(signature))
	}
	globalScope.Set("hmacSHA256", value.NewValue(hmacSHA256Func))

	// hmacSHA1 - HMAC с SHA1
	hmacSHA1Func := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: hmacSHA1() requires 2 arguments (key, data)")
		}
		
		key, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: hmacSHA1() first argument must be string (key)")
		}
		
		data, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: hmacSHA1() second argument must be string (data)")
		}
		
		mac := hmac.New(sha1.New, []byte(key))
		mac.Write([]byte(data))
		signature := mac.Sum(nil)
		
		return value.NewString(hex.EncodeToString(signature))
	}
	globalScope.Set("hmacSHA1", value.NewValue(hmacSHA1Func))

	// hmacMD5 - HMAC с MD5
	hmacMD5Func := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: hmacMD5() requires 2 arguments (key, data)")
		}
		
		key, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: hmacMD5() first argument must be string (key)")
		}
		
		data, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: hmacMD5() second argument must be string (data)")
		}
		
		mac := hmac.New(md5.New, []byte(key))
		mac.Write([]byte(data))
		signature := mac.Sum(nil)
		
		return value.NewString(hex.EncodeToString(signature))
	}
	globalScope.Set("hmacMD5", value.NewValue(hmacMD5Func))

	// ============ ГЕНЕРАЦИЯ СЛУЧАЙНЫХ ДАННЫХ ============
	
	// randomBytes - генерация случайных байтов в hex формате
	randomBytesFunc := func(args []*value.Value) *value.Value {
		if len(args) != 1 {
			return value.NewString("Error: randomBytes() requires 1 argument (length)")
		}
		
		length, ok := args[0].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[0].Any().(float64); ok {
				length = int64(floatVal)
			} else {
				return value.NewString("Error: randomBytes() requires numeric argument")
			}
		}
		
		if length <= 0 || length > 1024 {
			return value.NewString("Error: randomBytes() length must be between 1 and 1024")
		}
		
		bytes := make([]byte, length)
		_, err := rand.Read(bytes)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error generating random bytes: %v", err))
		}
		
		return value.NewString(hex.EncodeToString(bytes))
	}
	globalScope.Set("randomBytes", value.NewValue(randomBytesFunc))

	// randomString - генерация случайной строки из заданного алфавита
	randomStringFunc := func(args []*value.Value) *value.Value {
		if len(args) < 1 || len(args) > 2 {
			return value.NewString("Error: randomString() requires 1-2 arguments (length, [alphabet])")
		}
		
		length, ok := args[0].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[0].Any().(float64); ok {
				length = int64(floatVal)
			} else {
				return value.NewString("Error: randomString() first argument must be numeric")
			}
		}
		
		if length <= 0 || length > 1024 {
			return value.NewString("Error: randomString() length must be between 1 and 1024")
		}
		
		// Алфавит по умолчанию
		alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		
		// Пользовательский алфавит
		if len(args) == 2 {
			if userAlphabet, ok := args[1].Any().(string); ok {
				if len(userAlphabet) == 0 {
					return value.NewString("Error: alphabet cannot be empty")
				}
				alphabet = userAlphabet
			} else {
				return value.NewString("Error: alphabet must be a string")
			}
		}
		
		result := make([]byte, length)
		alphabetLen := big.NewInt(int64(len(alphabet)))
		
		for i := range result {
			randomIndex, err := rand.Int(rand.Reader, alphabetLen)
			if err != nil {
				return value.NewString(fmt.Sprintf("Error generating random string: %v", err))
			}
			result[i] = alphabet[randomIndex.Int64()]
		}
		
		return value.NewString(string(result))
	}
	globalScope.Set("randomString", value.NewValue(randomStringFunc))

	// randomInt - генерация случайного числа в диапазоне
	randomIntFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: randomInt() requires 2 arguments (min, max)")
		}
		
		min, ok := args[0].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[0].Any().(float64); ok {
				min = int64(floatVal)
			} else {
				return value.NewString("Error: randomInt() first argument must be numeric")
			}
		}
		
		max, ok := args[1].Any().(int64)
		if !ok {
			// Пробуем float
			if floatVal, ok := args[1].Any().(float64); ok {
				max = int64(floatVal)
			} else {
				return value.NewString("Error: randomInt() second argument must be numeric")
			}
		}
		
		if min >= max {
			return value.NewString("Error: min must be less than max")
		}
		
		// Генерируем случайное число в диапазоне [min, max)
		rangeSize := max - min
		randomBig, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
		if err != nil {
			return value.NewString(fmt.Sprintf("Error generating random int: %v", err))
		}
		
		result := min + randomBig.Int64()
		return value.NewInt64(result)
	}
	globalScope.Set("randomInt", value.NewValue(randomIntFunc))

	// ============ UUID ГЕНЕРАЦИЯ ============
	
	// randomUUID - генерация UUID v4
	randomUUIDFunc := func(args []*value.Value) *value.Value {
		if len(args) != 0 {
			return value.NewString("Error: randomUUID() requires 0 arguments")
		}
		
		// Генерируем 16 случайных байтов
		bytes := make([]byte, 16)
		_, err := rand.Read(bytes)
		if err != nil {
			return value.NewString(fmt.Sprintf("Error generating UUID: %v", err))
		}
		
		// Устанавливаем биты версии (4) и варианта
		bytes[6] = (bytes[6] & 0x0f) | 0x40 // Версия 4
		bytes[8] = (bytes[8] & 0x3f) | 0x80 // Вариант RFC 4122
		
		// Форматируем как UUID строку
		uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			bytes[0:4],
			bytes[4:6],
			bytes[6:8],
			bytes[8:10],
			bytes[10:16])
			
		return value.NewString(uuid)
	}
	globalScope.Set("randomUUID", value.NewValue(randomUUIDFunc))

	// ============ УТИЛИТЫ СРАВНЕНИЯ ============
	
	// constantTimeCompare - безопасное сравнение строк (защита от timing атак)
	constantTimeCompareFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: constantTimeCompare() requires 2 arguments (a, b)")
		}
		
		a, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: constantTimeCompare() first argument must be string")
		}
		
		b, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: constantTimeCompare() second argument must be string")
		}
		
		// Используем hmac.Equal для constant-time сравнения
		result := hmac.Equal([]byte(a), []byte(b))
		return value.NewBool(result)
	}
	globalScope.Set("constantTimeCompare", value.NewValue(constantTimeCompareFunc))

	// ============ ФУНКЦИИ ДЛЯ ПАРОЛЕЙ ============
	
	// passwordHash - простое хеширование пароля (SHA256 + соль)
	passwordHashFunc := func(args []*value.Value) *value.Value {
		if len(args) < 1 || len(args) > 2 {
			return value.NewString("Error: passwordHash() requires 1-2 arguments (password, [salt])")
		}
		
		password, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: passwordHash() first argument must be string")
		}
		
		// Генерируем соль если не предоставлена
		salt := ""
		if len(args) == 2 {
			if userSalt, ok := args[1].Any().(string); ok {
				salt = userSalt
			} else {
				return value.NewString("Error: salt must be a string")
			}
		} else {
			// Генерируем случайную соль
			saltBytes := make([]byte, 16)
			_, err := rand.Read(saltBytes)
			if err != nil {
				return value.NewString(fmt.Sprintf("Error generating salt: %v", err))
			}
			salt = hex.EncodeToString(saltBytes)
		}
		
		// Хешируем пароль с солью
		combined := password + salt
		hash := sha256.Sum256([]byte(combined))
		
		// Возвращаем в формате salt:hash
		result := salt + ":" + hex.EncodeToString(hash[:])
		return value.NewString(result)
	}
	globalScope.Set("passwordHash", value.NewValue(passwordHashFunc))

	// passwordVerify - проверка пароля против хеша
	passwordVerifyFunc := func(args []*value.Value) *value.Value {
		if len(args) != 2 {
			return value.NewString("Error: passwordVerify() requires 2 arguments (password, hash)")
		}
		
		password, ok := args[0].Any().(string)
		if !ok {
			return value.NewString("Error: passwordVerify() first argument must be string")
		}
		
		storedHash, ok := args[1].Any().(string)
		if !ok {
			return value.NewString("Error: passwordVerify() second argument must be string")
		}
		
		// Разбираем сохраненный хеш (salt:hash)
		parts := splitString(storedHash, ":")
		if len(parts) != 2 {
			return value.NewString("Error: invalid hash format")
		}
		
		salt := parts[0]
		expectedHash := parts[1]
		
		// Вычисляем хеш для проверяемого пароля
		combined := password + salt
		hash := sha256.Sum256([]byte(combined))
		actualHash := hex.EncodeToString(hash[:])
		
		// Используем constant-time сравнение
		result := hmac.Equal([]byte(actualHash), []byte(expectedHash))
		return value.NewBool(result)
	}
	globalScope.Set("passwordVerify", value.NewValue(passwordVerifyFunc))
}

// splitString разбивает строку по разделителю (простая реализация)
func splitString(s, sep string) []string {
	if sep == "" {
		return []string{s}
	}
	
	var result []string
	start := 0
	
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	
	// Добавляем последнюю часть
	if start < len(s) {
		result = append(result, s[start:])
	}
	
	return result
}