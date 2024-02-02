package osc_test

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/crgimenes/go-osc"
	"github.com/stretchr/testify/assert"
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
		wait.Add(1)

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
		err = app1.NewConn(addr2, addr1)

		err = d1.AddMsgHandler(ping, func(msg *osc.Message) {
			pingF64 = msg.Arguments[0].(float64)

			err = app1.SendMsg(pong, cBoolTrue, cBoolFalse, cI32, cI64, cF32, cF64, cStrTest, cStrEmpty, cI, nil, cArray, cTimetag)

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
		err = d2.AddMsgHandler(pong, func(msg *osc.Message) {

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

			wait.Done()
		})
		if err != nil {
			t.Error(err)
		}

		app2 := osc.NewServerAndClient(d2)
		err = app2.NewConn(addr1, addr2)
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
