package test

import (
	"foo_lang/ast"
	"foo_lang/parser"
	"foo_lang/scope"
	"strings"
	"testing"
)

func TestResultOk(t *testing.T) {
	InitTestEnvironment()

	const code = `let result = Ok(42)`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	resultVal, ok := val.Any().(*ast.ResultValue)
	if !ok {
		t.Errorf("expected ResultValue, got %T", val.Any())
		return
	}

	if !resultVal.IsOk() {
		t.Errorf("expected Ok result")
	}

	if resultVal.IsErr() {
		t.Errorf("expected not Err result")
	}

	unwrapped := resultVal.Unwrap()
	if unwrapped.Int64() != 42 {
		t.Errorf("expected 42, got %d", unwrapped.Int64())
	}
}

func TestResultErr(t *testing.T) {
	InitTestEnvironment()

	const code = `let result = Err("error message")`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	val, ok := scope.GlobalScope.Get("result")
	if !ok {
		t.Errorf("result not found")
		return
	}

	resultVal, ok := val.Any().(*ast.ResultValue)
	if !ok {
		t.Errorf("expected ResultValue, got %T", val.Any())
		return
	}

	if resultVal.IsOk() {
		t.Errorf("expected not Ok result")
	}

	if !resultVal.IsErr() {
		t.Errorf("expected Err result")
	}

	errorVal := resultVal.GetValue()
	if errorVal.String() != "error message" {
		t.Errorf("expected 'error message', got %s", errorVal.String())
	}
}

func TestResultMethods(t *testing.T) {
	InitTestEnvironment()

	const code = `
		let okResult = Ok(100)
		let errResult = Err("failed")
		
		let isOkTest = okResult.isOk()
		let isErrTest = errResult.isErr()
		let unwrapTest = okResult.unwrap()
		let unwrapOrTest1 = okResult.unwrapOr(0)
		let unwrapOrTest2 = errResult.unwrapOr(-1)
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	tests := []struct {
		name     string
		variable string
		expected interface{}
	}{
		{"isOk on Ok", "isOkTest", true},
		{"isErr on Err", "isErrTest", true},
		{"unwrap on Ok", "unwrapTest", int64(100)},
		{"unwrapOr on Ok", "unwrapOrTest1", int64(100)},
		{"unwrapOr on Err", "unwrapOrTest2", float64(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := scope.GlobalScope.Get(tt.variable)
			if !ok {
				t.Errorf("variable %s not found", tt.variable)
				return
			}

			result := val.Any()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestResultUnwrapPanic(t *testing.T) {
	InitTestEnvironment()

	const code = `
		let errResult = Err("something went wrong")
		errResult.unwrap()
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	
	// This should panic
	defer func() {
		if r := recover(); r != nil {
			errorStr := r.(string)
			if !strings.Contains(errorStr, "called unwrap on Err value") {
				t.Errorf("expected unwrap panic error, got: %s", errorStr)
			}
		} else {
			t.Errorf("expected panic due to unwrap on Err")
		}
	}()

	for _, expr := range exprs {
		expr.Eval()
	}
}

func TestResultInFunction(t *testing.T) {
	InitTestEnvironment()

	const code = `
		fn safeDivide(a, b) {
			if b == 0 {
				return Err("division by zero")
			}
			return Ok(a / b)
		}
		
		let result1 = safeDivide(10, 2)
		let result2 = safeDivide(10, 0)
		
		let value1 = result1.unwrap()
		let value2 = result2.unwrapOr(-1)
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test successful division
	value1Val, ok := scope.GlobalScope.Get("value1")
	if !ok {
		t.Errorf("value1 not found")
		return
	}

	if value1Val.Int64() != 5 {
		t.Errorf("expected 5, got %d", value1Val.Int64())
	}

	// Test division by zero with unwrapOr
	value2Val, ok := scope.GlobalScope.Get("value2")
	if !ok {
		t.Errorf("value2 not found")
		return
	}

	if value2Val.Int64() != -1 {
		t.Errorf("expected -1, got %d", value2Val.Int64())
	}
}

func TestResultChaining(t *testing.T) {
	InitTestEnvironment()

	const code = `
		fn parseInt(str) {
			if str == "42" {
				return Ok(42)
			}
			return Err("invalid number")
		}
		
		fn processNumber(numStr) {
			let parseResult = parseInt(numStr)
			if parseResult.isErr() {
				return parseResult
			}
			
			let num = parseResult.unwrap()
			if num > 40 {
				return Ok("big number: " + num)
			}
			return Ok("small number: " + num)
		}
		
		let result1 = processNumber("42")
		let result2 = processNumber("invalid")
		
		let final1 = result1.unwrapOr("no result")
		let final2 = result2.unwrapOr("no result")
	`
	
	exprs := parser.NewParser(code).ParseWithoutScopeInit()
	for _, expr := range exprs {
		expr.Eval()
	}

	// Test successful processing
	final1Val, ok := scope.GlobalScope.Get("final1")
	if !ok {
		t.Errorf("final1 not found")
		return
	}

	if final1Val.String() != "big number: 42" {
		t.Errorf("expected 'big number: 42', got %s", final1Val.String())
	}

	// Test error processing
	final2Val, ok := scope.GlobalScope.Get("final2")
	if !ok {
		t.Errorf("final2 not found")
		return
	}

	if final2Val.String() != "no result" {
		t.Errorf("expected 'no result', got %s", final2Val.String())
	}
}