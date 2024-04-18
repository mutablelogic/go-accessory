/*
 * ast package builds abstract syntax trees (ast) from tokens
 */
package ast

/*
  Current grammar:
    <expr> := number | string | ident | <unary_op> <expr> | <expr> <binary_op> <expr> | '(' <expr> ')' | func
    <func> := ident '(' <expr>, ... ')'
    <unary_op> := '!' | 'NOT'
    <binary_op> := '+' | '-' | '*' | '/' | '==' | '!=' | '<' | '<=' | '>' | '>=' | '~' | '!~' | 'OR' | 'AND'
*/
