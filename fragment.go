package html

import (
	"bytes"
	"io"
	"strings"
)

type HTMLElement interface {
	io.WriterTo
	element()
	String() string
	Bytes() []byte
	Reader() io.Reader
}

type Fragment []any

var _ HTMLElement = Fragment(nil)
var _ HTML = Fragment(nil)

func (frag Fragment) Render() HTMLElement { return frag }
func (frag Fragment) element()            {}
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
