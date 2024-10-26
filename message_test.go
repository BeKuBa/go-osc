package osc_test

import (
	"testing"

	"bekuba/go-osc"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {

	var oscAddress string = "/address"

	t.Run("should append data to message", func(t *testing.T) {
		message := osc.NewMessage(oscAddress)

		assert.Equal(t, oscAddress, message.Address)

		err := message.Append("string argument")
		assert.Nil(t, err)
		err = message.Append(int32(123456789))
		assert.Nil(t, err)
		err = message.Append(true)
		assert.Nil(t, err)

		assert.Equal(t, 3, len(message.Arguments))
	})

	t.Run("should message equal to another message", func(t *testing.T) {
		msg1 := osc.NewMessage(oscAddress)
		msg2 := osc.NewMessage(oscAddress)
		err := msg1.Append(int64(1234))
		assert.Nil(t, err)
		err = msg2.Append(int64(1234))
		assert.Nil(t, err)
		err = msg1.Append("test string")
		assert.Nil(t, err)
		err = msg2.Append("test string")
		assert.Nil(t, err)

		assert.True(t, msg1.Equals(msg2))
	})

	t.Run("unsupported type int throws error", func(t *testing.T) {
		msg1 := osc.NewMessage(oscAddress)
		err := msg1.Append(1234)
		assert.NotNil(t, err)
	})
}

func TestMessage_TypeTags(t *testing.T) {
	for _, tt := range []struct {
		desc string
		msg  *osc.Message
		tags string
		ok   bool
	}{
		{"addr_only", osc.NewMessage("/"), ",", true},
		{"nil", osc.NewMessage("/", nil), ",N", true},
		{"bool_true", osc.NewMessage("/", true), ",T", true},
		{"bool_false", osc.NewMessage("/", false), ",F", true},
		{"int32", osc.NewMessage("/", int32(1)), ",i", true},
		{"int64", osc.NewMessage("/", int64(2)), ",h", true},
		{"float32", osc.NewMessage("/", float32(3.0)), ",f", true},
		{"float64", osc.NewMessage("/", float64(4.0)), ",d", true},
		{"string", osc.NewMessage("/", "5"), ",s", true},
		{"[]byte", osc.NewMessage("/", []byte{'6'}), ",b", true},
		{"two_args", osc.NewMessage("/", "123", int32(456)), ",si", true},
	} {
		tags := tt.msg.TypeTags()
		if got, want := tags, tt.tags; got != want {
			t.Errorf("%s: TypeTags() = '%s', want = '%s'", tt.desc, got, want)
		}
	}
}

func TestMessage_String(t *testing.T) {
	for _, tt := range []struct {
		desc string
		msg  *osc.Message
		str  string
	}{
		{"nil message", nil, ""},
		{"message with 1 nil argument", osc.NewMessage("/foo/bar", nil), "/foo/bar ,N Nil"},
		{"addr_only", osc.NewMessage("/foo/bar"), "/foo/bar ,"},
		{"one_addr", osc.NewMessage("/foo/bar", "123"), "/foo/bar ,s \"123\""},
		{"two_args", osc.NewMessage("/foo/bar", "123", int32(456)), "/foo/bar ,si \"123\" 456"},
		{"timetag", osc.NewMessage("/foo/bar", osc.Timetag(16818286200017484014)), "/foo/bar ,t 16818286200017484014"},
		{"bytes", osc.NewMessage("/foo/bar", []byte{51, 52, 53}), "/foo/bar ,b [51 52 53]"},
	} {
		if got, want := tt.msg.String(), tt.str; got != want {
			t.Errorf("%s: String() = '%s', want = '%s'", tt.desc, got, want)
		}
	}
}

func TestTypeTagsString(t *testing.T) {
	msg := osc.NewMessage("/some/address")
	msg.Append(int32(100))
	msg.Append(true)
	msg.Append(false)

	typeTags := msg.TypeTags()

	if typeTags != ",iTF" {
		t.Errorf("Type tag string should be ',iTF' and is: %s", typeTags)
	}
}

func TestOscMessageMatch(t *testing.T) {
	tc := []struct {
		desc        string
		addr        string
		addrPattern string
		want        bool
	}{
		{
			"match everything",
			"*",
			"/a/b",
			true,
		},
		{
			"don't match",
			"/a/b",
			"/a",
			false,
		},
		{
			"don't match",
			"/a",
			"/a/b",
			false,
		},
		{
			"match alternatives",
			"/a/{foo,bar}",
			"/a/foo",
			true,
		},
		{
			"don't match if address is not part of the alternatives",
			"/a/{foo,bar}",
			"/a/bob",
			false,
		},
	}

	for _, tt := range tc {
		msg := osc.NewMessage(tt.addr)

		got := msg.Match(tt.addrPattern)
		if got != tt.want {
			t.Errorf("%s: msg('%v').Match('%s') = '%t', want = '%t'", tt.desc, tt.addr, tt.addrPattern, got, tt.want)
		}
	}
}

func TestArgumentGetter(t *testing.T) {

	//bool | int32 | int64 | float32 | float64 | string | []byte | Timetag | nil

	const (
		cInt32   int32   = 1
		cInt64   int64   = 2
		cFloat32 float32 = 3.0
		cFloat64 float64 = 4.0
		cString  string  = "5"

		cTimetag osc.Timetag = 16818286200017484014
	)
	var cBytes []byte = []byte{byte(7), byte(17), byte(27)}
	// true, false, nil

	var vInt32 int32
	var vInt64 int64
	var vFloat32 float32

	var vFloat64 float64
	var vString string
	var vTimetag osc.Timetag
	var vBytes []byte
	var vTrue bool
	var vFalse bool
	var vNil any // true as dummy for not nil

	msg := osc.NewMessage("/argtest", cInt32, cInt64, cFloat32, cFloat64, cString, cTimetag, cBytes, true, false, nil)

	//check values

	vInt32, err := msg.Arguments.Int32(0)
	assert.NoError(t, err)
	assert.Equal(t, cInt32, vInt32)

	vInt64, err = msg.Arguments.Int64(1)
	assert.NoError(t, err)
	assert.Equal(t, cInt64, vInt64)

	vFloat32, err = msg.Arguments.Float32(2)
	assert.NoError(t, err)
	assert.Equal(t, cFloat32, vFloat32)

	vFloat64, err = msg.Arguments.Float64(3)
	assert.NoError(t, err)
	assert.Equal(t, cFloat64, vFloat64)

	vString, err = msg.Arguments.Str(4)
	assert.NoError(t, err)
	assert.Equal(t, cString, vString)

	vTimetag, err = msg.Arguments.Timetag(5)
	assert.NoError(t, err)
	assert.Equal(t, cTimetag, vTimetag)

	vBytes, err = msg.Arguments.Bytes(6)
	assert.NoError(t, err)
	assert.Equal(t, cBytes, vBytes)

	vTrue, err = msg.Arguments.Bool(7)
	assert.NoError(t, err)
	assert.Equal(t, true, vTrue)

	vFalse, err = msg.Arguments.Bool(8)
	assert.NoError(t, err)
	assert.Equal(t, false, vFalse)

	vNil = msg.Arguments.Nil(9)
	assert.Equal(t, nil, vNil)

	// must throw error on wrong type
	for ix := 0; ix < 10; ix++ {
		if ix != 0 {
			_, err := msg.Arguments.Int32(ix)
			assert.Error(t, err)
		}
		if ix != 1 {
			_, err := msg.Arguments.Int64(ix)
			assert.Error(t, err)
		}
		if ix != 2 {
			_, err := msg.Arguments.Float32(ix)
			assert.Error(t, err)
		}
		if ix != 3 {
			_, err := msg.Arguments.Float64(ix)
			assert.Error(t, err)
		}
		if ix != 5 {
			_, err := msg.Arguments.Timetag(ix)
			assert.Error(t, err)
		}
		if ix != 6 {
			_, err := msg.Arguments.Bytes(ix)
			assert.Error(t, err)
		}
		if (ix != 7) && (ix != 8) {
			_, err := msg.Arguments.Bool(ix)
			assert.Error(t, err)
		}
		if ix != 9 {
			n := msg.Arguments.Nil(ix)
			assert.NotNil(t, n)
		}
	}

}

func TestClearMessage(t *testing.T) {
	msg := osc.NewMessage("/msg", int32(4), "msg")
	assert.Equal(t, "/msg", msg.Address)
	assert.Equal(t, 2, len(msg.Arguments))
	msg.Clear()
	assert.Equal(t, "", msg.Address)
	assert.Equal(t, 0, len(msg.Arguments))
}

func TestMatchPanic(t *testing.T) {
	msg := osc.NewMessage("}/")
	assert.Panics(t, func() { _ = msg.Match("/msg") })

}
