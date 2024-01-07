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
  - 
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

import "github.com/crgimenes/go-osc"

func main() {
    addr := "127.0.0.1:8765"
    d := osc.NewStandardDispatcher()
    d.AddMsgHandler("/message/address", func(msg *osc.Message) {
        fmt.Println(msg)
    })

    server := &osc.Server{
        Addr: addr,
        Dispatcher:d,
    }
    server.ListenAndServe()
}
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
