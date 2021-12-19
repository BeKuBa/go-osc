package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/crgimenes/go-osc"
)

func main() {
	addr := "127.0.0.1:8765"
	server := &osc.Server{}
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Println("Couldn't listen: ", err)
	}
	defer conn.Close()

	fmt.Println("### Welcome to go-osc receiver demo")
	fmt.Println("Press \"q\" to exit")

	go func() {
		fmt.Println("Start listening on", addr)

		for {
			packet, err := server.ReceivePacket(conn)
			if err != nil {
				fmt.Println("Server error: " + err.Error())
				os.Exit(1)
			}

			if packet != nil {
				switch p := packet.(type) {
				default:
					fmt.Println("Unknow packet type!")

				case *osc.Message:
					fmt.Println("-- OSC Message:", p)

				case *osc.Bundle:
					fmt.Println("-- OSC Bundle:")

					for i, message := range p.Messages {
						fmt.Printf("  -- OSC Message #%d: ", i+1)
						fmt.Println(message)
					}
				}
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		c, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			os.Exit(1)
		}

		if c == 'q' {
			os.Exit(0)
		}
	}
}
