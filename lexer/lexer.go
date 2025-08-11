package lexer

import (
	"foo_lang/token"
	"maps"
	"slices"
	"unicode"
)

type MapType[T comparable] map[T]token.Token

func (m MapType[T]) IsOperator(key T) bool {
	_, ok := m[key]
	return ok
}

func (m MapType[T]) Get(key T) token.Token {
	tok, ok := m[key]

	if !ok {
		return token.ILLEGAL
	}

	return tok
}

var operators = MapType[string]{
	"(":   token.LPAREN,
	")":   token.RPAREN,
	"[":   token.LBRACK,
	"]":   token.RBRACK,
	"{":   token.LBRACE,
	"}":   token.RBRACE,
	"+":   token.ADD,
	"++":  token.INC,
	"-":   token.SUB,
	"--":  token.DEC,
	"*":   token.MUL,
	"/":   token.QUO,
	"%":   token.REM,
	"&":   token.AND,
	"|":   token.OR,
	"^":   token.XOR,
	"?":   token.QUESTION,
	"_":   token.UNDERSCORE,
	";":   token.SEMICOLON,
	":":   token.COLON,
	"~":   token.TILDE,
	"!":   token.NOT,
	".":   token.DOT,
	",":   token.COMMA,
	"=":   token.EQ,
	"@":   token.AT,
	"#":   token.Pound,
	"<<":  token.LT_LT,
	">>":  token.GT_GT,
	"=>":  token.EQ_GT,
	"&^":  token.AND_NOT,
	"&^=": token.AND_NOT_ASSIGN,
	"&&":  token.AND_AND,
	"||":  token.OR_OR,
	">":   token.GT,
	"<":   token.LT,
	"==":  token.EQ_EQ,
	">=":  token.GT_EQ,
	"<=":  token.LT_EQ,
	"!=":  token.NOT_EQ,
}

var operatorKeys = func() []string {
	keys := make([]string, 0, len(operators))
	{
		for key := range maps.Keys(operators) {
			keys = append(keys, key)
		}
	}

	slices.SortFunc(keys, func(a, b string) int {
		return len(b) - len(a)
	})

	return keys
}()

var keywords = MapType[string]{
	"print":   token.PRINT,
	"println": token.PRINTLN,
	"macro":   token.MACRO,
	"quote":   token.QUOTE,
	"unquote": token.UNQUOTE,
	"typeof":  token.TYPEOF,
	"struct":  token.STRUCT,
	"type":    token.TYPE,
	"let":     token.LET,
	"const":   token.CONST,
	"match":   token.MATCH,
	"return":  token.RETURN,
	"yield":   token.YIELD,
	"break":   token.BREAK,
	"if":      token.IF,
	"else":    token.ELSE,
	"for":     token.FOR,
	"fn":      token.FN,
	"enum":    token.ENUM,
	"Result":  token.RESULT,
	"Ok":      token.OK,
	"Err":     token.ERR,
	"true":    token.TRUE,
	"false":   token.FALSE,
	"import":  token.IMPORT,
	"export":  token.EXPORT,
	"from":    token.FROM,
	"as":      token.AS,
}

type Lexer struct {
	input []rune
	line  int
	col   int
	pos   int
}

func NewLexer[T []rune | string | []byte](input T) *Lexer {
	return &Lexer{
		input: []rune(string(input)),
	}
}

func (l *Lexer) HasNext() bool {
	return l.pos < len(l.input)
}

func (l *Lexer) Next() rune {
	return l.NextN(1)
}

func (l *Lexer) NextN(pos int) rune {
	if !l.HasNext() {
		return '\000'
	}
	ch := l.input[l.pos]
	l.pos += pos

	if ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}

	return ch
}

func (l *Lexer) Match(symbol rune) bool {
	return l.Peek(0) == symbol
}

func (l *Lexer) MatchN(symbol rune, n int) bool {
	return l.Peek(n) == symbol
}

func (l *Lexer) MatchAllNext(symbol string) bool {
	for i, s := range symbol {
		if !l.MatchN(s, i) {
			return false
		}
	}

	l.NextN(len(symbol))

	return true
}

func (l *Lexer) MatchAndNext(symbol rune) bool {
	if l.Match(symbol) {
		l.Next()
		return true
	}
	return false
}

func (l *Lexer) MatchNAndNext(symbol rune, n int) bool {
	if l.MatchN(symbol, n) {
		l.NextN(n)
		return true
	}
	return false
}

func (l *Lexer) Peek(pos int) rune {
	pos = l.pos + pos

	if pos >= len(l.input) {
		return '\000'
	}

	return l.input[pos]
}

func (l *Lexer) Get(start, end int) []rune {
	if start == end {
		return nil
	}

	if end > len(l.input) {
		return nil
	}

	return l.input[start:end]
}

func (l *Lexer) SkipSpace() {
	for {
		// Пропускаем пробелы
		if unicode.IsSpace(l.Peek(0)) {
			l.Next()
			continue
		}
		
		// Пропускаем однострочные комментарии //
		if l.Peek(0) == '/' && l.Peek(1) == '/' {
			l.Next() // /
			l.Next() // /
			// Читаем до конца строки
			for l.HasNext() && l.Peek(0) != '\n' {
				l.Next()
			}
			continue
		}
		
		// Пропускаем многострочные комментарии /* */
		if l.Peek(0) == '/' && l.Peek(1) == '*' {
			l.Next() // /
			l.Next() // *
			// Читаем до */
			for l.HasNext() {
				if l.Peek(0) == '*' && l.Peek(1) == '/' {
					l.Next() // *
					l.Next() // /
					break
				}
				l.Next()
			}
			continue
		}
		
		// Если не пробел и не комментарий - выходим
		break
	}
}

func (l *Lexer) Tokens() []token.TokenType {
	var tokens []token.TokenType

	for l.HasNext() {
		t := l.Token()
		if t.Token == token.EOF {
			break
		}
		tokens = append(tokens, t)
	}

	tokens = append(tokens, token.TokenType{
		Token: token.EOF,
		Value: "EOF",
	})

	return tokens
}

func (l *Lexer) Token() token.TokenType {
	l.SkipSpace()

	ch := l.Peek(0)

	switch {

	case ch == '"':
		return l.ReadString()
	case ch == '`':
		return l.ReadStringFromat()
	case unicode.IsDigit(ch):
		return l.ReadNumber()

	case unicode.IsLetter(ch):
		return l.ReadIdentifier()

	case IsAnySymbol(ch):
		return l.ReadOperator()
	}

	if l.Next() == '\000' {
		return token.NewTokenType(token.EOF, "EOF", l.line, l.col)
	}

	return token.NewTokenType(token.ILLEGAL, "ILLEGAL", l.line, l.col)
}

func (l *Lexer) ReadStringFromat() token.TokenType {
	start := l.pos
	l.Next()

	for l.HasNext() {

		if l.Peek(0) == '`' {
			l.Next()
			break
		}

		l.Next()
	}

	str := l.Get(start+1, l.pos-1)
	{
		if len(str) == 0 {
			return token.NewTokenType(token.ILLEGAL, "ILLEGAL", l.line, l.col)
		}
	}

	return token.NewTokenType(token.STRING, string(str), l.line, l.col)
}

func (l *Lexer) ReadString() token.TokenType {
	start := l.pos
	l.Next()
	
	hasInterpolation := false

	for l.HasNext() {

		if l.Peek(0) == '"' {
			l.Next()
			break
		}
		
		// Проверяем на наличие ${
		if l.Peek(0) == '$' && l.Peek(1) == '{' {
			hasInterpolation = true
		}

		l.Next()
	}

	if start+1 == l.pos-1 {
		tokenType := token.STRING
		if hasInterpolation {
			tokenType = token.INTERP_STRING
		}
		return token.NewTokenType(tokenType, "", l.line, l.col)
	}

	str := l.Get(start+1, l.pos-1)
	{
		if len(str) == 0 {
			return token.NewTokenType(token.ILLEGAL, "ILLEGAL", l.line, l.col)
		}
	}

	tokenType := token.STRING
	if hasInterpolation {
		tokenType = token.INTERP_STRING
	}

	return token.NewTokenType(tokenType, string(str), l.line, l.col)
}

func (l *Lexer) ReadNumber() token.TokenType {
	start := l.pos
	var isDot bool
	var last rune = -1

	for l.HasNext() {
		ch := l.Peek(0)

		// Запрещённые последовательности: __, .., ._, _.
		if (last == '.' || last == '_') && (ch == '.' || ch == '_') {
			return token.NewTokenType(token.ILLEGAL, "ILLEGAL", l.line, l.col)
		}

		if ch == '.' {
			if isDot {
				return token.NewTokenType(token.ILLEGAL, "ILLEGAL", l.line, l.col)
			}
			isDot = true
			last = ch
			l.Next()
			continue
		}

		if ch == '_' {
			last = ch
			l.Next()
			continue
		}

		if !unicode.IsDigit(ch) {
			break
		}

		last = ch
		l.Next()
	}

	numbers := l.Get(start, l.pos)

	if len(numbers) == 0 || last == '_' || last == '.' {
		return token.NewTokenType(token.ILLEGAL, "ILLEGAL", l.line, l.col)
	}

	if isDot {
		return token.NewTokenType(token.FLOAT, string(numbers), l.line, l.col)
	}

	return token.NewTokenType(token.INT, string(numbers), l.line, l.col)
}

func (l *Lexer) ReadOperator() token.TokenType {
	for _, operator := range operatorKeys {
		if l.MatchAllNext(operator) {
			return token.NewTokenType(operators.Get(operator), operator, l.line, l.col)
		}
	}

	l.Next()

	return token.NewTokenType(token.ILLEGAL, "ILLEGAL", l.line, l.col)
}

func (l *Lexer) ReadIdentifier() token.TokenType {
	start := l.pos

	for l.HasNext() {

		if !unicode.IsLetter(l.Peek(0)) && !unicode.IsDigit(l.Peek(0)) && !(l.Peek(0) == '_') {
			break
		}

		l.Next()
	}

	symbols := l.Get(start, l.pos)
	{
		if len(symbols) == 0 {
			return token.NewTokenType(token.ILLEGAL, "ILLEGAL", l.line, l.col)
		}
	}

	tok := keywords.Get(string(symbols))
	{
		if tok == token.ILLEGAL {
			return token.NewTokenType(token.IDENT, string(symbols), l.line, l.col)
		}
	}

	return token.NewTokenType(tok, string(symbols), l.line, l.col)
}

func IsAnySymbol(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}
