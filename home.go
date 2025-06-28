package main

import (
	"net/http"
	"time"

	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func HomePage(msgs []PubSubMessage) Node {
	return page(PageProps{
		Title: "RestWatch",
	},
		Div(
			x.Data("{ openDrawer: false }"),
			Class("drawer drawer-end"),
			Input(
				ID("my-drawer"),
				Type("checkbox"),
				Class("drawer-toggle"),
				x.Bind("checked", "openDrawer"),
			),
			Div(
				Button(
					Class("btn btn-primary btn-sm"),
					Text("Clear Messages"),
					hx.Post("/clear-messages"),
				),
				Div(
					Class("drawer-content rounded-box border border-base-content/5 bg-base-100 overflow-x-auto"),
					Table(
						Attr("sse-connect", "/sse-events?stream=all"),
						Class("table table-auto table-sm border-collapse"),
						THead(
							Tr(
								Th(
									Class("text-xs overflow-auto w-3/4"),
									Text("Message"),
								),
								Th(
									Class("text-xs overflow-auto w-1/4"),
									Text("Received At"),
								),
							),
						),
						TBody(
							Attr("sse-swap", "incoming-messages"),
							Attr("hx-swap", "beforeend"),
							Tr(
								Class("no-data-row"),
								Td(
									ColSpan("2"),
									Text("No messages yet."),
								),
							),
							If(len(msgs) > 0,
								Map(msgs, func(msg PubSubMessage) Node {
									return MessageRow(msg)
								}),
							),
						),
					),
				),
			),
			Div(
				Class("drawer-side"),
				Label(
					Class("drawer-overlay"),
					x.On("click", "openDrawer = false"),
				),
				Ul(
					Class("menu min-h-full w-1/2 p-4 text-base-content bg-base-200 border-l border-base-100"),
					Li(Text("hello")),
					Li(Text("shitbird")),
				),
			),
		),
	)
}

func MessageRow(msg PubSubMessage) Node {
	return Tr(
		Class("hover:bg-base-200 cursor-pointer"),
		x.On("click", "openDrawer = true"),
		Td(
			Text(msg.RawMessage),
		),
		Td(
			Text(msg.ReceivedAt.Format(time.RFC3339)),
		),
	)
}

func (s *Server) clearMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.log.Info("Clearing messages")
		s.messages = []PubSubMessage{}
	}
}
