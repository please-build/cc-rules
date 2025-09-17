# please\_cc

please\_cc is a tool for invoking C/C++ compilers and linkers using command line options that are dynamically selected
depending on the compiler or linker being invoked.

please\_cc operates in one of two modes, indicated by the value of the first argument passed to it:

- `cc` ("compiler mode") invokes the compiler - and, implicitly, the linker - in a C or C++ toolchain.
- `ld` ("linker mode") invokes a linker independently of a C or C++ toolchain.

The second argument identifies the (relative or absolute) path to the compiler or linker. If the path is relative,
please\_cc follows the usual conventions for converting relative paths into absolute paths before invoking the compiler
or linker.

The remaining arguments are command line options to be passed through to the compiler or linker. Arguments surrounded
by `{{ ` and ` }}` are evaluated as [expressions](#Expressions) after those delimiters have been stripped; the result of
each evaluation is substituted into the list of command line options passed to the compiler or linker. The identity of
the compiler and/or linker is described in the environment in which the expression is evaluated, which allows for
different command line options to be passed through for different compilers and/or linkers.

# Expressions

please\_cc features an expression language that allows any number of compiler or linter command line options to be
substituted in place of an expression specified on please\_cc's command line.

## Types and identifiers

The language consists of five types:

- The nil type.
- Version numbers, represented in [dot-decimal notation](https://en.wikipedia.org/wiki/Dot-decimal_notation#Version_numbers)
  (e.g. `1`, `1.5`, `14.0.6`).
- Booleans.
- Strings, delimited with either `'` or `"` (e.g. `'single'`, `"double"`).
- String arrays, delimited with `,` and enclosed by `[` and `]` (e.g. `['single', "double"]`).

Identifiers are sequences of letters or digits (e.g. `gcc`, `ld64`) and may have nil or version number values. The
expression's environment is populated with identifiers representing the various compilers and linkers supported by
please\_cc; a version number value for one of these identifiers indicates that please\_cc detected the presence of that
version of that compiler/linker during invocation, while a nil value indicates that please\_cc did not detect the
presence of that compiler/linker. The identifier by which each supported compiler or linker is known is listed in the
[tool compatibility table](#Compatibility).

The language supports limited type coercion: the nil type evaluates to false in boolean contexts, and version numbers
evaluate to true in boolean contexts.

## Expression evaluation

Intermediate expressions may evaluate to any of the types above, although it is not possible to explicitly refer to
all of them - there are no `nil`, `true` or `false` constants, for example.

The overall expression must evaluate to either a string (in which case a single option is substituted into the
arguments list for the compiler/linker invocation in place of the expression), or a string array (in which case a
sequence of options is substituted, in the given order). If the overall expression evaluates to the empty string
array (`[]`), no arguments are substituted in place of the expression.

## Operators

The language consists of three types of operators:

- The comparison operators - `==`, `!=`, `<`, `<=`, `>`, and `>=` - which compare two version numbers (or the nil value)
  and return a boolean value. Version numbers are compared dot-wise; if version numbers of different lengths are
  compared with each other, the missing numbers in the shorter one are assumed to be zeroes - e.g., `1.2 > 1.1.3`,
  `1.2 == 1.2.0` and `1.2 < 1.2.3` all return true. The nil type is equal to itself, and is always less than any version
  number.
- The logical operators - `!`, `&&`, and `||` - which consume expressions as operands and return boolean values.
- The ternary operator - `? :` - which evaluates the expression on the left-hand side of `?` and returns the
  evaluation of the expression on the right-hand side of `?` if true and the evaluation of the expression on the
  right-hand side of `:` if false. If the false branch is omitted, it implicitly evaluates to the empty string array.

Operators earlier in the list above bind more tightly than those later in the list. This order of precedence can be
overridden with the use of parentheses (`()`).

# Compatibility

please\_cc is compatible with all operating systems and architectures with which Please itself is compatible, and is
known to correctly identify the following compilers and linkers:

| Tool        | Minimum supported version | Expression language identifier |
| ----------- | ------------------------- | ------------------------------ |
| GCC         | 9                         | `gcc`                          |
| Clang       | 11                        | `clang`                        |
| Apple Clang | 12                        | `aclang`                       |
| GNU ld      | 2.38                      | `gnuld`                        |
| GNU gold    | 1.15                      | `gold`                         |
| LLD         | 11                        | `lld`                          |
| ld64        | 609.8                     | `ld64`                         |
| Apple ld    | 1015.7                    | `appleld`                      |

# Examples

- Only make warnings fatal when compiling with GCC:

  ```sh
  please_cc cc cctool -o example '{{ gcc ? "-Werror" }}' example.c
  ```
- Compile C++ modules with GCC and Clang <= 15 with the `-fmodules-ts` option, or Clang >= 16 with the (equivalent,
  non-deprecated) `-std=c++20` option:

  ```sh
  please_cc cc cctool -o example '{{ gcc || (clang && clang <= 15) ? "-fmodules-ts" : "-std=c++20" }}' example.cc
  ```
- Link objects using the old ld64 code path if linking with Apple's new ld linker (enabled using `-ld64` prior to
  Xcode 15.1, and `-ld_classic` from Xcode 15.1 onwards):

  ```sh
  please_cc ld ldtool '{{ appleld ? (appleld >= 1022.1 ? "-ld_classic" : "-ld64") }}' obj1.o obj2.o -o example
  ```

# Tests

Test cases for all supported operating systems, compilers and linkers can be found in `cctool/test_data/`. The
`.test_data` files in this tree were automatically generated by `//tools/please_cc/cctool:generate_test_data`, which
consumes a JSON-encoded toolspec file describing the compilers and linkers for which test cases should be generated:

```sh
plz run //tools/please_cc/cctool:generate_test_data [-o OUTPUT_DIR] [TOOLSPEC_PATH]
```

The toolspec files in `cctool/test_data/` are designed to be used on specific platforms, as documented at the top of
each respective file.
