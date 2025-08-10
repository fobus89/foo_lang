package iter

import "testing"

func TestIterator_Number(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result, ok := iter.ReadMore('1', '2', '3'); !ok {
		t.Errorf("Expected true, got false %s", string(result))
	}
}

func TestIterator_SkipFunc(t *testing.T) {
	iter := NewIterator([]rune("123"))

	iter.SkipFunc(func(r rune) bool {
		return r == '1'
	})

	if result, ok := iter.ReadMore('2', '3'); !ok {
		t.Errorf("Expected true, got false %s", string(result))
	}
}

func TestIterator_Skip(t *testing.T) {
	iter := NewIterator([]rune("123"))

	iter.Skip('1')

	if result, ok := iter.ReadMore('2', '3'); !ok {
		t.Errorf("Expected true, got false %s", string(result))
	}
}

func TestIterator_Read(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result, ok := iter.Read('1', '2', '3'); !ok {
		t.Errorf("Expected true, got false %s", string(result))
	}
}

func TestIterator_ReadFunc(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result, ok := iter.ReadFunc(func(r rune) bool {
		return r == '1'
	}); !ok {
		t.Errorf("Expected true, got false %s", string(result))
	}
}

func TestIterator_ReadMore(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result, ok := iter.ReadMore('1', '2', '3'); !ok {
		t.Errorf("Expected true, got false %s", string(result))
	}
}

func TestIterator_ReadMoreFunc(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result, ok := iter.ReadMoreFunc(func(r rune) bool {
		return r == '1'
	}); !ok {
		t.Errorf("Expected true, got false %s", string(result))
	}
}

func TestIterator_Peek(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result := iter.Peek(0); result != '1' {
		t.Errorf("Expected true, got false %s", string(result))
	}
}

func TestIterator_Get(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result := iter.Get(); result != '1' {
		t.Errorf("Expected true, got false %s", string(result))
	}
}

func TestIterator_Match(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result := iter.Match('1'); !result {
		t.Errorf("Expected true, got false %v", result)
	}
}

func TestIterator_MatchN(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result := iter.MatchN('1', 0); !result {
		t.Errorf("Expected true, got false %v", result)
	}
}

func TestIterator_MatchAndNext(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result := iter.MatchAndNext('1'); !result {
		t.Errorf("Expected true, got false %v", result)
	}
}

func TestIterator_MatchNAndNext(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result := iter.MatchNAndNext('1', 0); !result {
		t.Errorf("Expected true, got false %v", result)
	}
}

func TestIterator_HasNext(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result := iter.HasNext(); !result {
		t.Errorf("Expected true, got false %v", result)
	}
}

func TestIterator_Next(t *testing.T) {
	iter := NewIterator([]rune("123"))

	if result := iter.Next(); result != '1' {
		t.Errorf("Expected true, got false %s", string(result))
	}
}
