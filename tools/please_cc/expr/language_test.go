package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/please-build/cc-rules/tools/please_cc/cctool"
)

func TestEvaluate(t *testing.T) {
	env := map[string]any{
		"a": cctool.Version{1, 2, 3},
		"b": cctool.Version{4, 5},
		"c": cctool.Version{1, 2, 3},
	}

	tests := map[string]map[string]struct {
		Expr     string
		Err      string
		Expected []string
	}{
		"Invalid return types": {
			"Empty expression": {
				Expr: "",
				Err:  "unexpected EOF while scanning extensions",
			},
			"Version": {
				Expr: "1.3.0",
				Err:  "expression must evaluate to string or string array",
			},
			"Empty string, single-quoted": {
				Expr: `''`,
				Err:  "string literal cannot be empty",
			},
			"Empty string, double-quoted": {
				Expr: `""`,
				Err:  "string literal cannot be empty",
			},
			"Bound identifier": {
				Expr: "a",
				Err:  "expression must evaluate to string or string array",
			},
			"Unbound identifier": {
				Expr: "z",
				Err:  "expression must evaluate to string or string array",
			},
			"True branch of ? operator evaluates to bound identifier": {
				Expr: "1 > 0 ? b : 'n'",
				Err:  "expression must evaluate to string or string array",
			},
			"True branch of ? operator evaluates to unbound identifier": {
				Expr: "1 > 0 ? z : 'n'",
				Err:  "expression must evaluate to string or string array",
			},
			"False branch of ? operator evaluates to bound identifier": {
				Expr: "0 > 1 ? 'y' : b",
				Err:  "expression must evaluate to string or string array",
			},
			"False branch of ? operator evaluates to unbound identifier": {
				Expr: "0 > 1 ? 'y' : z",
				Err:  "expression must evaluate to string or string array",
			},
		},
		"Constant return values": {
			"String, single-quoted": {
				Expr:     `'-an-option'`,
				Expected: []string{"-an-option"},
			},
			"String, double-quoted": {
				Expr:     `"-another-option"`,
				Expected: []string{"-another-option"},
			},
			"Zero-element string array": {
				Expr:     "[]",
				Expected: []string{},
			},
			"String array, mixed quotes": {
				Expr:     `['-an-option', "-another-option"]`,
				Expected: []string{"-an-option", "-another-option"},
			},
		},
		"Identifiers": {
			"Bound identifier evaluates to true": {
				Expr:     `a ? 'y' : 'n'`,
				Expected: []string{"y"},
			},
			"Unbound identifier evaluates to false": {
				Expr:     `z ? 'y' : 'n'`,
				Expected: []string{"n"},
			},
		},
		"== operator": {
			"Left-hand operand not a version, right-hand operand is a version": {
				Expr: "'2' == 2 ? 'y' : 'n'",
				Err:  "left-hand operand to == must be a version",
			},
			"Left-hand operand is a version, right-hand operand not a version": {
				Expr: "2 == '2' ? 'y' : 'n'",
				Err:  "right-hand operand to == must be a version",
			},
			"Left-hand operand not a version, right-hand operand is a bound identifier": {
				Expr: "'2' == a ? 'y' : 'n'",
				Err:  "left-hand operand to == must be a version",
			},
			"Left-hand operand is a bound identifier, right-hand operand not a version": {
				Expr: "a == '2' ? 'y' : 'n'",
				Err:  "right-hand operand to == must be a version",
			},
			"Left-hand operand not a version, right-hand operand is an unbound identifier": {
				Expr: "'2' == z ? 'y' : 'n'",
				Err:  "left-hand operand to == must be a version",
			},
			"Left-hand operand is an unbound identifier, right-hand operand not a version": {
				Expr: "z == '2' ? 'y' : 'n'",
				Err:  "right-hand operand to == must be a version",
			},
			"Equal (zero) versions, one component": {
				Expr:     "0 == 0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal (zero) versions, multiple components": {
				Expr:     "0.0.0 == 0.0.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal (non-zero) versions, one component": {
				Expr:     "3 == 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal (non-zero) versions, multiple components": {
				Expr:     "3.2.1 == 3.2.1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal versions, implied zero components": {
				Expr:     "3 == 3.0.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Unequal versions, implied zero components": {
				Expr:     "3 == 3.3.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Unequal versions, one component": {
				Expr:     "2 == 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Unequal versions, multiple components": {
				Expr:     "2.0.1 == 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version equal to identifier": {
				Expr:     "a == 1.2.3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version unequal to identifier": {
				Expr:     "a == 6.7.8 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal identifiers": {
				Expr:     "a == c ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Unequal identifiers": {
				Expr:     "a == b ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Left-hand operand is an unbound identifier": {
				Expr:     "z == 5.5 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Right-hand operand is an unbound identifier": {
				Expr:     "5.5 == z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
		},
		"!= operator": {
			"Left-hand operand not a version, right-hand operand is a version": {
				Expr: "'2' != 2 ? 'y' : 'n'",
				Err:  "left-hand operand to != must be a version",
			},
			"Left-hand operand is a version, right-hand operand not a version": {
				Expr: "2 != '2' ? 'y' : 'n'",
				Err:  "right-hand operand to != must be a version",
			},
			"Left-hand operand not a version, right-hand operand is a bound identifier": {
				Expr: "'2' != a ? 'y' : 'n'",
				Err:  "left-hand operand to != must be a version",
			},
			"Left-hand operand is a bound identifier, right-hand operand not a version": {
				Expr: "a != '2' ? 'y' : 'n'",
				Err:  "right-hand operand to != must be a version",
			},
			"Left-hand operand not a version, right-hand operand is an unbound identifier": {
				Expr: "'2' != z ? 'y' : 'n'",
				Err:  "left-hand operand to != must be a version",
			},
			"Left-hand operand is an unbound identifier, right-hand operand not a version": {
				Expr: "z != '2' ? 'y' : 'n'",
				Err:  "right-hand operand to != must be a version",
			},
			"Equal (zero) versions, one component": {
				Expr:     "0 != 0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal (zero) versions, multiple components": {
				Expr:     "0.0.0 != 0.0.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal (non-zero) versions, one component": {
				Expr:     "3 != 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal (non-zero) versions, multiple components": {
				Expr:     "3.2.1 != 3.2.1 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal versions, implied zero components": {
				Expr:     "3 != 3.0.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Unequal versions, implied zero components": {
				Expr:     "3 != 3.3.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Unequal versions, one component": {
				Expr:     "2 != 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Unequal versions, multiple components": {
				Expr:     "2.0.1 != 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version equal to identifier": {
				Expr:     "a != 1.2.3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version unequal to identifier": {
				Expr:     "a != 6.7.8 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal identifiers": {
				Expr:     "a != c ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Unequal identifiers": {
				Expr:     "a != b ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is an unbound identifier": {
				Expr:     "z != 5.5 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Right-hand operand is an unbound identifier": {
				Expr:     "5.5 != z ? 'y' : 'n'",
				Expected: []string{"y"},
			},
		},
		"> operator": {
			"Left-hand operand not a version, right-hand operand is a version": {
				Expr: "'2' > 2 ? 'y' : 'n'",
				Err:  "left-hand operand to > must be a version",
			},
			"Left-hand operand is a version, right-hand operand not a version": {
				Expr: "2 > '2' ? 'y' : 'n'",
				Err:  "right-hand operand to > must be a version",
			},
			"Left-hand operand not a version, right-hand operand is a bound identifier": {
				Expr: "'2' > a ? 'y' : 'n'",
				Err:  "left-hand operand to > must be a version",
			},
			"Left-hand operand is a bound identifier, right-hand operand not a version": {
				Expr: "a > '2' ? 'y' : 'n'",
				Err:  "right-hand operand to > must be a version",
			},
			"Left-hand operand not a version, right-hand operand is an unbound identifier": {
				Expr: "'2' > z ? 'y' : 'n'",
				Err:  "left-hand operand to > must be a version",
			},
			"Left-hand operand is an unbound identifier, right-hand operand not a version": {
				Expr: "z > '2' ? 'y' : 'n'",
				Err:  "right-hand operand to > must be a version",
			},
			"Equal (zero) versions, one component": {
				Expr:     "0 > 0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal (zero) versions, multiple components": {
				Expr:     "0.0.0 > 0.0.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal (non-zero) versions, one component": {
				Expr:     "3 > 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal (non-zero) versions, multiple components": {
				Expr:     "3.2.1 > 3.2.1 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal versions, implied zero components": {
				Expr:     "3 > 3.0.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"First version greater, implied zero components": {
				Expr:     "4 > 3.3.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"First version greater, one component": {
				Expr:     "4 > 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"First version greater, multiple components": {
				Expr:     "4.0.1 > 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Second version greater, implied zero components": {
				Expr:     "3 > 3.3.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Second version greater, one component": {
				Expr:     "2 > 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Second version greater, multiple components": {
				Expr:     "2.0.1 > 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version equal to identifier, identifier first": {
				Expr:     "a > 1.2.3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Identifier greater than version, identifier first": {
				Expr:     "a > 0.1.5 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version greater than identifier, identifier first": {
				Expr:     "a > 7.8.9 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version equal to identifier, version first": {
				Expr:     "1.2.3 > a ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Identifier greater than version, version first": {
				Expr:     "1.2.3 > b ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version greater than identifier, version first": {
				Expr:     "7.8.9 > a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal identifiers": {
				Expr:     "a > c ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"First identifier greater than second identifier": {
				Expr:     "b > a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Second identifier greater than first identifier": {
				Expr:     "c > b ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Left-hand operand is an unbound identifier": {
				Expr:     "z > 5.5 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Right-hand operand is an unbound identifier": {
				Expr:     "5.5 > z ? 'y' : 'n'",
				Expected: []string{"y"},
			},
		},
		">= operator": {
			"Left-hand operand not a version, right-hand operand is a version": {
				Expr: "'2' >= 2 ? 'y' : 'n'",
				Err:  "left-hand operand to >= must be a version",
			},
			"Left-hand operand is a version, right-hand operand not a version": {
				Expr: "2 >= '2' ? 'y' : 'n'",
				Err:  "right-hand operand to >= must be a version",
			},
			"Left-hand operand not a version, right-hand operand is a bound identifier": {
				Expr: "'2' >= a ? 'y' : 'n'",
				Err:  "left-hand operand to >= must be a version",
			},
			"Left-hand operand is a bound identifier, right-hand operand not a version": {
				Expr: "a >= '2' ? 'y' : 'n'",
				Err:  "right-hand operand to >= must be a version",
			},
			"Left-hand operand not a version, right-hand operand is an unbound identifier": {
				Expr: "'2' >= z ? 'y' : 'n'",
				Err:  "left-hand operand to >= must be a version",
			},
			"Left-hand operand is an unbound identifier, right-hand operand not a version": {
				Expr: "z >= '2' ? 'y' : 'n'",
				Err:  "right-hand operand to >= must be a version",
			},
			"Equal (zero) versions, one component": {
				Expr:     "0 >= 0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal (zero) versions, multiple components": {
				Expr:     "0.0.0 >= 0.0.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal (non-zero) versions, one component": {
				Expr:     "3 >= 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal (non-zero) versions, multiple components": {
				Expr:     "3.2.1 >= 3.2.1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal versions, implied zero components": {
				Expr:     "3 >= 3.0.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"First version greater, implied zero components": {
				Expr:     "4 >= 3.3.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"First version greater, one component": {
				Expr:     "4 >= 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"First version greater, multiple components": {
				Expr:     "4.0.1 >= 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Second version greater, implied zero components": {
				Expr:     "3 >= 3.3.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Second version greater, one component": {
				Expr:     "2 >= 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Second version greater, multiple components": {
				Expr:     "2.0.1 >= 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version equal to identifier, identifier first": {
				Expr:     "a >= 1.2.3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Identifier greater than version, identifier first": {
				Expr:     "a >= 0.1.5 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version greater than identifier, identifier first": {
				Expr:     "a >= 7.8.9 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version equal to identifier, version first": {
				Expr:     "1.2.3 >= a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Identifier greater than version, version first": {
				Expr:     "1.2.3 >= b ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version greater than identifier, version first": {
				Expr:     "7.8.9 >= a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal identifiers": {
				Expr:     "a >= c ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"First identifier greater than second identifier": {
				Expr:     "b >= a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Second identifier greater than first identifier": {
				Expr:     "c >= b ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Left-hand operand is an unbound identifier": {
				Expr:     "z >= 5.5 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Right-hand operand is an unbound identifier": {
				Expr:     "5.5 >= z ? 'y' : 'n'",
				Expected: []string{"y"},
			},
		},
		"< operator": {
			"Left-hand operand not a version, right-hand operand is a version": {
				Expr: "'2' < 2 ? 'y' : 'n'",
				Err:  "left-hand operand to < must be a version",
			},
			"Left-hand operand is a version, right-hand operand not a version": {
				Expr: "2 < '2' ? 'y' : 'n'",
				Err:  "right-hand operand to < must be a version",
			},
			"Left-hand operand not a version, right-hand operand is a bound identifier": {
				Expr: "'2' < a ? 'y' : 'n'",
				Err:  "left-hand operand to < must be a version",
			},
			"Left-hand operand is a bound identifier, right-hand operand not a version": {
				Expr: "a < '2' ? 'y' : 'n'",
				Err:  "right-hand operand to < must be a version",
			},
			"Left-hand operand not a version, right-hand operand is an unbound identifier": {
				Expr: "'2' < z ? 'y' : 'n'",
				Err:  "left-hand operand to < must be a version",
			},
			"Left-hand operand is an unbound identifier, right-hand operand not a version": {
				Expr: "z < '2' ? 'y' : 'n'",
				Err:  "right-hand operand to < must be a version",
			},
			"Equal (zero) versions, one component": {
				Expr:     "0 < 0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal (zero) versions, multiple components": {
				Expr:     "0.0.0 < 0.0.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal (non-zero) versions, one component": {
				Expr:     "3 < 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal (non-zero) versions, multiple components": {
				Expr:     "3.2.1 < 3.2.1 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal versions, implied zero components": {
				Expr:     "3 < 3.0.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"First version greater, implied zero components": {
				Expr:     "4 < 3.3.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"First version greater, one component": {
				Expr:     "4 < 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"First version greater, multiple components": {
				Expr:     "4.0.1 < 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Second version greater, implied zero components": {
				Expr:     "3 < 3.3.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Second version greater, one component": {
				Expr:     "2 < 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Second version greater, multiple components": {
				Expr:     "2.0.1 < 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version equal to identifier, identifier first": {
				Expr:     "a < 1.2.3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Identifier greater than version, identifier first": {
				Expr:     "a < 0.1.5 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version greater than identifier, identifier first": {
				Expr:     "a < 7.8.9 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version equal to identifier, version first": {
				Expr:     "1.2.3 < a ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Identifier greater than version, version first": {
				Expr:     "1.2.3 < b ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version greater than identifier, version first": {
				Expr:     "7.8.9 < a ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal identifiers": {
				Expr:     "a < c ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"First identifier greater than second identifier": {
				Expr:     "b < a ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Second identifier greater than first identifier": {
				Expr:     "c < b ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is an unbound identifier": {
				Expr:     "z < 5.5 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Right-hand operand is an unbound identifier": {
				Expr:     "5.5 < z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
		},
		"<= operator": {
			"Left-hand operand not a version, right-hand operand is a version": {
				Expr: "'2' <= 2 ? 'y' : 'n'",
				Err:  "left-hand operand to <= must be a version",
			},
			"Left-hand operand is a version, right-hand operand not a version": {
				Expr: "2 <= '2' ? 'y' : 'n'",
				Err:  "right-hand operand to <= must be a version",
			},
			"Left-hand operand not a version, right-hand operand is a bound identifier": {
				Expr: "'2' <= a ? 'y' : 'n'",
				Err:  "left-hand operand to <= must be a version",
			},
			"Left-hand operand is a bound identifier, right-hand operand not a version": {
				Expr: "a <= '2' ? 'y' : 'n'",
				Err:  "right-hand operand to <= must be a version",
			},
			"Left-hand operand not a version, right-hand operand is an unbound identifier": {
				Expr: "'2' <= z ? 'y' : 'n'",
				Err:  "left-hand operand to <= must be a version",
			},
			"Left-hand operand is an unbound identifier, right-hand operand not a version": {
				Expr: "z <= '2' ? 'y' : 'n'",
				Err:  "right-hand operand to <= must be a version",
			},
			"Equal (zero) versions, one component": {
				Expr:     "0 <= 0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal (zero) versions, multiple components": {
				Expr:     "0.0.0 <= 0.0.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal (non-zero) versions, one component": {
				Expr:     "3 <= 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal (non-zero) versions, multiple components": {
				Expr:     "3.2.1 <= 3.2.1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Equal versions, implied zero components": {
				Expr:     "3 <= 3.0.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"First version greater, implied zero components": {
				Expr:     "4 <= 3.3.0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"First version greater, one component": {
				Expr:     "4 <= 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"First version greater, multiple components": {
				Expr:     "4.0.1 <= 3 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Second version greater, implied zero components": {
				Expr:     "3 <= 3.3.0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Second version greater, one component": {
				Expr:     "2 <= 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Second version greater, multiple components": {
				Expr:     "2.0.1 <= 3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version equal to identifier, identifier first": {
				Expr:     "a <= 1.2.3 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Identifier greater than version, identifier first": {
				Expr:     "a <= 0.1.5 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Version greater than identifier, identifier first": {
				Expr:     "a <= 7.8.9 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version equal to identifier, version first": {
				Expr:     "1.2.3 <= a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Identifier greater than version, version first": {
				Expr:     "1.2.3 <= b ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Version greater than identifier, version first": {
				Expr:     "7.8.9 <= a ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Equal identifiers": {
				Expr:     "a <= c ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"First identifier greater than second identifier": {
				Expr:     "b <= a ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Second identifier greater than first identifier": {
				Expr:     "c <= b ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is an unbound identifier": {
				Expr:     "z <= 5.5 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Right-hand operand is an unbound identifier": {
				Expr:     "5.5 <= z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
		},
		"! operator": {
			"Operand not a version": {
				Expr: "!'3.3'",
				Err:  "operand to ! must be a boolean expression or a version",
			},
			"Negation of bound identifier": {
				Expr:     "!a ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Negation of unbound identifier": {
				Expr:     "!z ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Double negation of bound identifier": {
				Expr:     "!!a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Double negation of unbound identifier": {
				Expr:     "!!z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Negation of zero version": {
				Expr:     "!0 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Double negation of zero version": {
				Expr:     "!!0 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Negation of non-zero version": {
				Expr:     "!1.4.5 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Double negation of non-zero version": {
				Expr:     "!!1.4.5 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
		},
		"|| operator": {
			"Left-hand operand not a version, right-hand operand is a version": {
				Expr: "'2' || 2 ? 'y' : 'n'",
				Err:  "left-hand operand to || must be a boolean expression or a version",
			},
			"Left-hand operand is a version, right-hand operand not a version": {
				Expr: "2 || '2' ? 'y' : 'n'",
				Err:  "right-hand operand to || must be a boolean expression or a version",
			},
			"Left-hand operand not a version, right-hand operand is a bound identifier": {
				Expr: "'2' || a ? 'y' : 'n'",
				Err:  "left-hand operand to || must be a boolean expression or a version",
			},
			"Left-hand operand is a bound identifier, right-hand operand not a version": {
				Expr: "a || '2' ? 'y' : 'n'",
				Err:  "right-hand operand to || must be a boolean expression or a version",
			},
			"Left-hand operand not a version, right-hand operand is an unbound identifier": {
				Expr: "'2' || z ? 'y' : 'n'",
				Err:  "left-hand operand to || must be a boolean expression or a version",
			},
			"Left-hand operand is an unbound identifier, right-hand operand not a version": {
				Expr: "z || '2' ? 'y' : 'n'",
				Err:  "right-hand operand to || must be a boolean expression or a version",
			},
			"Both operands are versions": {
				Expr:     "0 || 1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is bound identifier": {
				Expr:     "a || 1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Right-hand operand is bound identifier": {
				Expr:     "2 || a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Both operands are bound identifiers": {
				Expr:     "a || b ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is unbound identifier": {
				Expr:     "z || 1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Right-hand operand is unbound identifier": {
				Expr:     "2 || z ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Both operands are unbound identifiers": {
				Expr:     "y || z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Both operands are true boolean expressions": {
				Expr:     "1 > 0 || 2 > 1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is true boolean expression, right-hand operand is false boolean expression": {
				Expr:     "1 > 0 || 2 < 1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is false boolean expression, right-hand operand is true boolean expression": {
				Expr:     "1 < 0 || 2 > 1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Both operands are false boolean expressions": {
				Expr:     "1 < 0 || 2 < 1 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Chained operators, true || true || true": {
				Expr:     "a || b || c ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Chained operators, true || true || false": {
				Expr:     "a || b || z ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Chained operators, false || true || true": {
				Expr:     "z || b || c ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Chained operators, false || true || false": {
				Expr:     "y || b || z ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Chained operators, false || false || false": {
				Expr:     "x || y || z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
		},
		"&& operator": {
			"Left-hand operand not a version, right-hand operand is a version": {
				Expr: "'2' && 2 ? 'y' : 'n'",
				Err:  "left-hand operand to && must be a boolean expression or a version",
			},
			"Left-hand operand is a version, right-hand operand not a version": {
				Expr: "2 && '2' ? 'y' : 'n'",
				Err:  "right-hand operand to && must be a boolean expression or a version",
			},
			"Left-hand operand not a version, right-hand operand is a bound identifier": {
				Expr: "'2' && a ? 'y' : 'n'",
				Err:  "left-hand operand to && must be a boolean expression or a version",
			},
			"Left-hand operand is a bound identifier, right-hand operand not a version": {
				Expr: "a && '2' ? 'y' : 'n'",
				Err:  "right-hand operand to && must be a boolean expression or a version",
			},
			"Left-hand operand not a version, right-hand operand is an unbound identifier": {
				Expr: "'2' && z ? 'y' : 'n'",
				Err:  "left-hand operand to && must be a boolean expression or a version",
			},
			"Left-hand operand is an unbound identifier, right-hand operand not a version": {
				Expr: "z && '2' ? 'y' : 'n'",
				Err:  "right-hand operand to && must be a boolean expression or a version",
			},
			"Both operands are versions": {
				Expr:     "0 && 1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is bound identifier": {
				Expr:     "a && 1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Right-hand operand is bound identifier": {
				Expr:     "2 && a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Both operands are bound identifiers": {
				Expr:     "a && b ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is unbound identifier": {
				Expr:     "z && 1 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Right-hand operand is unbound identifier": {
				Expr:     "2 && z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Both operands are unbound identifiers": {
				Expr:     "y && z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Both operands are true boolean expressions": {
				Expr:     "1 > 0 && 2 > 1 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Left-hand operand is true boolean expression, right-hand operand is false boolean expression": {
				Expr:     "1 > 0 && 2 < 1 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Left-hand operand is false boolean expression, right-hand operand is true boolean expression": {
				Expr:     "1 < 0 && 2 > 1 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Both operands are false boolean expressions": {
				Expr:     "1 < 0 && 2 < 1 ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Chained operators, true && true && true": {
				Expr:     "a && b && c ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Chained operators, true && true && false": {
				Expr:     "a && b && z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Chained operators, false && true && true": {
				Expr:     "z && b && c ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Chained operators, false && true && false": {
				Expr:     "y && b && z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Chained operators, false && false && false": {
				Expr:     "x && y && z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
		},
		"? operator": {
			"Implicit false branch evaluates to empty array": {
				Expr:     "0 > 1 ? 'y'",
				Expected: []string{},
			},
			"Chained operators, true ? ... : true ? ... : ...": {
				Expr:     "a ? 'y1' : b ? 'y2' : 'n2'",
				Expected: []string{"y1"},
			},
			"Chained operators, true ? ... : false ? ... : ...": {
				Expr:     "a ? 'y1' : z ? 'y2' : 'n2'",
				Expected: []string{"y1"},
			},
			"Chained operators, false ? ... : true ? ... : ...": {
				Expr:     "y ? 'y1' : a ? 'y2' : 'n2'",
				Expected: []string{"y2"},
			},
			"Chained operators, false ? ... : false ? ... : ...": {
				Expr:     "y ? 'y1' : z ? 'y2' : 'n2'",
				Expected: []string{"n2"},
			},
		},
		"Order of precedence": {
			"! over &&": {
				Expr:     "!z && a ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"! over ||": {
				Expr:     "!y || z ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"&& over ||": {
				Expr:     "y || a && z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"|| over ==": {
				Expr:     "a == 1.2.3 || z == 5 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"&& over ==": {
				Expr:     "a == 1.2.3 && b == 4.5 ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Parentheses give && precedence over !": {
				Expr:     "!(y && z) ? 'y' : 'n'",
				Expected: []string{"y"},
			},
			"Parentheses give || precedence over !": {
				Expr:     "!(a || z) ? 'y' : 'n'",
				Expected: []string{"n"},
			},
			"Parentheses give || precedence over &&": {
				Expr:     "(y || a) && z ? 'y' : 'n'",
				Expected: []string{"n"},
			},
		},
	}

	for category, tests := range tests {
		t.Run(category, func(t *testing.T) {
			for desc, tc := range tests {
				t.Run(desc, func(t *testing.T) {
					actual, err := Evaluate(tc.Expr, env)
					if tc.Err == "" {
						assert.EqualValues(t, tc.Expected, actual, "Expression evaluates to expected value")
						assert.NoError(t, err, "No error returned")
					} else {
						assert.Nil(t, actual, "No value returned")
						assert.ErrorContains(t, err, tc.Err, "Expected error returned")
					}
				})
			}
		})
	}
}
