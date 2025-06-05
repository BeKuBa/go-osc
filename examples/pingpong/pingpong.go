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
	done.Add(3)

	addr1 := "127.0.0.1:8000"
	addr2 := "127.0.0.1:9000"
	addr3 := "127.0.0.1:9001"

	app1, err := osc.NewNode(addr2)
	if err != nil {
		log.Fatal(err)
	}
	defer app1.Close()

	d1 := osc.NewStandardDispatcher()
	err = d1.AddMsgHandlerExt("*", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("%v -> %v: %v \n", addr, addr2, msg)
		err = app1.SendMsgTo(addr.String(), "/pong", 2)
		if err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		err := app1.ListenAndServe(d1)
		if err != nil {
			fmt.Println(err)
		}
	}()

	d2 := osc.NewStandardDispatcher()
	err = d2.AddMsgHandlerExt("*", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("%v -> %v: %v \n", addr, addr1, msg)
		done.Done()
	})
	if err != nil {
		fmt.Println(err)
	}

	app2, err := osc.NewNode(addr1)
	if err != nil {
		log.Fatal(err)
	}
	defer app2.Close()

	go func() {
		err := app2.ListenAndServe(d2)
		if err != nil {
			fmt.Println(err)
		}
	}()

	d3 := osc.NewStandardDispatcher()
	err = d3.AddMsgHandler("*", func(msg *osc.Message) {
		fmt.Printf("?? -> 127.0.0.1:9001 %v \n", msg)
		done.Done()
	})
	if err != nil {
		fmt.Println(err)
	}

	app3, err := osc.NewNode(addr3)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err := app3.ListenAndServe(d3)
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = app2.SendMsgTo(addr2, "/ping", 1.0)
	if err != nil {
		fmt.Println(err)
	}

	err = app3.SendMsgTo(addr1, "/pong", 3)
	if err != nil {
		fmt.Println(err)
	}

	err = app2.SendMsgTo(addr3, "/pong", 4.0)
	if err != nil {
		fmt.Println(err)
	}

	done.Wait()
}

// output:
// 127.0.0.1:9001 -> 127.0.0.1:8000: /pong ,i 3
// 127.0.0.1:8000 -> 127.0.0.1:9000: /ping ,d 1
// ?? -> 127.0.0.1:9001 /pong ,d 4
// 127.0.0.1:9000 -> 127.0.0.1:8000: /pong ,i 2
