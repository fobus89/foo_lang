package token

type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT     // //
	MUL_COMMENT // /* */
	IDENT       // main
	INT         // 12345
	FLOAT       // 123.45
	CHAR        // 'a'
	STRING      // "abc"
	INTERP_STRING // "hello ${name}"
	MACRO         // macro
	QUOTE         // quote
	UNQUOTE       // unquote
	TYPEOF        // typeof
	STRUCT        // struct
	TYPE          // type
	EXPR          // Expr
	FNTYPE        // FnType
	STRUCTTYPE    // StructType
	ENUMTYPE      // EnumType
	INT_TYPE      // int (как тип параметра)
	STRING_TYPE   // string (как тип параметра)
	FLOAT_TYPE    // float (как тип параметра)  
	BOOL_TYPE     // bool (как тип параметра)
	PRINT
	PRINTLN
	LET
	CONST
	MATCH
	IF
	ELSE
	FOR
	FN
	ENUM
	RESULT
	OK
	ERR
	TRUE
	FALSE
	RETURN
	YIELD
	BREAK
	IMPORT
	EXPORT
	FROM
	AS
	EXTENSION
	INTERFACE
	IMPL
	ASYNC
	AWAIT
	SLEEP
	PROMISE

	operator_beg
	ADD        // +
	SUB        // -
	MUL        // *
	QUO        // /
	REM        // %
	QUESTION   // ?
	AND        // &
	AND_AND    // &&
	OR         // |
	OR_OR      // ||
	XOR        // ^
	SHL        // <<
	SHR        // >>
	AND_NOT    // &^
	GT         // >
	EQ         // =
	LT         // <
	EQ_GT      // =>
	GT_GT      // >>
	AT         // @
	Pound      // @
	LT_LT      // <<
	GT_EQ      // >=
	LT_EQ      // <=
	EQ_LT      // <=
	EQ_EQ      // ==
	NOT_EQ     // !=
	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL // ==
	LSS // <
	GTR // >
	NOT // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	DOT    // .

	RPAREN     // )
	RBRACK     // ]
	RBRACE     // }
	SEMICOLON  // ;
	COLON      // :
	UNDERSCORE // _
	TILDE      // ~
	operator_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND:  "&&",
	LOR:   "||",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQL: "==",
	LSS: "<",
	GTR: ">",
	EQ:  "=",
	NOT: "!",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	DOT:    ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",
}

func (tok Token) IsOperator() bool {
	return (operator_beg < tok && tok < operator_end)
}

func (tok Token) String() string {
	var s string
	{
		if 0 <= tok && tok < Token(len(tokens)) {
			s = tokens[tok]
		}
	}

	return s
}
