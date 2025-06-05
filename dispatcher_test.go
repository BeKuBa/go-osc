package osc_test

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"bekuba.de/go-osc"

	"github.com/stretchr/testify/assert"
)

// Open question: is this desired behavior, or should server.serve return
// successfully in cases where it would otherwise throw this error?
func serveUntilInterrupted(server *osc.Node) error {
	if err := server.ListenAndServe(nil); err != nil &&
		!strings.Contains(err.Error(), "use of closed network connection") {
		return err
	}

	return nil
}

func TestDispatch(t *testing.T) {
	d := osc.NewStandardDispatcher()

	handlerName := [4]string{"/message", "/message/01", "/message/03", "*"}
	b := [4]bool{false, false, false, false}

	err := d.AddMsgHandler(handlerName[0], func(msg *osc.Message) {
		b[0] = true
	})
	assert.Nil(t, err)
	err = d.AddMsgHandler(handlerName[1], func(msg *osc.Message) {
		b[1] = true
	})
	assert.Nil(t, err)
	err = d.AddMsgHandler(handlerName[2], func(msg *osc.Message) {
		b[2] = true
	})
	assert.Nil(t, err)
	err = d.AddMsgHandler(handlerName[3], func(msg *osc.Message) {
		b[3] = true
	})
	assert.Nil(t, err)

	// error ERROR_OSC_ADDRESS_EXISTS
	err = d.AddMsgHandler(handlerName[1], func(msg *osc.Message) {
		b[1] = true
	})
	assert.NotNil(t, err)

	t.Run("dispatch message", func(t *testing.T) {

		tc := []struct {
			desc string
			msg  string
			b    [4]bool
			err  bool
		}{
			{
				"match everything",
				"*",
				[4]bool{true, true, true, true},
				false,
			},
			{
				"match /message",
				"/message",
				[4]bool{true, false, false, true},
				false,
			},
			{
				"match [1-3]",
				"/message/0[1-3]",
				[4]bool{false, true, true, true},
				false,
			},
			{
				"don't match",
				"/message/01/01",
				[4]bool{false, false, false, true},
				false,
			},
			{
				"Regex error",
				"}/",
				[4]bool{false, false, false, false},
				true,
			},
		}

		err = nil
		for _, tt := range tc {
			msg := osc.NewMessage(tt.msg)
			err = d.Dispatch(msg, nil)
			if tt.err {
				assert.NotNil(t, err, "%s: msgPath = '%s', expect error", tt.desc, tt.msg)
			} else {
				assert.Nil(t, err, "%s: msgPath = '%s', expect no error", tt.desc, tt.msg)
			}

			for i, got := range b {
				if got != tt.b[i] {
					t.Errorf("%s: msgPath='%v', handlerFunc='%s', got  = '%t', want = '%t'", tt.desc, tt.msg, handlerName[i], got, tt.b[i])
				}
			}

			b = [4]bool{false, false, false, false}
			err = nil
		}
	})
	t.Run("dispatch bundle", func(t *testing.T) {

		b = [4]bool{false, false, false, false}

		bundle := osc.NewBundle(time.Now())

		// 1 bundle, 2 messages
		err := bundle.Append(osc.NewMessage(handlerName[1], ""))
		assert.Nil(t, err)
		err = bundle.Append(osc.NewMessage(handlerName[2], "test2"))
		assert.Nil(t, err)

		err = d.Dispatch(bundle, nil)
		assert.Nil(t, err)

		assert.False(t, b[0], "check handlerFunc %v", handlerName[0])
		assert.True(t, b[1], "check handlerFunc %v", handlerName[1])
		assert.True(t, b[2], "check handlerFunc %v", handlerName[2])
		assert.True(t, b[3], "check handlerFunc %v", handlerName[3])

		// bundle2 with 2 bundles: bundle(2 messages), bundle3(1 message)
		bundle2 := osc.NewBundle(bundle.Timetag.Time())
		err = bundle2.Append(bundle)
		assert.Nil(t, err)

		bundle3 := osc.NewBundle(time.Now())
		err = bundle2.Append(bundle3)
		assert.Nil(t, err)
		err = bundle3.Append(osc.NewMessage(handlerName[0]))
		assert.Nil(t, err)

		err = d.Dispatch(bundle2, nil)
		assert.Nil(t, err)

		assert.True(t, b[0], "check handlerFunc %v", handlerName[0])
		assert.True(t, b[1], "check handlerFunc %v", handlerName[1])
		assert.True(t, b[2], "check handlerFunc %v", handlerName[2])
		assert.True(t, b[3], "check handlerFunc %v", handlerName[3])

		// bundle: test error handling
		err = bundle.Append(osc.NewMessage("}/"))
		assert.Nil(t, err)
		err = d.Dispatch(bundle, nil)
		assert.NotNil(t, err)

		// bundle3: test error handling
		err = bundle3.Append(osc.NewMessage("}/"))
		assert.Nil(t, err)
		err = d.Dispatch(bundle2, nil)
		assert.NotNil(t, err)

	})
}

func TestAddMsgHandler(t *testing.T) {
	d := osc.NewStandardDispatcher()
	err := d.AddMsgHandler("/address/test", func(msg *osc.Message) {})
	if err != nil {
		t.Error("Expected that OSC address '/address/test' is valid")
	}
}

func TestAddMsgHandlerWithInvalidAddress(t *testing.T) {
	d := osc.NewStandardDispatcher()
	err := d.AddMsgHandler("/address*/test", func(msg *osc.Message) {})
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

	addr := fmt.Sprintf("localhost:%v", port)

	server, err := osc.NewNode(addr)
	if err != nil {
		fmt.Println(err)
	}
	defer server.Close()

	d := osc.NewStandardDispatcher()
	if err := d.AddMsgHandlerExt(
		"/address/test",
		func(msg *osc.Message, addr net.Addr) {
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

			addr1 := "localhost:0"
			client, err := osc.NewNode(addr1)
			if err != nil {
				fmt.Println(err)
			}
			defer client.Close()

			msg := osc.NewMessage("/address/test")
			msg.Append(int32(1122))
			if err := client.SendTo(addr, msg); err != nil {
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
