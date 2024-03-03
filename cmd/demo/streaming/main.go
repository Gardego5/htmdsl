package main

import (
	"io"
	"net/http"
	"time"

	. "github.com/Gardego5/htmdsl"
	"github.com/Gardego5/htmdsl/suspense"
)

func layout(title string, children ...HTMLElement) HTMLElement {
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
	}}
}

func streamingDemo() io.WriterTo {
	return layout("Streaming demo",
		suspense.Suspense{
			Fallback: Span{"Loading 1..."},
			Children: func() any {
				time.Sleep(1 * time.Second)
				return Div{"Hello, world!"}
			},
		},
		suspense.Suspense{
			Fallback: Span{"Loading 2..."},
			Children: func() any {
				time.Sleep(1 * time.Second)
				return Div{"Hello, again!"}
			},
		},
	)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /normal", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		streamingDemo().WriteTo(w)
	})

	mux.HandleFunc("GET /chunked", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		dw := suspense.NewDeferrableWriter(w)
		streamingDemo().WriteTo(dw)
		dw.WriteDeferred()
	})

	http.ListenAndServe(":8080", mux)
}
