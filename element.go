package html

import (
	"context"
	"fmt"
	"html"
	"io"
	"slices"
)

func Element(tag string, children ...any) RenderedHTML {
	el := el{tag: tag, attrs: make(map[string]string), children: []any{}}
	return el.addChildren(children...)
}
func AttrsElement(tag string, attrs Attrs) RenderedHTML {
	el := el{tag: tag, attrs: make(map[string]string)}
	return el.addAttrs(attrs)
}

type el struct {
	tag        string
	attrs      map[string]string
	nakedAttrs []string
	children   Fragment
}

type attrIf struct {
	cond    bool
	yes, no any
}

func AttrIf(cond bool, then ...any) any {
	switch len(then) {
	case 0:
		return attrIf{cond: cond}
	case 1:
		return attrIf{cond: cond, yes: then[0]}
	case 2:
		return attrIf{cond: cond, yes: then[0], no: then[1]}
	default:
		panic("invalid number of arguments, expected 0, 1, or 2, found " + fmt.Sprint(len(then)))
	}
}

func (e *el) addAttrs(attrs Attrs) *el {
	for key, val := range attrs {
		switch val := val.(type) {
		// Nil case... naked attributes
		case nil:
			e.nakedAttrs = append(e.nakedAttrs, key)

		// Normal case
		case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, []byte, []rune:
			if key != "class" {
				e.attrs[key] = fmt.Sprint(val)
			} else if _, ok := e.attrs[key]; ok && key == "class" {
				e.attrs[key] = fmt.Sprint(e.attrs[key], " ", val)
			} else {
				e.attrs[key] = fmt.Sprint(val)
			}

		// Conditionally render attributes
		case func() any:
			if val := val(); val != nil {
				e.addAttrs(Attrs{key: val})
			}

		// Conditionally render naked attributes
		case func() bool:
			if val := val(); val {
				e.nakedAttrs = append(e.nakedAttrs, key)
			}

		case attrIf:
			if val.cond {
				if val.yes != nil {
					e.addAttrs(Attrs{key: val.yes})
				} else {
					e.nakedAttrs = append(e.nakedAttrs, key)
				}
			} else {
				if val.no != nil {
					e.addAttrs(Attrs{key: val.no})
				}
			}

		default:
			panic(fmt.Sprintf("unsupported type %T for attribute %s", val, key))
		}
	}

	return e
}

func (e *el) addChildren(children ...any) *el {
	for _, child := range children {
		switch inner := child.(type) {
		case nil:
			// Do nothing

		case ATTRS:
			attrs := inner.Attrs()
			if attrs != nil {
				e.addAttrs(attrs)
			}

			switch inner := child.(type) {
			case Fragment:
				e.addChildren(inner...)
			case HTML, RenderedHTML, IndirectHTML:
				e.addChildren(inner)
			}
		case Attrs:
			e.addAttrs(inner)
		case *Attrs:
			e.addAttrs(*inner)

		case IndirectHTML:
			e.addChildren(inner.Render()...)
		case []any:
			e.addChildren(inner...)
		case Fragment:
			e.addChildren(inner...)
		case *Fragment:
			e.addChildren(*inner...)

		default:
			e.children = append(e.children, inner)
		}
	}

	return e
}

var (
	_ RenderedHTML = (*el)(nil)
	_ HTML         = (*el)(nil)
)

func (e *el) Render(context.Context) RenderedHTML { return e }
func (e *el) WriteTo(w io.Writer) (int64, error) {
	// Start of tag
	nn := int64(0)
	n, err := fmt.Fprintf(w, "<%s", e.tag)
	nn += int64(n)
	if err != nil {
		return nn, err
	}

	// Gather set of unique attribute keys
	keySet := make(map[string]struct{}, len(e.attrs)+len(e.nakedAttrs))
	for key := range e.attrs {
		keySet[key] = struct{}{}
	}
	for _, key := range e.nakedAttrs {
		keySet[key] = struct{}{}
	}

	// Sort keys
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, key := range keys {
		if val, ok := e.attrs[key]; ok {
			n, err = fmt.Fprintf(w, " %s=\"%s\"", html.EscapeString(key), html.EscapeString(val))
			nn += int64(n)
			if err != nil {
				return nn, err
			}
		} else {
			n, err = fmt.Fprintf(w, " %s", html.EscapeString(key))
			nn += int64(n)
			if err != nil {
				return nn, err
			}
		}
	}

	if e.children == nil {
		n, err = fmt.Fprint(w, `/>`)
		nn += int64(n)
		if err != nil {
			return nn, err
		}
		return nn + int64(n), nil
	} else {
		n, err = fmt.Fprint(w, `>`)
		nn += int64(n)
		if err != nil {
			return nn, err
		}
		_n, err := Render(w, e.children)
		nn += _n
		if err != nil {
			return nn, err
		}
		n, err = fmt.Fprintf(w, `</%s>`, e.tag)
		nn += int64(n)
		return nn, err
	}
}
