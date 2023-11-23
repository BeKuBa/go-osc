package osc

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// These tests stop the server by forcibly closing the connection, which causes
// a "use of closed network connection" error the next time we try to read from
// the connection. As a workaround, this wraps server.ListenAndServe() in an
// error-handling layer that doesn't consider "use of closed network connection"
// an error.
//
// Open question: is this desired behavior, or should server.serve return
// successfully in cases where it would otherwise throw this error?
func serveUntilInterrupted(server *Server) error {
	if err := server.ListenAndServe(); err != nil &&
		!strings.Contains(err.Error(), "use of closed network connection") {
		return err
	}

	return nil
}

func TestServerMessageReceiving(t *testing.T) {
	port := 6677

	finish := make(chan bool)
	start := make(chan bool)
	done := sync.WaitGroup{}
	done.Add(2)

	// Start the server in a go-routine
	go func() {
		server := &Server{}

		c, err := net.ListenPacket("udp", "localhost:"+strconv.Itoa(port))
		if err != nil {
			t.Error(err)
			return
		}
		defer c.Close()

		// Start the client
		start <- true

		packet, err := server.Read(c)
		if err != nil {
			t.Errorf("server error: %v", err)
			return
		}
		if packet == nil {
			t.Error("nil packet")
			return
		}

		msg := packet.(*Message)
		lenArg := len(msg.Arguments)
		if lenArg != 2 {
			t.Errorf("Argument length should be 2 and is: %d\n", lenArg)
		}
		if msg.Arguments[0].(int32) != 1122 {
			t.Error("Argument should be 1122 and is: " + string(msg.Arguments[0].(int32)))
		}
		if msg.Arguments[1].(int32) != 3344 {
			t.Error("Argument should be 3344 and is: " + string(msg.Arguments[1].(int32)))
		}

		c.Close()
		finish <- true
	}()

	go func() {
		timeout := time.After(5 * time.Second)
		select {
		case <-timeout:
		case <-start:
			client := NewClient("localhost", port)
			msg := NewMessage("/address/test")
			msg.Append(int32(1122))
			msg.Append(int32(3344))
			time.Sleep(500 * time.Millisecond)
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

func TestReadTimeout(t *testing.T) {
	start := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		select {
		case <-time.After(5 * time.Second):
			t.Error("timed out")
			wg.Done()
		case <-start:
			client := NewClient("localhost", 6677)
			msg := NewMessage("/address/test1")
			err := client.Send(msg)
			if err != nil {
				t.Error(err)
			}

			time.Sleep(150 * time.Millisecond)
			msg = NewMessage("/address/test2")
			err = client.Send(msg)
			if err != nil {
				t.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()

		server := &Server{ReadTimeout: 100 * time.Millisecond}
		c, err := net.ListenPacket("udp", "localhost:6677")
		if err != nil {
			t.Error(err)
		}
		defer c.Close()

		// Start the client
		start <- true
		p, err := server.Read(c)
		if err != nil {
			t.Errorf("server error: %v", err)
			return
		}
		if got, want := p.(*Message).Address, "/address/test1"; got != want {
			t.Errorf("wrong address; got = %s, want = %s", got, want)
			return
		}

		// Second receive should time out since client is delayed 150 milliseconds
		if _, err = server.Read(c); err == nil {
			t.Errorf("expected error")
			return
		}

		// Next receive should get it
		p, err = server.Read(c)
		if err != nil {
			t.Errorf("server error: %v", err)
			return
		}
		if got, want := p.(*Message).Address, "/address/test2"; got != want {
			t.Errorf("wrong address; got = %s, want = %s", got, want)
			return
		}
	}()

	wg.Wait()
}
