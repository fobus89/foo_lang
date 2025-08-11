package ast

import (
	"foo_lang/scope"
	"foo_lang/value"
	"fmt"
	"strings"
)

// TypeInfo представляет информацию о типе
type TypeInfo struct {
	Kind   string                 // "struct", "fn", "enum", "primitive"
	Name   string                 // имя типа
	Fields map[string]*TypeInfo   // поля для struct
	Params []*TypeInfo            // параметры для fn
	Return *TypeInfo              // возвращаемый тип для fn
	Values []string               // значения для enum
	Data   interface{}            // дополнительная информация
}

// GetProperty возвращает свойство TypeInfo как объект
func (ti *TypeInfo) GetProperty(name string) *value.Value {
	switch name {
	case "Kind":
		return value.NewValue(ti.Kind)
	case "Name":
		return value.NewValue(ti.Name)
	case "String":
		return value.NewValue(&TypeStringMethod{typeInfo: ti})
	case "GetFieldNames":
		return value.NewValue(&GetFieldNamesMethod{typeInfo: ti})
	case "GetFieldType":
		return value.NewValue(&GetFieldTypeMethod{typeInfo: ti})
	case "HasField":
		return value.NewValue(&HasFieldMethod{typeInfo: ti})
	default:
		return nil
	}
}

// TypeStringMethod представляет метод String для TypeInfo
type TypeStringMethod struct {
	typeInfo *TypeInfo
}

func (m *TypeStringMethod) Call(args []*value.Value) *value.Value {
	return value.NewValue(m.typeInfo.String())
}

// GetFieldNamesMethod представляет метод GetFieldNames для TypeInfo
type GetFieldNamesMethod struct {
	typeInfo *TypeInfo
}

func (m *GetFieldNamesMethod) Call(args []*value.Value) *value.Value {
	names := m.typeInfo.GetFieldNames()
	values := make([]*value.Value, len(names))
	for i, name := range names {
		values[i] = value.NewValue(name)
	}
	return value.NewValue(values)
}

// GetFieldTypeMethod представляет метод GetFieldType для TypeInfo
type GetFieldTypeMethod struct {
	typeInfo *TypeInfo
}

func (m *GetFieldTypeMethod) Call(args []*value.Value) *value.Value {
	if len(args) != 1 {
		panic("GetFieldType expects 1 argument")
	}
	fieldName := args[0].String()
	fieldType := m.typeInfo.GetFieldType(fieldName)
	if fieldType == nil {
		return value.NewValue(nil)
	}
	return value.NewValue(fieldType)
}

// HasFieldMethod представляет метод HasField для TypeInfo
type HasFieldMethod struct {
	typeInfo *TypeInfo
}

func (m *HasFieldMethod) Call(args []*value.Value) *value.Value {
	if len(args) != 1 {
		panic("HasField expects 1 argument")
	}
	fieldName := args[0].String()
	return value.NewValue(m.typeInfo.HasField(fieldName))
}

// NewStructTypeInfo создает информацию о структуре
func NewStructTypeInfo(name string, fields map[string]*TypeInfo) *TypeInfo {
	return &TypeInfo{
		Kind:   "struct",
		Name:   name,
		Fields: fields,
	}
}

// NewFunctionTypeInfo создает информацию о функции
func NewFunctionTypeInfo(name string, params []*TypeInfo, returnType *TypeInfo) *TypeInfo {
	return &TypeInfo{
		Kind:   "function",
		Name:   name,
		Params: params,
		Return: returnType,
	}
}

// NewEnumTypeInfo создает информацию об enum
func NewEnumTypeInfo(name string, values []string) *TypeInfo {
	return &TypeInfo{
		Kind:   "enum",
		Name:   name,
		Values: values,
	}
}

// NewPrimitiveTypeInfo создает информацию о примитивном типе
func NewPrimitiveTypeInfo(name string) *TypeInfo {
	return &TypeInfo{
		Kind: "primitive",
		Name: name,
	}
}

// GetFieldNames возвращает имена полей структуры
func (ti *TypeInfo) GetFieldNames() []string {
	if ti.Kind != "struct" {
		return []string{}
	}
	
	names := make([]string, 0, len(ti.Fields))
	for name := range ti.Fields {
		names = append(names, name)
	}
	return names
}

// GetFieldType возвращает тип поля
func (ti *TypeInfo) GetFieldType(fieldName string) *TypeInfo {
	if ti.Kind != "struct" {
		return nil
	}
	return ti.Fields[fieldName]
}

// HasField проверяет наличие поля
func (ti *TypeInfo) HasField(fieldName string) bool {
	if ti.Kind != "struct" {
		return false
	}
	_, exists := ti.Fields[fieldName]
	return exists
}

// String возвращает строковое представление типа
func (ti *TypeInfo) String() string {
	switch ti.Kind {
	case "struct":
		var fields []string
		for name, fieldType := range ti.Fields {
			fields = append(fields, fmt.Sprintf("%s: %s", name, fieldType.String()))
		}
		return fmt.Sprintf("struct %s { %s }", ti.Name, strings.Join(fields, ", "))
	case "function":
		var params []string
		for _, param := range ti.Params {
			params = append(params, param.String())
		}
		returnStr := "void"
		if ti.Return != nil {
			returnStr = ti.Return.String()
		}
		return fmt.Sprintf("fn %s(%s) -> %s", ti.Name, strings.Join(params, ", "), returnStr)
	case "enum":
		return fmt.Sprintf("enum %s { %s }", ti.Name, strings.Join(ti.Values, ", "))
	case "primitive":
		return ti.Name
	default:
		return "unknown"
	}
}

// TypeofExpr представляет выражение typeof для получения типа
type TypeofExpr struct {
	Expr Expr
}

func NewTypeofExpr(expr Expr) *TypeofExpr {
	return &TypeofExpr{Expr: expr}
}

func (t *TypeofExpr) Eval() *value.Value {
	// Анализируем тип выражения
	val := t.Expr.Eval()
	
	var typeInfo *TypeInfo
	switch val.Any().(type) {
	case int64:
		typeInfo = NewPrimitiveTypeInfo("int")
	case float64:
		typeInfo = NewPrimitiveTypeInfo("float")
	case string:
		typeInfo = NewPrimitiveTypeInfo("string")
	case bool:
		typeInfo = NewPrimitiveTypeInfo("bool")
	case map[string]*value.Value:
		// Объект - анализируем его структуру
		obj := val.Any().(map[string]*value.Value)
		fields := make(map[string]*TypeInfo)
		for name, fieldVal := range obj {
			switch fieldVal.Any().(type) {
			case int64:
				fields[name] = NewPrimitiveTypeInfo("int")
			case float64:
				fields[name] = NewPrimitiveTypeInfo("float")
			case string:
				fields[name] = NewPrimitiveTypeInfo("string")
			case bool:
				fields[name] = NewPrimitiveTypeInfo("bool")
			default:
				fields[name] = NewPrimitiveTypeInfo("any")
			}
		}
		typeInfo = NewStructTypeInfo("object", fields)
	case []*value.Value:
		typeInfo = NewPrimitiveTypeInfo("array")
	default:
		typeInfo = NewPrimitiveTypeInfo("unknown")
	}
	
	return value.NewValue(typeInfo)
}

// StructDefExpr представляет определение структуры
type StructDefExpr struct {
	Name   string
	Fields map[string]Expr // имя поля -> тип или значение по умолчанию
}

func NewStructDefExpr(name string, fields map[string]Expr) *StructDefExpr {
	return &StructDefExpr{
		Name:   name,
		Fields: fields,
	}
}

func (s *StructDefExpr) Eval() *value.Value {
	// Создаем информацию о типе структуры
	fieldTypes := make(map[string]*TypeInfo)
	for name, fieldExpr := range s.Fields {
		// Вычисляем тип поля
		fieldTypeValue := fieldExpr.Eval()
		if typeInfo, ok := fieldTypeValue.Any().(*TypeInfo); ok {
			fieldTypes[name] = typeInfo
		} else {
			// Для простоты, определяем тип как any если не можем определить
			fieldTypes[name] = NewPrimitiveTypeInfo("any")
		}
	}
	
	typeInfo := NewStructTypeInfo(s.Name, fieldTypes)
	
	// Сохраняем структуру в scope для использования в макросах
	scope.GlobalScope.Set(s.Name, value.NewValue(typeInfo))
	
	return value.NewValue(typeInfo)
}

// TypeExpr представляет выражение для передачи типа в макрос
type TypeExpr struct {
	TypeName string
}

func NewTypeExpr(typeName string) *TypeExpr {
	return &TypeExpr{TypeName: typeName}
}

func (t *TypeExpr) Eval() *value.Value {
	// Получаем информацию о типе из scope
	typeValue, found := scope.GlobalScope.Get(t.TypeName)
	if !found {
		// Если не найден пользовательский тип, проверяем примитивные типы
		switch t.TypeName {
		case "int", "integer":
			return value.NewValue(NewPrimitiveTypeInfo("int"))
		case "float", "double":
			return value.NewValue(NewPrimitiveTypeInfo("float"))
		case "string":
			return value.NewValue(NewPrimitiveTypeInfo("string"))
		case "bool", "boolean":
			return value.NewValue(NewPrimitiveTypeInfo("bool"))
		case "array":
			return value.NewValue(NewPrimitiveTypeInfo("array"))
		default:
			panic(fmt.Sprintf("unknown type: %s", t.TypeName))
		}
	}
	
	return typeValue
}