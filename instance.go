package path

import (
	"encoding/json"
	"math"

	"github.com/pkg/errors"
)

// P is a slice-based representation of a JSON path. Each slice element
// represents a single path component: either an object key (if a string) or
// an array index (if an integer).
//
// Unlike JSONPath (https://goessner.net/articles/JsonPath/) and JSON pointers
// (https://datatracker.ietf.org/doc/html/rfc6901), paths are serialized as
// slices rather than strings. This helps avoid a class of errors around improperly
// (de-)escaping path components.
//
// Unlike []interface{}, paths are type-safe via the Append* methods.
//
// The zero value is safe to use.
type P struct {
	// components is the underlying representation of this path. It is not exposed
	// as part of path's API since clients could mistakenly insert non-string/int
	// components into the path.
	components []interface{}
}

// Str is a helper for constructing a path from one or more strings.
func Str(s ...string) P {
	return P{}.Str(s...)
}

// Int is a helper for constructing a path from one or more ints.
func Int(i ...int) P {
	return P{}.Int(i...)
}

// Path is a helper for constructing a path from one or more paths.
func Path(paths ...P) P {
	return P{}.Path(paths...)
}

// Sub(start, end) returns a subpath from [start, end).
//
// If not provided, end defaults to len(p).
// If not provided, start defaults to 0.
//
// For example, ["foo", "bar"].Sub(1) would return ["bar"].
func (p P) Sub(idxs ...int) P {
	start := 0
	end := len(p.components)
	if len(idxs) > 0 {
		start = idxs[0]
	}
	if len(idxs) > 1 {
		end = idxs[1]
	}

	return P{
		components: append([]interface{}{}, (p.components)[start:end]...),
	}
}

// Str appends one or more strings onto the provided path.
//
// It has the same semantics as the built-in append: the returned path
// may re-allocate the underlying array as needed.
func (p P) Str(s ...string) P {
	for _, si := range s {
		p.components = append(p.components, si)
	}
	return p
}

// Int appends one or more ints onto the provided path.
//
// It has the same semantics as the built-in append: the returned path
// may re-allocate the underlying array as needed.
func (p P) Int(i ...int) P {
	for _, ii := range i {
		p.components = append(p.components, ii)
	}
	return p
}

// Path appends one or more paths onto the provided path.
//
// It has the same semantics as the built-in append: the returned path
// may re-allocate the underlying array as needed.
func (p P) Path(paths ...P) P {
	for _, pi := range paths {
		if pi.components != nil {
			p.components = append(p.components, pi.components...)
		}
	}
	return p
}

func (p P) Len() int {
	return len(p.components)
}

func (p P) At(i int) interface{} {
	return p.components[i]
}

func (p P) Components() []interface{} {
	return append([]interface{}{}, p.components...)
}

var _ json.Marshaler = &P{}
var _ json.Unmarshaler = &P{}

func (p P) MarshalJSON() ([]byte, error) {
	components := []interface{}{}
	if p.components != nil {
		components = p.components
	}
	return json.Marshal(components)
}

func (p *P) UnmarshalJSON(b []byte) error {
	comps := []interface{}{}
	if err := json.Unmarshal(b, &comps); err != nil {
		return err
	}

	p.components = []interface{}{}
	for _, c := range comps {
		switch v := c.(type) {
		case float64:
			// JSON always unmarshals as float64s. Convert to ints.
			if math.Floor(v) != v {
				return errors.Errorf("unexpected float: %v", v)
			}
			p.components = append(p.components, int(v))
		case string:
			p.components = append(p.components, c)
		default:
			return errors.Errorf("unexpected value: %v", v)
		}
	}

	return nil
}

func (p P) String() string {
	components := []interface{}{}
	if p.components != nil {
		components = p.components
	}
	b, _ := json.Marshal(components)
	return string(b)
}
