package main

import (
	"log"
	"net/http"

	. "github.com/Gardego5/htmdsl"
	"github.com/Gardego5/htmdsl/util"
)

type Page [2]HTML

func (page Page) Render() RenderedHTML {
	return Fragment{DOCTYPE, Html{
		Head{
			Meta{{"charset", "UTF-8"}},
			Meta{{"name", "viewport"}, {"content", "width=device-width, initial-scale=1.0"}},
			Meta{{"http-equiv", "X-UA-Compatible"}, {"content", "ie=edge"}},
			page[0],
		},
		Body{
			page[1],
		},
	}}
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

func (list List) Render() RenderedHTML {
	items := make([]HTML, len(list.Items))
	for i, item := range list.Items {
		items[i] = Li{
			util.Match(list.BulletStyle).
				When(Bullet, "• ").
				When(Dash, "- ").
				When(Plus, "+ "),
			item,
		}
	}
	return Ul{items}.Render()
}

// we need to implement the HTMLComponent interface for a type component.
var _ HTML = (*Page)(nil)
var _ HTML = (*List)(nil)

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
						"hello",
						"fair",
						"world",
					},
				},
			},
		})
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
