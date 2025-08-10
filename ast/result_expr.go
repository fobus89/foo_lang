package ast

// ResultValue представляет Result<T, E> тип
type ResultValue struct {
	isOk  bool
	value *Value
	error *Value
}

// NewResultOk создает Ok(value)
func NewResultOk(value *Value) *ResultValue {
	return &ResultValue{
		isOk:  true,
		value: value,
		error: nil,
	}
}

// NewResultErr создает Err(error)
func NewResultErr(error *Value) *ResultValue {
	return &ResultValue{
		isOk:  false,
		value: nil,
		error: error,
	}
}

// IsOk возвращает true если Result содержит Ok
func (r *ResultValue) IsOk() bool {
	return r.isOk
}

// IsErr возвращает true если Result содержит Err
func (r *ResultValue) IsErr() bool {
	return !r.isOk
}

// Unwrap возвращает значение или panic
func (r *ResultValue) Unwrap() *Value {
	if r.isOk {
		return r.value
	}
	panic("called unwrap on Err value: " + FormatValue(r.error.Any()))
}

// UnwrapOr возвращает значение или default
func (r *ResultValue) UnwrapOr(defaultValue *Value) *Value {
	if r.isOk {
		return r.value
	}
	return defaultValue
}

// GetValue возвращает внутреннее значение (Ok или Err)
func (r *ResultValue) GetValue() *Value {
	if r.isOk {
		return r.value
	}
	return r.error
}

// OkExpr представляет Ok(value)
type OkExpr struct {
	Value Expr
}

func NewOkExpr(value Expr) *OkExpr {
	return &OkExpr{Value: value}
}

func (o *OkExpr) Eval() *Value {
	val := o.Value.Eval()
	result := NewResultOk(val)
	return NewValue(result)
}

// ErrExpr представляет Err(error)
type ErrExpr struct {
	Error Expr
}

func NewErrExpr(error Expr) *ErrExpr {
	return &ErrExpr{Error: error}
}

func (e *ErrExpr) Eval() *Value {
	err := e.Error.Eval()
	result := NewResultErr(err)
	return NewValue(result)
}