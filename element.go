package html

import (
	"io"
	"slices"
)

func Element(tag string, children ...any) RenderedHTML {
	childAttrs, childEls := make(map[string]string), []any{}
	addChildren(childAttrs, &childEls, children...)
	return el{tag, childAttrs, childEls}
}

func AttrsElement(tag string, attrs ...Attr) RenderedHTML {
	el := el{tag, make(map[string]string), nil}
	for _, attr := range attrs {
		addAttr(el.attrs, attr[0], attr[1])
	}
	return el
}

type el struct {
	tag      string
	attrs    map[string]string
	children Fragment
}

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

	keys := make([]string, 0, len(e.attrs))
	for key := range e.attrs {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, key := range keys {
		n, err = w.Write([]byte(" " + key + "=\"" + e.attrs[key] + "\""))
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

func addChildren(
	attributeElements map[string]string, childrenElements *[]any,
	children ...any,
) {
	for _, child := range children {
		switch child := child.(type) {
		case Attr:
			addAttr(attributeElements, child[0], child[1])
		case Attrs:
			for _, attr := range child {
				addAttr(attributeElements, attr[0], attr[1])
			}
		case Fragment:
			addChildren(attributeElements, childrenElements, child...)
		case []any:
			addChildren(attributeElements, childrenElements, child...)
		default:
			addChild(childrenElements, child)
		}
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
