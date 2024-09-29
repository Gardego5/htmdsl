package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	. "github.com/Gardego5/htmdsl"
)

func layout(title string, children ...HTML) HTML {
	return Fragment{DOCTYPE, Html{Attrs{{"hidden"}, {"lang", "en"}},
		Head{
			Title{title},
			Meta{{"charset", "UTF-8"}},
			Meta{{"name", "viewport"}, {"content", "width=device-width, initial-scale=1.0"}},
			Meta{{"http-equiv", "X-UA-Compatible"}, {"content", "ie=edge"}},

			// Tailwind
			Script{Attrs{{"type", "module"}, {"src", "https://cdn.skypack.dev/twind/shim"}}},
			Script{Attrs{{"type", "twind-config"}}, PreEscaped(`{"theme":{"fontFamily":{"sans":["Rokkitt","sans-serif"]}}}`)},

			// Google Fonts
			Link{{"rel", "preconnect"}, {"href", "https://fonts.googleapis.com"}},
			Link{{"rel", "preconnect"}, {"href", "https://fonts.gstatic.com"}, {"crossorigin", ""}},
			Link{{"rel", "stylesheet"}, {"href", "https://fonts.googleapis.com/css2?family=Rokkitt&display=swap"}},

			// Alpine.js
			Script{Attrs{{"src", "https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"}, {"defer"}}},
		},
		Body{Main{
			Class("font-sans max-w-xl mx-auto flex flex-col gap-2"),
			children,
		}},
	}}
}

func index(w http.ResponseWriter, r *http.Request) {
	const NAME, VALUE = "~name", "~value"

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		name, value := r.FormValue(NAME), r.FormValue(VALUE)
		r.Form.Del(NAME)
		r.Form.Del(VALUE)
		if name != "" && value != "" {
			r.Form.Add(name, value)
		}
	}

	keys := make([]string, 0, len(r.Form))
	for key := range r.Form {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	tabledata := map[string]string{}
	for _, key := range keys {
		tabledata[key] = key
	}
	data, err := json.Marshal(tabledata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	table := make(Table, 0, len(keys)+3)
	table = append(table, Class("w-full"),
		Attrs{{"x-data", strings.ReplaceAll(string(data), "\"", "'")}},
		Tr{Class("text-left"), Th{"Name"}, Th{"Value"}})
	for _, key := range keys {
		values := r.Form[key]
		slices.Sort(values)
		for i, value := range values {
			tr := make(Tr, 0, 2)
			if i == 0 {
				tr = append(tr, Td{Attrs{Class("align-top"), {"rowspan", fmt.Sprint(len(values))}},
					Input{Class("w-full rounded px-2"), {"type", "text"}, {"x-model", key}}})
			}
			tr = append(tr, Td{Input{Class("w-full rounded px-2"), {"type", "text"}, {":name", key}, {"value", value}}})
			table = append(table, tr)
		}
	}

	w.Header().Add("Content-Type", "text/html")
	page := layout("Form Demo",
		H1{Class("text-center text-2xl"), "Form Demo"},
		P{"You can add query parameters to see them displayed on this page."},
		Form{Class("grid gap-2"), Attrs{{"method", "POST"}, {"action", "/"}},
			Div{Class("p-2 rounded border bg-gray-100"), table},
			Div{Class("grid grid-cols-2 gap-2"),
				Label{Attrs{{"for", NAME}}, "Name"},
				Input{Class("rounded border px-2"), {"type", "text"}, {"name", NAME}, {"id", NAME}},
				Label{Attrs{{"for", VALUE}}, "Value"},
				Input{Class("rounded border px-2"), {"type", "text"}, {"name", VALUE}, {"id", VALUE}},
				Button{Attrs{Class("rounded border px-2 col-span-full"), {"type", "submit"}}, "Submit"},
			},
		},
	)

	Render(w, page)
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(index))

	http.ListenAndServe(":8080", mux)
}
