package suspense

import (
	"fmt"
	"io"
	"net/http"

	html "github.com/Gardego5/htmdsl"
)

type (
	Id struct{ int }

	DeferrableWriter interface {
		io.Writer
		deferredChannel() chan DeferredWriterTo
		expectDeferred() Id
		WriteDeferred() (int, error)
	}

	deferrableWriter struct {
		*http.ResponseController
		io.Writer
		ch        chan DeferredWriterTo
		expecting map[Id]struct{}
		counter   int
	}
)

func NewDeferrableWriter(rw http.ResponseWriter) DeferrableWriter {
	rc := http.NewResponseController(rw)
	return &deferrableWriter{rc, rw, make(chan DeferredWriterTo), make(map[Id]struct{}), 0}
}

func (w *deferrableWriter) deferredChannel() chan DeferredWriterTo {
	return w.ch
}
func (w *deferrableWriter) expectDeferred() Id {
	id := Id{w.counter}
	w.expecting[id] = struct{}{}
	w.counter++
	return id
}
func (w *deferrableWriter) WriteDeferred() (int, error) {
	w.Flush()
	n := 0

	for len(w.expecting) > 0 {
		child := <-w.ch
		id := child.Id()
		if _, exists := w.expecting[id]; exists {
			delete(w.expecting, id)
			nn, err := html.Template{
				html.Attrs{{"data-promised", fmt.Sprint(id.int)}},
				child,
			}.WriteTo(w)
			n += int(nn)
			if err != nil {
				return n, err
			}
			w.Flush()
		} else {
			return n, fmt.Errorf("unexpected promise id %d", id.int)
		}
	}

	return n, nil
}
