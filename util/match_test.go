package util_test

import (
	"bytes"
	"testing"

	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
)

func TestMatch(t *testing.T) {
	for _, test := range []struct {
		name, expect string
		input        any
	}{
		{
			name:   "match number",
			expect: "<div>You say three.</div>",
			input: Div{
				"You say ",
				Match(3).
					When(1, "one").
					When(2, "two").
					When(3, "three").
					Default("otherwise"),
				".",
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			_, err := Render(buf, test.input)
			if err != nil {
				t.Fatal(err)
			}

			got := buf.String()
			if got != test.expect {
				t.Errorf("got %q; expected %q", got, test.expect)
			}
		})
	}
}
