package token

import "fmt"

type TokenType struct {
	Token
	Value string
	Line  int
	Col   int
	Pos   int    // Позиция в исходном тексте
}

func NewTokenType(token Token, value string, line, col int) TokenType {
	return TokenType{
		Col:   col,
		Line:  line,
		Token: token,
		Value: value,
		Pos:   -1, // По умолчанию
	}
}

func NewTokenTypeWithPos(token Token, value string, line, col, pos int) TokenType {
	return TokenType{
		Col:   col,
		Line:  line,
		Token: token,
		Value: value,
		Pos:   pos,
	}
}

func (t TokenType) String() string {
	return fmt.Sprintf("%s", t.Value)
}
