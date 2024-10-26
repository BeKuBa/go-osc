package main

import (
	"fmt"
	"net"
	"time"

	"bekuba/go-osc"
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

	addr := fmt.Sprintf("%v:%v", xrIP, xr18Port)

	app, err := osc.NewServerAndClient(":0")
	if err != nil {
		fmt.Println(err)
	}
	defer app.Close()

	d := osc.NewStandardDispatcher()
	err = d.AddMsgHandlerExt("*", func(msg *osc.Message, addr net.Addr) {
		fmt.Printf("xr %v: %v  \n", addr, msg)
	})
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		err = app.ListenAndServe(d)
		if err != nil {
			fmt.Println(err)
		}
	}()

	app.SendMsgTo(addr, "/xinfo")

	for {
		// keepp connection alive (for multi client usage)
		app.SendMsgTo(addr, "/xremote")
		// show status of xair
		app.SendMsgTo(addr, "/status")

		time.Sleep(1 * time.Second)
	}

	// Output:
	//	xr 10.0.1.174:10024: /xinfo ,ssss 10.0.1.174 XR18-35-54-8A XR18 1.22
	//	xr 10.0.1.174:10024: /status ,sss active 10.0.1.174 XR18-35-54-8A
	//	xr 10.0.1.174:10024: /status ,sss active 10.0.1.174 XR18-35-54-8A
	//	xr 10.0.1.174:10024: /status ,sss active 10.0.1.174 XR18-35-54-8A
	// ...

}
