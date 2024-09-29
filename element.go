package html

import "io"

func Element(tag string, children ...any) RenderedHTML {
	childAttrs, childEls := make(map[string]string), []any{}

	for _, child := range children {
		if attrs, ok := child.(Attrs); ok {
			for _, attr := range attrs {
				addAttr(childAttrs, attr[0], attr[1])
			}
		} else if attr, ok := child.(Attr); ok {
			addAttr(childAttrs, attr[0], attr[1])
		} else {
			addChild(&childEls, child)
		}
	}

	el := el{tag, childAttrs, childEls}

	return el
}

func AttrsElement(tag string, attrs ...Attr) RenderedHTML {
	el := el{tag, make(map[string]string), nil}
	for _, attr := range attrs {
		addAttr(el.attrs, attr[0], attr[1])
	}
	return el
}

type (
	el struct {
		tag      string
		attrs    map[string]string
		children Fragment
	}
)

var (
	_ RenderedHTML = (*el)(nil)
	_ HTML         = (*el)(nil)
)

func (e el) Render() RenderedHTML { return e }
func (e el) WriteTo(w io.Writer) (int64, error) {
	nn := int64(0)
	n, err := w.Write([]byte("<" + e.tag))
	nn += int64(n)
	if err != nil {
		return nn, err
	}

	for key, value := range e.attrs {
		n, err = w.Write([]byte(" " + key + "=\"" + value + "\""))
		nn += int64(n)
		if err != nil {
			return nn, err
		}
	}

	if e.children == nil {
		n, err = w.Write([]byte("/>"))
		nn += int64(n)
		if err != nil {
			return nn, err
		}
		return nn + int64(n), nil
	} else {
		n, err = w.Write([]byte(">"))
		nn += int64(n)
		if err != nil {
			return nn, err
		}
		_n, err := Render(w, e.children)
		nn += _n
		if err != nil {
			return nn, err
		}
		n, err = w.Write([]byte("</" + e.tag + ">"))
		nn += int64(n)
		return nn, err
	}
}

func addAttr(attrs map[string]string, key string, value string) {
	if _, ok := attrs[key]; ok && key == "class" {
		current := attrs[key]
		attrs[key] = current + " " + value
	} else {
		attrs[key] = value
	}
}

func addChild(children *[]any, child any) {
	*children = append(*children, child)
}
