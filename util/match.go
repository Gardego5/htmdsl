package util

import (
	"context"
	"reflect"

	html "github.com/Gardego5/htmdsl"
)

type match struct {
	val     any
	matched bool
}

var _ html.HTML = (*match)(nil)

func Match(val any) *match { return &match{val: val} }

func (m *match) When(val any, then ...any) *match {
	if !m.matched && reflect.DeepEqual(m.val, val) {
		m.matched = true
		m.val = then
	}
	return m
}
func (m *match) Default(then ...any) html.Fragment {
	if m.matched {
		return html.Fragment(m.val.([]any))
	} else {
		return then
	}
}
func (m match) Render(context.Context) html.RenderedHTML {
	if m.matched {
		return html.Fragment(m.val.([]any))
	} else {
		return html.Fragment{}
	}
}
