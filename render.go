package html

import (
	"fmt"
	"io"

	"github.com/sym01/htmlsanitizer"
)

type HTML interface{ Render() HTMLElement }

func Render(w io.Writer, child any) (int64, error) {
	if child == nil {
		return 0, nil
	}

	switch child := child.(type) {
	case HTMLElement:
		return child.WriteTo(w)
	case HTML:
		if child := child.Render(); child != nil {
			return child.WriteTo(w)
		} else {
			return 0, nil
		}
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
	case []HTML:
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
