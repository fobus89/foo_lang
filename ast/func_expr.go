package ast

import (
	"fmt"
	"foo_lang/scope"
	"foo_lang/value"
)

// FuncParam представляет типизированный параметр функции
type FuncParam struct {
	Name     string
	TypeName string // "int", "string", "float", "bool", или пустая строка для нетипизированного
	Default  Expr   // Значение по умолчанию (может быть nil)
}

type FuncStatment struct {
	funcName string
	args     []map[string]Expr
	body     Expr
	isMacro  bool
}

// TypedFuncStatement представляет функцию с типизированными параметрами
type TypedFuncStatement struct {
	FuncName   string
	Params     []FuncParam
	Body       Expr
	ReturnType string // Тип возвращаемого значения
}

// TypeConstraint представляет ограничение типа в generic функции
type TypeConstraint struct {
	TypeName   string   // Имя параметра типа, например "T"
	Constraints []string // Ограничения, например ["Drawable", "Moveable"]
}

// GenericFuncStatement представляет generic функцию с параметрами типов
type GenericFuncStatement struct {
	FuncName      string
	TypeParams    []TypeConstraint // Параметры типов с ограничениями, например [{TypeName: "T", Constraints: ["Drawable"]}]
	Params        []FuncParam      // Обычные параметры функции
	ReturnType    string           // Тип возвращаемого значения
	Body          Expr
}

func NewGenericFuncStatement(funcName string, typeParams []TypeConstraint, params []FuncParam, returnType string, body Expr) *GenericFuncStatement {
	return &GenericFuncStatement{
		FuncName:   funcName,
		TypeParams: typeParams,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
	}
}

// Создает простой параметр типа без ограничений (для обратной совместимости)
func NewSimpleTypeParam(typeName string) TypeConstraint {
	return TypeConstraint{
		TypeName:   typeName,
		Constraints: []string{},
	}
}

// Создает параметр типа с ограничениями
func NewConstrainedTypeParam(typeName string, constraints []string) TypeConstraint {
	return TypeConstraint{
		TypeName:   typeName,
		Constraints: constraints,
	}
}

func NewTypedFuncStatement(funcName string, params []FuncParam, body Expr, returnType string) *TypedFuncStatement {
	return &TypedFuncStatement{
		FuncName:   funcName,
		Params:     params,
		Body:       body,
		ReturnType: returnType,
	}
}

func (f *TypedFuncStatement) Eval() *Value {
	// Создаем замыкание с типизированными параметрами
	closure := NewTypedClosure(f.FuncName, f.Params, f.Body)
	
	// Создаем сигнатуру для перегрузки
	signature := CreateSignatureFromFunction(f.FuncName, f.Params, "")
	
	// Пытаемся зарегистрировать как перегруженную функцию
	err := RegisterOverloadedMethod(signature, closure)
	if err != nil {
		// Если это метод интерфейса, не выдаем ошибку
		// Методы интерфейсов могут иметь одинаковые имена и сигнатуры
		if !isInterfaceMethod(f.FuncName) {
			panic(fmt.Sprintf("Function overload error: %v", err))
		}
	}
	
	// Также регистрируем в обычном scope для обратной совместимости
	// (если это первая функция с таким именем)
	if existingValue, exists := scope.GlobalScope.Get(f.FuncName); !exists {
		scope.GlobalScope.Set(f.FuncName, NewValue(closure))
	} else if existingValue == nil {
		scope.GlobalScope.Set(f.FuncName, NewValue(closure))
	}
	
	return NewValue(nil)
}

func NewFuncStatment(funcName string, args []map[string]Expr, body Expr, isMacro bool) *FuncStatment {
	f := &FuncStatment{
		funcName: funcName,
		args:     args,
		body:     body,
		isMacro:  isMacro,
	}

	// Создаем замыкание для захвата переменных из текущей области видимости
	closure := NewClosure(funcName, args, body, isMacro)
	
	// Регистрируем замыкание в области видимости
	scope.GlobalScope.Set(funcName, NewValue(closure))

	return f
}

func (f *FuncStatment) Name() string {
	return f.funcName
}

func (f *FuncStatment) Params() []string {
	var args []string
	{
		for _, v := range f.args {
			for k := range v {
				args = append(args, k)
			}
		}
	}
	return args
}

func (f *FuncStatment) Eval() *Value {
	// Function definitions don't return values, they register the function in scope
	// The function is already registered in NewFuncStatment constructor
	return nil
}

func (f *FuncStatment) IsMacro() bool {
	return f.isMacro
}

func (f *FuncStatment) String() string {
	return fmt.Sprintf("func %s(%s) { %s }", f.funcName, f.args, f.body)
}

// Реализация интерфейса Callable для FuncStatment
func (f *FuncStatment) Call(args []*Value) *Value {
	bodyStm := f.body.(*BodyExpr)

	expected := len(f.args)
	passed := len(args)

	if passed > expected {
		panic(fmt.Sprintf("too many arguments: expected %d, got %d", expected, passed))
	}

	// Создаем новую область видимости для функции с проверкой рекурсии
	err := scope.GlobalScope.PushFunction()
	if err != nil {
		panic(err.Error())
	}
	defer scope.GlobalScope.PopFunction()

	// Устанавливаем параметры функции в локальной области
	for i, arg := range f.args {
		for name, expr := range arg {
			if i < len(args) {
				scope.GlobalScope.Set(name, args[i])
			} else if expr != nil {
				defaultValue := expr.Eval()
				scope.GlobalScope.Set(name, defaultValue)
			} else {
				panic(fmt.Sprintf("missing required argument: %s", name))
			}
		}
	}

	// Выполняем тело функции
	for _, stm := range bodyStm.Statments {
		if stm == nil {
			continue
		}

		result := stm.Eval()
		
		// Проверяем на return
		if result != nil && result.IsReturn() {
			return result
		}
	}

	return nil
}

// Eval для GenericFuncStatement - регистрирует generic функцию
func (g *GenericFuncStatement) Eval() *Value {
	// Сохраняем generic функцию в глобальной области видимости
	scope.GlobalScope.Set(g.FuncName, NewValue(g))
	return NewValue(g)
}

// Call для GenericFuncStatement с поддержкой типов
func (g *GenericFuncStatement) Call(args []*Value) *Value {
	// Создаем новую область видимости для функции
	scope.GlobalScope.Push()
	defer scope.GlobalScope.Pop()

	// Проверяем количество аргументов
	if len(args) != len(g.Params) {
		panic(fmt.Sprintf("function '%s' expects %d arguments, got %d", 
			g.FuncName, len(g.Params), len(args)))
	}

	// Устанавливаем параметры в области видимости
	for i, param := range g.Params {
		argValue := args[i]
		
		// Проверяем ограничения типов, если они есть
		if param.TypeName != "" {
			if err := g.validateTypeConstraints(param.TypeName, argValue); err != nil {
				panic(fmt.Sprintf("Type constraint violation in function '%s': %v", g.FuncName, err))
			}
		}
		
		scope.GlobalScope.Set(param.Name, argValue)
	}

	// Выполняем тело функции
	if bodyStm, ok := g.Body.(*BodyExpr); ok {
		for _, stm := range bodyStm.Statments {
			if stm == nil {
				continue
			}

			result := stm.Eval()
			
			// Проверяем на return
			if result != nil && result.IsReturn() {
				return result
			}
		}
	}

	return nil
}

// Name для интерфейса Callable
func (g *GenericFuncStatement) Name() string {
	return g.FuncName
}

// validateTypeConstraints проверяет, соответствует ли значение ограничениям типа
func (g *GenericFuncStatement) validateTypeConstraints(paramType string, val *Value) error {
	// Находим ограничения для данного типа
	var constraints []string
	for _, typeParam := range g.TypeParams {
		if typeParam.TypeName == paramType {
			constraints = typeParam.Constraints
			break
		}
	}
	
	// Если ограничений нет, проверка пройдена
	if len(constraints) == 0 {
		return nil
	}
	
	// Проверяем каждое ограничение
	for _, constraint := range constraints {
		if !g.checkInterfaceConstraint(constraint, val) {
			return fmt.Errorf("value does not satisfy interface constraint '%s'", constraint)
		}
	}
	
	return nil
}

// checkInterfaceConstraint проверяет, реализует ли тип данный интерфейс
func (g *GenericFuncStatement) checkInterfaceConstraint(interfaceName string, val *Value) bool {
	// Получаем определение интерфейса
	interfaceDef := GetInterface(interfaceName)
	if interfaceDef == nil {
		// Если интерфейс не найден, считаем ограничение не выполненным
		return false
	}
	
	// Определяем тип значения
	var typeName string
	valueType := value.GetValueTypeName(val)
	
	if valueType == "struct" {
		if structObj, ok := val.Any().(*StructObject); ok {
			typeName = structObj.TypeInfo.Name
		} else {
			return false
		}
	} else {
		// Для примитивных типов используем их тип
		typeName = valueType
	}
	
	// Проверяем, есть ли реализация интерфейса для данного типа
	return TypeImplementsInterface(typeName, interfaceName)
}

// isInterfaceMethod проверяет, является ли функция методом интерфейса
// Методы интерфейсов имеют общие имена и могут дублироваться для разных типов
func isInterfaceMethod(methodName string) bool {
	// Проверяем все зарегистрированные интерфейсы
	for _, interfaceDef := range registeredInterfaces {
		if interfaceDef.HasMethod(methodName) {
			return true
		}
	}
	return false
}
