package ast

import (
	"foo_lang/value"
	"strings"
)

// UnionTypeExpr представляет union тип: string | number | null
type UnionTypeExpr struct {
	Types []string // Список типов в union
}

func NewUnionTypeExpr(types []string) *UnionTypeExpr {
	return &UnionTypeExpr{Types: types}
}

func (u *UnionTypeExpr) Eval() *value.Value {
	// Создаем специальный UnionTypeInfo
	unionInfo := NewUnionTypeInfo(u.Types)
	return value.NewValue(unionInfo)
}

// UnionTypeInfo хранит информацию о union типе
type UnionTypeInfo struct {
	TypeName string   // Имя для отображения, например "string | number | null"
	Types    []string // Список типов
}

func NewUnionTypeInfo(types []string) *UnionTypeInfo {
	return &UnionTypeInfo{
		TypeName: strings.Join(types, " | "),
		Types:    types,
	}
}

func (u *UnionTypeInfo) String() string {
	return u.TypeName
}

func (u *UnionTypeInfo) GetTypeName() string {
	return u.TypeName
}

// IsValidType проверяет, соответствует ли значение одному из типов union
func (u *UnionTypeInfo) IsValidType(val *value.Value) bool {
	valueType := value.GetValueTypeName(val)
	
	for _, unionType := range u.Types {
		switch unionType {
		case "int", "number":
			if valueType == "int" {
				return true
			}
		case "float":
			if valueType == "float" {
				return true
			}
		case "string":
			if valueType == "string" {
				return true
			}
		case "bool", "boolean":
			if valueType == "bool" {
				return true
			}
		case "null":
			if val.Any() == nil {
				return true
			}
		case "array":
			if valueType == "array" {
				return true
			}
		case "object":
			if valueType == "object" {
				return true
			}
		default:
			// Поддержка пользовательских типов
			if valueType == unionType {
				return true
			}
		}
	}
	
	return false
}