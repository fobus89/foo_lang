package ast

import (
	"fmt"
	"foo_lang/scope"
)

// TypedLetExpr представляет типизированную переменную
type TypedLetExpr struct {
	name     string
	varType  string // Ожидаемый тип переменной
	expr     Expr
}

func NewTypedLetExpr(name string, varType string, expr Expr) *TypedLetExpr {
	return &TypedLetExpr{
		name:    name,
		varType: varType,
		expr:    expr,
	}
}

func (n *TypedLetExpr) Eval() *Value {
	if scope.GlobalScope.Has(n.name) {
		panic("variable " + n.name + " is already defined")
	}

	val := n.expr.Eval()
	
	// Проверяем тип, если он указан
	if n.varType != "" {
		if err := validateVariableType(val, n.varType); err != nil {
			panic(fmt.Sprintf("variable '%s' type error: %s", n.name, err.Error()))
		}
	}

	scope.GlobalScope.Set(n.name, val)

	return nil
}

// validateVariableType проверяет соответствие значения ожидаемому типу переменной
func validateVariableType(value *Value, expectedType string) error {
	switch expectedType {
	case "int":
		if !value.IsInt64() {
			return fmt.Errorf("expected int, got %T", value.Any())
		}
		return nil
	case "string":
		if !value.IsString() {
			return fmt.Errorf("expected string, got %T", value.Any())
		}
		return nil
	case "float":
		if !value.IsFloat64() && !value.IsInt64() { // int может быть приведен к float
			return fmt.Errorf("expected float, got %T", value.Any())
		}
		return nil
	case "bool":
		if !value.IsBool() {
			return fmt.Errorf("expected bool, got %T", value.Any())
		}
		return nil
	default:
		return fmt.Errorf("unknown type constraint: %s", expectedType)
	}
}

// GetName возвращает имя переменной
func (n *TypedLetExpr) GetName() string {
	return n.name
}

// GetType возвращает тип переменной
func (n *TypedLetExpr) GetType() string {
	return n.varType
}

// GetExpr возвращает выражение
func (n *TypedLetExpr) GetExpr() Expr {
	return n.expr
}