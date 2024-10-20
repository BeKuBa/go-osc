package osc

import "errors"

// OSC Errors
var (
	ErrorOscInvalidCharacter = errors.New("OSC Address string may not contain any characters in \"*?,[]{}#")
	ErrorOscAddress          = errors.New("invalid OSC address")
	ErrorOscAddressFormat    = errors.New("invalid OSC address format")
	ErrorOscAddressExists    = errors.New("OSC address exists already")
	ErrorUnsuportedPackage   = errors.New("unsupported OSC packet type: only Bundle and Message are supported")
	ErrorInvalidPacked       = errors.New("invalid OSC packet")
)
