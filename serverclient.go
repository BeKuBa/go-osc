package osc

import (
	"fmt"
	"math"
	"net"
)

// ServerAndClient structure
type ServerAndClient struct {
	conn   *net.UDPConn
	RAddr  *net.UDPAddr // default remote adr (for Send and SendMsg)
	server *Server
}

// NewServerAndClient create a new ServerandClient
func NewServerAndClient(dispatcher Dispatcher) *ServerAndClient {
	return &ServerAndClient{server: &Server{Dispatcher: dispatcher}}
}

// NewConn create a new UDP Connection for Server and Client
func (sc *ServerAndClient) NewConn(laddr *net.UDPAddr, raddr *net.UDPAddr) error {
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}

	sc.conn = conn
	sc.RAddr = raddr

	return err
}

// SendTo sends an OSC Bundle or an OSC Message (as OSC Client) to a given address.
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

// Send sends an OSC Bundle or an OSC Message (as OSC Client).
func (sc *ServerAndClient) Send(packet Packet) error {
	return sc.SendTo(sc.RAddr, packet)
}

// SendMsgTo sends a OSC Message to a given address(all int types converted to int32)
// Default int is int32, include int values in range of int32
// If you need a int value in range of int64 convert the arg to int64
func (sc *ServerAndClient) SendMsgTo(addr net.Addr, path string, args ...any) error {
	var a []any

	for _, arg := range args {
		switch t := arg.(type) {
		case int8:
			a = append(a, int32(t))
		case uint8:
			a = append(a, int32(t))
		case int16:
			a = append(a, int32(t))
		case uint16:
			a = append(a, int32(t))
		case int:
			if (t <= math.MaxInt32) && (t >= math.MinInt32) {
				a = append(a, int32(t))
			} else {
				return fmt.Errorf("int32 %d out of range", t)
			}
		case bool, int64, int32, float32, float64, string, nil, []byte, Timetag:
			a = append(a, t)
		default:
			return fmt.Errorf("wrong datatype, can't send OSC packet")
		}

	}

	return sc.SendTo(addr, NewMessage(path, a...))
}

// SendMsg sends a OSC Message to a given address(all int types converted to int32)
// Default int is int32, include int values in range of int32
// If you need a int value in range of int64 convert the arg to int64
func (sc *ServerAndClient) SendMsg(path string, args ...any) error {
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
	}
	return fmt.Errorf("ServerAndClient connection is not created")
}

// Close close ServerAndClient connection
func (sc *ServerAndClient) Close() error {
	conn := sc.conn
	// for handle return server.serve error
	sc.conn = nil

	err := conn.Close()

	return err
}

// Conn ServerAndClient conn
func (sc *ServerAndClient) Conn() *net.UDPConn {
	return sc.conn
}

// Server ServerAndClient conn
func (sc *ServerAndClient) Server() *Server {
	return sc.server
}
