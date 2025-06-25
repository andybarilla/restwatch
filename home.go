package main

import (
	"net/http"
	"time"

	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func HomePage(msgs []PubSubMessage) Node {
	return page(PageProps{
		Title: "RestWatch",
	},
		Button(
			Class("btn btn-primary btn-sm"),
			Text("Clear Messages"),
			hx.Post("/clear-messages"),
		),
		Div(
			Class("rounded-box border border-base-content/5 bg-base-100 overflow-x-auto"),
			Table(
				Attr("sse-connect", "/sse-events?stream=all"),
				Class("table table-zebra table-pin-cols table-xs"),
				THead(
					Tr(
						Th(Text("Message")),
						Th(Text("Received At")),
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
								return StaticMessageRow(msg)
							}),
						),
					),
				),
			),
		),
	)
}

func StaticMessageRow(msg PubSubMessage) Node {
	return Tr(
		Td(
			Class("bg-accent"),
			Text(msg.RawMessage),
		),
		Td(
			Text(msg.ReceivedAt.Format(time.RFC3339)),
		),
	)
}

func MessageRow(msg PubSubMessage) Node {
	return Tr(
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
