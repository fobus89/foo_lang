package ast

import (
	"foo_lang/scope"
	"foo_lang/value"
)

// TypeAliasExpr представляет определение псевдонима типа: type UserId = int
type TypeAliasExpr struct {
	AliasName string // Имя нового типа (например, "UserId")
	BaseType  string // Базовый тип (например, "int")
}

func NewTypeAliasExpr(aliasName, baseType string) *TypeAliasExpr {
	return &TypeAliasExpr{
		AliasName: aliasName,
		BaseType:  baseType,
	}
}

func (t *TypeAliasExpr) Eval() *value.Value {
	// Создаем псевдоним типа в глобальной области видимости
	
	// Получаем информацию о базовом типе
	var baseTypeInfo *TypeInfo
	
	// Сначала проверяем примитивные типы
	switch t.BaseType {
	case "int":
		baseTypeInfo = NewPrimitiveTypeInfo("int")
	case "string":
		baseTypeInfo = NewPrimitiveTypeInfo("string")
	case "float":
		baseTypeInfo = NewPrimitiveTypeInfo("float")
	case "bool":
		baseTypeInfo = NewPrimitiveTypeInfo("bool")
	case "array":
		baseTypeInfo = NewPrimitiveTypeInfo("array")
	case "object":
		baseTypeInfo = NewPrimitiveTypeInfo("object")
	default:
		// Проверяем другие псевдонимы (цепочка псевдонимов)
		if aliasValue, exists := scope.GlobalScope.Get(t.BaseType); exists {
			if existingAlias, ok := aliasValue.Any().(*TypeAliasInfo); ok {
				// Ссылка на псевдоним - используем его базовый тип
				baseTypeInfo = existingAlias.BaseTypeInfo
			} else if baseInfo, ok := aliasValue.Any().(*TypeInfo); ok {
				baseTypeInfo = baseInfo
			} else {
				panic("Invalid alias value for type: " + t.BaseType)
			}
		} else if typeInfoValue, exists := scope.GlobalScope.Get(t.BaseType + "__TypeInfo"); exists {
			if baseInfo, ok := typeInfoValue.Any().(*TypeInfo); ok {
				baseTypeInfo = baseInfo
			} else {
				panic("Invalid type info for type: " + t.BaseType)
			}
		} else {
			panic("Unknown base type: " + t.BaseType)
		}
	}
	
	// Создаем TypeAliasInfo
	aliasInfo := NewTypeAliasInfo(t.AliasName, baseTypeInfo)
	
	// Сохраняем псевдоним в scope
	scope.GlobalScope.Set(t.AliasName, value.NewValue(aliasInfo))
	scope.GlobalScope.Set(t.AliasName+"__TypeInfo", value.NewValue(aliasInfo))
	
	return value.NewString("type alias '" + t.AliasName + "' = '" + t.BaseType + "' defined")
}

// TypeAliasInfo хранит информацию о псевдониме типа
type TypeAliasInfo struct {
	AliasName    string    // Имя псевдонима (например, "UserId")
	BaseTypeInfo *TypeInfo // Информация о базовом типе
}

func NewTypeAliasInfo(aliasName string, baseTypeInfo *TypeInfo) *TypeAliasInfo {
	return &TypeAliasInfo{
		AliasName:    aliasName,
		BaseTypeInfo: baseTypeInfo,
	}
}

func (t *TypeAliasInfo) String() string {
	return "type " + t.AliasName + " = " + t.BaseTypeInfo.String()
}

func (t *TypeAliasInfo) GetTypeName() string {
	return t.AliasName
}

// IsValidType проверяет валидность значения для псевдонима типа
func (t *TypeAliasInfo) IsValidType(val *value.Value) bool {
	// Псевдоним считается валидным, если значение соответствует базовому типу
	baseTypeName := t.BaseTypeInfo.Name
	valueType := value.GetValueTypeName(val)
	
	switch baseTypeName {
	case "int":
		return valueType == "int"
	case "string":
		return valueType == "string"
	case "float":
		return valueType == "float"
	case "bool":
		return valueType == "bool"
	case "array":
		return valueType == "array"
	case "object":
		return valueType == "object"
	default:
		// Для пользовательских типов сравниваем имена типов
		return valueType == baseTypeName
	}
}