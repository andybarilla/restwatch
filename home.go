package main

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HomePage(msgs []PubSubMessage) Node {
	return page(PageProps{
		Title: "RestWatch",
	},
		Div(
			Table(
				Attr("sse-connect", "/sse-events?stream=all"),
				Class("table"),
				THead(
					Tr(
						Th(Text("Message")),
					),
					TBody(
						Attr("sse-swap", "incoming-messages"),
						Attr("hx-swap", "beforeend"),
						If(len(msgs) == 0,
							Tr(
								Td(
									Text("No messages yet."),
								),
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
	)
}

func MessageRow(msg PubSubMessage) Node {
	return Tr(
		Td(
			Text(msg.RawMessage),
		),
	)
}
