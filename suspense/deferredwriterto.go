package suspense

import (
	"io"

	html "github.com/Gardego5/htmdsl"
)

type DeferredWriterTo interface {
	io.WriterTo
	Id() Id
}

type settledSuspense struct {
	id    Id
	value any
}

func (s settledSuspense) Id() Id { return s.id }
func (s settledSuspense) WriteTo(w io.Writer) (int64, error) {
	return html.Render(w, s.value)
}
