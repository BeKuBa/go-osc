package osc

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestListenAndServe(t *testing.T) {
	done := sync.WaitGroup{}
	done.Add(2)

	dispatcher := NewStandardDispatcher()
	dispatcher.AddMsgHandler("/osc/address", func(msg *Message) {
		assert.Equal(t, "/osc/address", msg.Address)
		assert.Equal(t, 3, len(msg.Arguments))
		assert.Equal(t, int32(111), msg.Arguments[0].(int32))
		assert.Equal(t, true, msg.Arguments[1].(bool))
		assert.Equal(t, "hello", msg.Arguments[2].(string))

		done.Done()
	})
	addr := "127.0.0.1:8765"
	server := &Server{
		Addr:       addr,
		Dispatcher: dispatcher,
	}
	defer func() {
		err := server.Close()
		assert.NoError(t, err)
	}()

	go server.ListenAndServe()

	go func() {
		client := NewClient("localhost", 8765)
		msg := NewMessage("/osc/address", int32(111), true, "hello")
		client.Send(msg)

		done.Done()
	}()

	done.Wait()
}

func TestReadTimeout(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	dispatcher := NewStandardDispatcher()
	dispatcher.AddMsgHandler("/address/test", func(msg *Message) {
		assert.Equal(t, "/address/test", msg.Address)
		assert.Equal(t, 0, len(msg.Arguments))

		wg.Done()
	})
	addr := "127.0.0.1:6677"
	server := &Server{
		Addr:        addr,
		Dispatcher:  dispatcher,
		ReadTimeout: 100 * time.Millisecond,
	}
	defer func() {
		err := server.Close()
		assert.NoError(t, err)
	}()

	go server.ListenAndServe()

	go func() {
		time.Sleep(150 * time.Millisecond)
		client := NewClient("localhost", 6677)
		msg := NewMessage("/address/test")
		client.Send(msg)

		wg.Done()
	}()

	wg.Wait()
}
