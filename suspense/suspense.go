package suspense

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	html "github.com/Gardego5/htmdsl"
)

type Suspense struct {
	html.ExoticElement
	Fallback html.PushAttrs
	Children func() any
}

func (s Suspense) WriteTo(w io.Writer) (int64, error) {
	switch w := w.(type) {
	case DeferrableWriter:
		id := w.expectDeferred()

		// Potentially blocking, do any async tasks here first, then send result
		go func() { w.deferredChannel() <- settledSuspense{id, s.Children()} }()

		// Write the Fallback and return
		return s.Fallback.
			PushAttrs(html.Attr{"data-deferred", fmt.Sprint(id.int)}).
			WriteTo(w)

	default:
		// If this isn't called in the context of a DeferStream writer, we can
		// must render the children now.
		return html.Render(w, s.Children())
	}
}
func (Suspense) element() {}
func (s Suspense) String() string {
	b := strings.Builder{}
	s.WriteTo(&b)
	return b.String()
}
func (s Suspense) Bytes() []byte {
	b := bytes.Buffer{}
	s.WriteTo(&b)
	return b.Bytes()
}
