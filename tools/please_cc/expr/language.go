package expr

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"text/scanner"

	"github.com/PaesslerAG/gval"

	"github.com/please-build/cc-rules/tools/please_cc/cctool"
)

// scanMode is a scanner mode that recognises identifiers, integers, characters and strings (excluding raw strings) in
// Go's grammar. These are the only token types used in please_cc's expression language.
var scanMode uint =
	scanner.ScanIdents |
	scanner.ScanInts |
	scanner.ScanChars |
	scanner.ScanStrings

var language = gval.NewLanguage(
	// Restrict the tokens recognised by the parser to those used by please_cc's expression language:
	gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
		p.SetMode(scanMode)
		return p.ParseExpression(ctx)
	}),

	gval.PrefixExtension(scanner.Int, parseVersion),
	gval.PrefixExtension(scanner.String, parseString),
	gval.PrefixExtension(scanner.Char, parseSingleQuotedString),
	gval.PrefixExtension(scanner.Ident, parseIdent),

	gval.PrefixExtension('[', parseArray),

	gval.PrefixOperator("!", negationOperator),
	gval.InfixOperator("||", boolOperator("||")),
	gval.InfixOperator("&&", boolOperator("&&")),

	gval.InfixOperator("==", compareVersions("==", []int{0})),
	gval.InfixOperator("!=", compareVersions("!=", []int{-1, 1})),
	gval.InfixOperator(">", compareVersions(">", []int{1})),
	gval.InfixOperator(">=", compareVersions(">=", []int{0, 1})),
	gval.InfixOperator("<", compareVersions("<", []int{-1})),
	gval.InfixOperator("<=", compareVersions("<=", []int{-1, 0})),

	gval.PrefixExtension('(', parseParentheses),

	gval.PostfixOperator("?", parseIf),

	gval.Precedence("||", 20),
	gval.Precedence("&&", 21),

	gval.Precedence("==", 40),
	gval.Precedence("!=", 40),
	gval.Precedence(">", 40),
	gval.Precedence(">=", 40),
	gval.Precedence("<", 40),
	gval.Precedence("<=", 40),
)

// ErrInvalidReturnType is an error indicating that an expression did not evaluate to a legal return type.
var ErrInvalidReturnType = errors.New("expression must evaluate to string or string array")

// Evaluate evaluates the string representation of an expression in the given environment (mapping identifiers to
// values) and returns a string slice representing the evaluation of the overall expression.
func Evaluate(expr string, env any) ([]string, error) {
	out, err := language.EvaluateWithContext(context.Background(), expr, env)
	if err != nil {
		return nil, err
	}
	tout := reflect.TypeOf(out)
	if tout == nil {
		return nil, ErrInvalidReturnType
	} else if tout.Kind() == reflect.String {
		return []string{out.(string)}, nil
	} else if tout.Kind() == reflect.Slice && tout.Elem().Kind() == reflect.String {
		return out.([]string), nil
	}
	return nil, ErrInvalidReturnType
}

func negationOperator(c context.Context, a any) (any, error) {
	if a == nil {
		return true, nil
	}
	ab, isBool := a.(bool)
	if isBool {
		return !ab, nil
	}
	_, isVer := a.(cctool.Version)
	if !isVer {
		return nil, errors.New("operand to ! must be a boolean expression or a version")
	}
	return false, nil
}

func boolOperator(operator string) func(any, any) (any, error) {
	return func(a, b any) (any, error) {
		ab, aIsBool := a.(bool)
		_, aIsVer := a.(cctool.Version)
		bb, bIsBool := b.(bool)
		_, bIsVer := b.(cctool.Version)
		if !aIsBool && !aIsVer && a != nil {
			return nil, fmt.Errorf("left-hand operand to %s must be a boolean expression or a version", operator)
		}
		if !bIsBool && !bIsVer && b != nil {
			return nil, fmt.Errorf("right-hand operand to %s must be a boolean expression or a version", operator)
		}
		aTest := (aIsBool && ab) || (aIsVer && a != nil)
		bTest := (bIsBool && bb) || (bIsVer && b != nil)
		if operator == "&&" {
			return aTest && bTest, nil
		}
		return aTest || bTest, nil
	}
}

func compareVersions(operator string, trueValues []int) func(any, any) (any, error) {
	return func(a, b any) (any, error) {
		aVer, aIsVer := a.(cctool.Version)
		bVer, bIsVer := b.(cctool.Version)
		if !aIsVer && a != nil {
			return nil, fmt.Errorf("left-hand operand to %s must be a version", operator)
		}
		if !bIsVer && b != nil {
			return nil, fmt.Errorf("right-hand operand to %s must be a version", operator)
		}
		// This is slightly hacky - if either of the operands are nil, (i.e. unbound identifiers), pretend they're the version
		// number -1 for the purpose of comparison with the other one. Since the parser doesn't permit negative integers as
		// part of version numbers, these "fake" version numbers will always compare lower than a genuine version number.
		if a == nil {
			aVer = cctool.Version{-1}
		}
		if b == nil {
			bVer = cctool.Version{-1}
		}
		return slices.Contains(trueValues, aVer.Compare(bVer)), nil
	}
}
