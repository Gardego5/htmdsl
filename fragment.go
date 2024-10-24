package html

import (
	"context"
	"io"
)

type Fragment []any

var (
	_ RenderedHTML = Fragment{}
	_ HTML         = Fragment{}
)

func (f Fragment) Render(context.Context) RenderedHTML { return f }
func (f Fragment) WriteTo(w io.Writer) (int64, error) {
	var n int64
	for _, a := range f {
		nn, err := Render(w, a)
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
