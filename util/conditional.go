package util

import html "github.com/Gardego5/htmdsl"

type (
	conditional struct {
		cond bool
		then any
	}
)

var _ html.HTML = (*conditional)(nil)

func Switch() conditional                { return conditional{} }
func If(bool bool, then any) conditional { return conditional{cond: bool, then: then} }

func (i conditional) Default(then any) any {
	if i.cond {
		return i.then
	} else {
		return then
	}
}
func (i conditional) Else(then any) any {
	if i.cond {
		return i.then
	} else {
		return then
	}
}
func (i conditional) Case(cond bool, then any) conditional {
	if i.cond {
		return i
	} else {
		return conditional{cond: cond, then: then}
	}
}
func (i conditional) ElseIf(cond bool, then any) conditional {
	if i.cond {
		return i
	} else {
		return conditional{cond: cond, then: then}
	}
}

func (i conditional) Render() html.HTMLElement {
	if i.cond {
		return html.Fragment{i.then}
	}
	return nil
}
