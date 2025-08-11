package ast

import (
	"fmt"
	"foo_lang/scope"
)

// InterfaceMethod представляет сигнатуру метода в интерфейсе
type InterfaceMethod struct {
	Name       string      // Имя метода
	Params     []FuncParam // Параметры метода
	ReturnType string      // Тип возвращаемого значения
}

// InterfaceDefinition представляет определение интерфейса
type InterfaceDefinition struct {
	Name    string            // Имя интерфейса
	Methods []InterfaceMethod // Методы интерфейса
}

// NewInterfaceDefinition создает новое определение интерфейса
func NewInterfaceDefinition(name string, methods []InterfaceMethod) *InterfaceDefinition {
	return &InterfaceDefinition{
		Name:    name,
		Methods: methods,
	}
}

// Eval регистрирует интерфейс в области видимости
func (i *InterfaceDefinition) Eval() *Value {
	// Регистрируем интерфейс в глобальной области видимости
	RegisterInterface(i.Name, i)
	scope.GlobalScope.Set(i.Name, NewValue(i))
	
	return NewValue(i)
}

// String возвращает строковое представление интерфейса
func (i *InterfaceDefinition) String() string {
	return fmt.Sprintf("interface %s { ... }", i.Name)
}

// GetMethodByName находит метод интерфейса по имени
func (i *InterfaceDefinition) GetMethodByName(methodName string) *InterfaceMethod {
	for _, method := range i.Methods {
		if method.Name == methodName {
			return &method
		}
	}
	return nil
}

// HasMethod проверяет, есть ли метод с данным именем в интерфейсе
func (i *InterfaceDefinition) HasMethod(methodName string) bool {
	return i.GetMethodByName(methodName) != nil
}

// MethodsMatch проверяет, совпадает ли сигнатура метода с требованиями интерфейса
func (i *InterfaceDefinition) MethodMatches(methodName string, params []FuncParam, returnType string) bool {
	interfaceMethod := i.GetMethodByName(methodName)
	if interfaceMethod == nil {
		return false
	}
	
	// Проверяем количество параметров
	if len(interfaceMethod.Params) != len(params) {
		return false
	}
	
	// Проверяем типы параметров
	for j, interfaceParam := range interfaceMethod.Params {
		if interfaceParam.TypeName != params[j].TypeName {
			return false
		}
	}
	
	// Проверяем тип возвращаемого значения
	if interfaceMethod.ReturnType != returnType {
		return false
	}
	
	return true
}

// ImplBlock представляет блок реализации интерфейса для типа
type ImplBlock struct {
	InterfaceName string            // Имя реализуемого интерфейса
	TypeName      string            // Имя типа, для которого реализуется интерфейс
	Methods       []*TypedFuncStatement // Методы реализации
}

// NewImplBlock создает новый блок реализации
func NewImplBlock(interfaceName, typeName string, methods []*TypedFuncStatement) *ImplBlock {
	return &ImplBlock{
		InterfaceName: interfaceName,
		TypeName:      typeName,
		Methods:       methods,
	}
}

// Eval выполняет блок реализации интерфейса
func (impl *ImplBlock) Eval() *Value {
	// Получаем интерфейс
	interfaceDef := GetInterface(impl.InterfaceName)
	if interfaceDef == nil {
		panic(fmt.Sprintf("Interface '%s' is not defined", impl.InterfaceName))
	}
	
	// Проверяем, что все методы интерфейса реализованы
	err := impl.ValidateImplementation(interfaceDef)
	if err != nil {
		panic(fmt.Sprintf("Implementation error: %v", err))
	}
	
	// Регистрируем реализацию
	RegisterImplementation(impl.TypeName, impl.InterfaceName, impl)
	
	// Регистрируем методы в области видимости
	for _, method := range impl.Methods {
		method.Eval()
	}
	
	return NewValue(impl)
}

// ValidateImplementation проверяет, что все методы интерфейса правильно реализованы
func (impl *ImplBlock) ValidateImplementation(interfaceDef *InterfaceDefinition) error {
	// Проверяем, что каждый метод интерфейса реализован
	for _, interfaceMethod := range interfaceDef.Methods {
		found := false
		
		for _, method := range impl.Methods {
			if method.FuncName == interfaceMethod.Name {
				// Проверяем совпадение сигнатур
				if !interfaceDef.MethodMatches(interfaceMethod.Name, method.Params, method.ReturnType) {
					return fmt.Errorf("method '%s' signature doesn't match interface", interfaceMethod.Name)
				}
				found = true
				break
			}
		}
		
		if !found {
			return fmt.Errorf("method '%s' is not implemented", interfaceMethod.Name)
		}
	}
	
	return nil
}

// String возвращает строковое представление блока реализации
func (impl *ImplBlock) String() string {
	return fmt.Sprintf("impl %s for %s", impl.InterfaceName, impl.TypeName)
}

// Глобальный реестр интерфейсов
var registeredInterfaces = make(map[string]*InterfaceDefinition)

// Глобальный реестр реализаций интерфейсов
var interfaceImplementations = make(map[string]map[string]*ImplBlock) // [typeName][interfaceName] = implBlock

// RegisterInterface регистрирует новый интерфейс
func RegisterInterface(name string, interfaceDef *InterfaceDefinition) {
	registeredInterfaces[name] = interfaceDef
}

// GetInterface получает интерфейс по имени
func GetInterface(name string) *InterfaceDefinition {
	return registeredInterfaces[name]
}

// RegisterImplementation регистрирует реализацию интерфейса для типа
func RegisterImplementation(typeName, interfaceName string, impl *ImplBlock) {
	if interfaceImplementations[typeName] == nil {
		interfaceImplementations[typeName] = make(map[string]*ImplBlock)
	}
	interfaceImplementations[typeName][interfaceName] = impl
}

// GetImplementation получает реализацию интерфейса для типа
func GetImplementation(typeName, interfaceName string) *ImplBlock {
	if typeImpls, exists := interfaceImplementations[typeName]; exists {
		return typeImpls[interfaceName]
	}
	return nil
}

// TypeImplementsInterface проверяет, реализует ли тип интерфейс
func TypeImplementsInterface(typeName, interfaceName string) bool {
	return GetImplementation(typeName, interfaceName) != nil
}

// ClearInterfaces очищает реестр интерфейсов (для тестов)
func ClearInterfaces() {
	registeredInterfaces = make(map[string]*InterfaceDefinition)
	interfaceImplementations = make(map[string]map[string]*ImplBlock)
}