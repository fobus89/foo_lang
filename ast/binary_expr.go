package ast

import (
	"fmt"
	"foo_lang/token"
)

type BinaryExpr struct {
	Left  Expr
	Op    token.Token
	Right Expr
}

func NewBinaryExpr(left Expr, op token.Token, right Expr) *BinaryExpr {
	return &BinaryExpr{Left: left, Op: op, Right: right}
}

func (b *BinaryExpr) Eval() *Value {

	left := b.Left.Eval()
	right := b.Right.Eval()

	switch b.Op {

	case token.ADD:
		return b.add(left, right)
	case token.SUB:
		return b.subtract(left, right)
	case token.MUL:
		return b.multiply(left, right)
	case token.QUO:
		return b.divide(left, right)
	case token.REM:
		return b.modulo(left, right)

	case token.GT:
		return NewValue(left.Float64() > right.Float64())
	case token.LT:
		return NewValue(left.Float64() < right.Float64())
	case token.NOT_EQ:
		return NewValue(left.Float64() != right.Float64())
	case token.LT_EQ:
		return NewValue(left.Float64() <= right.Float64())
	case token.GT_EQ:
		return NewValue(left.Float64() >= right.Float64())
	case token.EQ_EQ:
		return NewValue(left.Float64() == right.Float64())

	case token.AND:
		return NewValue(left.Int64() & right.Int64())
	case token.OR:
		return NewValue(left.Int64() | right.Int64())
	case token.XOR:
		return NewValue(left.Int64() ^ right.Int64())
	case token.AND_NOT:
		return NewValue(left.Int64() &^ right.Int64())
	case token.LT_LT:
		return NewValue(left.Int64() << right.Int64())
	case token.GT_GT:
		return NewValue(left.Int64() >> right.Int64())
	case token.AND_AND:
		return b.AndAnd(left, right)
	case token.OR_OR:
		return b.OrOr(left, right)

	}

	panic(fmt.Sprintf("unknown token: %s", b.Op))
}

func (b *BinaryExpr) add(left, right *Value) *Value {

	if left.IsNumber() && right.IsNumber() {
		// Сохраняем типы: int + int = int, иначе float
		if left.IsInt64() && right.IsInt64() {
			return NewValue(left.Int64() + right.Int64())
		} else {
			return NewValue(left.Float64() + right.Float64())
		}
	}

	if left.IsString() && right.IsString() {
		return NewValue(left.String() + right.String())
	}

	if left.IsString() && right.IsNumber() {
		return NewValue(left.String() + right.String())
	}

	if left.IsNumber() && right.IsString() {
		return NewValue(left.String() + right.String())
	}

	if left.IsBool() && right.IsBool() {
		return NewValue(left.Int64() + right.Int64())
	}

	if left.IsBool() && right.IsNumber() {
		return NewValue(left.Int64() + right.Int64())
	}

	if left.IsNumber() && right.IsBool() {
		return NewValue(left.Int64() + right.Int64())
	}

	if left.IsBool() && right.IsString() {
		return NewValue(left.String() + right.String())
	}

	if left.IsString() && right.IsBool() {
		return NewValue(left.String() + right.String())
	}

	// Конкатенация строки с любым другим типом (включая массивы)
	if left.IsString() {
		return NewValue(left.String() + FormatValue(right.Any()))
	}

	if right.IsString() {
		return NewValue(FormatValue(left.Any()) + right.String())
	}

	return nil
}

func (b *BinaryExpr) OrOr(left, right *Value) *Value {
	if left.Bool() {
		return left
	}
	return right
}

func (b *BinaryExpr) AndAnd(left, right *Value) *Value {
	if left.Bool() {
		return right
	}
	return left
}

// Арифметические операции с сохранением типов

func (b *BinaryExpr) subtract(left, right *Value) *Value {
	if left.IsNumber() && right.IsNumber() {
		// int - int = int, иначе float
		if left.IsInt64() && right.IsInt64() {
			return NewValue(left.Int64() - right.Int64())
		} else {
			return NewValue(left.Float64() - right.Float64())
		}
	}
	panic("subtract operation requires numeric operands")
}

func (b *BinaryExpr) multiply(left, right *Value) *Value {
	if left.IsNumber() && right.IsNumber() {
		// int * int = int, иначе float
		if left.IsInt64() && right.IsInt64() {
			return NewValue(left.Int64() * right.Int64())
		} else {
			return NewValue(left.Float64() * right.Float64())
		}
	}
	panic("multiply operation requires numeric operands")
}

func (b *BinaryExpr) divide(left, right *Value) *Value {
	if left.IsNumber() && right.IsNumber() {
		// Деление всегда возвращает float для точности
		return NewValue(left.Float64() / right.Float64())
	}
	panic("divide operation requires numeric operands")
}

func (b *BinaryExpr) modulo(left, right *Value) *Value {
	if left.IsNumber() && right.IsNumber() {
		// int % int = int, float % float = float
		if left.IsInt64() && right.IsInt64() {
			return NewValue(left.Int64() % right.Int64())
		} else {
			// Нужен import math для Mod
			return NewValue(left.Float64() - right.Float64()*float64(int64(left.Float64()/right.Float64())))
		}
	}
	panic("modulo operation requires numeric operands")
}
