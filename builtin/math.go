package builtin

import (
	"math"
	"foo_lang/value"
)

// MathFunction представляет математическую функцию
type MathFunction struct {
	name string
	fn   func([]*value.Value) *value.Value
}

func (mf *MathFunction) Eval() *value.Value {
	// Математические функции не вызываются напрямую через Eval
	return value.NewValue(mf)
}

func (mf *MathFunction) Call(args []*value.Value) *value.Value {
	return mf.fn(args)
}

func (mf *MathFunction) Name() string {
	return mf.name
}

func (mf *MathFunction) String() string {
	return "builtin function " + mf.name
}

// CreateMathFunctions создает все встроенные математические функции
func CreateMathFunctions() map[string]*value.Value {
	functions := make(map[string]*value.Value)
	
	// sin(x) - синус
	functions["sin"] = value.NewValue(&MathFunction{
		name: "sin",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("sin() takes exactly 1 argument")
			}
			x := args[0].Float64()
			return value.NewValue(math.Sin(x))
		},
	})
	
	// cos(x) - косинус
	functions["cos"] = value.NewValue(&MathFunction{
		name: "cos", 
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("cos() takes exactly 1 argument")
			}
			x := args[0].Float64()
			return value.NewValue(math.Cos(x))
		},
	})
	
	// tan(x) - тангенс
	functions["tan"] = value.NewValue(&MathFunction{
		name: "tan",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("tan() takes exactly 1 argument")
			}
			x := args[0].Float64()
			return value.NewValue(math.Tan(x))
		},
	})
	
	// sqrt(x) - квадратный корень
	functions["sqrt"] = value.NewValue(&MathFunction{
		name: "sqrt",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("sqrt() takes exactly 1 argument")
			}
			x := args[0].Float64()
			if x < 0 {
				panic("sqrt() argument must be non-negative")
			}
			return value.NewValue(math.Sqrt(x))
		},
	})
	
	// abs(x) - абсолютное значение
	functions["abs"] = value.NewValue(&MathFunction{
		name: "abs",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("abs() takes exactly 1 argument")
			}
			x := args[0].Float64()
			return value.NewValue(math.Abs(x))
		},
	})
	
	// pow(x, y) - возведение в степень
	functions["pow"] = value.NewValue(&MathFunction{
		name: "pow",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 2 {
				panic("pow() takes exactly 2 arguments")
			}
			x := args[0].Float64()
			y := args[1].Float64()
			return value.NewValue(math.Pow(x, y))
		},
	})
	
	// floor(x) - округление вниз
	functions["floor"] = value.NewValue(&MathFunction{
		name: "floor",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("floor() takes exactly 1 argument")
			}
			x := args[0].Float64()
			return value.NewValue(math.Floor(x))
		},
	})
	
	// ceil(x) - округление вверх
	functions["ceil"] = value.NewValue(&MathFunction{
		name: "ceil", 
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("ceil() takes exactly 1 argument")
			}
			x := args[0].Float64()
			return value.NewValue(math.Ceil(x))
		},
	})
	
	// round(x) - округление к ближайшему целому
	functions["round"] = value.NewValue(&MathFunction{
		name: "round",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("round() takes exactly 1 argument")
			}
			x := args[0].Float64()
			return value.NewValue(math.Round(x))
		},
	})
	
	// min(x, y) - минимум
	functions["min"] = value.NewValue(&MathFunction{
		name: "min",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 2 {
				panic("min() takes exactly 2 arguments")
			}
			x := args[0].Float64()
			y := args[1].Float64()
			return value.NewValue(math.Min(x, y))
		},
	})
	
	// max(x, y) - максимум
	functions["max"] = value.NewValue(&MathFunction{
		name: "max",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 2 {
				panic("max() takes exactly 2 arguments")
			}
			x := args[0].Float64()
			y := args[1].Float64()
			return value.NewValue(math.Max(x, y))
		},
	})
	
	// log(x) - натуральный логарифм
	functions["log"] = value.NewValue(&MathFunction{
		name: "log",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("log() takes exactly 1 argument")
			}
			x := args[0].Float64()
			if x <= 0 {
				panic("log() argument must be positive")
			}
			return value.NewValue(math.Log(x))
		},
	})
	
	// log10(x) - логарифм по основанию 10
	functions["log10"] = value.NewValue(&MathFunction{
		name: "log10",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("log10() takes exactly 1 argument")
			}
			x := args[0].Float64()
			if x <= 0 {
				panic("log10() argument must be positive")
			}
			return value.NewValue(math.Log10(x))
		},
	})
	
	// exp(x) - e^x
	functions["exp"] = value.NewValue(&MathFunction{
		name: "exp",
		fn: func(args []*value.Value) *value.Value {
			if len(args) != 1 {
				panic("exp() takes exactly 1 argument")
			}
			x := args[0].Float64()
			return value.NewValue(math.Exp(x))
		},
	})
	
	return functions
}

// InitializeMathFunctions добавляет математические функции в глобальную область видимости
func InitializeMathFunctions(scopeStack ScopeStack) {
	functions := CreateMathFunctions()
	for name, fn := range functions {
		scopeStack.Set(name, fn)
	}
}

