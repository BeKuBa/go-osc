package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"bekuba.de/go-osc"
)

func main() {

	addr1 := "localhost:0"
	addr2 := "localhost:8765"

	client, err := osc.NewNode(addr1)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("### Welcome to go-osc transmitter demo")
	fmt.Println("Please, select the OSC packet type you would like to send:")
	fmt.Println("\tm: OSCMessage")
	fmt.Println("\tb: OSCBundle")
	fmt.Println("\tPress \"q\" to exit")
	fmt.Printf("# ")

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}

		sline := strings.TrimRight(string(line), "\n")
		if sline == "m" {
			message := osc.NewMessage("/message/address")
			message.Append(int32(12345))
			message.Append("teststring")
			message.Append(true)
			message.Append(false)
			client.SendTo(addr1, message)
		} else if sline == "b" {
			bundle := osc.NewBundle(time.Now())
			message1 := osc.NewMessage("/bundle/message/1")
			message1.Append(int32(12345))
			message1.Append("teststring")
			message1.Append(true)
			message1.Append(false)
			message2 := osc.NewMessage("/bundle/message/2")
			message2.Append(int32(3344))
			message2.Append(float32(101.9))
			message2.Append("string1")
			message2.Append("string2")
			message2.Append(true)
			bundle.Append(message1)
			bundle.Append(message2)
			client.SendTo(addr2, bundle)
		} else if sline == "q" {
			fmt.Println("Exit!")
			os.Exit(0)
		}

		fmt.Printf("# ")
	}
}
