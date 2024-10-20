package main

import (
	"fmt"
	"log"

	"github.com/bekuba/go-osc"
)

func main() {
	addr := "127.0.0.1:8765"

	server, err := osc.NewServerAndClient(addr)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("/message/address", func(msg *osc.Message) {
		fmt.Printf("Received message from %v: %v", addr, msg)
	})

	server.ListenAndServe(d)
}
