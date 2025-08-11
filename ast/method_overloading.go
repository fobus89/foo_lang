package ast

import (
	"fmt"
	"strings"
	"foo_lang/value"
)

// MethodSignature представляет сигнатуру метода для системы перегрузок
type MethodSignature struct {
	Name        string   // Имя метода
	ParamTypes  []string // Типы параметров (int, string, float, bool, any)
	ParamCount  int      // Количество параметров
	ReturnType  string   // Тип возвращаемого значения
}

// String возвращает строковое представление сигнатуры
func (s *MethodSignature) String() string {
	params := strings.Join(s.ParamTypes, ", ")
	if s.ReturnType == "" {
		return fmt.Sprintf("%s(%s)", s.Name, params)
	}
	return fmt.Sprintf("%s(%s) -> %s", s.Name, params, s.ReturnType)
}

// Equals проверяет идентичность сигнатур
func (s *MethodSignature) Equals(other *MethodSignature) bool {
	if s.Name != other.Name || s.ParamCount != other.ParamCount {
		return false
	}
	
	for i, paramType := range s.ParamTypes {
		if paramType != other.ParamTypes[i] {
			return false
		}
	}
	
	return true
}

// IsCompatible проверяет совместимость сигнатуры с переданными аргументами
func (s *MethodSignature) IsCompatible(argTypes []string) bool {
	if len(argTypes) != s.ParamCount {
		return false
	}
	
	for i, paramType := range s.ParamTypes {
		argType := argTypes[i]
		
		// Точное совпадение типов
		if paramType == argType {
			continue
		}
		
		// any принимает любой тип
		if paramType == "any" {
			continue
		}
		
		// Автоматическое приведение числовых типов
		if (paramType == "float" && argType == "int") {
			continue
		}
		
		// Несовместимые типы
		return false
	}
	
	return true
}

// MatchScore возвращает оценку соответствия сигнатуры аргументам (выше = лучше)
func (s *MethodSignature) MatchScore(argTypes []string) int {
	if !s.IsCompatible(argTypes) {
		return -1 // Несовместимо
	}
	
	score := 1000 // Базовый балл за совместимость
	
	for i, paramType := range s.ParamTypes {
		argType := argTypes[i]
		
		if paramType == argType {
			score += 100 // Точное совпадение типа
		} else if paramType == "any" {
			score += 10  // any совместим, но менее предпочтителен
		} else if paramType == "float" && argType == "int" {
			score += 50  // Автоматическое приведение int->float
		}
	}
	
	return score
}

// OverloadedMethod представляет перегруженный метод с несколькими вариантами
type OverloadedMethod struct {
	Name      string
	Overloads map[string]*MethodOverload // ключ - строковое представление сигнатуры
}

// MethodOverload представляет одну перегрузку метода
type MethodOverload struct {
	Signature   *MethodSignature
	Function    Callable // Функция или метод для вызова
}

// NewOverloadedMethod создает новый перегруженный метод
func NewOverloadedMethod(name string) *OverloadedMethod {
	return &OverloadedMethod{
		Name:      name,
		Overloads: make(map[string]*MethodOverload),
	}
}

// AddOverload добавляет новую перегрузку метода
func (om *OverloadedMethod) AddOverload(signature *MethodSignature, function Callable) error {
	sigString := signature.String()
	
	// Проверяем на дублирование сигнатур
	if _, exists := om.Overloads[sigString]; exists {
		return fmt.Errorf("method overload with signature '%s' already exists", sigString)
	}
	
	om.Overloads[sigString] = &MethodOverload{
		Signature: signature,
		Function:  function,
	}
	
	return nil
}

// ResolveOverload находит наиболее подходящую перегрузку для данных аргументов
func (om *OverloadedMethod) ResolveOverload(argTypes []string) (*MethodOverload, error) {
	var bestOverload *MethodOverload
	bestScore := -1
	
	// Ищем лучшую перегрузку
	for _, overload := range om.Overloads {
		score := overload.Signature.MatchScore(argTypes)
		
		if score > bestScore {
			bestScore = score
			bestOverload = overload
		}
	}
	
	if bestOverload == nil {
		availableSigs := make([]string, 0, len(om.Overloads))
		for sig := range om.Overloads {
			availableSigs = append(availableSigs, sig)
		}
		argTypesStr := strings.Join(argTypes, ", ")
		return nil, fmt.Errorf("no matching overload for '%s(%s)'. Available: %s", 
			om.Name, argTypesStr, strings.Join(availableSigs, ", "))
	}
	
	return bestOverload, nil
}

// GetArgTypesFromValues извлекает типы из значений аргументов
func GetArgTypesFromValues(args []*Value) []string {
	types := make([]string, len(args))
	for i, arg := range args {
		types[i] = value.GetValueTypeName(arg)
	}
	return types
}

// CreateSignatureFromFunction создает сигнатуру из TypedFuncStatement
func CreateSignatureFromFunction(name string, params []FuncParam, returnType string) *MethodSignature {
	paramTypes := make([]string, len(params))
	for i, param := range params {
		if param.TypeName == "" {
			paramTypes[i] = "any"
		} else {
			paramTypes[i] = param.TypeName
		}
	}
	
	return &MethodSignature{
		Name:       name,
		ParamTypes: paramTypes,
		ParamCount: len(params),
		ReturnType: returnType,
	}
}

// Глобальный реестр перегруженных методов
var overloadedMethods = make(map[string]*OverloadedMethod)

// RegisterOverloadedMethod регистрирует перегруженный метод
func RegisterOverloadedMethod(signature *MethodSignature, function Callable) error {
	methodName := signature.Name
	
	// Получаем или создаем перегруженный метод
	overloadedMethod, exists := overloadedMethods[methodName]
	if !exists {
		overloadedMethod = NewOverloadedMethod(methodName)
		overloadedMethods[methodName] = overloadedMethod
	}
	
	// Добавляем перегрузку
	return overloadedMethod.AddOverload(signature, function)
}

// ResolveMethodOverload разрешает вызов перегруженного метода
func ResolveMethodOverload(methodName string, argTypes []string) (Callable, error) {
	overloadedMethod, exists := overloadedMethods[methodName]
	if !exists {
		return nil, fmt.Errorf("method '%s' is not overloaded", methodName)
	}
	
	overload, err := overloadedMethod.ResolveOverload(argTypes)
	if err != nil {
		return nil, err
	}
	
	return overload.Function, nil
}

// IsOverloadedMethod проверяет, является ли метод перегруженным
func IsOverloadedMethod(methodName string) bool {
	_, exists := overloadedMethods[methodName]
	return exists
}

// ClearOverloadedMethods очищает реестр перегруженных методов (для тестов)
func ClearOverloadedMethods() {
	overloadedMethods = make(map[string]*OverloadedMethod)
}