package osc

import (
	"regexp"
	"strings"
)

// addressExists returns true if the OSC address `addr` is found in `handlers`.
func addressExists(addr string, handlers map[string]Handler) bool {
	for h := range handlers {
		if h == addr {
			return true
		}
	}
	return false
}

// getRegEx compiles and returns a regular expression object for the given
// address `pattern`.
func getRegEx(pattern string) (*regexp.Regexp, error) {
	for _, trs := range []struct {
		old, new string
	}{
		{".", `\.`}, // Escape all '.' in the pattern
		{"(", `\(`}, // Escape all '(' in the pattern
		{")", `\)`}, // Escape all ')' in the pattern
		{"*", ".*"}, // Replace a '*' with '.*' that matches zero or more chars
		{"{", "("},  // Change a '{' to '('
		{",", "|"},  // Change a ',' to '|'
		{"}", ")"},  // Change a '}' to ')'
		{"?", "."},  // Change a '?' to '.'
	} {
		pattern = strings.Replace(pattern, trs.old, trs.new, -1)
		pattern = "^" + pattern + "$"
	}

	return regexp.Compile(pattern)
}

// getTypeTag returns the OSC type tag for the given argument.
func getTypeTag(arg any) byte {
	switch t := arg.(type) {
	case bool:
		if t {
			return 'T'
		}
		return 'F'
	case nil:
		return 'N'
	case int32:
		return 'i'
	case float32:
		return 'f'
	case string:
		return 's'
	case []byte:
		return 'b'
	case int64:
		return 'h'
	case float64:
		return 'd'
	case Timetag:
		return 't'
	default:
		return '\xff'
	}
}
