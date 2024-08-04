package html

import (
	"io"
	"strings"
)

type literal string

const DOCTYPE literal = "<!DOCTYPE html>"

func (lit literal) element()       {}
func (lit literal) String() string { return string(lit) }
func (lit literal) Bytes() []byte  { return []byte(lit) }
func (lit literal) Reader() io.Reader { return strings.NewReader(string(lit)) }
func (lit literal) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(lit))
	return int64(n), err
}
