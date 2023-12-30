package main

import (
	"fmt"
	"github.com/crgimenes/go-osc"
	"net"
	"sync"
)

func main() {
	done := sync.WaitGroup{}
	done.Add(2)

	addr1, err := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
	}

	addr2, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println(err)
	}

	addr3, err := net.ResolveUDPAddr("udp", "127.0.0.1:9001")
	if err != nil {
		fmt.Println(err)
	}

	d1 := osc.NewStandardDispatcher()
	app1 := osc.NewServerAndClient(d1)
	err = app1.NewConn(addr2, addr1)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := app1.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = d1.AddMsgHandler("*", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("%v -> %v: %v \n", addr, addr2, msg)
		err = app1.SendMsg("/pong", 2)
		if err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		err := app1.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	d2 := osc.NewStandardDispatcher()
	err = d2.AddMsgHandler("*", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("%v -> %v: %v \n", addr, addr1, msg)
		done.Done()
	})
	if err != nil {
		fmt.Println(err)
	}

	app2 := osc.NewServerAndClient(d2)
	err = app2.NewConn(addr1, addr2)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := app2.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := app2.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	d3 := osc.NewStandardDispatcher()
	err = d3.AddMsgHandler("*", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("d3 %v \n", msg)
		done.Done()
	})
	if err != nil {
		fmt.Println(err)
	}

	app3 := osc.NewServerAndClient(d3)
	err = app3.NewConn(addr3, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := app3.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := app3.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = app2.SendMsg("/ping", 1.0)
	if err != nil {
		fmt.Println(err)
	}

	err = app3.SendMsgTo(addr1, "/pong", 3)
	if err != nil {
		fmt.Println(err)
	}

	done.Wait()
}

// output:
//127.0.0.1:9001 d2 /pong ,i 3
//127.0.0.1:8000 d1 /ping ,d 1
//127.0.0.1:9000 d2 /pong ,i 2
