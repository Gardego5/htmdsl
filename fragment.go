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
	Reader() io.Reader
}

type Fragment []any

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
		nn, err := Render(w, child)
		n += nn
		if err != nil {
			return 0, err
		}
	}

	return n, nil
}
func (frag Fragment) Reader() io.Reader {
	r, w := io.Pipe()
	go func(w *io.PipeWriter) {
		_, err := frag.WriteTo(w)
		if err != nil {
			w.CloseWithError(err)
		} else {
			w.Close()
		}
	}(w)
	return r
}

type HTMLComponent interface{ Render() HTMLElement }

func Render(w io.Writer, child any) (int64, error) {
	switch child := child.(type) {
	case nil:
		return 0, nil
	case HTMLComponent:
		return child.Render().WriteTo(w)
	case HTMLElement:
		return child.WriteTo(w)
	case string:
		n, err := htmlsanitizer.NewWriter(w).Write([]byte(child))
		return int64(n), err
	case *string:
		if child == nil {
			return 0, nil
		} else {
			return Render(w, *child)
		}
	case io.WriterTo:
		return child.WriteTo(htmlsanitizer.NewWriter(w))
	case io.Reader:
		return io.Copy(htmlsanitizer.NewWriter(w), child)
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
