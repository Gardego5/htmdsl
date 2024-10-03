package html

import (
	"fmt"
	"io"
	"reflect"

	"github.com/sym01/htmlsanitizer"
)

type (
	HTML         interface{ Render() RenderedHTML }
	RenderedHTML interface{ io.WriterTo }
)

func Render(w io.Writer, child any) (int64, error) {
	if child == nil {
		return 0, nil
	}

	switch child := child.(type) {
	case RenderedHTML:
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
	case io.WriterTo:
		return child.WriteTo(htmlsanitizer.NewWriter(w))
	case io.Reader:
		return io.Copy(htmlsanitizer.NewWriter(w), child)
	default:
		ty := reflect.ValueOf(child)
		if ty.Kind() == reflect.Slice {
			if ty.IsNil() {
				return 0, nil
			}

			nn := int64(0)
			for i := 0; i < ty.Len(); i++ {
				n, err := Render(w, ty.Index(i).Interface())
				nn += n
				if err != nil {
					return nn, err
				}
			}
			return nn, nil
		} else if ty.Kind() == reflect.Pointer {
			if ty.IsNil() {
				return 0, nil
			}
			return Render(w, ty.Elem().Interface())
		} else {
			n, err := fmt.Fprint(htmlsanitizer.NewWriter(w), child)
			return int64(n), err
		}
	}
}
