package value

import (
	"fmt"
	"strconv"
)

type Value struct {
	isConst  bool
	isReturn bool
	isYield  bool
	isBreak  bool
	data     any
}

type Number interface {
	int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func NewValue(data any, isConst ...bool) *Value {
	return &Value{
		data: data,
		isConst: func() bool {
			if len(isConst) > 0 {
				return isConst[0]
			}
			return false
		}(),
	}
}

func (n *Value) Bool() bool {
	switch v := n.Any().(type) {
	case bool:
		return v
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	case string:
		return v != ""
	default:
		return v != nil
	}
}

func (n *Value) IsBool() bool {
	switch n.Any().(type) {
	case bool:
		return true
	default:
		return false
	}
}

func (n *Value) IsNumber() bool {
	switch n.Any().(type) {
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func (n *Value) IsFloat64() bool {
	switch n.Any().(type) {
	case float64:
		return true
	default:
		return false
	}
}

func (n *Value) IsInt64() bool {
	switch n.Any().(type) {
	case int64:
		return true
	default:
		return false
	}
}

func (n *Value) IsString() bool {
	switch n.Any().(type) {
	case string:
		return true
	default:
		return false
	}
}

func (n *Value) String() string {
	switch v := n.Any().(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	}
	return ""
}

func (n *Value) Int() int {
	switch v := n.Any().(type) {
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case string:
		x, _ := strconv.ParseInt(v, 10, 64)
		return int(x)
	case bool:
		if v {
			return 1
		}
		return 0
	}
	return 0
}

func (n *Value) Int64() int64 {
	switch v := n.Any().(type) {
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case string:
		x, _ := strconv.ParseInt(v, 10, 64)
		return int64(x)
	case bool:
		if v {
			return 1
		}
		return 0
	}
	return 0
}

func (n *Value) Float64() float64 {
	switch v := n.Any().(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case string:
		x, _ := strconv.ParseFloat(v, 64)
		return float64(x)
	case bool:
		if v {
			return 1
		}
		return 0
	}
	return 0
}

func (n *Value) Float32() float32 {
	switch v := n.Any().(type) {
	case float64:
		return float32(v)
	case float32:
		return float32(v)
	case int8:
		return float32(v)
	case int16:
		return float32(v)
	case int32:
		return float32(v)
	case int64:
		return float32(v)
	case uint8:
		return float32(v)
	case uint16:
		return float32(v)
	case uint32:
		return float32(v)
	case uint64:
		return float32(v)
	case string:
		x, _ := strconv.ParseFloat(v, 32)
		return float32(x)
	case bool:
		if v {
			return 1
		}
		return 0
	}
	return 0
}

func (n *Value) Any() any {
	if n == nil {
		return nil
	}
	return n.data
}

func (n *Value) IsConst() bool {
	return n.isConst
}

func (n *Value) SetConst(constant bool) {
	n.isConst = constant
}


func (n *Value) SetYield(yield bool) {
	n.isYield = yield
}

func (n *Value) IsYield() bool {
	return n.isYield
}

func (n *Value) SetBreak(brk bool) {
	n.isBreak = brk
}

func (n *Value) IsBreak() bool {
	return n.isBreak
}

// Extension methods registry
var extensionMethods = make(map[string]map[string]interface{})

// RegisterExtensionMethod регистрирует метод расширения для типа
func RegisterExtensionMethod(typeName string, methodName string, method interface{}) {
	if extensionMethods[typeName] == nil {
		extensionMethods[typeName] = make(map[string]interface{})
	}
	extensionMethods[typeName][methodName] = method
}

// GetExtensionMethod возвращает метод расширения для типа
func GetExtensionMethod(typeName string, methodName string) (interface{}, bool) {
	if methods, ok := extensionMethods[typeName]; ok {
		if method, ok := methods[methodName]; ok {
			return method, true
		}
	}
	return nil, false
}

// GetValueTypeName возвращает имя типа для значения
func GetValueTypeName(v *Value) string {
	if v == nil || v.data == nil {
		return "nil"
	}
	
	switch v.data.(type) {
	case int64:
		return "int"
	case float64:
		return "float"
	case string:
		return "string"
	case bool:
		return "bool"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		// Проверяем на StructObject через рефлексию
		typeName := fmt.Sprintf("%T", v.data)
		if typeName == "*ast.StructObject" {
			return "struct"
		}
		return "unknown"
	}
}
