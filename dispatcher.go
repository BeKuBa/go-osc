package osc

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Dispatcher is an interface for an OSC message dispatcher. A dispatcher is
// responsible for dispatching received OSC messages.
type Dispatcher interface {
	Dispatch(packet Packet, addr net.Addr) error
}

// Handler is an interface for message handlers. Every handler implementation
// for an OSC message must implement this interface.
type Handler interface {
	HandleMessage(msg *Message, addr net.Addr)
}

// HandlerFuncExt implements the Handler interface. Type definition for an OSC
// handler function(with msg and addr).
type HandlerFuncExt func(msg *Message, addr net.Addr)

// HandleMessage calls itself with the given OSC Message. Implements the
// Handler interface for HandlerFuncExt(with msg and addr ).
func (f HandlerFuncExt) HandleMessage(msg *Message, addr net.Addr) {
	f(msg, addr)
}

// HandlerFuncExt implements the Handler interface. Type definition for an OSC
// handler function.
type HandlerFunc func(msg *Message)

// NewStandardDispatcher returns an /*
// HandleMessage calls itself with the given OSC Message. Implements the
// Handler interface for HandlerFunc.
func (f HandlerFunc) HandleMessage(msg *Message, addr net.Addr) {
	f(msg)
}

// StandardDispatcher is a dispatcher for OSC packets. It handles the dispatching of
// received OSC packets to Handlers for their given address.
type StandardDispatcher struct {
	handlers       map[string]Handler
	defaultHandler Handler
}

func NewStandardDispatcher() *StandardDispatcher {
	return &StandardDispatcher{
		handlers:       make(map[string]Handler),
		defaultHandler: nil,
	}
}

// AddMsgHandlerExt adds a new message handler (HandlerFuncExt) for the given OSC address.
func (s *StandardDispatcher) AddMsgHandlerExt(addr string, handler HandlerFuncExt) error {

	if addr == "*" {
		s.defaultHandler = handler
		return nil
	}

	for _, chr := range "*?,[]{}# " {
		if strings.Contains(addr, fmt.Sprintf("%c", chr)) {
			return ERROR_OSC_INVALID_CHARACTER
		}
	}

	if addressExists(addr, s.handlers) {
		return ERROR_OSC_ADDRESS_EXISTS
	}

	s.handlers[addr] = handler

	return nil
}

// AddMsgHandler adds a new message handler (HandlerFunc) for the given OSC address.
func (s *StandardDispatcher) AddMsgHandler(addr string, handler HandlerFunc) error {
	return s.AddMsgHandlerExt(addr, func(msg *Message, addr net.Addr) { handler(msg) })
}

// Dispatch dispatches OSC packets. Implements the Dispatcher interface.
func (s *StandardDispatcher) Dispatch(packet Packet, raddr net.Addr) (err error) {
	switch p := packet.(type) {
	case *Message:
		regex, err := getRegEx(p.Address)
		if err != nil {
			return err
		}

		for addr, handler := range s.handlers {
			if regex.MatchString(addr) {
				handler.HandleMessage(p, raddr)
			}
		}

		if s.defaultHandler != nil {
			s.defaultHandler.HandleMessage(p, raddr)
		}

	case *Bundle:
		timer := time.NewTimer(p.Timetag.ExpiresIn())

		errChan := make(chan error)
		go func() {
			<-timer.C

			for _, message := range p.Messages {
				regex, err := getRegEx(message.Address)
				if err != nil {
					errChan <- err
					return
				}
				for address, handler := range s.handlers {
					if regex.MatchString(address) {
						handler.HandleMessage(message, raddr)
					}
				}

				if s.defaultHandler != nil {
					s.defaultHandler.HandleMessage(message, raddr)
				}
			}

			// Process all bundles
			for _, b := range p.Bundles {
				err := s.Dispatch(b, raddr)
				if err != nil {
					errChan <- err
					return
				}
			}
			errChan <- nil
		}()
		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}
