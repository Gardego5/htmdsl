package html

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/sym01/htmlsanitizer"
)

type HTMLElement interface {
	io.WriterTo
	element()
	String() string
	Bytes() []byte
}

type PushAttrs interface {
	HTMLElement
	PushAttrs(attrs ...Attr) PushAttrs
}

type Fragment []HTMLElement

func (frag Fragment) element() {}
func (frag Fragment) String() string {
	b := strings.Builder{}
	frag.WriteTo(&b)
	return b.String()
}
func (frag Fragment) Bytes() []byte {
	b := bytes.Buffer{}
	frag.WriteTo(&b)
	return b.Bytes()
}
func (frag Fragment) WriteTo(w io.Writer) (int64, error) {
	n := int64(0)

	for _, child := range frag {
		nn, err := child.WriteTo(w)
		n += nn
		if err != nil {
			return 0, err
		}
	}

	return n, nil
}

func Render(w io.Writer, child any) (int64, error) {
	switch child := child.(type) {
	case HTMLElement:
		n, err := child.WriteTo(w)
		return n, err
	case string:
		n, err := htmlsanitizer.NewWriter(w).Write([]byte(child))
		return int64(n), err
	case io.WriterTo:
		n, err := child.WriteTo(htmlsanitizer.NewWriter(w))
		return n, err
	case []HTMLElement:
		n := int64(0)
		for _, child := range child {
			nn, err := Render(w, child)
			n += nn
			if err != nil {
				return n, err
			}
		}
		return n, nil
	default:
		n, err := fmt.Fprint(htmlsanitizer.NewWriter(w), child)
		return int64(n), err
	}
}
