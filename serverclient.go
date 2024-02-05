package osc

import (
	"fmt"
	"math"
	"net"
)

type ServerAndClient struct {
	conn   *net.UDPConn
	RAddr  *net.UDPAddr // default remote adr (for Send and SendMsg)
	server *Server
}

func NewServerAndClient(dispatcher Dispatcher) *ServerAndClient {
	return &ServerAndClient{server: &Server{Dispatcher: dispatcher}}
}

// New UDP Connection for Server and Client
func (sc *ServerAndClient) NewConn(laddr *net.UDPAddr, raddr *net.UDPAddr) error {
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}

	sc.conn = conn
	sc.RAddr = raddr

	return err
}

// Send sends an OSC Bundle or an OSC Message (as OSC Client).
func (sc *ServerAndClient) SendTo(raddr net.Addr, packet Packet) (err error) {
	if sc.conn != nil {
		data, err := packet.MarshalBinary()
		if err != nil {
			return err
		}
		if _, err = sc.conn.WriteTo(data, raddr); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("can't send OSC packet! %s", err.Error())
	}
	return err
}

func (sc *ServerAndClient) Send(packet Packet) error {
	return sc.SendTo(sc.RAddr, packet)
}

// SendMsg sends a OSC Message (all int types konverted to int32)
func (sc *ServerAndClient) SendMsgTo(addr net.Addr, path string, args ...interface{}) error {
	var a []interface{}

	for _, arg := range args {
		switch t := arg.(type) {
		case int8:
			a = append(a, int32(t))
		case uint8:
			a = append(a, int32(t))
		case int:
			if (arg.(int) <= math.MaxInt32) && (arg.(int) >= math.MinInt32) {
				a = append(a, int32(t))
			} else {
				return fmt.Errorf("int32 %d out of range", arg.(int))
			}

		default:
			a = append(a, arg)
		}

	}

	return sc.SendTo(addr, NewMessage(path, a...))
}

func (sc *ServerAndClient) SendMsg(path string, args ...interface{}) error {
	return sc.SendMsgTo(sc.RAddr, path, args...)
}

// ListenAndServe listen and serve as an OSC Server
func (sc *ServerAndClient) ListenAndServe() error {
	if sc.conn != nil {
		if sc.server.Dispatcher == nil {
			sc.server.Dispatcher = NewStandardDispatcher()
		}

		err := sc.server.serve(sc.conn)

		// serve is a go routine with a loop that only ends on error
		// so can now sc.conn (e.g. after close connection)  maybe nil
		if sc.conn == nil {
			err = nil
		}

		return err
	} else {
		return fmt.Errorf("ServerAndClient connection is not created")
	}
}

func (sc *ServerAndClient) Close() error {
	conn := sc.conn
	// for handle return server.serve error
	sc.conn = nil

	err := conn.Close()

	return err
}

func (sc *ServerAndClient) Conn() *net.UDPConn {
	return sc.conn
}
