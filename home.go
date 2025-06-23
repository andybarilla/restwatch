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
				Attr("sse-connect", "/sse-events"),
				Attr("sse-swap", "incoming-messages"),
				ID("incoming-messages"),
				Class("table"),
				THead(
					Tr(
						Th(Text("Message")),
					),
					TBody(
						Map(msgs, func(msg PubSubMessage) Node {
							return MessageRow(msg)
						}),
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
