package html

import (
	"fmt"
	"strings"
)

type Attrs map[string]any

func Class(list ...string) Attrs { return Attrs{"class": strings.Join(list, " ")} }
func Id(id any) Attrs            { return Attrs{"id": fmt.Sprint(id)} }
