package html

import (
	"fmt"
	"io"
	"strings"
)

type PreEscaped string

var (
	_ RenderedHTML = PreEscaped("")
	_ HTML         = PreEscaped("")
)

func (text PreEscaped) Render() RenderedHTML { return text }
func (text PreEscaped) element()             {}
func (text PreEscaped) WriteTo(w io.Writer) (int64, error) {
	nn, err := fmt.Fprint(w, text)
	return int64(nn), err
}
func (text PreEscaped) String() string    { return string(text) }
func (text PreEscaped) Bytes() []byte     { return []byte(text) }
func (text PreEscaped) Reader() io.Reader { return strings.NewReader(string(text)) }
