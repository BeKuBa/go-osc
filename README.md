# GoOSC

[![Build Status](https://travis-ci.org/crgimenes/go-osc.png?branch=master)](https://travis-ci.org/crgimenes/go-osc) [![GoDoc](https://godoc.org/github.com/crgimenes/go-osc/osc?status.svg)](https://godoc.org/github.com/crgimenes/go-osc/osc) [![Coverage Status](https://coveralls.io/repos/github/crgimenes/go-osc/badge.svg?branch=master)](https://coveralls.io/github/crgimenes/go-osc?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/crgimenes/go-osc)](https://goreportcard.com/report/github.com/crgimenes/go-osc)

[Open Sound Control (OSC)](http://opensoundcontrol.org/introduction-osc) library for Golang. Implemented in pure Go.

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

import "github.com/crgimenes/go-osc"

func main() {
    addr := "127.0.0.1:8765"
    d := osc.NewStandardDispatcher()
    d.AddMsgHandler("/message/address", func(msg *osc.Message) {
        osc.PrintMessage(msg)
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
make test
```
