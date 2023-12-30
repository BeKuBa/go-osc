package osc_test

import (
	"github.com/crgimenes/go-osc"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
)

const (
	ping = "/ping"
	pong = "/pong"
)

func TestServerAndClient(t *testing.T) {

	timeout := time.After(1 * time.Second)
	done := make(chan bool)

	go func() {

		wait := sync.WaitGroup{}
		wait.Add(3)

		var d float64 = 0.0
		var i int32 = 0

		addr1, err := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
		if err != nil {
			t.Error(err)
		}

		addr2, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
		if err != nil {
			t.Error(err)
		}

		d1 := osc.NewStandardDispatcher()
		app1 := osc.NewServerAndClient(d1)
		err = app1.NewConn(addr1, addr2)
		if err != nil {
			t.Error(err)
		}
		defer func() {
			err := app1.Close()
			if err != nil {
				t.Error(err)
			}
		}()

		err = d1.AddMsgHandler(ping, func(msg *osc.Message, addr net.Addr) {
			d = msg.Arguments[0].(float64)
			assert.Equal(t, 1.0, d)
			err = app1.SendMsg(pong, 2)
			if err != nil {
				t.Error(err)
			}
		})
		if err != nil {
			t.Error(err)
		}

		go func() {
			err := app1.ListenAndServe()
			if err != nil {
				t.Error(err)
			}
		}()

		d2 := osc.NewStandardDispatcher()
		err = d2.AddMsgHandler(pong, func(msg *osc.Message, addr net.Addr) {
			i = msg.Arguments[0].(int32)
			assert.Equal(t, int32(2), i)
			wait.Done()
		})
		if err != nil {
			t.Error(err)
		}

		app2 := osc.NewServerAndClient(d2)
		err = app2.NewConn(addr2, addr1)
		if err != nil {
			t.Error(err)
		}
		defer func() {
			err := app2.Close()
			if err != nil {
				t.Error(err)
			}
		}()
		go func() {
			err := app2.ListenAndServe()
			if err != nil {
				t.Error(err)
			}
		}()

		err = app2.SendMsg(ping, 1.0)
		if err != nil {
			t.Error(err)
		}

		err = app1.SendMsg(pong, 2)
		if err != nil {
			t.Error(err)
		}

		app3 := osc.NewServerAndClient(nil)
		err = app3.NewConn(nil, nil)
		err = app3.SendMsgTo(addr2, pong, 2)
		if err != nil {
			t.Error(err)
		}

		wait.Wait()

		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("test didn't finish in time")
	case <-done:
	}
}
