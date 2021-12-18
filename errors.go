package osc

import "errors"

var (
	ERROR_OSC_INVALID_CHARACTER = errors.New("OSC Address string may not contain any characters in \"*?,[]{}#")
	ERROR_OSC_ADDRESS_EXISTS    = errors.New("OSC address exists already")
	ERROR_MESSAGE_IS_NIL        = errors.New("message is nil")
	ERROR_UNSUPORTED_PACKAGE    = errors.New("unsupported OSC packet type: only Bundle and Message are supported")
)
