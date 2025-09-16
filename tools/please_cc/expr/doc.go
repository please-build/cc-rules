// Package expr implements please_cc's expression language.
//
// # Types and identifiers
//
// The language consists of five types:
//
// - The nil type.
// - [Version] numbers, represented in [dot-decimal notation] (e.g. `1`, `1.5`, `14.0.6`).
// - Booleans.
// - Strings, delimited with either `'` or `"` (e.g. `'single'`, `"double"`).
// - String arrays, delimited with `,` and enclosed by `[` and `]` (e.g. `['single', "double"]`).
//
// Identifiers are sequences of one or more letters or digits (e.g. `gcc`, `ld64`) and may have nil or version number
// values. The expression's environment is populated with identifiers representing the various compilers and linkers
// supported by please_cc; a version number value for one of these identifiers indicates that please_cc detected the
// presence of that version of that compiler/linker during invocation, while a nil value indicates that please_cc did
// not detect the presence of that compiler/linker.
//
// The language supports limited type coercion: the nil type evaluates to false in boolean contexts.
//
// # Expressions
//
// Intermediate expressions may evaluate to any of the types above, although it is not possible to explicitly refer to
// all of them - there are no `nil`, `true` or `false` constants, for example.
//
// The overall expression must evaluate to either a string (in which case a single option is substituted into the
// arguments list for the compiler/linker invocation in place of the expression), or a string array (in which case a
// sequence of options is substituted, in the given order). If the overall expression evaluates to the empty string
// array (`[]`), no arguments are substituted in place of the expression.
//
// # Operators
//
// The language consists of three types of operators:
//
// - The comparison operators - `==`, `!=`, `<`, `<=`, `>`, and `>=` - which compare version numbers according to the
//   rules described in [Version.Compare] and return a boolean value. The nil type is equal to itself, and is always
//   less than any version number.
// - The logical operators - `!`, `&&`, and `||` - which consume expressions as operands and return boolean values.
// - The ternary operator - `? :` - which evaluates the expression on the left-hand side of `?` and returns the
//   evaluation of the expression on the right-hand side of `?` if true and the evaluation of the expression on the
//   right-hand side of `:` if false. If the false branch is omitted, it implicitly evaluates to the empty array.
//
// Operators earlier in the list above bind more tightly than those later in the list. This order of precedence can be
// overridden with the use of parentheses (`()`).
//
// [Dot-decimal notation]: https://en.wikipedia.org/wiki/Dot-decimal_notation#Version_numbers
package expr
