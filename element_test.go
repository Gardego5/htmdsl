package html_test

import (
	"bytes"
	"context"
	"testing"

	. "github.com/Gardego5/htmdsl"
)

func TestElement(t *testing.T) {
	for _, test := range []struct {
		name, expected string
		result         any
		context        context.Context
	}{
		{
			name:     "div with text",
			expected: "<div>hello</div>",
			result:   Element("div", "hello"),
		},
		{
			name:     "div with nested div",
			expected: "<div><div></div></div>",
			result:   Element("div", Element("div")),
		},
		{
			name:     "div with nested div and text",
			expected: "<div><div>hello</div></div>",
			result:   Element("div", Element("div", "hello")),
		},
		{
			name:     "web components",
			expected: "<my-component></my-component>",
			result:   Element("my-component"),
		},
		{
			name:     "naked attributes",
			expected: `<select><option value="1">one</option><option selected value="2">two</option></select>`,
			result: Select{
				Option{Attrs{"value": 1}, "one"},
				Option{Attrs{"selected": nil, "value": 2}, "two"},
			},
		},
		{
			name:     "attributes",
			expected: `<div class="container" id="main"></div>`,
			result:   Element("div", Attrs{"class": "container", "id": "main"}),
		},
		{
			name:     "multiple classes",
			expected: `<div class="container main" id="main"></div>`,
			result:   Element("div", Class("container"), Attrs{"class": "main", "id": "main"}),
		},
		{
			name:     "multiple classes with children",
			expected: `<div class="container main" id="main"><div class="inner">hello</div></div>`,
			result:   Element("div", Class("container"), Attrs{"class": "main"}, Attrs{"id": "main"}, Element("div", Class("inner"), "hello")),
		},
		{
			name:     "input",
			expected: "<input type=\"text\" value=\"hello\"/>",
			result:   Input{"type": "text", "value": "hello"},
		},
		{
			name:     "defined tags: a, b, button, div, h1, h2, h3, h4, h5, h6, img, input, label, li, ol, p, span, ul",
			expected: `<a href="https://example.com">click me</a><b>bold</b><button>click me</button><div>hello</div><h1>hello</h1><h2>hello</h2><h3>hello</h3><h4>hello</h4><h5>hello</h5><h6>hello</h6><img src="https://example.com"/><input type="text" value="hello"/>`,
			result: Fragment{
				A{"click me", Attrs{"href": "https://example.com"}},
				B{"bold"},
				Button{"click me"}, Div{"hello"},
				H1{"hello"}, H2{"hello"}, H3{"hello"},
				H4{"hello"}, H5{"hello"}, H6{"hello"},
				Img{"src": "https://example.com"},
				Input{"type": "text", "value": "hello"},
			},
		},
		{
			name:     "attrs hoisting",
			expected: `<div class="container" id="main"><div class="inner">hello</div></div>`,
			result: Div{
				Attrs{"id": "main"},
				Fragment{
					Class("container"), // hoisted
					Div{Class("inner"), "hello"},
				},
			},
		},
		{
			name:     "conditional attributes",
			expected: `<select id="choices"><option selected>One</option><option>Two</option></select>`,
			result: Select{
				Attrs{
					"class": AttrIf(false, "container"),
					"id":    AttrIf(1 == 2, "options", "choices"),
				},
				Option{Attrs{"selected": AttrIf("One" == "One")}, "One"},
				Option{Attrs{"selected": AttrIf("One" == "Two")}, "Two"},
			},
		},
		{
			name:     "invalid hoisting",
			expected: `map[class:container]<div class="inner">hello</div>`,
			result: Fragment{
				Class("container"),
				Div{Class("inner"), "hello"},
			},
		},
		{
			name:     "context value",
			expected: "<div>context-data</div>",
			result:   Div{func(ctx context.Context) any { return "context-data" }},
			context:  context.WithValue(context.Background(), "key", "context-data"),
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := test.context
			if ctx == nil {
				ctx = context.Background()
			}

			buf := new(bytes.Buffer)
			_, err := RenderContext(buf, ctx, test.result)
			if err != nil {
				t.Fatal(err)
			}

			got := buf.String()
			if got != test.expected {
				t.Errorf("expected %q", test.expected)
				t.Errorf("got      %q", got)
			}
		})
	}
}
