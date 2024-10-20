# GoOSC

[![Go Report Card](https://goreportcard.com/badge/github.com/bekuba/go-osc)](https://goreportcard.com/report/github.com/bekuba/go-osc)

[Open Sound Control (OSC)](https://opensoundcontrol.stanford.edu) library for Golang. Implemented in pure Go.



This repository is a heavily modified fork of the [original go-osc](https://github.com/hypebeast/go-osc) started by https://github.com/crgimenes. Please consider using the original project [original go-osc](https://github.com/hypebeast/go-osc).

Version 0.0.1 is not compatible with further versions. But it is easy to migrate.

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
go get github.com/bekuba/go-osc
```

## Usage

See exambles.</br>
There is only one object that can act as "OSC server" AND/OR "OSC client".
( OSC Spec. 1.0: "Any application that sends OSC Packets is an OSC Client; any application that receives OSC Packets is an OSC Server.")

### Server and Client Ping Pong

```go
func main() {
	done := sync.WaitGroup{}
	done.Add(1)

	addr1 := "127.0.0.1:8000"
	addr2 := "127.0.0.1:9000"

	// OSC app 1 with AddMsgHandlerExt

	app1, err := osc.NewServerAndClient(addr2)
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

	app2, err := osc.NewServerAndClient(addr1)
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
```
