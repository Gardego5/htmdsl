package main

import (
	"context"
	"net/http"

	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
)

type key int

const keySession key = iota

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), keySession, "session-id")))
	})
}

func getSessionId(c context.Context) string {
	if val, ok := c.Value(keySession).(string); ok {
		return val
	} else {
		return ""
	}
}

type SessionComponent struct{}

var _ ContextHTML = (*SessionComponent)(nil)

func (c SessionComponent) Render() RenderedHTML { return c.RenderWithContext(nil) }
func (SessionComponent) RenderWithContext(c context.Context) RenderedHTML {
	return P{
		"Session ID: ",
		If(c != nil, func() any { return getSessionId(c) }).Else("nil"),
	}.Render()
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RenderWithContext(w, r.Context(), Fragment{DOCTYPE, Html{
			Head{},
			Body{
				SessionComponent{},
			},
		}})
	}))

	http.ListenAndServe(":8080", sessionMiddleware(mux))
}
