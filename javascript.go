package path

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var canDot = regexp.MustCompile(`^\w+$`)

// ToJS converts a path to a JSONPath-style string representation.
//
// For example:
//   ["foo", 0, "bar"] -> "foo[0].bar"
//
// The produced value can be parsed by FromJS and will produce an identical path.
func (p P) ToJS() string {
	var b strings.Builder
	for _, c := range p.components {
		switch v := c.(type) {
		case string:
			if canDot.MatchString(v) {
				if b.Len() > 0 {
					b.WriteRune('.')
				}
				b.WriteString(v)
			} else {
				b.WriteRune('[')
				b.WriteString(strconv.Quote(v))
				b.WriteRune(']')
			}
		case int:
			b.WriteRune('[')
			b.WriteString(strconv.FormatInt(int64(v), 10))
			b.WriteRune(']')
		}
	}

	return b.String()
}

const dotString = `\.\w+`
const bracketNumber = `\[\d+\]`
const bracketString = `\["(?:[^"\\]|(?:\\.))*"\]`

// pieceReg is used with FindAllStringIndex to break the string into
// non-overlapping components.
var pieceReg = regexp.MustCompile(`(?:^\w+)|(?:` + dotString + `)|(?:` + bracketNumber + `)|(?:` + bracketString + `)`)

// FromJSPartial parses a JSONPath-style path string from the start greedily,
// and returns the index of how far it managed to parse.
//
// For example:
//   "foo[0].bar asdf" -> ["foo", 0, "bar"], index for " asdf"
func FromJSPartial(s string) (p P, idx int, err error) {
	p = P{components: []interface{}{}}

	if len(s) == 0 {
		return p, 0, nil
	}

	// There is one specific case that the given regex doesn't handle well - we
	// don't want to accept an initial dot at the start of the string (e.g.
	// ".x.y"), but without lookahead/lookbehind tricks, we can't easily filter
	// this out, so we handle this explicitly.
	if s[0] == '.' {
		return p, 0, nil
	}
	// Note that we have to manually check the indices to make sure it covers
	// the whole string - FindAllStringIndex only guarantees that the matches
	// are non-overlapping, not that they are exhaustive.
	var lastIndex = 0
	for _, idx := range pieceReg.FindAllStringIndex(s, -1) {
		if idx[0] != lastIndex {
			return p, lastIndex, nil
		}
		lastIndex = idx[1]
		m := s[idx[0]:idx[1]]
		if m[0] == '.' {
			p.components = append(p.components, m[1:])
		} else if m[0] == '[' {
			if v, err := strconv.Atoi(m[1 : len(m)-1]); err == nil {
				p.components = append(p.components, v)
			} else if v, err := strconv.Unquote(m[1 : len(m)-1]); err == nil {
				p.components = append(p.components, v)
			} else {
				return p, 0, errors.New("could not parse portion in square brackets")
			}
		} else {
			p.components = append(p.components, m)
		}
	}
	if lastIndex != len(s) {
		return p, lastIndex, nil
	}

	return p, len(s), nil
}

// FromJS parses a JSONPath-style path string and returns an error
// if it is unable to parse the string.
//
// For example:
//   "foo[0].bar" -> ["foo", 0, "bar"]
func FromJS(s string) (P, error) {
	p, idx, err := FromJSPartial(s)
	if err != nil {
		return p, err
	}
	if idx != len(s) {
		return p, errors.New("could not parse entire string")
	}

	return p, nil
}
