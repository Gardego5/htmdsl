package util

import (
	"context"

	html "github.com/Gardego5/htmdsl"
)

type gconditional[T any] struct {
	cond bool
	then T
}

func GSwitch[T any]() gconditional[T]             { return gconditional[T]{} }
func GIf[T any](cond bool, val T) gconditional[T] { return gconditional[T]{cond: cond, then: val} }

func (i *gconditional[T]) Default(then T) T {
	if i.cond {
		return i.then
	} else {
		return then
	}
}

func (i *gconditional[T]) Else(then T) T {
	if i.cond {
		return i.then
	} else {
		return then
	}
}

func (i *gconditional[T]) Case(cond bool, then T) *gconditional[T] {
	if !i.cond {
		i.cond = cond
		i.then = then
	}

	return i
}

func (i *gconditional[T]) ElseIf(cond bool, then T) *gconditional[T] {
	if !i.cond {
		i.cond = cond
		i.then = then
	}

	return i
}

func (i *gconditional[T]) Render(context.Context) html.RenderedHTML {
	if i.cond {
		return html.Fragment{i.then}
	}

	return nil
}
