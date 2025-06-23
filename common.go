package main

import (
	"crypto/sha256"
	"fmt"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"

	"os"
	"path/filepath"
	"strings"
	"sync"
)

type PageProps struct {
	Title       string
	Description string
}

var hashOnce sync.Once
var appCSSPath string
var htmxJSPath string
var htmxSseExtJSPath string
var debugJSPath string

func page(props PageProps, children ...Node) Node {
	hashOnce.Do(func() {
		appCSSPath = getHashedPath("public/styles/app.css")
		htmxJSPath = getHashedPath("public/scripts/htmx.min.js")
		htmxSseExtJSPath = getHashedPath("public/scripts/htmx-ext-sse.min.js")
		debugJSPath = getHashedPath("public/scripts/debug.js")
	})

	return HTML5(HTML5Props{
		Title:       props.Title,
		Description: props.Description,
		Language:    "en",
		Head: []Node{
			Link(Rel("stylesheet"), Href(appCSSPath)),
			Script(Src(htmxJSPath), Defer()),
			Script(Src(htmxSseExtJSPath), Defer()),
			Script(Src(debugJSPath), Defer()),
		},
		Body: []Node{
			hx.Ext("sse"),
			Div(
				Class("container mx-auto p-4"),
				header(),
				Group(children),
			),
		},
	},
	)
}

func header() Node {
	return Div(Text("RestWatch"))
}

func getHashedPath(path string) string {
	externalPath := strings.TrimPrefix(path, "public")
	ext := filepath.Ext(path)
	if ext == "" {
		panic("no extension found")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("%v%v", strings.TrimSuffix(externalPath, ext), ext)
	}

	return fmt.Sprintf("%v%v?%x", strings.TrimSuffix(externalPath, ext), ext, sha256.Sum256(data))
}
