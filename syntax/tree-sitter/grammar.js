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
      $.function_declaration,
      $.expression_statement
    ),

    // Comments
    comment: $ => token(seq('//', /.*/)),
    block_comment: $ => token(seq('/*', /[^*]*\*+([^/*][^*]*\*+)*/, '/')),

    // Variables
    variable_declaration: $ => seq(
      choice('let', 'const'),
      field('name', $.identifier),
      optional(seq(':', field('type', $.identifier))),
      '=',
      field('value', $._expression)
    ),

    // Functions
    function_declaration: $ => seq(
      'fn',
      field('name', $.identifier),
      '(',
      optional($.parameter_list),
      ')',
      field('body', $.block)
    ),

    parameter_list: $ => sep1($.parameter, ','),

    parameter: $ => seq(
      field('name', $.identifier),
      optional(seq(':', field('type', $.identifier)))
    ),

    // Expressions
    _expression: $ => choice(
      $.binary_expression,
      $.call_expression,
      $.identifier,
      $.integer_literal,
      $.string_literal,
      $.boolean_literal,
      $.parenthesized_expression
    ),

    binary_expression: $ => choice(
      prec.left(10, seq($._expression, '+', $._expression)),
      prec.left(10, seq($._expression, '-', $._expression)),
      prec.left(11, seq($._expression, '*', $._expression)),
      prec.left(11, seq($._expression, '/', $._expression))
    ),

    call_expression: $ => prec(15, seq(
      field('function', $.identifier),
      '(',
      optional($.argument_list),
      ')'
    )),

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

    integer_literal: $ => /\d+/,

    string_literal: $ => seq('"', repeat(choice(/[^"\\]/, /\\./)), '"'),

    boolean_literal: $ => choice('true', 'false')
  }
});

function sep1(rule, separator) {
  return seq(rule, repeat(seq(separator, rule)));
}