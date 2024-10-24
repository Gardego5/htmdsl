package html

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/sym01/htmlsanitizer"
)

type (
	RenderedHTML interface{ io.WriterTo }
	HTML         interface {
		Render(context.Context) RenderedHTML
	}
	contextWriter struct {
		writer  io.Writer
		context context.Context
	}
	ContextWriter interface {
		io.Writer
		context.Context
	}
)

var _ ContextWriter = (*contextWriter)(nil)

func (w *contextWriter) Write(p []byte) (int, error)             { return w.writer.Write(p) }
func (w *contextWriter) Value(key interface{}) interface{}       { return w.context.Value(key) }
func (w *contextWriter) Deadline() (deadline time.Time, ok bool) { return w.context.Deadline() }
func (w *contextWriter) Done() <-chan struct{}                   { return w.context.Done() }
func (w *contextWriter) Err() error                              { return w.context.Err() }

func RenderContext(w io.Writer, c context.Context, child any) (int64, error) {
	return Render(&contextWriter{writer: w, context: c}, child)
}

var htmlEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`'`, "&#39;", // "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "&#34;", // "&#34;" is shorter than "&quot;".
)

func Render(w io.Writer, child any) (int64, error) {
	if child == nil {
		return 0, nil
	}

	cw, ok := w.(ContextWriter)
	if !ok {
		cw = &contextWriter{writer: w, context: context.Background()}
	}

	switch child := child.(type) {
	case Fragment:
		nn := int64(0)
		for _, child := range child {
			n, err := Render(cw, child)
			nn += n
			if err != nil {
				return nn, err
			}
		}
		return nn, nil
	case RenderedHTML:
		return child.WriteTo(cw)
	case HTML:
		if child := child.Render(cw); child != nil {
			return child.WriteTo(cw)
		} else {
			return 0, nil
		}
	case func() any:
		return Render(cw, child())
	case func(context.Context) any:
		return Render(cw, child(cw))
	case string:
		n, err := htmlEscaper.WriteString(htmlsanitizer.NewWriter(cw), child)
		return int64(n), err
	case io.Reader:
		return io.Copy(htmlsanitizer.NewWriter(cw), child)
	default:
		ty := reflect.ValueOf(child)
		if ty.Kind() == reflect.Slice {
			if ty.IsNil() {
				return 0, nil
			}

			nn := int64(0)
			for i := 0; i < ty.Len(); i++ {
				n, err := Render(cw, ty.Index(i).Interface())
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
			return Render(cw, ty.Elem().Interface())
		} else {
			n, err := fmt.Fprint(htmlsanitizer.NewWriter(cw), child)
			return int64(n), err
		}
	}
}
