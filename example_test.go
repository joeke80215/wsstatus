package wsstatus

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// ExampleHandler example use this package
func ExampleHandler() {
	wsh := NewHandler()
	wsh.Registry("start", []byte(`{"state":"start"}`))
	wsh.Registry("aStep", []byte(`{"state":"a Step"}`))
	wsh.Registry("bStep", []byte(`{"state":"b Step"}`))
	wsh.Registry("cStep", []byte(`{"state":"c Step"}`))
	wsh.Registry("dStep", []byte(`{"state":"d Step","timestamp": "%s"}`))
	wsh.Registry("completed", []byte(`{"state":"completed"}`))

	// execute websocket client
	u := &url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws/go"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	go func() {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("receive: %s", message)

		wsh.Next()
		// NextTo("cStep")
	}()

	for {
		if wsh.GetStatus() == "dStep" {
			err = c.WriteMessage(websocket.TextMessage, wsh.GetPayload(time.Now().Format(time.RFC3339)))
		} else {
			err = c.WriteMessage(websocket.TextMessage, wsh.GetPayload())
		}

		if err != nil {
			log.Fatal(err)
		}
	}
}
