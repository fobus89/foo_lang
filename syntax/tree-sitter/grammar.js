module.exports = grammar({
  name: 'foo',

  extras: $ => [
    /\s+/,
    $.comment,
    $.block_comment
  ],

  rules: {
    source_file: $ => repeat($._statement),

    _statement: $ => choice(
      $.variable_declaration,
      $.constant_declaration,
      $.function_declaration,
      $.struct_declaration,
      $.enum_declaration,
      $.interface_declaration,
      $.impl_declaration,
      $.extension_declaration,
      $.macro_declaration,
      $.import_statement,
      $.export_statement,
      $.expression_statement,
      $.if_statement,
      $.for_statement,
      $.match_statement,
      $.return_statement,
      $.yield_statement,
      $.break_statement
    ),

    // Comments
    comment: $ => token(seq('//', /.*/)),
    block_comment: $ => token(seq('/*', /[^*]*\*+([^/*][^*]*\*+)*/, '/')),

    // Variables and Constants
    variable_declaration: $ => seq(
      'let',
      field('name', $.identifier),
      optional(seq(':', field('type', $._type))),
      '=',
      field('value', $._expression)
    ),

    constant_declaration: $ => seq(
      'const',
      field('name', $.identifier),
      optional(seq(':', field('type', $._type))),
      '=',
      field('value', $._expression)
    ),

    // Functions
    function_declaration: $ => seq(
      'fn',
      field('name', $.identifier),
      optional($.generic_parameters),
      '(',
      optional($.parameter_list),
      ')',
      optional(seq('->', field('return_type', $._type))),
      field('body', $.block)
    ),

    generic_parameters: $ => seq(
      '<',
      sep1($.generic_parameter, ','),
      '>'
    ),

    generic_parameter: $ => seq(
      $.identifier,
      optional(seq(':', sep1($.identifier, '+')))
    ),

    parameter_list: $ => sep1($.parameter, ','),

    parameter: $ => seq(
      field('name', $.identifier),
      optional(seq(':', field('type', $._type))),
      optional(seq('=', field('default', $._expression)))
    ),

    // Structures, Enums, Interfaces
    struct_declaration: $ => seq(
      'struct',
      field('name', $.identifier),
      '{',
      optional($.field_list),
      '}'
    ),

    field_list: $ => sep1($.field, ','),

    field: $ => seq(
      field('name', $.identifier),
      ':',
      field('type', $._type)
    ),

    enum_declaration: $ => seq(
      'enum',
      field('name', $.identifier),
      '{',
      sep1($.identifier, ','),
      '}'
    ),

    interface_declaration: $ => seq(
      'interface',
      field('name', $.identifier),
      '{',
      repeat($.interface_method),
      '}'
    ),

    interface_method: $ => seq(
      'fn',
      field('name', $.identifier),
      '(',
      optional($.parameter_list),
      ')',
      optional(seq('->', field('return_type', $._type)))
    ),

    impl_declaration: $ => seq(
      'impl',
      field('interface', $.identifier),
      'for',
      field('type', $.identifier),
      '{',
      repeat($.function_declaration),
      '}'
    ),

    extension_declaration: $ => seq(
      'extension',
      field('type', $.identifier),
      '{',
      repeat($.function_declaration),
      '}'
    ),

    // Macros
    macro_declaration: $ => seq(
      'macro',
      field('name', $.identifier),
      '(',
      optional($.parameter_list),
      ')',
      field('body', $.block)
    ),

    macro_call: $ => seq(
      '@',
      field('name', $.identifier),
      '(',
      optional($.argument_list),
      ')'
    ),

    // Import/Export
    import_statement: $ => choice(
      seq('import', field('path', $.string_literal)),
      seq('import', '{', sep1($.identifier, ','), '}', 'from', field('path', $.string_literal)),
      seq('import', '*', 'as', field('alias', $.identifier), 'from', field('path', $.string_literal))
    ),

    export_statement: $ => seq(
      'export',
      choice(
        $.variable_declaration,
        $.constant_declaration,
        $.function_declaration,
        $.struct_declaration,
        $.enum_declaration
      )
    ),

    // Control Flow
    if_statement: $ => seq(
      'if',
      field('condition', $._expression),
      field('then', $.block),
      optional(seq('else', field('else', choice($.block, $.if_statement))))
    ),

    for_statement: $ => seq(
      'for',
      choice(
        seq($.variable_declaration, ';', $._expression, ';', $._expression),
        $._expression
      ),
      field('body', $.block)
    ),

    match_statement: $ => seq(
      'match',
      field('value', $._expression),
      '{',
      repeat($.match_arm),
      '}'
    ),

    match_arm: $ => seq(
      field('pattern', $._expression),
      '=>',
      field('body', choice($.block, $._expression))
    ),

    return_statement: $ => seq(
      'return',
      optional(sep1($._expression, ','))
    ),

    yield_statement: $ => seq(
      'yield',
      $._expression
    ),

    break_statement: $ => 'break',

    // Expressions
    _expression: $ => choice(
      $.binary_expression,
      $.unary_expression,
      $.ternary_expression,
      $.call_expression,
      $.member_expression,
      $.index_expression,
      $.assignment_expression,
      $.async_expression,
      $.await_expression,
      $.macro_call,
      $.array_literal,
      $.object_literal,
      $.function_literal,
      $.identifier,
      $.integer_literal,
      $.float_literal,
      $.string_literal,
      $.boolean_literal,
      $.null_literal,
      $.parenthesized_expression
    ),

    binary_expression: $ => choice(
      prec.left(10, seq($._expression, '+', $._expression)),
      prec.left(10, seq($._expression, '-', $._expression)),
      prec.left(11, seq($._expression, '*', $._expression)),
      prec.left(11, seq($._expression, '/', $._expression)),
      prec.left(11, seq($._expression, '%', $._expression)),
      prec.left(7, seq($._expression, '==', $._expression)),
      prec.left(7, seq($._expression, '!=', $._expression)),
      prec.left(8, seq($._expression, '<', $._expression)),
      prec.left(8, seq($._expression, '>', $._expression)),
      prec.left(8, seq($._expression, '<=', $._expression)),
      prec.left(8, seq($._expression, '>=', $._expression)),
      prec.left(4, seq($._expression, '&&', $._expression)),
      prec.left(3, seq($._expression, '||', $._expression))
    ),

    unary_expression: $ => choice(
      prec(12, seq('!', $._expression)),
      prec(12, seq('-', $._expression)),
      prec(12, seq('++', $._expression)),
      prec(12, seq('--', $._expression))
    ),

    ternary_expression: $ => prec.right(2, seq(
      field('condition', $._expression),
      '?',
      field('true', $._expression),
      ':',
      field('false', $._expression)
    )),

    assignment_expression: $ => prec.right(1, seq(
      field('left', $._expression),
      choice('=', '+=', '-=', '*=', '/=', '%='),
      field('right', $._expression)
    )),

    call_expression: $ => seq(
      field('function', $._expression),
      '(',
      optional($.argument_list),
      ')'
    ),

    member_expression: $ => seq(
      field('object', $._expression),
      '.',
      field('property', $.identifier)
    ),

    index_expression: $ => seq(
      field('object', $._expression),
      '[',
      field('index', $._expression),
      ']'
    ),

    async_expression: $ => seq(
      'async',
      $._expression
    ),

    await_expression: $ => seq(
      'await',
      $._expression
    ),

    // Literals
    array_literal: $ => seq(
      '[',
      optional(sep1($._expression, ',')),
      ']'
    ),

    object_literal: $ => seq(
      '{',
      optional(sep1($.object_field, ',')),
      '}'
    ),

    object_field: $ => seq(
      choice($.identifier, $.string_literal),
      ':',
      $._expression
    ),

    function_literal: $ => seq(
      'fn',
      '(',
      optional($.parameter_list),
      ')',
      choice(
        seq('=>', $._expression),
        $.block
      )
    ),

    // Types
    _type: $ => choice(
      $.primitive_type,
      $.union_type,
      $.optional_type,
      $.tuple_type,
      $.identifier
    ),

    primitive_type: $ => choice('int', 'float', 'string', 'bool'),

    union_type: $ => seq(
      $._type,
      repeat1(seq('|', $._type))
    ),

    optional_type: $ => seq(
      $._type,
      '?'
    ),

    tuple_type: $ => seq(
      '(',
      sep1($._type, ','),
      ')'
    ),

    // Block and Expression Statement
    block: $ => seq(
      '{',
      repeat($._statement),
      '}'
    ),

    expression_statement: $ => $._expression,

    parenthesized_expression: $ => seq('(', $._expression, ')'),

    argument_list: $ => sep1($._expression, ','),

    // Terminals
    identifier: $ => /[a-zA-Z_][a-zA-Z0-9_]*/,

    integer_literal: $ => /\\d+/,

    float_literal: $ => /\\d+\\.\\d+/,

    string_literal: $ => choice(
      seq('"', repeat(choice(/[^"\\\\]/, /\\\\./)), '"'),
      seq("'", repeat(choice(/[^'\\\\]/, /\\\\./)), "'")
    ),

    boolean_literal: $ => choice('true', 'false'),

    null_literal: $ => 'null'
  }
});

function sep1(rule, separator) {
  return seq(rule, repeat(seq(separator, rule)));
}