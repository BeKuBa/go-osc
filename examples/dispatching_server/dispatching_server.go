package main

import (
	"fmt"

	"github.com/crgimenes/go-osc"
)

func main() {
	addr := "127.0.0.1:8765"

	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("/message/address", func(msg *osc.Message) {
		fmt.Println("Received message:", msg)
	})

	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}

	server.ListenAndServe()
}
