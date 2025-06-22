package main

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HomePage() Node {
	return page(PageProps{
		Title: "RestWatch",
	},
		Div(
			Class("p-4 bg-gray-100 rounded shadow"),
			Text("hello world"),
		),
	)
}
