package osc

import (
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestAddMsgHandler(t *testing.T) {
	d := NewStandardDispatcher()
	err := d.AddMsgHandler("/address/test", func(msg *Message) {})
	if err != nil {
		t.Error("Expected that OSC address '/address/test' is valid")
	}
}

func TestAddMsgHandlerWithInvalidAddress(t *testing.T) {
	d := NewStandardDispatcher()
	err := d.AddMsgHandler("/address*/test", func(msg *Message) {})
	if err == nil {
		t.Error("Expected error with '/address*/test'")
	}
}

func TestServerMessageDispatching(t *testing.T) {
	finish := make(chan bool)
	start := make(chan bool)
	done := sync.WaitGroup{}
	done.Add(2)

	port := 6677
	addr := "localhost:" + strconv.Itoa(port)

	d := NewStandardDispatcher()
	server := &Server{Addr: addr, Dispatcher: d}

	defer func() {
		err := server.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	if err := d.AddMsgHandlerExt(
		"/address/test",
		func(msg *Message, addr net.Addr) {
			lenArgs := len(msg.Arguments)
			if lenArgs != 1 {
				t.Errorf("Argument length should be 1 and is: %d", lenArgs)
			}

			if msg.Arguments[0].(int32) != 1122 {
				t.Error("Argument should be 1122 and is: " + string(msg.Arguments[0].(int32)))
			}

			finish <- true
		},
	); err != nil {
		t.Error("Error adding message handler")
	}

	// Server goroutine
	go func() {
		start <- true

		if err := serveUntilInterrupted(server); err != nil {
			t.Errorf("error during Serve: %s", err.Error())
		}
	}()

	// Client goroutine
	go func() {
		timeout := time.After(5 * time.Second)
		select {
		case <-timeout:
		case <-start:
			time.Sleep(500 * time.Millisecond)
			client := NewClient("localhost", port)
			msg := NewMessage("/address/test")
			msg.Append(int32(1122))
			if err := client.Send(msg); err != nil {
				t.Error(err)
				done.Done()
				done.Done()
				return
			}
		}

		done.Done()

		select {
		case <-timeout:
		case <-finish:
		}
		done.Done()
	}()

	done.Wait()
}
