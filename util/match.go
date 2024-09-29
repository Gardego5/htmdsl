package util

import (
	"reflect"

	html "github.com/Gardego5/htmdsl"
)

type match struct {
	val     any
	matched bool
}

var _ html.HTML = (*match)(nil)

func Match(val any) match { return match{val: val} }

func (m match) When(val any, then any) match {
	if !m.matched && reflect.DeepEqual(m.val, val) {
		return match{matched: true, val: then}
	} else {
		return m
	}
}
func (m match) Default(then any) any {
	if m.matched {
		return m.val
	} else {
		return then
	}
}
func (m match) Render() html.RenderedHTML { return html.Fragment{m.val} }
