package util

import html "github.com/Gardego5/htmdsl"

type conditional struct {
	cond bool
	then html.Fragment
}

var _ html.IndirectHTML = (*conditional)(nil)

func Switch() *conditional                   { return &conditional{} }
func If(bool bool, then ...any) *conditional { return &conditional{cond: bool, then: then} }

func (i *conditional) Default(then ...any) any {
	if i.cond {
		return i.then
	} else {
		return then
	}
}
func (i *conditional) Else(then ...any) any {
	if i.cond {
		return i.then
	} else {
		return then
	}
}
func (i *conditional) Case(cond bool, then ...any) *conditional {
	if !i.cond {
		i.cond = cond
		i.then = then
	}

	return i
}
func (i *conditional) ElseIf(cond bool, then ...any) *conditional {
	if !i.cond {
		i.cond = cond
		i.then = then
	}

	return i
}
func (i *conditional) Render() html.Fragment {
	if i.cond {
		return i.then
	}

	return nil
}
