package token

import "fmt"

type TokenType struct {
	Token
	Value string
	Line  int
	Col   int
}

func NewTokenType(token Token, value string, line, col int) TokenType {
	return TokenType{
		Col:   col,
		Line:  line,
		Token: token,
		Value: value,
	}
}

func (t TokenType) String() string {
	return fmt.Sprintf("%s", t.Value)
}
