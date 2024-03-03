package suspense

import (
	"io"

	html "github.com/Gardego5/htmdsl"
)

type syncArea struct {
	html.Fragment
}

func SyncArea(children ...html.HTMLElement) syncArea {
	return syncArea{Fragment: html.Fragment(children)}
}
func (sa syncArea) WriteTo(w io.Writer) (int64, error) {
	dw := NewDeferrableWriter(w)
	n, err := sa.Fragment.WriteTo(dw)
	if err != nil {
		return n, err
	}
	nn, err := dw.WriteDeferred()
	n += int64(nn)
	return n, err
}
