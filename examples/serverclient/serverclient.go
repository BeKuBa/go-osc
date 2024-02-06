package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/crgimenes/go-osc"
)

func main() {
	done := sync.WaitGroup{}
	done.Add(1)

	addr1, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
	addr2, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9000")

	// OSC app 1 with AddMsgHandlerExt
	d1 := osc.NewStandardDispatcher()
	app1 := osc.NewServerAndClient(d1)
	app1.NewConn(addr2, nil)

	d1.AddMsgHandlerExt("*", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("%v -> %v: %v \n", addr, addr2, msg)
		app1.SendMsgTo(addr, "/pong", 2)
	})

	go app1.ListenAndServe()

	// OSC app 2 with AddMsgHandler
	d2 := osc.NewStandardDispatcher()
	d2.AddMsgHandler("*", func(msg *osc.Message) {
		fmt.Printf("-> %v: %v \n", addr1, msg)
		done.Done()
	})

	app2 := osc.NewServerAndClient(d2)
	app2.NewConn(addr1, addr2)

	go app2.ListenAndServe()

	app2.SendMsg("/ping", 1.0)

	done.Wait()
}

// output:
// 127.0.0.1:8000 -> 127.0.0.1:9000: /ping ,d 1
// -> 127.0.0.1:8000: /pong ,i 2
