package osc_test

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"

	"bekuba.de/go-osc"
)

func TestParsePacket(t *testing.T) {
	for _, tt := range []struct {
		desc string
		msg  string
		pkt  osc.Packet
		ok   bool
	}{
		{
			"no_args",
			"/a/b/c" + nulls(2) + "," + nulls(3),
			makePacket("/a/b/c", nil),
			true,
		},
		{
			"string_arg",
			"/d/e/f" + nulls(2) + ",s" + nulls(2) + "foo" + nulls(1),
			makePacket("/d/e/f", []string{"foo"}),
			true,
		},
		{
			"empty_string_arg",
			"/d/e/g" + nulls(2) + ",s" + nulls(2) + "" + nulls(4),
			makePacket("/d/e/g", []string{""}),
			true,
		},
		{"empty", "", nil, false},
	} {
		var start int
		pkt, err := osc.ReadPacket(bufio.NewReader(bytes.NewBufferString(tt.msg)), &start, len(tt.msg))
		if err != nil && tt.ok {
			t.Errorf("%s: readPacket() returned unexpected error; %s", tt.desc, err)
		}
		if err == nil && !tt.ok {
			t.Errorf("%s: readPacket() expected error", tt.desc)
		}
		if !tt.ok {
			continue
		}

		pktBytes, err := pkt.MarshalBinary()
		if err != nil {
			t.Errorf("%s: failure converting pkt to byte array; %s", tt.desc, err)
			continue
		}
		ttpktBytes, err := tt.pkt.MarshalBinary()
		if err != nil {
			t.Errorf("%s: failure converting tt.pkt to byte array; %s", tt.desc, err)
			continue
		}
		if got, want := pktBytes, ttpktBytes; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: readPacket() as bytes = '%s', want = '%s'", tt.desc, got, want)
			continue
		}
	}
}

// makePacket creates a fake Message Packet.
func makePacket(addr string, args []string) osc.Packet {
	msg := osc.NewMessage(addr)
	for _, arg := range args {
		msg.Append(arg)
	}
	return msg
}

const zero = string(byte(0))

// nulls returns a string of `i` nulls.
func nulls(i int) string {
	s := ""
	for j := 0; j < i; j++ {
		s += zero
	}
	return s
}
