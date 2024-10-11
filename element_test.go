package html_test

import (
	"bytes"
	"testing"

	. "github.com/Gardego5/htmdsl"
)

func TestElement(t *testing.T) {
	for _, test := range []struct {
		name, expected string
		result         RenderedHTML
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
			name:     "attributes",
			expected: `<div class="container" id="main"></div>`,
			result:   Element("div", Attr{"class", "container"}, Attr{"id", "main"}),
		},
		{
			name:     "multiple classes",
			expected: `<div class="container main" id="main"></div>`,
			result:   Element("div", Class("container"), Attr{"class", "main"}, Attr{"id", "main"}),
		},
		{
			name:     "multiple classes with children",
			expected: `<div class="container main" id="main"><div class="inner">hello</div></div>`,
			result:   Element("div", Class("container"), Attr{"class", "main"}, Attr{"id", "main"}, Element("div", Class("inner"), "hello")),
		},
		{
			name:     "input",
			expected: "<input type=\"text\" value=\"hello\"/>",
			result:   Input{{"type", "text"}, {"value", "hello"}}.Render(),
		},
		{
			name:     "defined tags: a, b, button, div, h1, h2, h3, h4, h5, h6, img, input, label, li, ol, p, span, ul",
			expected: `<a href="https://example.com">click me</a><b>bold</b><button>click me</button><div>hello</div><h1>hello</h1><h2>hello</h2><h3>hello</h3><h4>hello</h4><h5>hello</h5><h6>hello</h6><img src="https://example.com"/><input type="text" value="hello"/>`,
			result: Fragment{
				A{"click me", Attr{"href", "https://example.com"}},
				B{"bold"},
				Button{"click me"}, Div{"hello"},
				H1{"hello"}, H2{"hello"}, H3{"hello"},
				H4{"hello"}, H5{"hello"}, H6{"hello"},
				Img{Attr{"src", "https://example.com"}},
				Input{Attr{"type", "text"}, Attr{"value", "hello"}},
			},
		},
		{
			name:     "attrs hoisting",
			expected: `<div class="container" id="main"><div class="inner">hello</div></div>`,
			result: Div{
				Attrs{{"id", "main"}},
				Fragment{
					Class("container"), // hoisted
					Div{Class("inner"), "hello"},
				},
			}.Render(),
		},
		{
			name:     "invalid hoisting",
			expected: `[class container]<div class="inner">hello</div>`,
			result: Fragment{
				Class("container"),
				Div{Class("inner"), "hello"},
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			_, err := test.result.WriteTo(buf)
			if err != nil {
				t.Fatal(err)
			}

			got := buf.String()
			if got != test.expected {
				t.Errorf("got %q; expected %q", got, test.expected)
			}
		})
	}
}
