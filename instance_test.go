package path

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	require := require.New(t)

	p := Int(0).Int(1).Int(2).Int(3, 4).Int(5)
	require.Equal([]interface{}{
		0, 1, 2, 3, 4, 5,
	}, p.components)

	// Path should work across multiple types of component.
	p2 := Str("hello").Int(0).Str("world", "friends").Int(1)
	require.Equal([]interface{}{
		"hello", 0, "world", "friends", 1,
	}, p2.components)

	// Path should not mutate the path parameter:
	p2 = p2.Path(p)
	require.Equal([]interface{}{
		0, 1, 2, 3, 4, 5,
	}, p.components)
	// But it should mutate the chained path:
	require.Equal([]interface{}{
		"hello", 0, "world", "friends", 1, 0, 1, 2, 3, 4, 5,
	}, p2.components)
}

func TestZeroValue(t *testing.T) {
	require := require.New(t)

	p := P{}
	p = p.Int(0).Str("foo", "bar").Path(P{})
	require.Equal([]interface{}{0, "foo", "bar"}, p.components)
}

func TestSliceMethods(t *testing.T) {
	require := require.New(t)

	// Len
	p := P{}
	require.Equal(0, p.Len())
	p = p.Str("hello").Int(0)
	require.Equal(2, p.Len())

	// Sub
	require.Equal([]interface{}{"hello", 0}, p.Sub().components)
	require.Equal([]interface{}{"hello", 0}, p.Sub(0).components)
	require.Equal([]interface{}{0}, p.Sub(1).components)
	require.Equal([]interface{}{}, p.Sub(2).components)
	require.Equal([]interface{}{"hello"}, p.Sub(0, 1).components)
	require.Equal([]interface{}{"hello", 0}, p.Sub(0, 2).components)
	require.Equal([]interface{}{0}, p.Sub(1, 2).components)
	require.Equal([]interface{}{}, p.Sub(1, 1).components)

	// At
	require.Equal("hello", p.At(0))
	require.Equal(0, p.At(1))

	// Components
	require.Equal([]interface{}{"hello", 0}, p.Components())
}

func TestString(t *testing.T) {
	require := require.New(t)

	require.Equal("[]", P{}.String())
	require.Equal(`["hello",1]`, Str("hello").Int(1).String())
}

func TestMarshalJSON(t *testing.T) {
	require := require.New(t)

	p := &P{}
	require.NoError(json.Unmarshal([]byte(`["hello", "world", 1, 2, "three"]`), p))
	require.Equal([]interface{}{"hello", "world", 1, 2, "three"}, p.components)

	require.Equal("[]", P{}.String())
	require.Equal(`["hello",1]`, Str("hello").Int(1).String())

	require.Error(json.Unmarshal([]byte(`[1.23]`), &P{}))  // float
	require.Error(json.Unmarshal([]byte(`"hello"`), &P{})) // not an array

	v, err := json.Marshal(p)
	require.NoError(err)
	require.Equal(`["hello","world",1,2,"three"]`, string(v))
}
