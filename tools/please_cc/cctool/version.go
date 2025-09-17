package cctool

import (
	"strconv"
	"strings"
)

// Version represents a version number in [dot-decimal notation].
//
// [Dot-decimal notation]: https://en.wikipedia.org/wiki/Dot-decimal_notation#Version_numbers
type Version []int64

// MustParseVersion returns a Version representing the given version number string. It panics if given the empty string
// or a string that is not a version number in dot-decimal notation.
func MustParseVersion(version string) Version {
	if version == "" {
		panic("Version number must not be an empty string")
	}
	ver := Version{}
	for _, num := range strings.Split(version, ".") {
		n, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			panic("Invalid version number: " + version)
		}
		ver = append(ver, n)
	}
	return ver
}

// Compare compares this Version with another. It returns -1 if the version number represented by this Version is less
// than that of the other one, 1 if it is greater than the other one, or 0 if they are equivalent.
//
// If one of the version numbers contains more dots than the other, the missing numbers in the shorter of the two
// version numbers are assumed to be zeroes for comparison purposes - so, for example:
//
// - 1.2 > 1.1.3
// - 1.2 = 1.2.0
// - 1.2 < 1.2.3
func (v Version) Compare(to Version) int {
	vlen, tolen := len(v), len(to)
	var vel, toel int64
	end := max(vlen, tolen)
	for i := 0; i < end; i++ {
		vel, toel = 0, 0
		if i < vlen {
			vel = v[i]
		}
		if i < tolen {
			toel = to[i]
		}
		if vel < toel {
			return -1
		}
		if vel > toel {
			return 1
		}
	}
	return 0
}

// String returns the version number as a string in dot-decimal notation.
func (v Version) String() string {
	var s strings.Builder
	for i, el := range v {
		if i != 0 {
			s.WriteString(".")
		}
		s.WriteString(strconv.FormatInt(el, 10))
	}
	return s.String()
}
