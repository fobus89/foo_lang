package ast

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
	"strings"
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
func validateVariableType(val *Value, expectedType string) error {
	switch expectedType {
	case "int":
		if !val.IsInt64() {
			return fmt.Errorf("expected int, got %T", val.Any())
		}
		return nil
	case "string":
		if !val.IsString() {
			return fmt.Errorf("expected string, got %T", val.Any())
		}
		return nil
	case "float":
		if !val.IsFloat64() && !val.IsInt64() { // int может быть приведен к float
			return fmt.Errorf("expected float, got %T", val.Any())
		}
		return nil
	case "bool":
		if !val.IsBool() {
			return fmt.Errorf("expected bool, got %T", val.Any())
		}
		return nil
	case "array":
		if _, ok := val.Any().([]interface{}); !ok {
			return fmt.Errorf("expected array, got %T", val.Any())
		}
		return nil
	case "object":
		if _, ok := val.Any().(map[string]*value.Value); !ok {
			return fmt.Errorf("expected object, got %T", val.Any())
		}
		return nil
	default:
		// Проверяем Tuple типы (начинаются с '(' и заканчиваются ')')
		if strings.HasPrefix(expectedType, "(") && strings.HasSuffix(expectedType, ")") {
			return validateTupleType(val, expectedType)
		}
		
		// Проверяем Union типы (содержат символ |)
		if strings.Contains(expectedType, "|") {
			unionTypes := strings.Split(expectedType, "|")
			for i := range unionTypes {
				unionTypes[i] = strings.TrimSpace(unionTypes[i])
			}
			
			// Проверяем, соответствует ли значение одному из Union типов
			for _, unionType := range unionTypes {
				if unionType == "null" && val.Any() == nil {
					return nil
				}
				if err := validateVariableType(val, unionType); err == nil {
					return nil // Найден подходящий тип
				}
			}
			
			return fmt.Errorf("expected one of [%s], got %s", 
				strings.Join(unionTypes, " | "), value.GetValueTypeName(val))
		}
		
		// Проверяем псевдонимы типов
		if aliasValue, exists := scope.GlobalScope.Get(expectedType); exists {
			if aliasInfo, ok := aliasValue.Any().(*TypeAliasInfo); ok {
				// Для псевдонимов проверяем соответствие базовому типу
				return validateVariableType(val, aliasInfo.BaseTypeInfo.Name)
			}
		}
		
		// Проверяем пользовательские типы  
		if _, exists := scope.GlobalScope.Get(expectedType + "__TypeInfo"); exists {
			// Для пользовательских типов пока возвращаем nil (считаем валидными)
			return nil
		}
		
		return fmt.Errorf("unknown type constraint: %s", expectedType)
	}
}

// validateTupleType проверяет соответствие значения Tuple типу
func validateTupleType(val *Value, expectedType string) error {
	// Парсим Tuple тип: "(string,int,float)" -> ["string", "int", "float"]
	tupleContent := strings.TrimPrefix(expectedType, "(")
	tupleContent = strings.TrimSuffix(tupleContent, ")")
	expectedTypes := strings.Split(tupleContent, ",")
	
	// Очищаем пробелы
	for i := range expectedTypes {
		expectedTypes[i] = strings.TrimSpace(expectedTypes[i])
	}
	
	// Проверяем, что значение является массивом
	arr, ok := val.Any().([]interface{})
	if !ok {
		return fmt.Errorf("expected tuple (array), got %s", value.GetValueTypeName(val))
	}
	
	// Проверяем количество элементов
	if len(arr) != len(expectedTypes) {
		return fmt.Errorf("tuple length mismatch: expected %d elements, got %d", 
			len(expectedTypes), len(arr))
	}
	
	// Проверяем каждый элемент Tuple
	for i, element := range arr {
		elementVal := value.NewValue(element)
		expectedElementType := expectedTypes[i]
		
		if err := validateVariableType(elementVal, expectedElementType); err != nil {
			return fmt.Errorf("tuple element %d: %s", i, err.Error())
		}
	}
	
	return nil
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