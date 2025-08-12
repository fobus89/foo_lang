package builtin

import "foo_lang/value"

// ScopeStack интерфейс для области видимости (чтобы избежать циклических импортов)
type ScopeStack interface {
	Set(name string, val *value.Value)
}