package parser

import (
	"fmt"
	"foo_lang/ast"
	"foo_lang/lexer"
	"foo_lang/token"
	"slices"
	"strconv"
	"strings"
)

type Parser struct {
	tokens []token.TokenType
	pos    int
}

func (p *Parser) error(msg string, tok token.TokenType) {
	panic(fmt.Sprintf("Parse error at line %d, column %d: %s (got '%s')", tok.Line, tok.Col, msg, tok.Value))
}

func NewParser[T []rune | string | []byte](input T) *Parser {
	tokens := lexer.NewLexer(input).Tokens()
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Peek(offset int) token.TokenType {
	pos := p.pos + offset
	if pos < 0 || pos >= len(p.tokens) {
		return token.TokenType{Token: token.EOF, Value: "EOF"}
	}
	return p.tokens[pos]
}

func (p *Parser) Next() token.TokenType {
	tok := p.Peek(0)
	p.pos++
	return tok
}

func (p *Parser) NextN(n int) token.TokenType {
	tok := p.Peek(n)
	p.pos += n
	return tok
}

func (p *Parser) Match(t token.Token) bool {
	return p.MatchN(t, 0)
}

func (p *Parser) MatchN(t token.Token, n int) bool {
	return p.Peek(n).Token == t
}

func (p *Parser) MatchAndNext(token token.Token) bool {
	if !p.Match(token) {
		return false
	}
	p.Next()
	return true
}

func (p *Parser) MatchAndNextN(token token.Token, n int) bool {
	if !p.MatchN(token, n) {
		return false
	}
	p.NextN(n)
	return true
}

func (p *Parser) MatchAllNext(tokens ...token.Token) bool {
	for i, t := range tokens {
		if !p.MatchN(t, i) {
			return false
		}
	}
	p.NextN(len(tokens))
	return true
}

func (p *Parser) MatchAll(tokens ...token.Token) bool {
	for i, t := range tokens {
		if !p.MatchN(t, i) {
			return false
		}
	}
	return true
}

func (p *Parser) MatchAnyNext(tokens ...token.Token) bool {
	if slices.ContainsFunc(tokens, p.Match) {
		p.NextN(1)
		return true
	}
	return false
}

func (p *Parser) Parse() []ast.Expr {
	var exprs []ast.Expr

	for !p.Match(token.EOF) {
		expr := p.Statement()
		exprs = append(exprs, expr)
	}

	return exprs
}

func (p *Parser) Statement() ast.Expr {
	//+= -= *= /= %=
	if p.MatchAllNext(token.IDENT, token.ADD, token.EQ) ||
		p.MatchAllNext(token.IDENT, token.SUB, token.EQ) ||
		p.MatchAllNext(token.IDENT, token.MUL, token.EQ) ||
		p.MatchAllNext(token.IDENT, token.QUO, token.EQ) ||
		p.MatchAllNext(token.IDENT, token.REM, token.EQ) {

		tok := p.Peek(-3)
		ident := tok.Value

		op := p.Peek(-2).Token

		expr := ast.NewBinaryExpr(ast.NewVarExpr(ident, nil), op, p.Expression())

		return ast.NewVarExpr(ident, expr)
	}

	//=
	if p.MatchAllNext(token.IDENT, token.EQ) {
		tok := p.Peek(-2)
		ident := tok.Value
		return ast.NewVarExpr(ident, p.Expression())
	}

	//@generateOther("1+1+1")
	if p.MatchAllNext(token.AT, token.IDENT, token.LPAREN) {
		return p.CallMacro()
	}

	if p.MatchAllNext(token.Pound, token.IDENT, token.LPAREN) {

		tok := p.Peek(-2)
		ident := tok.Value

		macros, ok := ast.Container[ident]
		if !ok {
			p.error("expected macro", p.Peek(0))
		}

		macro, ok := macros.Any().(*ast.FuncStatment)
		if !ok || !macro.IsMacro() {
			p.error("expected macro", p.Peek(0))
		}

		p.MatchAndNext(token.RPAREN)

		fn := p.FunctionStatement().(*ast.FuncStatment)

		arg := fmt.Sprintf("{name:%s, args:[%s]}", fn.Name(), strings.Join(fn.Params(), ","))

		call := ast.NewMacrosCallExpr(ident,
			[]ast.Expr{ast.NewLiteralString(arg)},
		)

		fmt.Println(call.Eval().String())

		return fn
	}

	if p.Match(token.FN) {
		return p.FunctionStatement()
	}

	if p.Match(token.MACROS) {
		return p.MacrosStatement()
	}

	if p.MatchAndNext(token.RETURN) {
		// Parse multiple return values: return a, b, c
		var returnValues []ast.Expr
		returnValues = append(returnValues, p.Expression())
		
		// Check for additional return values
		for p.MatchAndNext(token.COMMA) {
			returnValues = append(returnValues, p.Expression())
		}
		
		if len(returnValues) == 1 {
			// Single return value
			return ast.NewReturnExpr(returnValues[0])
		} else {
			// Multiple return values
			return ast.NewReturnExpr(ast.NewMultiReturnExpr(returnValues))
		}
	}

	if p.MatchAndNext(token.BREAK) {
		return ast.NewBreakExpr(nil)
	}

	if p.MatchAndNext(token.YIELD) {
		return ast.NewYieldExpr(p.Statement())
	}

	if p.MatchAllNext(token.CONST, token.IDENT, token.EQ) {
		tok := p.Peek(-2)
		ident := tok.Value

		if p.MatchAllNext(token.AT, token.IDENT, token.LPAREN) {
			return ast.NewLetExpr(ident, p.CallMacro())
		}

		if p.MatchAndNext(token.IF) {
			return ast.NewConstExpr(ident, p.IfStatement())
		}

		if p.Match(token.MATCH) {
			return ast.NewConstExpr(ident, p.MatchStatement())
		}

		if p.MatchAndNext(token.FOR) {
			return ast.NewConstExpr(ident, p.ForStatement())
		}

		return ast.NewConstExpr(ident, p.Expression())
	}

	// Check for multiple variable assignment: let a, b, c = expr
	if p.MatchAndNext(token.LET) {
		var names []string
		
		// Parse first identifier
		if !p.Match(token.IDENT) {
			p.error("expected identifier after let", p.Peek(0))
		}
		names = append(names, p.Next().Value)
		
		// Check for additional identifiers
		for p.MatchAndNext(token.COMMA) {
			if !p.Match(token.IDENT) {
				p.error("expected identifier after comma", p.Peek(0))
			}
			names = append(names, p.Next().Value)
		}
		
		if !p.MatchAndNext(token.EQ) {
			p.error("expected '=' after variable names", p.Peek(0))
		}
		
		// Multiple assignment case
		if len(names) > 1 {
			return ast.NewMultiAssignExpr(names, p.Expression())
		}
		
		// Single assignment case - keep existing logic
		ident := names[0]

		if p.MatchAllNext(token.AT, token.IDENT, token.LPAREN) {
			return ast.NewLetExpr(ident, p.CallMacro())
		}

		if p.MatchAndNext(token.IF) {
			return ast.NewLetExpr(ident, p.IfStatement())
		}

		if p.Match(token.MATCH) {
			return ast.NewLetExpr(ident, p.MatchStatement())
		}

		if p.MatchAndNext(token.FOR) {
			return ast.NewLetExpr(ident, p.ForStatement())
		}

		return ast.NewLetExpr(ident, p.Expression())
	}

	if p.MatchAnyNext(token.PRINT, token.PRINTLN) {
		isPrint := p.MatchN(token.PRINT, -1)
		return ast.NewPrintExpr(p.Expression(), isPrint)
	}

	if p.MatchAndNext(token.FOR) {
		return p.ForStatement()
	}

	if p.MatchAllNext(token.IF) {
		return p.IfStatement()
	}

	if p.Match(token.MATCH) {
		return p.MatchStatement()
	}

	if p.MatchAndNext(token.ENUM) {
		return p.EnumStatement()
	}

	if p.MatchAndNext(token.IMPORT) {
		return p.ImportStatement()
	}

	if p.MatchAndNext(token.EXPORT) {
		return p.ExportStatement()
	}

	return p.Expression()
}

func (p *Parser) ForStatement() ast.Expr {
	init := p.Statement()
	{
		if !p.MatchAndNext(token.SEMICOLON) {
			p.error("expected expression", p.Peek(0))
		}
	}

	condition := p.Expression()
	{
		if !p.MatchAndNext(token.SEMICOLON) {
			p.error("expected expression", p.Peek(0))
		}
	}

	step := p.Expression()

	body := p.BlockStatement()

	return ast.NewForExpr(init, condition, step, body)
}

func (p *Parser) CallMacroWithParams(args []ast.Expr) string {
	tok := p.Peek(-2)
	ident := tok.Value

	macros, ok := ast.Container[ident]
	if !ok {
		p.error("expected macro", p.Peek(0))
	}

	macro, ok := macros.Any().(*ast.FuncStatment)
	if !ok || !macro.IsMacro() {
		p.error("expected macro", p.Peek(0))
	}

	call := ast.NewMacrosCallExpr(ident, args)

	return call.Eval().String()
}

func (p *Parser) CallMacro() ast.Expr {
	tok := p.Peek(-2)
	ident := tok.Value

	macros, ok := ast.Container[ident]
	if !ok {
		p.error("expected macro", p.Peek(0))
	}

	macro, ok := macros.Any().(*ast.FuncStatment)
	if !ok || !macro.IsMacro() {
		p.error("expected macro", p.Peek(0))
	}

	start := p.pos - 3

	var args []ast.Expr

	for !p.MatchAndNext(token.RPAREN) {
		args = append(args, p.Statement())
		p.MatchAndNext(token.COMMA)
	}

	call := ast.NewMacrosCallExpr(ident, args)

	outCode := call.Eval().String()
	genTokens := lexer.NewLexer(outCode).Tokens()

	end := p.pos // после ')'

	before := slices.Clone(p.tokens[:start])
	after := p.tokens[end:]

	p.tokens = append(append(before, genTokens[:len(genTokens)-1]...), after...)
	p.pos = start

	return p.Statement()
}

func (p *Parser) FunctionStatement() ast.Expr {
	if !p.MatchAllNext(token.FN, token.IDENT, token.LPAREN) {
		p.error("expected fn", p.Peek(0))
	}

	identTok := p.Peek(-2)

	var args []map[string]ast.Expr

	for !p.MatchAndNext(token.RPAREN) {

		argToken := p.Next()
		{
			if argToken.Token != token.IDENT {
				p.error("expected identifier", argToken)
			}
		}

		key := argToken.Value
		{
			if p.MatchAndNext(token.EQ) {
				args = append(args, map[string]ast.Expr{
					key: p.Expression(),
				})
			} else {
				args = append(args, map[string]ast.Expr{
					key: nil,
				})
			}
		}

		p.MatchAllNext(token.COMMA)
	}

	body := p.BlockStatement()

	return ast.NewFuncStatment(identTok.Value, args, body, false)
}

func (p *Parser) MacrosStatement() ast.Expr {
	if !p.MatchAllNext(token.MACROS, token.IDENT, token.LPAREN) {
		p.error("expected fn", p.Peek(0))
	}

	identTok := p.Peek(-2)

	var args []map[string]ast.Expr

	for !p.MatchAndNext(token.RPAREN) {

		argToken := p.Next()
		{
			if argToken.Token != token.IDENT {
				p.error("expected identifier", argToken)
			}
		}

		key := argToken.Value
		{
			if p.Match(token.EQ) {
				args = append(args, map[string]ast.Expr{
					key: p.Expression(),
				})
			} else {
				args = append(args, map[string]ast.Expr{
					key: nil,
				})
			}
		}

		p.MatchAllNext(token.COMMA)
	}

	body := p.BlockStatement()

	return ast.NewFuncStatment(identTok.Value, args, body, true)
}

func (p *Parser) IfStatement() ast.Expr {
	conditions := []ast.Expr{p.Expression()}
	then := []ast.Expr{p.BlockStatement()}

	for p.MatchAllNext(token.ELSE, token.IF) {
		conditions = append(conditions, p.Expression())
		then = append(then, p.BlockStatement())
	}

	if p.MatchAllNext(token.ELSE, token.LBRACE) {
		p.NextN(-1)
		return ast.NewIfExpr(conditions, then, p.BlockStatement())
	}

	return ast.NewIfExpr(conditions, then, nil)
}

func (p *Parser) MatchStatement() ast.Expr {
	if !p.MatchAndNext(token.MATCH) {
		p.error("expected match", p.Peek(0))
	}

	p.MatchAndNext(token.LPAREN)
	condition := p.Expression()
	p.MatchAndNext(token.RPAREN)

	if !p.MatchAndNext(token.LBRACE) {
		p.error("expected '{'", p.Peek(0))
	}

	var arms []ast.MatchArm

	var underscore *ast.MatchArm

	for !p.MatchAllNext(token.RBRACE) {

		if p.MatchAndNext(token.UNDERSCORE) {
			if underscore != nil {
				p.error("only one underscore is allowed", p.Peek(0))
			}

			if !p.MatchAllNext(token.EQ_GT) {
				p.error("expected '=>'", p.Peek(0))
			}

			var resultExpr ast.Expr
			{
				if !p.Match(token.LBRACE) {
					resultExpr = p.Statement()
				} else {
					resultExpr = p.BlockStatement()
				}
			}

			tmp := ast.NewMatchArm(ast.NewLiteralString("_"), resultExpr)
			underscore = &tmp
			p.MatchAllNext(token.COMMA)
			continue
		}

		caseExpr := p.Expression()

		if !p.MatchAllNext(token.EQ_GT) {
			p.error("expected '=>'", p.Peek(0))
		}

		var resultExpr ast.Expr

		if !p.Match(token.LBRACE) {
			resultExpr = p.Statement()
		} else {
			resultExpr = p.BlockStatement()
		}

		p.MatchAllNext(token.COMMA)

		arms = append(arms, ast.NewMatchArm(caseExpr, resultExpr))
	}

	if underscore != nil {
		arms = append(arms, *underscore)
	}

	return ast.NewMatchExpr(condition, arms)
}

func (p *Parser) BlockStatement() ast.Expr {
	if !p.MatchAndNext(token.LBRACE) {
		p.error("expected {", p.Peek(0))
	}

	var statments []ast.Expr

	for !p.MatchAndNext(token.RBRACE) {
		statments = append(statments, p.Statement())
	}

	return ast.NewBodyStatment(statments)
}

func (p *Parser) Expression() ast.Expr {
	return p.Conditional()
}

func (p *Parser) Conditional() ast.Expr {
	expr := p.Logical()

	// true ? 1 : 22
	if p.MatchAndNext(token.QUESTION) {
		thenBranch := p.Expression()
		if !p.MatchAndNext(token.COLON) {
			p.error("expected ':' in conditional expression", p.Peek(0))
		}
		elseBranch := p.Expression()
		return ast.NewConditionalExpr(expr, thenBranch, elseBranch)
	}

	return expr
}

func (p *Parser) Logical() ast.Expr {
	expr := p.Comparison()

	for {
		if p.MatchAndNext(token.AND_AND) {
			expr = ast.NewBinaryExpr(expr, token.AND_AND, p.Comparison())
		} else if p.MatchAndNext(token.OR_OR) {
			expr = ast.NewBinaryExpr(expr, token.OR_OR, p.Comparison())
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) Comparison() ast.Expr {
	expr := p.Addition()

	for {
		tok := p.Peek(0)
		if p.MatchAnyNext(token.GT, token.LT, token.EQ_EQ, token.GT_EQ, token.LT_EQ, token.NOT_EQ) {
			expr = ast.NewBinaryExpr(expr, tok.Token, p.Addition())
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) Addition() ast.Expr {
	expr := p.Multiplication()

	for {
		tok := p.Peek(0)

		if p.MatchAnyNext(token.ADD, token.SUB) {
			expr = ast.NewBinaryExpr(expr, tok.Token, p.Multiplication())
		} else if p.MatchAnyNext(token.INC, token.DEC) {

			var op token.Token
			{
				if tok.Token == token.INC {
					op = token.ADD
				} else {
					op = token.SUB
				}
			}

			identTok := p.Peek(-2)

			if identTok.Token == token.IDENT {
				ident := identTok.Value
				expr = ast.NewBinaryExpr(ast.NewVarExpr(ident, nil), op, ast.NewInt64Expr(1))
				expr = ast.NewVarExpr(ident, expr)
			} else if identTok.Token == token.INT || identTok.Token == token.FLOAT {
				expr = ast.NewBinaryExpr(expr, op, ast.NewInt64Expr(1))
			}

		} else {
			break
		}
	}

	return expr
}

func (p *Parser) Multiplication() ast.Expr {
	expr := p.Unary()
	for {
		tok := p.Peek(0)
		if p.MatchAnyNext(token.MUL, token.QUO, token.REM) {
			expr = ast.NewBinaryExpr(expr, tok.Token, p.Unary())
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) Unary() ast.Expr {
	if p.MatchAndNext(token.SUB) {
		return ast.NewUnaryOpExpr('-', p.Postfix(), 0)
	}

	if p.MatchAndNext(token.NOT) {
		count := 1 // First NOT is already consumed
		{
			for p.MatchAndNext(token.NOT) {
				count++
			}
		}
		return ast.NewUnaryOpExpr('!', p.Postfix(), count)
	}

	return p.Postfix()
}

func (p *Parser) Postfix() ast.Expr {
	expr := p.Primary()
	
	for {
		if p.MatchAndNext(token.DOT) {
			// Точечная нотация: obj.property или obj.method()
			propTok := p.Next()
			if propTok.Token != token.IDENT {
				p.error("expected property name after '.'", propTok)
			}
			
			// Проверяем, это вызов метода или доступ к свойству
			if p.MatchAndNext(token.LPAREN) {
				// Это вызов метода: obj.method(args)
				var args []ast.Expr
				for !p.MatchAndNext(token.RPAREN) {
					args = append(args, p.Expression())
					p.MatchAndNext(token.COMMA)
				}
				expr = ast.NewMethodCallExpr(expr, propTok.Value, args)
			} else {
				// Это доступ к свойству: obj.property
				expr = ast.NewMemberExpr(expr, propTok.Value)
			}
		} else if p.MatchAndNext(token.LBRACK) {
			// Индексация: arr[index] или obj["key"]
			indexExpr := p.Expression()
			if !p.MatchAndNext(token.RBRACK) {
				p.error("expected ']'", p.Peek(0))
			}
			expr = ast.NewIndexExpr(expr, indexExpr)
		} else if p.MatchAndNext(token.LPAREN) {
			// Вызов функции только для VarExpr (переменных)
			if varExpr, ok := expr.(*ast.VarExpr); ok {
				var args []ast.Expr
				for !p.MatchAndNext(token.RPAREN) {
					args = append(args, p.Expression())
					p.MatchAndNext(token.COMMA)
				}
				expr = ast.NewFuncCallExpr(varExpr.Name, args)
			} else {
				p.error("cannot call non-function", p.Peek(-1))
			}
		} else {
			break
		}
	}
	
	return expr
}

func (p *Parser) Primary() ast.Expr {
	tok := p.Peek(0)

	switch tok.Token {

	case token.EOF:
		p.Next()
		return ast.NewLiteralString("")
	case token.INT:
		p.Next()
		value, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			p.error(fmt.Sprintf("invalid int literal: %s", tok.Value), p.Peek(0))
		}
		return ast.NewInt64Expr(value)

	case token.FLOAT:
		p.Next()
		value, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			p.error(fmt.Sprintf("invalid float literal: %s", tok.Value), p.Peek(0))
		}
		return ast.NewFloat64Expr(value)

	case token.STRING:
		raw := p.Peek(0).Value
		p.Next()
		return ast.NewLiteralString(raw)

	case token.INTERP_STRING:
		return p.InterpolatedString()

	case token.TRUE, token.FALSE:
		p.Next()
		b, _ := strconv.ParseBool(tok.Value)
		return ast.NewBoolExpr(b)

	case token.LPAREN:
		p.Next()
		expr := p.Expression()
		if !p.Match(token.RPAREN) {
			p.error("expected ')'", p.Peek(0))
		}
		p.Next()
		return expr

	case token.LBRACE:
		return p.ObjectLiteral()

	case token.LBRACK:
		return p.ArrayLiteral()

	case token.OK:
		p.Next()
		if !p.MatchAndNext(token.LPAREN) {
			p.error("expected '(' after Ok", p.Peek(0))
		}
		value := p.Expression()
		if !p.MatchAndNext(token.RPAREN) {
			p.error("expected ')' after Ok value", p.Peek(0))
		}
		return ast.NewOkExpr(value)

	case token.ERR:
		p.Next()
		if !p.MatchAndNext(token.LPAREN) {
			p.error("expected '(' after Err", p.Peek(0))
		}
		error := p.Expression()
		if !p.MatchAndNext(token.RPAREN) {
			p.error("expected ')' after Err value", p.Peek(0))
		}
		return ast.NewErrExpr(error)

	case token.IDENT:
		p.Next()
		return ast.NewVarExpr(tok.Value, nil)
	}

	p.error("unexpected token", tok)
	return nil
}

// ObjectLiteral парсит объектные литералы {key: value, key2: value2}
func (p *Parser) ObjectLiteral() ast.Expr {
	if !p.MatchAndNext(token.LBRACE) {
		p.error("expected '{'", p.Peek(0))
	}

	fields := make(map[string]ast.Expr)

	// Пустой объект {}
	if p.MatchAndNext(token.RBRACE) {
		return ast.NewObjectExpr(fields)
	}

	for {
		// Ожидаем идентификатор или строку как ключ
		keyTok := p.Next()
		var key string
		
		if keyTok.Token == token.IDENT {
			key = keyTok.Value
		} else if keyTok.Token == token.STRING {
			key = keyTok.Value
		} else {
			p.error("expected property name", keyTok)
		}

		// Ожидаем двоеточие
		if !p.MatchAndNext(token.COLON) {
			p.error("expected ':'", p.Peek(0))
		}

		// Значение
		value := p.Expression()
		fields[key] = value

		// Проверяем запятую или закрывающую скобку
		if p.MatchAndNext(token.RBRACE) {
			break
		}
		
		if !p.MatchAndNext(token.COMMA) {
			p.error("expected ',' or '}'", p.Peek(0))
		}

		// Если после запятой сразу закрывающая скобка - это нормально
		if p.MatchAndNext(token.RBRACE) {
			break
		}
	}

	return ast.NewObjectExpr(fields)
}

func (p *Parser) format() ast.Expr {
	raw := p.Peek(0).Value
	p.Next()

	var parts []ast.Expr
	var buf strings.Builder
	i := 0

	for i < len(raw) {
		if raw[i] == '{' {

			// Добавляем накопленный литерал
			if buf.Len() > 0 {
				parts = append(parts, ast.NewLiteralString(buf.String()))
				buf.Reset()
			}

			// Найти конец выражения
			j := i + 1
			for j < len(raw) && raw[j] != '}' {
				j++
			}
			if j >= len(raw) {
				p.error("unclosed { in string", p.Peek(0))
			}

			// Вырезаем выражение между { и }
			exprStr := raw[i+1 : j]
			subParser := NewParser(exprStr)
			exprs := subParser.Parse()

			if len(exprs) != 1 {
				p.error("invalid embedded expression in string", p.Peek(0))
			}

			parts = append(parts, exprs...)

			i = j + 1 // после }
			continue
		}

		buf.WriteByte(raw[i])
		i++
	}

	// Остаток как литерал
	if buf.Len() > 0 {
		parts = append(parts, ast.NewLiteralString(buf.String()))
	}

	return ast.NewStringFormatExpr(parts)
}

func (p *Parser) EnumStatement() ast.Expr {
	if !p.Match(token.IDENT) {
		p.error("expected enum name", p.Peek(0))
	}
	
	enumName := p.Next().Value
	
	if !p.MatchAndNext(token.LBRACE) {
		p.error("expected '{'", p.Peek(0))
	}
	
	var values []string
	
	for !p.MatchAndNext(token.RBRACE) {
		if !p.Match(token.IDENT) {
			p.error("expected enum value", p.Peek(0))
		}
		
		values = append(values, p.Next().Value)
		
		if p.MatchAndNext(token.RBRACE) {
			break
		}
		
		if !p.MatchAndNext(token.COMMA) {
			p.error("expected ',' or '}'", p.Peek(0))
		}
		
		if p.MatchAndNext(token.RBRACE) {
			break
		}
	}
	
	return ast.NewEnumExpr(enumName, values)
}

func (p *Parser) ArrayLiteral() ast.Expr {
	if !p.MatchAndNext(token.LBRACK) {
		p.error("expected '['", p.Peek(0))
	}
	
	var elements []ast.Expr
	
	if p.MatchAndNext(token.RBRACK) {
		return ast.NewArrayExpr(elements)
	}
	
	for {
		elements = append(elements, p.Expression())
		
		if p.MatchAndNext(token.RBRACK) {
			break
		}
		
		if !p.MatchAndNext(token.COMMA) {
			p.error("expected ',' or ']'", p.Peek(0))
		}
		
		if p.MatchAndNext(token.RBRACK) {
			break
		}
	}
	
	return ast.NewArrayExpr(elements)
}

func (p *Parser) InterpolatedString() ast.Expr {
	raw := p.Peek(0).Value
	p.Next()

	var parts []ast.Expr
	var buf strings.Builder
	i := 0

	for i < len(raw) {
		if i < len(raw)-1 && raw[i] == '$' && raw[i+1] == '{' {
			// Добавляем накопленный литерал
			if buf.Len() > 0 {
				parts = append(parts, ast.NewLiteralString(buf.String()))
				buf.Reset()
			}

			// Пропускаем ${
			i += 2

			// Найти конец выражения }
			j := i
			braceCount := 1
			for j < len(raw) && braceCount > 0 {
				if raw[j] == '{' {
					braceCount++
				} else if raw[j] == '}' {
					braceCount--
				}
				if braceCount > 0 {
					j++
				}
			}

			if braceCount > 0 {
				p.error("unclosed ${ in string", p.Peek(0))
			}

			// Вырезаем выражение между ${ и }
			exprStr := raw[i:j]
			if len(exprStr) == 0 {
				p.error("empty expression in string interpolation", p.Peek(0))
			}

			subParser := NewParser(exprStr)
			exprs := subParser.Parse()

			if len(exprs) != 1 {
				p.error("invalid embedded expression in string", p.Peek(0))
			}

			parts = append(parts, exprs[0])

			i = j + 1 // после }
			continue
		}

		buf.WriteByte(raw[i])
		i++
	}

	// Остаток как литерал
	if buf.Len() > 0 {
		parts = append(parts, ast.NewLiteralString(buf.String()))
	}

	// Если нет частей для интерполяции, возвращаем простую строку
	if len(parts) == 1 {
		if strExpr, ok := parts[0].(*ast.LiteralString); ok {
			return strExpr
		}
	}

	return ast.NewStringFormatExpr(parts)
}

// ImportStatement parses import statements:
// import "./module.foo"
// import { func1, var1 } from "./module.foo"  
// import * as ModuleName from "./module.foo"
func (p *Parser) ImportStatement() ast.Expr {
	// import "./path"
	if p.Match(token.STRING) {
		path := p.Peek(0).Value
		p.Next()
		return ast.NewImportExpr(path)
	}

	// import { item1, item2 } from "./path"
	if p.MatchAndNext(token.LBRACE) {
		var items []string
		
		for !p.Match(token.RBRACE) {
			if !p.Match(token.IDENT) {
				p.error("expected identifier in import list", p.Peek(0))
			}
			items = append(items, p.Peek(0).Value)
			p.Next()
			
			if !p.MatchAndNext(token.COMMA) {
				break
			}
		}
		
		if !p.MatchAndNext(token.RBRACE) {
			p.error("expected '}' after import list", p.Peek(0))
		}
		
		if !p.MatchAndNext(token.FROM) {
			p.error("expected 'from' after import list", p.Peek(0))
		}
		
		if !p.Match(token.STRING) {
			p.error("expected module path string", p.Peek(0))
		}
		
		path := p.Peek(0).Value
		p.Next()
		
		return ast.NewSelectiveImportExpr(path, items)
	}

	// import * as Name from "./path"
	if p.MatchAllNext(token.MUL, token.AS) {
		if !p.Match(token.IDENT) {
			p.error("expected identifier after 'as'", p.Peek(0))
		}
		
		alias := p.Peek(0).Value
		p.Next()
		
		if !p.MatchAndNext(token.FROM) {
			p.error("expected 'from' after alias", p.Peek(0))
		}
		
		if !p.Match(token.STRING) {
			p.error("expected module path string", p.Peek(0))
		}
		
		path := p.Peek(0).Value
		p.Next()
		
		return ast.NewAliasImportExpr(path, alias)
	}

	p.error("invalid import syntax", p.Peek(0))
	return nil
}

// ExportStatement parses export statements:
// export fn name() { }
// export let variable = value
// export enum Color { RED, GREEN, BLUE }
func (p *Parser) ExportStatement() ast.Expr {
	var declaration ast.Expr
	var name string
	
	// export fn name() { }
	if p.Match(token.FN) {
		declaration = p.FunctionStatement()
		// Extract function name from FuncStatment
		if funcStmt, ok := declaration.(*ast.FuncStatment); ok {
			name = funcStmt.Name()
		} else {
			p.error("invalid function declaration", p.Peek(0))
		}
	} else if p.MatchAll(token.LET, token.IDENT, token.EQ) {
		// export let variable = value
		p.NextN(2) // Skip LET and IDENT
		name = p.Peek(-1).Value
		p.Next() // Skip EQ
		declaration = ast.NewLetExpr(name, p.Expression())
	} else if p.MatchAll(token.CONST, token.IDENT, token.EQ) {
		// export const variable = value
		p.NextN(2) // Skip CONST and IDENT
		name = p.Peek(-1).Value 
		p.Next() // Skip EQ
		declaration = ast.NewConstExpr(name, p.Expression())
	} else if p.MatchAndNext(token.ENUM) {
		// export enum Name { }
		if !p.Match(token.IDENT) {
			p.error("expected enum name", p.Peek(0))
		}
		name = p.Peek(0).Value // Get enum name before parsing
		declaration = p.EnumStatement()
	} else {
		p.error("invalid export declaration", p.Peek(0))
		return nil
	}
	
	return ast.NewExportExpr(declaration, name)
}
