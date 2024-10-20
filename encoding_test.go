package osc_test

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/bekuba/go-osc"
)

func TestPadBytesNeeded(t *testing.T) {
	var n int
	n = osc.PadBytesNeeded(4)
	if n != 0 {
		t.Errorf("Number of pad bytes should be 0 and is: %d", n)
	}

	n = osc.PadBytesNeeded(3)
	if n != 1 {
		t.Errorf("Number of pad bytes should be 1 and is: %d", n)
	}

	n = osc.PadBytesNeeded(2)
	if n != 2 {
		t.Errorf("Number of pad bytes should be 2 and is: %d", n)
	}

	n = osc.PadBytesNeeded(1)
	if n != 3 {
		t.Errorf("Number of pad bytes should be 3 and is: %d", n)
	}

	n = osc.PadBytesNeeded(0)
	if n != 0 {
		t.Errorf("Number of pad bytes should be 0 and is: %d", n)
	}

	n = osc.PadBytesNeeded(5)
	if n != 3 {
		t.Errorf("Number of pad bytes should be 3 and is: %d", n)
	}

	n = osc.PadBytesNeeded(7)
	if n != 1 {
		t.Errorf("Number of pad bytes should be 1 and is: %d", n)
	}

	n = osc.PadBytesNeeded(32)
	if n != 0 {
		t.Errorf("Number of pad bytes should be 0 and is: %d", n)
	}

	n = osc.PadBytesNeeded(63)
	if n != 1 {
		t.Errorf("Number of pad bytes should be 1 and is: %d", n)
	}

	n = osc.PadBytesNeeded(10)
	if n != 2 {
		t.Errorf("Number of pad bytes should be 2 and is: %d", n)
	}
}

func TestWritePaddedString(t *testing.T) {
	for _, tt := range []struct {
		s   string // string
		buf []byte // resulting buffer
		n   int    // bytes expected
	}{
		{"testString", []byte{'t', 'e', 's', 't', 'S', 't', 'r', 'i', 'n', 'g', 0, 0}, 12},
		{"testers", []byte{'t', 'e', 's', 't', 'e', 'r', 's', 0}, 8},
		{"tests", []byte{'t', 'e', 's', 't', 's', 0, 0, 0}, 8},
		{"test", []byte{'t', 'e', 's', 't', 0, 0, 0, 0}, 8},
		{"tes", []byte{'t', 'e', 's', 0}, 4},
		{"tes\x00", []byte{'t', 'e', 's', 0}, 4},                 // Don't add a second null terminator if one is already present
		{"tes\x00\x00\x00\x00\x00", []byte{'t', 'e', 's', 0}, 4}, // Skip extra nulls
		{"tes\x00\x00\x00", []byte{'t', 'e', 's', 0}, 4},         // Even if they don't fall on a 4 byte padding boundary
		{"", []byte{0, 0, 0, 0}, 4},                              // OSC uses null terminated strings, padded to the 4 byte boundary
	} {
		buf := []byte{}
		bytesBuffer := bytes.NewBuffer(buf)

		n, err := osc.WritePaddedString(tt.s, bytesBuffer)
		if err != nil {
			t.Errorf(err.Error())
		}
		if got, want := n, tt.n; got != want {
			t.Errorf("%q: Count of bytes written don't match; got = %d, want = %d", tt.s, got, want)
		}
		if got, want := bytesBuffer, tt.buf; !bytes.Equal(got.Bytes(), want) {
			t.Errorf("%q: Buffers don't match; got = %q, want = %q", tt.s, got.Bytes(), want)
		}
	}
}

func TestReadPaddedString(t *testing.T) {
	for _, tt := range []struct {
		buf []byte // buffer
		n   int    // bytes needed
		s   string // resulting string
		e   error  // expected error
	}{
		{[]byte{'t', 'e', 's', 't', 'S', 't', 'r', 'i', 'n', 'g', 0, 0}, 12, "testString", nil},
		{[]byte{'t', 'e', 's', 't', 'e', 'r', 's', 0}, 8, "testers", nil},
		{[]byte{'t', 'e', 's', 't', 's', 0, 0, 0}, 8, "tests", nil},
		{[]byte{'t', 'e', 's', 't', 0, 0, 0, 0}, 8, "test", nil},
		{[]byte{}, 0, "", io.EOF},
		{[]byte{'t', 'e', 's', 0}, 4, "tes", nil},             // OSC uses null terminated strings
		{[]byte{'t', 'e', 's', 0, 0, 0, 0, 0}, 4, "tes", nil}, // Additional nulls should be ignored
		{[]byte{'t', 'e', 's', 0, 0, 0}, 4, "tes", nil},       // Whether or not the nulls fall on a 4 byte padding boundary
		{[]byte{'t', 'e', 's', 't'}, 0, "", io.EOF},           // if there is no null byte at the end, it doesn't work.
	} {
		buf := bytes.NewBuffer(tt.buf)
		s, n, err := osc.ReadPaddedString(bufio.NewReader(buf))
		if got, want := err, tt.e; got != want {
			t.Errorf("%q: Unexpected error reading padded string; got = %q, want = %q", tt.s, got, want)
		}
		if got, want := n, tt.n; got != want {
			t.Errorf("%q: Bytes needed don't match; got = %d, want = %d", tt.s, got, want)
		}
		if got, want := s, tt.s; got != want {
			t.Errorf("%q: Strings don't match; got = %q, want = %q", tt.s, got, want)
		}
	}
}

func TestReadBlob(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    []byte
		want    []byte
		want1   int
		wantErr bool
	}{
		{"negative value", []byte{255, 255, 255, 255}, nil, 0, true},
		{"large value", []byte{0, 1, 17, 112}, nil, 0, true},
		{"regular value", []byte{0, 0, 0, 1, 10, 0, 0, 0}, []byte{10}, 8, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := osc.ReadBlob(bufio.NewReader(bytes.NewBuffer(tt.args)))
			if (err != nil) != tt.wantErr {
				t.Errorf("readBlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readBlob() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readBlob() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
