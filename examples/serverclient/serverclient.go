package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"bekuba.de/go-osc"
)

func main() {
	done := sync.WaitGroup{}
	done.Add(1)

	addr1 := "127.0.0.1:8000"
	addr2 := "127.0.0.1:9000"

	// OSC app 1 with AddMsgHandlerExt

	app1, err := osc.NewNode(addr2)
	if err != nil {
		log.Fatal(err)
	}
	defer app1.Close()

	d1 := osc.NewStandardDispatcher()
	d1.AddMsgHandlerExt("*", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("%v -> %v: %v \n", addr, addr2, msg)
		app1.SendMsgTo(addr1, "/pong", 2)
	})

	go app1.ListenAndServe(d1)

	// OSC app 2 with AddMsgHandler

	d2 := osc.NewStandardDispatcher()
	d2.AddMsgHandler("*", func(msg *osc.Message) {
		fmt.Printf("-> %v: %v \n", addr1, msg)
		done.Done()
	})

	app2, err := osc.NewNode(addr1)
	if err != nil {
		log.Fatal(err)
	}
	defer app2.Close()

	go app2.ListenAndServe(d2)

	app2.SendMsgTo(addr2, "/ping", 1.0)

	done.Wait()
}

// output:
// 127.0.0.1:8000 -> 127.0.0.1:9000: /ping ,d 1
// -> 127.0.0.1:8000: /pong ,i 2
