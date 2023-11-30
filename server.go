package osc

import (
	"bufio"
	"bytes"
	"net"
	"time"
)

// Server represents an OSC server. The server listens on Address and Port for
// incoming OSC packets and bundles.
type Server struct {
	Addr        string
	Dispatcher  Dispatcher
	ReadTimeout time.Duration
	close       func() error
}

// ListenAndServe retrieves incoming OSC packets and dispatches the retrieved
// OSC packets.
func (s *Server) ListenAndServe() error {
	defer s.Close()

	if s.Dispatcher == nil {
		s.Dispatcher = NewStandardDispatcher()
	}

	ln, err := net.ListenPacket("udp", s.Addr)
	if err != nil {
		return err
	}

	s.close = ln.Close

	return s.serve(ln)
}

// Serve retrieves incoming OSC packets from the given connection and dispatches
// retrieved OSC packets. If something goes wrong an error is returned.
func (s *Server) serve(c net.PacketConn) error {
	tempDelay := 25 + time.Millisecond

	for {
		msg, err := s.Read(c)
		if err != nil {
			ne, ok := err.(net.Error)

			if ok && ne.Temporary() {
				time.Sleep(tempDelay)
				continue
			}

			return err
		}

		go s.Dispatcher.Dispatch(msg)
	}
}

// Close forcibly closes a server's connection.
//
// This causes a "use of closed network connection" error the next time the
// server attempts to read from the connection.
func (s *Server) Close() error {
	if s.close == nil {
		return nil
	}

	return s.close()
}

// Read retrieves OSC packets.
func (s *Server) Read(c net.PacketConn) (Packet, error) {
	if s.ReadTimeout != 0 {
		err := c.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		if err != nil {
			return nil, err
		}
	}

	data := make([]byte, 65535)

	n, _, err := c.ReadFrom(data)
	if err != nil {
		return nil, err
	}

	var start int
	p, err := readPacket(bufio.NewReader(bytes.NewBuffer(data)), &start, n)

	return p, err
}
