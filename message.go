package osc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
)

// Datatype for Arguments
type ArgumentsType []any

// Message represents a single OSC message. An OSC message consists of an OSC
// address pattern and zero or more arguments.
type Message struct {
	Address   string
	Arguments ArgumentsType
}

// Verify that Messages implements the Packet interface.
// var _ Packet = (*Message)(nil)

// Append appends the given arguments to the arguments list.
func (msg *Message) Append(args ...any) error {
	// check types of args

	for _, arg := range args {
		switch t := arg.(type) {

		// OSC types are ok
		case bool, int32, int64, float32, float64, string, nil, []byte, Timetag: // do nothing
		// type is not an OSC type
		default:
			return fmt.Errorf("unsupported type: %T", t)
		}
	}

	msg.Arguments = append(msg.Arguments, args...)
	return nil
}

// Equals returns true if the given OSC Message `m` is equal to the current OSC
// Message. It checks if the OSC address and the arguments are equal. Returns
// true if the current object and `m` are equal.
func (msg *Message) Equals(m *Message) bool {
	return reflect.DeepEqual(msg, m)
}

// Clear clears the OSC address and all arguments.
func (msg *Message) Clear() {
	msg.Address = ""
	msg.ClearData()
}

// ClearData removes all arguments from the OSC Message.
func (msg *Message) ClearData() {
	msg.Arguments = msg.Arguments[len(msg.Arguments):]
}

// Match returns true, if the OSC address pattern of the OSC Message matches the given
// address. The match is case sensitive!
func (msg *Message) Match(addr string) bool {
	regex, err := getRegEx(msg.Address)
	if err != nil {
		if err != nil {
			panic("regexp: Compile(msg.Address): " + err.Error())
		}
	}
	return regex.MatchString(addr)
}

// typeTags returns the type tag string.
func (msg *Message) typeTags() string {
	if len(msg.Arguments) == 0 {
		return ","
	}

	var tags strings.Builder
	_ = tags.WriteByte(',')

	for _, m := range msg.Arguments {
		tags.WriteByte(getTypeTag(m))
	}

	return tags.String()
}

// String implements the fmt.Stringer interface.
func (msg *Message) String() string {
	if msg == nil {
		return ""
	}

	var s strings.Builder
	tags := msg.typeTags()
	s.WriteString(fmt.Sprintf("%s %s", msg.Address, tags))

	for _, arg := range msg.Arguments {
		switch argType := (arg).(type) {
		case bool, int32, int64, float32, float64:
			s.WriteString(fmt.Sprintf(" %v", argType))
		case string:
			s.WriteString(fmt.Sprintf(" %q", argType))
		case nil:
			s.WriteString(" Nil")
		case []byte:
			s.WriteString(fmt.Sprintf(" %d", argType))

		case Timetag:
			s.WriteString(fmt.Sprintf(" %d", Timetag(argType)))
		}
	}

	return s.String()
}

// MarshalBinary serializes the OSC message to a byte buffer. The byte buffer
// has the following format:
// 1. OSC Address Pattern
// 2. OSC Type Tag String
// 3. OSC Arguments.
func (msg *Message) MarshalBinary() ([]byte, error) {
	// We can start with the OSC address and add it to the buffer
	data := new(bytes.Buffer)

	_, err := writePaddedString(msg.Address, data)
	if err != nil {
		return nil, err
	}

	// Type tag string starts with ","
	lenArgs := len(msg.Arguments)
	typetags := make([]byte, lenArgs+1)
	typetags[0] = ','

	// Process the type tags and collect all arguments
	payload := new(bytes.Buffer)

	for i, arg := range msg.Arguments {
		switch t := arg.(type) {
		case bool:
			if t {
				typetags[i+1] = 'T'
				continue
			}

			typetags[i+1] = 'F'

		case nil:
			typetags[i+1] = 'N'

		case int32:
			typetags[i+1] = 'i'

			err = binary.Write(payload, binary.BigEndian, t)
			if err != nil {
				return nil, err
			}

		case float32:
			typetags[i+1] = 'f'

			err := binary.Write(payload, binary.BigEndian, t)
			if err != nil {
				return nil, err
			}

		case string:
			typetags[i+1] = 's'

			_, err = writePaddedString(t, payload)
			if err != nil {
				return nil, err
			}

		case []byte:
			typetags[i+1] = 'b'

			_, err = writeBlob(t, payload)
			if err != nil {
				return nil, err
			}

		case int64:
			typetags[i+1] = 'h'
			err = binary.Write(payload, binary.BigEndian, t)
			if err != nil {
				return nil, err
			}

		case float64:
			typetags[i+1] = 'd'

			err = binary.Write(payload, binary.BigEndian, t)
			if err != nil {
				return nil, err
			}

		case Timetag:
			typetags[i+1] = 't'

			b, err := t.MarshalBinary()
			if err != nil {
				return nil, err
			}

			_, err = payload.Write(b)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported type: %T", t)
		}
	}

	// Write the type tag string to the data buffer
	if _, err := writePaddedString(string(typetags), data); err != nil {
		return nil, err
	}

	// Write the payload (OSC arguments) to the data buffer
	if _, err := data.Write(payload.Bytes()); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

// NewMessage returns a new Message. The address parameter is the OSC address.
// if args has invalid types it return nil
func NewMessage(addr string, args ...any) *Message {
	msg := &Message{Address: addr}
	err := msg.Append(args...)
	if err != nil {
		panic(err)
	}

	return msg
}

// Help function for argument getter
func (args ArgumentsType) arg(ix int) (result any, err error) {
	if ix >= 0 && ix < len(args) {
		return args[ix], nil
	}
	return nil, fmt.Errorf("out of bounds")
}

// Argument getter for bool value
func (args *ArgumentsType) Bool(ix int) (bool, error) {
	v, err := args.arg(ix)
	if err == nil {
		switch t := v.(type) {
		case bool:
			return t, nil
		default:
			return false, fmt.Errorf("type(%T) is not bool", v)
		}
	}
	return false, err
}

// Argument getter for bool value
func (args *ArgumentsType) Int32(ix int) (int32, error) {
	v, err := args.arg(ix)
	if err == nil {
		switch t := v.(type) {
		case int32:
			return t, nil
		default:
			return 0, fmt.Errorf("type(%T) is not int32", v)
		}
	}
	return 0, err
}

// Argument getter for bool value
func (args *ArgumentsType) Int64(ix int) (int64, error) {
	v, err := args.arg(ix)
	if err == nil {
		switch t := v.(type) {
		case int64:
			return t, nil
		default:
			return 0, fmt.Errorf("type(%T) is not int64", v)
		}
	}
	return 0, err
}

// Argument getter for bool value
func (args *ArgumentsType) Float32(ix int) (float32, error) {
	v, err := args.arg(ix)
	if err == nil {
		switch t := v.(type) {
		case float32:
			return t, nil
		default:
			return 0.0, fmt.Errorf("type(%T) is not float32", v)
		}
	}
	return 0.0, err
}

// Argument getter for bool value
func (args *ArgumentsType) Float64(ix int) (float64, error) {
	v, err := args.arg(ix)
	if err == nil {
		switch t := v.(type) {
		case float64:
			return t, nil
		default:
			return 0.0, fmt.Errorf("type(%T) is not float64", v)
		}
	}
	return 0.0, err
}

// Argument getter for bool value
func (args *ArgumentsType) Str(ix int) (string, error) {
	v, err := args.arg(ix)
	if err == nil {
		switch t := v.(type) {
		case string:
			return t, nil
		default:
			return "", fmt.Errorf("type(%T) is not string", v)
		}
	}
	return "", err
}

// Argument getter for bool value
func (args *ArgumentsType) Bytes(ix int) ([]byte, error) {
	v, err := args.arg(ix)
	if err == nil {
		switch t := v.(type) {
		case []byte:
			return t, nil
		default:
			return nil, fmt.Errorf("type(%T) is not []byte", v)
		}
	}
	return nil, err
}

// Argument getter for bool value
func (args *ArgumentsType) Timetag(ix int) (Timetag, error) {
	v, err := args.arg(ix)
	if err == nil {
		switch t := v.(type) {
		case Timetag:
			return t, nil
		default:
			return 0, fmt.Errorf("type(%T) is not Timetag", v)
		}
	}
	return 0, err
}

// Argument getter for nil value
// also nil if ix out of range
func (args *ArgumentsType) Nil(ix int) any {
	var dummy any = true
	v, err := args.arg(ix)
	if (err == nil) && (v != nil) {
		return dummy
	}
	return nil
}
