package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"bekuba/go-osc"
)

func main() {

	addr := "localhost:8765"

	server, err := osc.NewServerAndClient(addr)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	fmt.Println("### Welcome to go-osc receiver demo")
	fmt.Println("Press \"q\" to exit")

	go func() {
		fmt.Println("Start listening on", addr)

		for {
			packet, _, err := server.Read()
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
