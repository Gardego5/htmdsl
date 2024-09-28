package main

import (
	"net/http"

	. "github.com/Gardego5/htmdsl"
	"github.com/Gardego5/htmdsl/util"
)

type Page [2]HTMLElement

func (page Page) Render() HTMLElement {
	return Fragment{DOCTYPE,
		Head{
			Meta{{"charset", "UTF-8"}},
			Meta{{"name", "viewport"}, {"content", "width=device-width, initial-scale=1.0"}},
			Meta{{"http-equiv", "X-UA-Compatible"}, {"content", "ie=edge"}},
			page[0],
		},
		Body{
			page[1],
		},
	}
}

type (
	BulletStyle int
	List        struct {
		BulletStyle BulletStyle
		Items       []string
	}
)

const (
	Bullet BulletStyle = iota
	Dash
	Plus
)

func (list List) Render() HTMLElement {
	items := make([]HTMLElement, len(list.Items))
	for i, item := range list.Items {
		items[i] = Li{
			util.Match(list.BulletStyle).
				When(Bullet, "• ").
				When(Dash, "- ").
				When(Plus, "+ "),
			item,
		}
	}
	return Ul{items}
}

// we need to implement the HTMLComponent interface for a type component.
var _ HTMLComponent = (*Page)(nil)
var _ HTMLComponent = (*List)(nil)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		title := "Hello, World!"
		Render(w, Page{
			Fragment{
				Title{title},
				Style{PreEscaped(`
ul {
	list-style-type: none;
}
`)},
			},
			Fragment{
				H1{title},
				List{
					BulletStyle: Dash,
					Items: []string{
						"Item 1",
						"Item 2",
						"Item 3",
					},
				},
			},
		})
	})
	http.ListenAndServe(":8080", mux)
}
