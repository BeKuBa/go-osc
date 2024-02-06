# GoOSC

[![Go Report Card](https://goreportcard.com/badge/github.com/crgimenes/go-osc)](https://goreportcard.com/report/github.com/crgimenes/go-osc)

[Open Sound Control (OSC)](https://opensoundcontrol.stanford.edu) library for Golang. Implemented in pure Go.

---

This repository is a heavily modified fork of the [original go-osc](https://github.com/hypebeast/go-osc). Please consider using the original project.

---


## Features

- OSC Bundles, including timetags
- OSC Messages
- OSC Client
- OSC Server
- Supports the following OSC argument types:
  - 'i' (Int32)
  - 'f' (Float32)
  - 's' (string)
  - 'b' (blob / binary data)
  - 'h' (Int64)
  - 't' (OSC timetag)
  - 'd' (Double/int64)
  - 'T' (True)
  - 'F' (False)
  - 'N' (Nil)
- Support for OSC address pattern including '\*', '?', '{,}' and '[]' wildcards

## Install

```shell
go get github.com/crgimenes/go-osc
```

## Usage

### Client

```go
import "github.com/crgimenes/go-osc"

func main() {
    client := osc.NewClient("localhost", 8765)
    msg := osc.NewMessage("/osc/address")
    msg.Append(int32(111))
    msg.Append(true)
    msg.Append("hello")
    client.Send(msg)
}
```

### Server

```go
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
        Addr: addr,
        Dispatcher:d,
    }
    server.ListenAndServe()
}
```
### Server and Client
```go
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
```

## Tests

```shell
go test ./...
```

# Contributing

- Fork the repo on GitHub
- Clone the project to your own machine
- Create a *branch* with your modifications `git checkout -b fantastic-feature`.
- Then _commit_ your changes `git commit -m 'Implementation of new fantastic feature'`
- Make a _push_ to your _branch_ `git push origin fantastic-feature`.
- Submit a **Pull Request** so that we can review your changes
