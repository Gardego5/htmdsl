package main

import (
	"io"
	"net/http"
	"time"

	. "github.com/Gardego5/htmdsl"
	"github.com/Gardego5/htmdsl/suspense"
)

func layout(title string, children ...HTMLElement) HTMLElement {
	return suspense.SyncArea(DOCTYPE, Html{Attrs{{"lang", "en"}},
		Head{
			Title{title},
			Meta{{"charset", "UTF-8"}},
			Meta{{"name", "viewport"}, {"content", "width=device-width, initial-scale=1.0"}},
			Meta{{"http-equiv", "X-UA-Compatible"}, {"content", "ie=edge"}},
		},
		Body{
			Script{PreEscaped(`
				const observer = new MutationObserver(function (mutations, observer) {
					for (const mutation of mutations) {
						if (mutation.type === 'childList') {
							for (const node of mutation.addedNodes) {
								const promised = node.dataset.promised;
								if (typeof promised !== 'undefined') {
									document.querySelector("[data-deferred='" + promised + "']").replaceWith(node.content);
									node.remove();
								}
							}
						}
					}
				});
				const el = document.querySelector('body');
				observer.observe(el, { subtree: true, childList: true, attributes: true, characterData: true });
			`)},
			Main{
				Class("font-sans max-w-xl mx-auto flex flex-col gap-2"),
				children,
			},
		},
	})
}

func slowComponent(name string, delay int32) HTMLElement {
	return suspense.Suspense{
		Fallback: Div{"Loading ", name, "..."},
		Children: func() any {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			return Div{
				H2{"Slow Component ", name},
				P{"This is a slow component. ",
					"I take ", delay, " milliseconds to load."},
			}
		},
	}
}

func streamingDemo() io.WriterTo {
	return layout("Streaming demo",
		slowComponent("A", 1000),
		slowComponent("B", 1500),
		slowComponent("C", 300),
		slowComponent("D", 400),
	)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		streamingDemo().WriteTo(w)
	})

	http.ListenAndServe(":8080", mux)
}
