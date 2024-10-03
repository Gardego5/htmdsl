package util

import "github.com/Gardego5/htmdsl"

type Block func() any

func (f Block) Render() html.RenderedHTML { return html.Fragment{f()} }
