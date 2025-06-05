package osc

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"net"
	"sync"
	"time"
)

// Node structure
type Node struct {
	conn *net.UDPConn
	//	Dispatcher  Dispatcher
	ReadTimeout time.Duration
}

// Node create a new OSC Server and/or Client connection
func NewNode(laddr string) (*Node, error) {

	addr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		return nil, ErrorOscAddressFormat
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, ErrorOscAddress
	}

	return &Node{conn: conn}, nil
}

// SendTo sends an OSC Bundle or an OSC Message (as OSC Client) to a given UDP address.
func (sc *Node) SendToUDPAddr(raddr *net.UDPAddr, packet Packet) (err error) {
	if sc.conn != nil {

		data, err := packet.MarshalBinary()
		if err != nil {
			return err
		}
		if _, err = sc.conn.WriteTo(data, raddr); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("can't send OSC packet! %s", err)
	}
	return err
}

// SendTo sends an OSC Bundle or an OSC Message (as OSC Client) to a given address.
func (sc *Node) SendTo(raddr string, packet Packet) (err error) {
	addr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		return ErrorOscAddressFormat
	}
	return sc.SendToUDPAddr(addr, packet)
}

// SendMsgTo sends a OSC Message to a given UDP address(all int types converted to int32)
// Default int is int32, include int values in range of int32
// If you need a int value in range of int64 convert the arg to int64
func (sc *Node) SendMsgToUDPAddr(addr *net.UDPAddr, path string, args ...any) error {
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

	return sc.SendToUDPAddr(addr, NewMessage(path, a...))
}

// SendMsgTo sends a OSC Message to a given address(all int types converted to int32)
// Default int is int32, include int values in range of int32
// If you need a int value in range of int64 convert the arg to int64
func (sc *Node) SendMsgTo(raddr string, path string, args ...any) error {
	addr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		return ErrorOscAddressFormat
	}
	return sc.SendMsgToUDPAddr(addr, path, args...)
}

// ListenAndServe listen and serve as an OSC Server
func (sc *Node) ListenAndServe(d Dispatcher) error {
	if sc.conn != nil {

		err := sc.serve(sc.conn, d)

		// serve is a go routine with a loop that only ends on error
		// so can now sc.conn (e.g. after close connection)  maybe nil
		if sc.conn == nil {
			err = nil
		}

		return err
	}
	return fmt.Errorf("ServerAndClient connection is not created")
}

/* ************************************** */

// Serve retrieves incoming OSC packets from the given connection and dispatches
// retrieved OSC packets. If something goes wrong an error is returned.
func (sc *Node) serve(c net.PacketConn, d Dispatcher) error {
	tempDelay := 25 + time.Millisecond

	for {
		if c == nil {
			return nil
		}
		msg, raddr, err := sc.Read()
		if err != nil {
			ne, ok := err.(net.Error)

			if ok && ne.Temporary() {
				time.Sleep(tempDelay)
				continue
			}

			return err
		}
		if d != nil {
			errChan := make(chan error)
			go func() {
				errChan <- d.Dispatch(msg, raddr)
			}()
			if err := <-errChan; err != nil {
				return err
			}
		}
	}
}

// Read retrieves OSC packets.
func (s *Node) Read() (Packet, net.Addr, error) {
	if s.ReadTimeout != 0 {
		err := s.conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		if err != nil {
			return nil, nil, err
		}
	}

	data := make([]byte, 65535)

	n, addr, err := s.conn.ReadFrom(data)
	if err != nil {
		return nil, nil, err
	}

	var start int
	p, err := readPacket(bufio.NewReader(bytes.NewBuffer(data)), &start, n)

	return p, addr, err
}

func (sc *Node) Close() {
	done := sync.WaitGroup{}
	done.Add(1)
	c := sc.conn
	sc.conn = nil
	done.Done()
	c.Close()
}

func (sc *Node) Conn() *net.UDPConn {
	return sc.conn
}
