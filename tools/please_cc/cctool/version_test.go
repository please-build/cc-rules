package cctool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	for desc, tc := range map[string]struct {
		One      Version
		Two      Version
		Expected int
	}{
		"Both empty": {
			One:      Version{},
			Two:      Version{},
			Expected: 0,
		},
		"One empty": {
			One:      Version{},
			Two:      Version{1},
			Expected: -1,
		},
		"Two empty": {
			One:      Version{1},
			Two:      Version{},
			Expected: 1,
		},
		"Same length, One greater by last component": {
			One:      Version{1, 5, 6},
			Two:      Version{1, 5, 5},
			Expected: 1,
		},
		"Same length, Two greater by last component": {
			One:      Version{1, 5, 6},
			Two:      Version{1, 5, 8},
			Expected: -1,
		},
		"Same length, One greater by non-last component": {
			One:      Version{1, 6, 5},
			Two:      Version{1, 5, 5},
			Expected: 1,
		},
		"Same length, Two greater by non-last component": {
			One:      Version{1, 5, 8},
			Two:      Version{1, 7, 8},
			Expected: -1,
		},
		"One longer, One greater by non-last component": {
			One:      Version{1, 6, 5, 3},
			Two:      Version{1, 5, 5},
			Expected: 1,
		},
		"One longer, Two greater by non-last component": {
			One:      Version{1, 6, 5, 3},
			Two:      Version{1, 7, 5},
			Expected: -1,
		},
		"Two longer, One greater by non-last component": {
			One:      Version{1, 8, 8},
			Two:      Version{1, 6, 8, 3},
			Expected: 1,
		},
		"Two longer, Two greater by non-last component": {
			One:      Version{1, 5, 8},
			Two:      Version{1, 7, 8, 3},
			Expected: -1,
		},
		"Same length, identical": {
			One:      Version{4, 6, 3, 2},
			Two:      Version{4, 6, 3, 2},
			Expected: 0,
		},
		"One longer, same version": {
			One:      Version{4, 6, 3, 0},
			Two:      Version{4, 6, 3},
			Expected: 0,
		},
		"Two longer, same version": {
			One:      Version{4, 6, 3},
			Two:      Version{4, 6, 3, 0},
			Expected: 0,
		},
	} {
		assert.Equal(t, tc.Expected, tc.One.Compare(tc.Two), desc)
	}
}
