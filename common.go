package main

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	"sync"
)

type PageProps struct {
	Title       string
	Description string
}

var hashOnce sync.Once
var appCSSPath string

func page(props PageProps, children ...Node) Node {
	hashOnce.Do(func() {
		appCSSPath = getHashedPath("static/styles/app.css")
	})

	return HTML5(HTML5Props{
		Title:       props.Title,
		Description: props.Description,
		Language:    "en",
		Head: []Node{
			Link(Rel("stylesheet"), Href(appCSSPath)),
		},
		Body: []Node{
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
	return "/styles/app.css"
	//externalPath := strings.TrimPrefix(path, "static")
	//ext := filepath.Ext(path)
	//if ext == "" {
	//	panic("no extension found")
	//}
	//
	//data, err := staticContent.ReadFile(path)
	//if err != nil {
	//	return fmt.Sprintf("%v.x%v", strings.TrimSuffix(externalPath, ext), ext)
	//}
	//
	//return fmt.Sprintf("%v.%x%v", strings.TrimSuffix(externalPath, ext), sha256.Sum256(data), ext)
}
