package main

import (
	"fmt"
	"net"

	"github.com/crgimenes/go-osc"
)

func main() {
	addr := "127.0.0.1:8765"

	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("/message/address", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("Received message from %v: %v", addr, msg)
	})

	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}

	server.ListenAndServe()
}
