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

func (e *el) addAttrs(attrs Attrs) *el {
	for key, val := range attrs {
		switch val := val.(type) {
		case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, []byte, []rune:
			if key != "class" {
				e.attrs[key] = fmt.Sprint(val)
			} else if _, ok := e.attrs[key]; ok && key == "class" {
				e.attrs[key] = fmt.Sprint(e.attrs[key], " ", val)
			} else {
				e.attrs[key] = fmt.Sprint(val)
			}

		case nil:
			e.nakedAttrs = append(e.nakedAttrs, key)

		default:
			panic(fmt.Sprintf("unsupported type %T for attribute %s", val, key))
		}
	}

	return e
}

func (e *el) addChildren(children ...any) *el {
	for _, child := range children {
		switch child := child.(type) {
		case Attrs:
			e.addAttrs(child)

		case *Attrs:
			e.addAttrs(*child)

		case []any:
		case Fragment:
			e.addChildren(child...)

		default:
			e.children = append(e.children, child)
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
