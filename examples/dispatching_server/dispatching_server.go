package main

import (
	"fmt"

	"github.com/crgimenes/go-osc"
)

func main() {
	addr := "127.0.0.1:8765"

	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("/message/address", func(msg *osc.Message) {
		fmt.Printf("Received message from %v: %v", addr, msg)
	})

	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}

	server.ListenAndServe()
}
