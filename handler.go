package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type PubSubMessage struct {
	Id           string    `json:"id"`
	Subscription string    `json:"subscription"`
	Message      Message   `json:"Message"`
	RawMessage   string    `json:"RawMessage"`
	ReceivedAt   time.Time `json:"receivedAt"`
}

type Message struct {
	PublishTime   string            `json:"publishTime"`
	Data          string            `json:"data"`
	MessageId     string            `json:"messageId"`
	Attributes    map[string]string `json:"attributes"`
	ExtractedData string
}

func messageHandler(statusChannel chan PubSubMessage, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("could not read body: %s", err)
		}

		msg := PubSubMessage{
			RawMessage: string(body),
			ReceivedAt: time.Now(),
		}
		logger.Info(fmt.Sprintf("Received message: %+v", msg))

		//err = json.Unmarshal(body, &msg)
		//if err != nil {
		//	fmt.Printf("could not unmarshal body")
		//}
		//
		//decodedData, err := base64.StdEncoding.DecodeString(msg.Message.Data)
		//if err != nil {
		//	fmt.Printf("could not decode data: %s\n", err)
		//	return
		//}
		//
		//if len(decodedData) != 0 {
		//	var data map[string]interface{}
		//	err = json.Unmarshal(decodedData, &data)
		//	if err != nil {
		//		fmt.Printf("could not unmarshal body")
		//	}
		//	msg.Message.ExtractedData = string(decodedData)
		//}

		statusChannel <- msg
	}
}
