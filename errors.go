package osc

import "errors"

// OSC Errors
var (
	ErrorOscInvalidCharacter = errors.New("OSC Address string may not contain any characters in \"*?,[]{}#")
	ErrorOscAddressExists    = errors.New("OSC address exists already")
	ErrorUnsuportedPackage   = errors.New("unsupported OSC packet type: only Bundle and Message are supported")
	ErrorInvalidPacked       = errors.New("invalid OSC packet")
)
