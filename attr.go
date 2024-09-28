package html

import (
	"fmt"
	"io"
	"strings"
)

// attr is a marker interface for attributes (plural and singular).
type attr interface{ attr() }

type Attr [2]string

func (attr Attr) attr() {}
func (attr Attr) WriteTo(w io.Writer) (int64, error) {
	if attr[1] == "" {
		n, err := fmt.Fprintf(w, " %s", attr[0])
		return int64(n), err
	} else {
		n, err := fmt.Fprintf(w, " %s=\"%s\"", attr[0], attr[1])
		return int64(n), err
	}
}
func Class(list ...string) Attr {
	return Attr{"class", strings.Join(list, " ")}
}

type baseAttrs []Attr

func (attrs baseAttrs) attr() {}
func (attrs baseAttrs) WriteTo(w io.Writer) (int64, error) {
	n := int64(0)

	classes := []string{}

	for _, attr := range attrs {
		if attr[0] == "class" {
			classes = append(classes, attr[1])
			continue
		}

		nn, err := attr.WriteTo(w)
		n += nn
		if err != nil {
			return n, err
		}
	}

	if len(classes) > 0 {
		nn, err := fmt.Fprintf(w, " %s=\"%s\"", "class", strings.Join(classes, " "))
		if err != nil {
			return n, err
		}
		n += int64(nn)
	}

	return n, nil
}

type Attrs baseAttrs

func (attrs Attrs) attr() {}
func (attrs Attrs) WriteTo(w io.Writer) (int64, error) {
	return baseAttrs(attrs).WriteTo(w)
}
