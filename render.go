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
	// RenderedHTML is a wrapper around io.WriterTo. It adds no new
	// functionality, but acts as a named target for HTML to return.
	//
	// This library uses writers as the underlying target to generate output to.
	// An `io.WriterTo` is a type that can write itself to an `io.Writer`. Once
	// the HTML is rendered, it can be written to the underlying writer.
	RenderedHTML interface{ io.WriterTo }

	// IndirectHTML is an interface that can be implemented by types that
	// represent HTML that is not yet rendered. This is useful for creating
	// things that can return attributes that should be inserted into the
	// enclosing element, or just a list of children.
	IndirectHTML interface {
		Render() Fragment
	}

	// HTML is an interface that can be implemented by types that represent
	// something that can be rendered to HTML. This is useful for components
	// that should be rendered to HTML.
	HTML interface {
		Render(context.Context) RenderedHTML
	}

	// Dynamically generated attributes. If this is a child of an element, it
	// will be rendered as attributes on that element.
	//
	// This, unlike other HTML types, can be rendered further to generate some
	// other content to add as children for the element. If the `ATTRS` also has
	// an `HTML`, `RenderedHTML`, or `IndirectHTML` implementation, or is a
	// `Fragment`, it will *additionally* be rendered as children of the element.
	ATTRS interface{ Attrs() Attrs }

	// A writer to render HTML to that also has a context for passing along
	// information.
	//
	// If you call `Render()` with a `ContextWriter`, it will use the context
	// provided as the context passed to the `HTML.Render` function. This lets
	// you use a template that dynamically renders different data based on the
	// context.
	//
	// An example where this might be useful: you can create a central context
	// that stores a user's authentication status. Then in any components
	// rendered, you can check the context to see if the user is authenticated
	// or not, and render different content based on that.
	//
	// A simpler approach to the same goal can be had by calling `RenderContext()`
	// instead of `Render()`. This lets you pass in a context as a param and
	// creates a `ContextWriter` for you.
	ContextWriter interface {
		io.Writer
		context.Context
	}

	contextWriter struct {
		writer  io.Writer
		context context.Context
	}
)

var _ ContextWriter = (*contextWriter)(nil)

func (w *contextWriter) Write(p []byte) (int, error)             { return w.writer.Write(p) }
func (w *contextWriter) Value(key interface{}) interface{}       { return w.context.Value(key) }
func (w *contextWriter) Deadline() (deadline time.Time, ok bool) { return w.context.Deadline() }
func (w *contextWriter) Done() <-chan struct{}                   { return w.context.Done() }
func (w *contextWriter) Err() error                              { return w.context.Err() }

// RenderContext renders the given HTML to the writer with the given context.
//
// You can forward arbitrary values to components you write using this context,
// allowing for dynamic rendering different content based on that context.
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

// Render renders the given HTML to the writer.
//
// If you want to customize how enclosed components are rendered, pass in a
// writer that implements `ContextWriter`, or use `RenderContext()` instead.
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
	case IndirectHTML:
		if child := child.Render(); child != nil {
			return child.WriteTo(cw)
		} else {
			return 0, nil
		}
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
