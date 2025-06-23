package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

//
////go:embed public
//var staticContent embed.FS
//var staticContentFS, _ = fs.Sub(staticContent, "public")

func Static(mux *chi.Mux) {
	h := http.FileServer(http.Dir("public"))
	mux.Handle(`/{:[^.]+\.[^.]+}`, h)
	mux.Handle(`/{:images|scripts|styles}/*`, h)
}
