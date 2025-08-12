; Tree-sitter highlighting queries for Foo language

; Keywords
[
  "let"
  "const"
  "fn" 
  "struct"
  "enum"
  "interface"
  "impl"
  "extension"
  "macro"
  "import"
  "export"
  "from"
  "as"
] @keyword

[
  "if"
  "else" 
  "for"
  "match"
  "return"
  "yield"
  "break"
  "async"
  "await"
] @keyword.control

; Types
[
  "int"
  "float"
  "string"
  "bool"
] @type.builtin

; Booleans
[
  "true"
  "false"
] @constant.builtin.boolean

; Null
"null" @constant.builtin

; Numbers
(integer_literal) @constant.numeric.integer
(float_literal) @constant.numeric.float

; Strings
(string_literal) @string

; Comments
(comment) @comment.line
(block_comment) @comment.block

; Functions
(function_declaration
  name: (identifier) @function)

(call_expression
  function: (identifier) @function.call)

; Macros
(macro_declaration
  name: (identifier) @function.macro)

(macro_call
  name: (identifier) @function.macro)

; Variables
(variable_declaration
  name: (identifier) @variable)

(constant_declaration
  name: (identifier) @constant)

; Parameters
(parameter
  name: (identifier) @variable.parameter)

; Fields
(field
  name: (identifier) @property)

; Operators
[
  "+"
  "-"
  "*"
  "/"
  "%"
  "++"
  "--"
] @operator.arithmetic

[
  "=="
  "!="
  "<"
  ">"
  "<="
  ">="
] @operator.comparison

[
  "&&"
  "||"
  "!"
] @operator.logical

[
  "="
  "+="
  "-="
  "*="
  "/="
  "%="
] @operator.assignment

[
  "=>"
  "->"
  "<-"
] @operator.arrow

"?" @operator.ternary

"|" @operator.union

; Punctuation
[
  "("
  ")"
  "["
  "]"
  "{"
  "}"
] @punctuation.bracket

[
  ","
  ";"
  ":"
  "."
] @punctuation.delimiter

; Error handling
"@" @punctuation.special

; Identifiers
(identifier) @variable