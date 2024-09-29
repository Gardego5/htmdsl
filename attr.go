package html

import "strings"

type (
	Attr  [2]string
	Attrs []Attr
)

func Class(list ...string) Attr { return Attr{"class", strings.Join(list, " ")} }
