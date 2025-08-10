package ast

import (
	"fmt"
	"foo_lang/token"
	"math"
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
		return NewValue(left.Float64() - right.Float64())
	case token.MUL:
		return NewValue(left.Float64() * right.Float64())
	case token.QUO:
		return NewValue(left.Float64() / right.Float64())
	case token.REM:
		return NewValue(math.Mod(left.Float64(), right.Float64()))

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
		return NewValue(left.Float64() + right.Float64())
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
