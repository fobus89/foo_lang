package iter

type Iterator[T comparable] interface {
	Peek(pos int) T
	Match(symbol T) bool
	MatchN(symbol T, n int) bool
	MatchAndNext(symbol T) bool
	MatchNAndNext(symbol T, n int) bool
	Next() T
	Get() T
	Data() []T
	Position() int
	HasNext() bool
	IsZero() bool

	Read(symbols ...T) (T, bool)
	ReadMoreFunc(fn func(T) bool) ([]T, bool)
	ReadMore(symbols ...T) ([]T, bool)
	ReadFunc(fn func(T) bool) (T, bool)

	Skip(symbols ...T)
	SkipFunc(fn func(T) bool)
}

type Iterable[T comparable] struct {
	pos  int
	zero T
	data []T
}

func NewIterator[T comparable](data []T) Iterator[T] {
	return &Iterable[T]{
		pos:  0,
		data: data,
	}
}

func (i *Iterable[T]) Skip(symbols ...T) {
	for _, symbol := range symbols {
		i.MatchAndNext(symbol)
	}
}

func (i *Iterable[T]) SkipFunc(fn func(T) bool) {
	for fn(i.Get()) {
		i.Next()
	}
}

func (i *Iterable[T]) ReadMore(symbols ...T) ([]T, bool) {
	var result []T

	for _, symbol := range symbols {

		if !i.MatchAndNext(symbol) {
			break
		}

		result = append(result, symbol)
	}

	return result, len(result) > 0
}

func (i *Iterable[T]) ReadMoreFunc(fn func(T) bool) ([]T, bool) {
	var result []T

	for {
		symbol := i.Get()

		if !fn(symbol) {
			break
		}

		result = append(result, symbol)

		i.Next()
	}

	return result, len(result) > 0
}

func (i *Iterable[T]) Read(symbols ...T) (T, bool) {
	for _, symbol := range symbols {
		if i.MatchAndNext(symbol) {
			return symbol, true
		}
	}
	return i.zero, false
}

func (i *Iterable[T]) ReadFunc(fn func(T) bool) (T, bool) {
	symbol := i.Get()

	if fn(symbol) {
		i.Next()
		return symbol, true
	}

	return i.zero, false
}

func (i *Iterable[T]) Peek(pos int) T {

	cursor := i.pos + pos

	if cursor >= len(i.data) {
		return i.zero
	}

	return i.data[cursor]
}

func (i *Iterable[T]) Match(symbol T) bool {
	return i.Peek(0) == symbol
}

func (i *Iterable[T]) MatchN(symbol T, n int) bool {
	return i.Peek(n) == symbol
}

func (i *Iterable[T]) MatchAndNext(symbol T) bool {
	if i.Match(symbol) {
		i.Next()
		return true
	}
	return false
}

func (i *Iterable[T]) MatchNAndNext(symbol T, n int) bool {
	if i.MatchN(symbol, n) {
		i.Next()
		return true
	}
	return false
}

func (i *Iterable[T]) Next() T {
	if !i.HasNext() {
		return i.zero
	}
	r := i.data[i.pos]
	i.pos++
	return r
}

func (i *Iterable[T]) HasNext() bool {
	return i.pos < len(i.data)
}

func (i *Iterable[T]) Data() []T {
	return i.data
}

func (i *Iterable[T]) Position() int {
	return i.pos
}

func (i *Iterable[T]) Get() T {
	return i.Peek(0)
}

func (i *Iterable[T]) IsZero() bool {
	return i.Peek(0) == i.zero
}
