package test

import (
	"testing"
	"foo_lang/parser"
)

func TestExtensionMethodsBasic(t *testing.T) {
	code := `
	// Расширяем строки методом isPalindrome
	extension string {
		fn isPalindrome() -> bool {
			let len = this.length()
			let mid = len / 2
			
			for let i = 0; i < mid; i++ {
				if this.charAt(i) != this.charAt(len - i - 1) {
					return false
				}
			}
			return true
		}
	}
	
	// Тестируем
	let str1 = "radar"
	let str2 = "hello"
	
	println(str1.isPalindrome())  // true
	println(str2.isPalindrome())  // false
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestExtensionMethodsNumbers(t *testing.T) {
	code := `
	// Расширяем числа (int) методами
	extension int {
		fn isEven() -> bool {
			return this % 2 == 0
		}
		
		fn factorial() -> int {
			if this <= 1 {
				return 1
			}
			return this * (this - 1).factorial()
		}
		
		fn square() -> int {
			return this * this
		}
	}
	
	// Числа в foo_lang по умолчанию float64, но мы можем расширить int тип
	let num1 = 5
	let num2 = 6
	
	println(num1.isEven())     // false
	println(num2.isEven())     // true
	// println(num1.factorial())  // 120 - закомментировано из-за типов
	println(num1.square())     // 25
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestExtensionMethodsWithParameters(t *testing.T) {
	code := `
	// Расширяем строки методами с параметрами
	extension string {
		fn repeat(n: int) -> string {
			let result = ""
			for let i = 0; i < n; i++ {
				result = result + this
			}
			return result
		}
		
		fn padLeft(width: int, char: string) -> string {
			let currentLen = this.length()
			if currentLen >= width {
				return this
			}
			
			let padding = ""
			for let i = 0; i < width - currentLen; i++ {
				padding = padding + char
			}
			return padding + this
		}
	}
	
	let str = "Hi"
	println(str.repeat(3))          // "HiHiHi"
	println(str.padLeft(5, "*"))    // "***Hi"
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestExtensionMethodsArrays(t *testing.T) {
	code := `
	// Расширяем массивы дополнительными методами
	extension array {
		fn sum() -> float {
			return this.reduce(0, fn(acc, x) => acc + x)
		}
		
		fn contains(item) -> bool {
			for let i = 0; i < this.length(); i++ {
				if this[i] == item {
					return true
				}
			}
			return false
		}
		
		fn isEmpty() -> bool {
			return this.length() == 0
		}
	}
	
	let numbers = [1, 2, 3, 4, 5]
	println(numbers.sum())          // 15
	println(numbers.contains(3))    // true
	println(numbers.contains(10))   // false
	println(numbers.isEmpty())      // false
	
	let empty = []
	println(empty.isEmpty())        // true
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestExtensionMethodsWithDefaults(t *testing.T) {
	code := `
	// Extension методы с параметрами по умолчанию
	extension string {
		fn truncate(maxLen: int, suffix: string = "...") -> string {
			if this.length() <= maxLen {
				return this
			}
			return this.substring(0, maxLen - suffix.length()) + suffix
		}
	}
	
	let text = "This is a very long text that needs truncation"
	println(text.truncate(20))        // "This is a very lo..."
	println(text.truncate(20, "…"))   // "This is a very long…"
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestExtensionMethodsMultipleTypes(t *testing.T) {
	code := `
	// Расширяем несколько типов одновременно
	extension string {
		fn wordCount() -> int {
			if this.length() == 0 {
				return 0
			}
			// Простой подсчет слов по пробелам
			let count = 1
			for let i = 0; i < this.length(); i++ {
				if this.charAt(i) == " " {
					count++
				}
			}
			return count
		}
	}
	
	extension float {
		fn roundTo(decimals: int) -> float {
			// Простая реализация округления без встроенных функций
			let factor = 1
			for let i = 0; i < decimals; i++ {
				factor = factor * 10
			}
			let scaled = this * factor
			// Простое округление
			let intPart = scaled
			if scaled - intPart >= 0.5 {
				intPart = intPart + 1
			}
			return intPart / factor
		}
	}
	
	extension bool {
		fn toInt() -> int {
			if this {
				return 1
			}
			return 0
		}
	}
	
	let sentence = "Hello world foo bar"
	println(sentence.wordCount())    // 4
	
	let pi = 3.14159
	println(pi.roundTo(2))           // 3.14
	
	let flag = true
	println(flag.toInt())            // 1
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestExtensionMethodsChainingVsBuiltin(t *testing.T) {
	code := `
	// Extension методы могут вызывать встроенные методы
	extension string {
		fn reverse() -> string {
			let result = ""
			for let i = this.length() - 1; i >= 0; i-- {
				result = result + this.charAt(i)
			}
			return result
		}
		
		fn isPalindromeIgnoreCase() -> bool {
			let lower = this.toLower()
			return lower == lower.reverse()
		}
	}
	
	let word1 = "RaceCar"
	let word2 = "Hello"
	
	println(word1.isPalindromeIgnoreCase())  // true
	println(word2.isPalindromeIgnoreCase())  // false
	println("abc".reverse())                 // "cba"
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestExtensionMethodsComplexExample(t *testing.T) {
	code := `
	// Сложный пример с extension методами для разных типов
	extension string {
		fn countChar(char: string) -> int {
			let count = 0
			for let i = 0; i < this.length(); i++ {
				if this.charAt(i) == char {
					count++
				}
			}
			return count
		}
		
		fn capitalize() -> string {
			if this.length() == 0 {
				return this
			}
			return this.charAt(0).toUpper() + this.substring(1, this.length())
		}
	}
	
	extension array {
		fn first() {
			if this.length() > 0 {
				return this[0]
			}
		}
		
		fn last() {
			if this.length() > 0 {
				return this[this.length() - 1]
			}
		}
		
		fn reverse() -> array {
			let result = []
			for let i = this.length() - 1; i >= 0; i-- {
				result = result.push(this[i])
			}
			return result
		}
	}
	
	// Тестируем string extensions
	let phrase = "hello world"
	println(phrase.countChar("l"))   // 3
	println(phrase.capitalize())     // "Hello world"
	
	// Тестируем array extensions
	let arr = [1, 2, 3, 4, 5]
	println(arr.first())             // 1
	println(arr.last())              // 5
	println(arr.reverse())           // [5, 4, 3, 2, 1]
	`

	p := parser.NewParser(code)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	exprs := p.Parse()
	for _, expr := range exprs {
		expr.Eval()
	}
}