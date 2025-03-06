package util_test

import (
	"bytes"
	"testing"

	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
)

func TestIfAndSwitch(t *testing.T) {
	for _, test := range []struct {
		name, expect string
		input        any
	}{
		{
			name:   "switch with case",
			expect: "<div>You say goodbye.</div>",
			input: Div{
				"You say ",
				Switch().
					Case(false, "hello").
					Case(true, "goodbye").
					Default("world"),
				".",
			},
		},
		{
			name:   "switch with default",
			expect: "yep",
			input: Switch().
				Case(1 == 2, "nope").
				Default("yep"),
		},
		{
			name:   "if attrs",
			expect: `<div class="container">inside</div>`,
			input:  Div{If(true, Attrs{"class": "container"}), "inside"},
		},
		{
			name:   "if attrs with else",
			expect: `<div id="div-1">inside</div>`,
			input:  Div{If(false, Id("container")).Else(Id("div-1")), "inside"},
		},
		{
			name:   "if with else if",
			expect: `<div>two</div>`,
			input:  Div{If(false, "one").ElseIf(true, "two").Else("three")},
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
