package main

import (
	"fmt"
	"github.com/crgimenes/go-osc"
	"net"
	"time"
)

const (
	xrIP = "10.0.1.174"
	//Port for XR18
	xr18Port = 10024
	//Port for XR32
	xr32Port = 10023
)

// XR18 example
// printout all mixer messages
// send /xinfo, /xremote, /status
func main() {

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", xrIP, xr18Port))
	if err != nil {
		fmt.Println(err)
	}

	d := osc.NewStandardDispatcher()
	err = d.AddMsgHandler("*", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("xr %v: %v  \n", addr, msg)
	})
	if err != nil {
		fmt.Println(err)
	}

	app := osc.NewServerAndClient(d)
	err = app.NewConn(nil, addr)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := app.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err = app.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	app.SendMsg("/xinfo")

	for {
		// keepp connection alive (for multi client usage)
		app.SendMsg("/xremote")
		// show status of xair
		app.SendMsg("/status")

		time.Sleep(1 * time.Second)
	}

	// Output:
	//	xr 10.0.1.174:10024: /xinfo ,ssss 10.0.1.174 XR18-35-54-8A XR18 1.18
	//	xr 10.0.1.174:10024: /status ,sss active 10.0.1.174 XR18-35-54-8A
	//	xr 10.0.1.174:10024: /status ,sss active 10.0.1.174 XR18-35-54-8A
	//	xr 10.0.1.174:10024: /status ,sss active 10.0.1.174 XR18-35-54-8A
	// ...

}
