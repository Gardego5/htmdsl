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
		ctx := r.Context()
		ctx = context.WithValue(ctx, keySession, "session-id")
		next.ServeHTTP(w, r.WithContext(ctx))
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

func (SessionComponent) Render(ctx context.Context) RenderedHTML {
	return P{
		"Session ID: ",
		If(ctx != nil, func() any { return getSessionId(ctx) }).Else("nil"),
	}.Render(ctx)
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RenderContext(w, r.Context(), Fragment{DOCTYPE, Html{
			Head{},
			Body{
				SessionComponent{},
			},
		}})
	}))

	http.ListenAndServe(":8080", sessionMiddleware(mux))
}
