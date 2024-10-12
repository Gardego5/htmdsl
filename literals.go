package html

import (
	"fmt"
	"io"
)

type literal string

const DOCTYPE literal = "<!DOCTYPE html>"

var (
	_ RenderedHTML = literal("")
	_ HTML         = literal("")
)

func (lit literal) Render() RenderedHTML { return lit }
func (lit literal) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, lit)
	return int64(n), err
}
