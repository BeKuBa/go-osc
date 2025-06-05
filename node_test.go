package osc_test

import (
	"net"
	"sync"
	"testing"
	"time"

	"bekuba.de/go-osc"

	"github.com/stretchr/testify/assert"
)

const (
	ping = "/ping"
	pong = "/pong"
)

func TestListenAndServe(t *testing.T) {
	done := sync.WaitGroup{}
	done.Add(2)

	addr := "127.0.0.1:0"

	server, err := osc.NewNode(addr)
	assert.NoError(t, err)
	addr = server.Conn().LocalAddr().String()
	defer server.Close()

	dispatcher := osc.NewStandardDispatcher()
	err = dispatcher.AddMsgHandler("/osc/address", func(msg *osc.Message) {
		assert.Equal(t, "/osc/address", msg.Address)
		assert.Equal(t, 3, len(msg.Arguments))
		assert.Equal(t, int32(111), msg.Arguments[0].(int32))
		assert.Equal(t, true, msg.Arguments[1].(bool))
		assert.Equal(t, "hello", msg.Arguments[2].(string))

		done.Done()
	})
	assert.NoError(t, err)
	go func() {
		err := server.ListenAndServe(dispatcher)
		assert.NoError(t, err)
	}()

	go func() {
		client, err := osc.NewNode("localhost:0")
		assert.NoError(t, err)
		defer client.Close()

		msg := osc.NewMessage("/osc/address", int32(111), true, "hello")
		client.SendTo(addr, msg)

		done.Done()
	}()

	done.Wait()
}

func TestServerAndClient(t *testing.T) {

	timeout := time.After(5 * time.Second)
	done := make(chan bool)

	go func() {

		wait := sync.WaitGroup{}
		wait.Add(2)

		var pingF64 = 0.0

		var boolTrue = false
		var boolFalse = true
		var i32 int32 = 0
		var i64 int64 = 0
		var f32 float32 = 0.0
		var f64 float64 = 0.0
		var strTest string = ""
		var strEmpty string = "e"
		var i int = 0
		var null any = 10
		var array []byte = []byte{10, 11, 12}
		var timetag osc.Timetag

		const (
			cBoolTrue          = true
			cBoolFalse         = false
			cI32       int32   = 2
			cI64       int64   = 3
			cF32       float32 = 4.0
			cF64       float64 = 5.0
			cStrTest   string  = "6test"
			cStrEmpty  string  = ""
			cI         int     = 8
			// nil
		)
		var cArray []byte = []byte{10, 48}
		const cTimetag osc.Timetag = 16818286200017484014

		addr1 := "127.0.0.1:0"
		addr2 := "127.0.0.1:0"

		// app1
		app1, err := osc.NewNode(addr1)
		assert.NoError(t, err)
		addr1 = app1.Conn().LocalAddr().String()
		defer app1.Close()

		d1 := osc.NewStandardDispatcher()
		err = d1.AddMsgHandler(ping, func(msg *osc.Message) {
			pingF64 = msg.Arguments[0].(float64)

			err = app1.SendMsgTo(addr2, pong, cBoolTrue, cBoolFalse, cI32, cI64, cF32, cF64, cStrTest, cStrEmpty, cI, nil, cArray, cTimetag)

			assert.NoError(t, err)
		})

		assert.NoError(t, err)

		go func() {
			err := app1.ListenAndServe(d1)
			assert.NoError(t, err)
		}()

		// app2
		app2, err := osc.NewNode(addr2)
		assert.NoError(t, err)
		addr2 = app2.Conn().LocalAddr().String()
		defer app2.Close()

		d2 := osc.NewStandardDispatcher()
		err = d2.AddMsgHandlerExt(pong, func(msg *osc.Message, raddr net.Addr) {

			if addr1 == raddr.String() {
				boolTrue = msg.Arguments[0].(bool)
				boolFalse = msg.Arguments[1].(bool)
				i32 = msg.Arguments[2].(int32)
				i64 = msg.Arguments[3].(int64)
				f32 = msg.Arguments[4].(float32)
				f64 = msg.Arguments[5].(float64)
				strTest = msg.Arguments[6].(string)
				strEmpty = msg.Arguments[7].(string)
				i = int(msg.Arguments[8].(int32))
				null = msg.Arguments[9]
				array = msg.Arguments[10].([]byte)
				timetag = msg.Arguments[11].(osc.Timetag)
			}
			wait.Done()
		})
		assert.NoError(t, err)

		go func() {
			err := app2.ListenAndServe(d2)
			assert.NoError(t, err)
		}()

		// app2 send ping, app1 send pong, app3 send pong
		err = app2.SendMsgTo(addr1, ping, 1.0)
		assert.NoError(t, err)

		app3, err := osc.NewNode(":0")
		assert.NoError(t, err)

		err = app3.SendMsgTo(addr2, pong, 2)
		assert.NoError(t, err)

		wait.Wait()

		// check if send and receive are the same
		assert.Equal(t, 1.0, pingF64)

		assert.Equal(t, cBoolTrue, boolTrue)
		assert.Equal(t, cBoolFalse, boolFalse)
		assert.Equal(t, cI32, i32)
		assert.Equal(t, cI64, i64)
		assert.Equal(t, cF32, f32)
		assert.Equal(t, cF64, f64)
		assert.Equal(t, cStrTest, strTest)
		assert.Equal(t, cStrEmpty, strEmpty)
		assert.Equal(t, cI, i)
		assert.Nil(t, null)
		assert.Equal(t, cArray, array)
		assert.Equal(t, cTimetag, timetag)

		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("test didn't finish in time")
	case <-done:
	}
}

func TestReadTimeout(t *testing.T) {

	addr1 := "127.0.0.1:0"
	addr2 := "127.0.0.1:0"

	// app1
	d1 := osc.NewStandardDispatcher()
	app1, err := osc.NewNode(addr1)
	assert.NoError(t, err)
	addr1 = app1.Conn().LocalAddr().String()
	app1.Conn().SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	defer app1.Close()

	get := false
	err = d1.AddMsgHandler(ping, func(msg *osc.Message) {
		get = true
	})

	assert.NoError(t, err)

	// app2
	app2, err := osc.NewNode(addr2)
	assert.NoError(t, err)
	addr2 = app2.Conn().LocalAddr().String()
	defer app2.Close()

	// In time
	wait := sync.WaitGroup{}
	wait.Add(1)

	go func() {

		p, addr, err := app1.Read()
		assert.NoError(t, err)
		if err == nil {
			d1.Dispatch(p, addr)
		}

		wait.Done()
	}()
	time.Sleep(50 * time.Millisecond)
	// app2 send ping
	err = app2.SendMsgTo(addr1, ping)
	assert.NoError(t, err)

	wait.Wait()
	assert.Equal(t, true, get)

	// Timeout
	get = false
	wait = sync.WaitGroup{}
	wait.Add(1)

	go func() {

		p, addr, err := app1.Read()
		assert.Error(t, err)
		if err == nil {
			d1.Dispatch(p, addr)
		}

		wait.Done()
	}()
	time.Sleep(150 * time.Millisecond)
	// app2 send ping
	err = app2.SendMsgTo(addr1, ping)
	assert.NoError(t, err)

	wait.Wait()
	assert.Equal(t, false, get)
}
