package main

import (
	"net/http"
)

//
////go:embed public
//var staticContent embed.FS
//var staticContentFS, _ = fs.Sub(staticContent, "public")

func Static(mux *http.ServeMux) {
	h := http.FileServer(http.Dir("public"))
	//mux.Handle(`/{:[^.]+\.[^.]+}`, h)
	//mux.Handle(`/{:images|scripts|styles}/*`, h)
	mux.Handle("/styles/", h)
	mux.Handle("/scripts/", h)
}
