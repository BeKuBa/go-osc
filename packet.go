package osc

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

const (
	bundleTagString = "#bundle"
)

// Packet is the interface for Message and Bundle.
type Packet interface {
	MarshalBinary() (data []byte, err error)
}

// receivePacket receives an OSC packet from the given reader.
func readPacket(reader *bufio.Reader, start *int, end int) (Packet, error) {
	// var buf []byte
	buf, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}

	switch buf[0] {
	case '/':
		return readMessage(reader, start)

	case '#':
		return readBundle(reader, start, end)
	}

	return nil, ERROR_INVALID_PACKET
}

// readBundle reads an Bundle from reader.
func readBundle(reader *bufio.Reader, start *int, end int) (*Bundle, error) {
	// Read the '#bundle' OSC string
	startTag, n, err := readPaddedString(reader)
	if err != nil {
		return nil, err
	}
	*start += n

	if startTag != bundleTagString {
		return nil, fmt.Errorf("Invalid bundle start tag: %s", startTag)
	}

	// Read the timetag
	var timeTag uint64
	err = binary.Read(reader, binary.BigEndian, &timeTag)
	if err != nil {
		return nil, err
	}

	*start += 8

	// Create a new bundle
	bundle := NewBundle(timetagToTime(Timetag(timeTag)))

	// Read until the end of the buffer
	//
	for (end - *start) > 4 {
		// Read the size of the bundle element
		var length int32

		err = binary.Read(reader, binary.BigEndian, &length)
		if err != nil {
			return nil, err
		}
		if length == 0 {
			break
		}

		*start += 4

		p, err := readPacket(reader, start, end)
		if err != nil {
			return nil, err
		}

		err = bundle.Append(p)
		if err != nil {
			return nil, err
		}
	}

	return bundle, nil
}

// readMessage from `reader`.
func readMessage(reader *bufio.Reader, start *int) (*Message, error) {
	// First, read the OSC address
	addr, n, err := readPaddedString(reader)
	if err != nil {
		return nil, err
	}
	*start += n

	// Read all arguments
	msg := NewMessage(addr)

	err = readArguments(msg, reader, start)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// readArguments from `reader` and add them to the OSC message `msg`.
func readArguments(msg *Message, reader *bufio.Reader, start *int) error {
	// Read the type tag string
	var n int
	typetags, n, err := readPaddedString(reader)
	if err != nil {
		return err
	}
	*start += n

	if len(typetags) == 0 {
		return nil
	}

	// If the typetag doesn't start with ',', it's not valid
	if typetags[0] != ',' {
		return fmt.Errorf("unsupported type tag string %s", typetags)
	}

	// Remove ',' from the type tag
	typetags = typetags[1:]

	for _, c := range typetags {
		switch c {
		case 'i': // int32
			var i int32

			err = binary.Read(reader, binary.BigEndian, &i)
			if err != nil {
				return err
			}

			*start += 4
			msg.Append(i)

		case 'h': // int64
			var i int64
			err = binary.Read(reader, binary.BigEndian, &i)
			if err != nil {
				return err
			}
			*start += 8
			msg.Append(i)

		case 'f': // float32
			var f float32
			err = binary.Read(reader, binary.BigEndian, &f)
			if err != nil {
				return err
			}
			*start += 4
			msg.Append(f)

		case 'd': // float64/double
			var d float64
			err = binary.Read(reader, binary.BigEndian, &d)
			if err != nil {
				return err
			}
			*start += 8
			msg.Append(d)

		case 's': // string
			// TODO: fix reading string value
			var s string
			s, _, err = readPaddedString(reader)
			if err != nil {
				return err
			}
			lenStr := len(s)
			*start += lenStr + padBytesNeeded(lenStr)
			msg.Append(s)

		case 'b': // blob
			var buf []byte
			var n int
			buf, n, err = readBlob(reader)
			if err != nil {
				return err
			}
			*start += n
			msg.Append(buf)

		case 't': // OSC time tag
			var tt uint64

			err = binary.Read(reader, binary.BigEndian, &tt)
			if err != nil {
				return nil
			}

			*start += 8
			msg.Append(Timetag(tt))

		case 'N': // nil
			msg.Append(nil)

		case 'T': // true
			msg.Append(true)

		case 'F': // false
			msg.Append(false)

		default:
			return fmt.Errorf("unsupported type tag: %c", c)
		}
	}

	return nil
}
