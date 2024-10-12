package html

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/sym01/htmlsanitizer"
)

type (
	HTML        interface{ Render() RenderedHTML }
	ContextHTML interface {
		HTML
		RenderWithContext(c context.Context) RenderedHTML
	}
	RenderedHTML interface{ io.WriterTo }
)

func is[T any](val any) bool { _, ok := val.(T); return ok }

type (
	contextWriter struct {
		w io.Writer
		c context.Context
	}
	ContextWriter interface {
		io.Writer
		context.Context
	}
)

var _ ContextWriter = (*contextWriter)(nil)

func (w *contextWriter) Write(p []byte) (int, error)             { return w.w.Write(p) }
func (w *contextWriter) Value(key interface{}) interface{}       { return w.c.Value(key) }
func (w *contextWriter) Deadline() (deadline time.Time, ok bool) { return w.c.Deadline() }
func (w *contextWriter) Done() <-chan struct{}                   { return w.c.Done() }
func (w *contextWriter) Err() error                              { return w.c.Err() }

func RenderWithContext(w io.Writer, c context.Context, child any) (int64, error) {
	return Render(&contextWriter{w: w, c: c}, child)
}

func Render(w io.Writer, child any) (int64, error) {
	if child == nil {
		return 0, nil
	}

	cw, ok := w.(ContextWriter)
	if !ok {
		cw = &contextWriter{w: w, c: context.Background()}
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
	case ContextHTML:
		if child := child.RenderWithContext(cw); child != nil {
			return child.WriteTo(cw)
		} else {
			return 0, nil
		}
	case HTML:
		if child := child.Render(); child != nil {
			return child.WriteTo(cw)
		} else {
			return 0, nil
		}
	case func() any:
		return Render(cw, child())
	case func(context.Context) any:
		return Render(cw, child(cw))
	case string:
		n, err := fmt.Fprint(htmlsanitizer.NewWriter(cw), child)
		return int64(n), err
	case io.WriterTo:
		return child.WriteTo(htmlsanitizer.NewWriter(cw))
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
