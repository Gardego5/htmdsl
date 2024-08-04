package html

import (
	"io"
	"strings"
)

type text []any

type PreEscaped string

func (text PreEscaped) element() {}
func (text PreEscaped) WriteTo(w io.Writer) (int64, error) {
	nn, err := w.Write([]byte(text))
	return int64(nn), err
}
func (text PreEscaped) String() string    { return string(text) }
func (text PreEscaped) Bytes() []byte     { return []byte(text) }
func (text PreEscaped) Reader() io.Reader { return strings.NewReader(string(text)) }
