package value

import (
	"fmt"
)

// FromInterface создает Value из interface{}
func FromInterface(data interface{}) *Value {
	return NewValue(data)
}

// Арифметические операции
func Add(a, b *Value) *Value {
	switch aVal := a.Any().(type) {
	case int64:
		if bVal, ok := b.Any().(int64); ok {
			return NewInt64(aVal + bVal)
		}
		if bVal, ok := b.Any().(float64); ok {
			return NewFloat64(float64(aVal) + bVal)
		}
	case float64:
		if bVal, ok := b.Any().(float64); ok {
			return NewFloat64(aVal + bVal)
		}
		if bVal, ok := b.Any().(int64); ok {
			return NewFloat64(aVal + float64(bVal))
		}
	case string:
		if bVal, ok := b.Any().(string); ok {
			return NewString(aVal + bVal)
		}
	}
	return NewString(fmt.Sprintf("Error: cannot add %T and %T", a.Any(), b.Any()))
}

func Subtract(a, b *Value) *Value {
	switch aVal := a.Any().(type) {
	case int64:
		if bVal, ok := b.Any().(int64); ok {
			return NewInt64(aVal - bVal)
		}
		if bVal, ok := b.Any().(float64); ok {
			return NewFloat64(float64(aVal) - bVal)
		}
	case float64:
		if bVal, ok := b.Any().(float64); ok {
			return NewFloat64(aVal - bVal)
		}
		if bVal, ok := b.Any().(int64); ok {
			return NewFloat64(aVal - float64(bVal))
		}
	}
	return NewString(fmt.Sprintf("Error: cannot subtract %T and %T", a.Any(), b.Any()))
}

func Multiply(a, b *Value) *Value {
	switch aVal := a.Any().(type) {
	case int64:
		if bVal, ok := b.Any().(int64); ok {
			return NewInt64(aVal * bVal)
		}
		if bVal, ok := b.Any().(float64); ok {
			return NewFloat64(float64(aVal) * bVal)
		}
	case float64:
		if bVal, ok := b.Any().(float64); ok {
			return NewFloat64(aVal * bVal)
		}
		if bVal, ok := b.Any().(int64); ok {
			return NewFloat64(aVal * float64(bVal))
		}
	}
	return NewString(fmt.Sprintf("Error: cannot multiply %T and %T", a.Any(), b.Any()))
}

func Divide(a, b *Value) *Value {
	switch aVal := a.Any().(type) {
	case int64:
		if bVal, ok := b.Any().(int64); ok {
			if bVal == 0 {
				return NewString("Error: division by zero")
			}
			return NewFloat64(float64(aVal) / float64(bVal))
		}
		if bVal, ok := b.Any().(float64); ok {
			if bVal == 0 {
				return NewString("Error: division by zero")
			}
			return NewFloat64(float64(aVal) / bVal)
		}
	case float64:
		if bVal, ok := b.Any().(float64); ok {
			if bVal == 0 {
				return NewString("Error: division by zero")
			}
			return NewFloat64(aVal / bVal)
		}
		if bVal, ok := b.Any().(int64); ok {
			if bVal == 0 {
				return NewString("Error: division by zero")
			}
			return NewFloat64(aVal / float64(bVal))
		}
	}
	return NewString(fmt.Sprintf("Error: cannot divide %T and %T", a.Any(), b.Any()))
}

func Modulo(a, b *Value) *Value {
	switch aVal := a.Any().(type) {
	case int64:
		if bVal, ok := b.Any().(int64); ok {
			if bVal == 0 {
				return NewString("Error: modulo by zero")
			}
			return NewInt64(aVal % bVal)
		}
	}
	return NewString(fmt.Sprintf("Error: cannot modulo %T and %T", a.Any(), b.Any()))
}

func Negate(a *Value) *Value {
	switch aVal := a.Any().(type) {
	case int64:
		return NewInt64(-aVal)
	case float64:
		return NewFloat64(-aVal)
	}
	return NewString(fmt.Sprintf("Error: cannot negate %T", a.Any()))
}

// Логические операции
func Not(a *Value) *Value {
	return NewBool(!a.IsTruthy())
}

func And(a, b *Value) *Value {
	if !a.IsTruthy() {
		return a
	}
	return b
}

func Or(a, b *Value) *Value {
	if a.IsTruthy() {
		return a
	}
	return b
}

// Операции сравнения
func Equal(a, b *Value) *Value {
	return NewBool(isEqual(a.Any(), b.Any()))
}

func NotEqual(a, b *Value) *Value {
	return NewBool(!isEqual(a.Any(), b.Any()))
}

func Greater(a, b *Value) *Value {
	result, err := compare(a.Any(), b.Any())
	if err != nil {
		return NewString(err.Error())
	}
	return NewBool(result > 0)
}

func GreaterEqual(a, b *Value) *Value {
	result, err := compare(a.Any(), b.Any())
	if err != nil {
		return NewString(err.Error())
	}
	return NewBool(result >= 0)
}

func Less(a, b *Value) *Value {
	result, err := compare(a.Any(), b.Any())
	if err != nil {
		return NewString(err.Error())
	}
	return NewBool(result < 0)
}

func LessEqual(a, b *Value) *Value {
	result, err := compare(a.Any(), b.Any())
	if err != nil {
		return NewString(err.Error())
	}
	return NewBool(result <= 0)
}

// Индексация
func Index(obj, index *Value) *Value {
	switch objVal := obj.Any().(type) {
	case []*Value:
		if idxVal, ok := index.Any().(int64); ok {
			if idxVal < 0 || idxVal >= int64(len(objVal)) {
				return NewString("Error: array index out of bounds")
			}
			return objVal[idxVal]
		}
	case map[string]*Value:
		if idxVal, ok := index.Any().(string); ok {
			if val, exists := objVal[idxVal]; exists {
				return val
			}
			return NewNil()
		}
	case string:
		if idxVal, ok := index.Any().(int64); ok {
			if idxVal < 0 || idxVal >= int64(len(objVal)) {
				return NewString("Error: string index out of bounds")
			}
			return NewString(string(objVal[idxVal]))
		}
	}
	return NewString(fmt.Sprintf("Error: cannot index %T with %T", obj.Any(), index.Any()))
}

// Вспомогательные функции
func isEqual(a, b interface{}) bool {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal == bVal
		}
		if bVal, ok := b.(float64); ok {
			return float64(aVal) == bVal
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal == bVal
		}
		if bVal, ok := b.(int64); ok {
			return aVal == float64(bVal)
		}
	case string:
		if bVal, ok := b.(string); ok {
			return aVal == bVal
		}
	case bool:
		if bVal, ok := b.(bool); ok {
			return aVal == bVal
		}
	case nil:
		return b == nil
	}
	return false
}

func compare(a, b interface{}) (int, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			if aVal < bVal {
				return -1, nil
			} else if aVal > bVal {
				return 1, nil
			}
			return 0, nil
		}
		if bVal, ok := b.(float64); ok {
			aFloat := float64(aVal)
			if aFloat < bVal {
				return -1, nil
			} else if aFloat > bVal {
				return 1, nil
			}
			return 0, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			if aVal < bVal {
				return -1, nil
			} else if aVal > bVal {
				return 1, nil
			}
			return 0, nil
		}
		if bVal, ok := b.(int64); ok {
			bFloat := float64(bVal)
			if aVal < bFloat {
				return -1, nil
			} else if aVal > bFloat {
				return 1, nil
			}
			return 0, nil
		}
	case string:
		if bVal, ok := b.(string); ok {
			if aVal < bVal {
				return -1, nil
			} else if aVal > bVal {
				return 1, nil
			}
			return 0, nil
		}
	}
	return 0, fmt.Errorf("Error: cannot compare %T and %T", a, b)
}

// Конструкторы для Value
func NewInt64(value int64) *Value {
	return NewValue(value)
}

func NewFloat64(value float64) *Value {
	return NewValue(value)
}

func NewString(value string) *Value {
	return NewValue(value)
}

func NewBool(value bool) *Value {
	return NewValue(value)
}

func NewNil() *Value {
	return NewValue(nil)
}

func NewArray(elements []*Value) *Value {
	return NewValue(elements)
}

func NewChannelValue(ch *Channel) *Value {
	return NewValue(ch)
}

// IsTruthy определяет, является ли значение истинным
func (v *Value) IsTruthy() bool {
	switch val := v.Any().(type) {
	case nil:
		return false
	case bool:
		return val
	case int64:
		return val != 0
	case float64:
		return val != 0.0
	case string:
		return val != ""
	case []*Value:
		return len(val) > 0
	case map[string]*Value:
		return len(val) > 0
	default:
		return true
	}
}

// IsReturn проверяет, является ли значение return-значением
func (v *Value) IsReturn() bool {
	return v.isReturn
}

// SetReturn устанавливает флаг return
func (v *Value) SetReturn(isReturn bool) {
	v.isReturn = isReturn
}