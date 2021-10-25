package path

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJS(tt *testing.T) {
	for _, test := range []struct {
		s string
		c []interface{}
	}{
		{"", []interface{}{}},
		{"outputs", []interface{}{"outputs"}},
		{`outputs[""]`, []interface{}{"outputs", ""}},
		{"[0]", []interface{}{0}},
		{`["an output"]`, []interface{}{"an output"}},
		{"outputs.foo", []interface{}{"outputs", "foo"}},
		{"outputs[0]", []interface{}{"outputs", 0}},
		{"outputs[0][1]", []interface{}{"outputs", 0, 1}},
		{`a["hello world"][10].foo.bar`, []interface{}{"a", "hello world", 10, "foo", "bar"}},
		{`["[]"]`, []interface{}{"[]"}},
		{`["]["]`, []interface{}{"]["}},
		{`["\""]`, []interface{}{`"`}},
		{`["\\"]`, []interface{}{`\`}},
		{`[""]`, []interface{}{``}},
		{`a["]\"\\["][10].b`, []interface{}{"a", `]"\[`, 10, "b"}},
		{`a[".b"][10]["[\".c\"]"]`, []interface{}{`a`, `.b`, 10, `[".c"]`}},
	} {
		tt.Run(test.s, func(t *testing.T) {
			require := require.New(t)

			inst, err := FromJS(test.s)
			require.NoError(err)
			require.Equal(test.c, inst.components)
			require.Equal(test.s, inst.ToJS())
		})
	}
}

func TestInvalidJS(tt *testing.T) {
	for _, test := range []struct {
		s string
	}{
		{".outputs"},
		{"outputs."},
		{"outputs.."},
		{"outputs.[0]"},
		{"["},
		{"]"},
		{"[[]"},
		{`outputs["""]`},
		{`outputs[.]`},
		{`outputs[abc]`},
	} {
		tt.Run(test.s, func(t *testing.T) {
			require := require.New(t)

			require.NotPanics(func() {
				_, err := FromJS(test.s)
				require.NotNil(err)
			})
		})
	}
}
