package html

import "strings"

type Attrs map[string]any

func Class(list ...string) Attrs { return Attrs{"class": strings.Join(list, " ")} }
func Id(id string) Attrs         { return Attrs{"id": id} }
