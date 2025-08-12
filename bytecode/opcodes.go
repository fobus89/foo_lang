package bytecode

import "fmt"

// OpCode представляет операцию в bytecode
type OpCode uint8

const (
	// Константы и литералы
	OP_CONSTANT OpCode = iota
	OP_NIL
	OP_TRUE
	OP_FALSE
	
	// Арифметические операции
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_MODULO
	OP_NEGATE
	
	// Логические операции
	OP_NOT
	OP_AND
	OP_OR
	
	// Операции сравнения
	OP_EQUAL
	OP_NOT_EQUAL
	OP_GREATER
	OP_GREATER_EQUAL
	OP_LESS
	OP_LESS_EQUAL
	
	// Переменные
	OP_GET_GLOBAL
	OP_SET_GLOBAL
	OP_GET_LOCAL
	OP_SET_LOCAL
	OP_DEFINE_GLOBAL
	
	// Управление потоком
	OP_JUMP
	OP_JUMP_IF_FALSE
	OP_LOOP
	OP_CALL
	OP_RETURN
	
	// Массивы и объекты
	OP_ARRAY
	OP_OBJECT
	OP_INDEX
	OP_SET_INDEX
	
	// Встроенные функции
	OP_PRINT
	OP_PRINTLN
	
	// Функции
	OP_CLOSURE
	OP_CALL_FUNCTION
	
	// Стек операции
	OP_POP
	OP_DUP
	
	// Async/await
	OP_ASYNC
	OP_AWAIT
	OP_PROMISE_ALL
	OP_PROMISE_ANY
	OP_SLEEP
	
	// HTTP операции
	OP_HTTP_GET
	OP_HTTP_POST
	OP_HTTP_PUT
	OP_HTTP_DELETE
	OP_HTTP_START_SERVER
	OP_HTTP_ROUTE
	
	// Файловая система
	OP_READ_FILE
	OP_WRITE_FILE
	OP_FILE_EXISTS
	
	// Математические функции
	OP_MATH_SIN
	OP_MATH_COS
	OP_MATH_SQRT
	OP_MATH_POW
	OP_MATH_ABS
	
	// Строковые операции
	OP_STRING_LEN
	OP_STRING_CONCAT
	OP_STRING_CHAR_AT
	OP_STRING_SUBSTRING
	
	// Методы типов
	OP_METHOD_CALL
	OP_PROPERTY_ACCESS
	
	// Макросы и метапрограммирование
	OP_MACRO_CALL
	OP_TYPE_OF
	OP_TYPE_CHECK
	
	// Интерфейсы и структуры
	OP_STRUCT_INSTANCE
	OP_INTERFACE_CALL
	OP_IMPL_METHOD
	
	// Extension methods
	OP_EXTENSION_CALL
	
	// Отладка и профилинг
	OP_DEBUG_TRACE
	OP_PROFILE_START
	OP_PROFILE_END
)

// Instruction представляет одну инструкцию bytecode
type Instruction struct {
	OpCode   OpCode
	Operands []int
	Line     int // Для отладки
}

// Chunk представляет блок bytecode инструкций
type Chunk struct {
	Code      []Instruction
	Constants []interface{}
	Lines     []int
}

// NewChunk создает новый chunk
func NewChunk() *Chunk {
	return &Chunk{
		Code:      make([]Instruction, 0),
		Constants: make([]interface{}, 0),
		Lines:     make([]int, 0),
	}
}

// WriteInstruction добавляет инструкцию в chunk
func (c *Chunk) WriteInstruction(opcode OpCode, operands []int, line int) int {
	instruction := Instruction{
		OpCode:   opcode,
		Operands: operands,
		Line:     line,
	}
	c.Code = append(c.Code, instruction)
	c.Lines = append(c.Lines, line)
	return len(c.Code) - 1
}

// AddConstant добавляет константу и возвращает её индекс
func (c *Chunk) AddConstant(value interface{}) int {
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

// DisassembleChunk выводит human-readable представление chunk'а
func DisassembleChunk(chunk *Chunk, name string) {
	println("== " + name + " ==")
	
	for i, instruction := range chunk.Code {
		DisassembleInstruction(&instruction, i)
	}
}

// DisassembleInstruction выводит human-readable представление инструкции
func DisassembleInstruction(instruction *Instruction, offset int) {
	print(fmt.Sprintf("%04d ", offset))
	
	switch instruction.OpCode {
	case OP_CONSTANT:
		print(fmt.Sprintf("OP_CONSTANT %d", instruction.Operands[0]))
	case OP_NIL:
		print("OP_NIL")
	case OP_TRUE:
		print("OP_TRUE")
	case OP_FALSE:
		print("OP_FALSE")
	case OP_ADD:
		print("OP_ADD")
	case OP_SUBTRACT:
		print("OP_SUBTRACT")
	case OP_MULTIPLY:
		print("OP_MULTIPLY")
	case OP_DIVIDE:
		print("OP_DIVIDE")
	case OP_MODULO:
		print("OP_MODULO")
	case OP_NEGATE:
		print("OP_NEGATE")
	case OP_NOT:
		print("OP_NOT")
	case OP_AND:
		print("OP_AND")
	case OP_OR:
		print("OP_OR")
	case OP_EQUAL:
		print("OP_EQUAL")
	case OP_NOT_EQUAL:
		print("OP_NOT_EQUAL")
	case OP_GREATER:
		print("OP_GREATER")
	case OP_GREATER_EQUAL:
		print("OP_GREATER_EQUAL")
	case OP_LESS:
		print("OP_LESS")
	case OP_LESS_EQUAL:
		print("OP_LESS_EQUAL")
	case OP_GET_GLOBAL:
		print(fmt.Sprintf("OP_GET_GLOBAL %d", instruction.Operands[0]))
	case OP_SET_GLOBAL:
		print(fmt.Sprintf("OP_SET_GLOBAL %d", instruction.Operands[0]))
	case OP_GET_LOCAL:
		print(fmt.Sprintf("OP_GET_LOCAL %d", instruction.Operands[0]))
	case OP_SET_LOCAL:
		print(fmt.Sprintf("OP_SET_LOCAL %d", instruction.Operands[0]))
	case OP_DEFINE_GLOBAL:
		print(fmt.Sprintf("OP_DEFINE_GLOBAL %d", instruction.Operands[0]))
	case OP_JUMP:
		print(fmt.Sprintf("OP_JUMP %d", instruction.Operands[0]))
	case OP_JUMP_IF_FALSE:
		print(fmt.Sprintf("OP_JUMP_IF_FALSE %d", instruction.Operands[0]))
	case OP_LOOP:
		print(fmt.Sprintf("OP_LOOP %d", instruction.Operands[0]))
	case OP_CALL:
		print(fmt.Sprintf("OP_CALL %d", instruction.Operands[0]))
	case OP_RETURN:
		print("OP_RETURN")
	case OP_ARRAY:
		print(fmt.Sprintf("OP_ARRAY %d", instruction.Operands[0]))
	case OP_OBJECT:
		print(fmt.Sprintf("OP_OBJECT %d", instruction.Operands[0]))
	case OP_INDEX:
		print("OP_INDEX")
	case OP_SET_INDEX:
		print("OP_SET_INDEX")
	case OP_PRINT:
		print("OP_PRINT")
	case OP_PRINTLN:
		print("OP_PRINTLN")
	case OP_CLOSURE:
		print(fmt.Sprintf("OP_CLOSURE %d", instruction.Operands[0]))
	case OP_CALL_FUNCTION:
		print(fmt.Sprintf("OP_CALL_FUNCTION %d", instruction.Operands[0]))
	case OP_POP:
		print("OP_POP")
	case OP_DUP:
		print("OP_DUP")
	case OP_ASYNC:
		print("OP_ASYNC")
	case OP_AWAIT:
		print("OP_AWAIT")
	case OP_PROMISE_ALL:
		print(fmt.Sprintf("OP_PROMISE_ALL %d", instruction.Operands[0]))
	case OP_PROMISE_ANY:
		print(fmt.Sprintf("OP_PROMISE_ANY %d", instruction.Operands[0]))
	case OP_SLEEP:
		print("OP_SLEEP")
	case OP_HTTP_GET:
		print("OP_HTTP_GET")
	case OP_HTTP_POST:
		print("OP_HTTP_POST")
	case OP_HTTP_PUT:
		print("OP_HTTP_PUT")
	case OP_HTTP_DELETE:
		print("OP_HTTP_DELETE")
	case OP_READ_FILE:
		print("OP_READ_FILE")
	case OP_WRITE_FILE:
		print("OP_WRITE_FILE")
	case OP_MATH_SIN:
		print("OP_MATH_SIN")
	case OP_MATH_COS:
		print("OP_MATH_COS")
	case OP_MATH_SQRT:
		print("OP_MATH_SQRT")
	case OP_STRING_LEN:
		print("OP_STRING_LEN")
	case OP_STRING_CONCAT:
		print("OP_STRING_CONCAT")
	case OP_METHOD_CALL:
		print(fmt.Sprintf("OP_METHOD_CALL %d", instruction.Operands[0]))
	case OP_PROPERTY_ACCESS:
		print(fmt.Sprintf("OP_PROPERTY_ACCESS %d", instruction.Operands[0]))
	case OP_DEBUG_TRACE:
		print("OP_DEBUG_TRACE")
	case OP_PROFILE_START:
		print(fmt.Sprintf("OP_PROFILE_START %d", instruction.Operands[0]))
	case OP_PROFILE_END:
		print(fmt.Sprintf("OP_PROFILE_END %d", instruction.Operands[0]))
	default:
		print(fmt.Sprintf("Unknown opcode %d", instruction.OpCode))
	}
	
	println()
}

