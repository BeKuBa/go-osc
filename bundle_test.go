package osc

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBundle(t *testing.T) {
	bundle := NewBundle(time.Now())

	bundle.Append(NewMessage("/a", "test"))
	bundle.Append(NewMessage("/b", "test2"))

	d, err := bundle.MarshalBinary()
	assert.Nil(t, err)

	t.Run("bundle as it is", func(t *testing.T) {
		io := bufio.NewReader(bytes.NewReader(d))
		start := 0
		b, err := readBundle(io, &start, len(d))

		assert.Nil(t, err)
		assert.Equal(t, 2, len(b.Messages))
	})

	t.Run("should append data(4 nulls) to bundle", func(t *testing.T) {

		d1 := append(d, 0, 0, 0, 0)

		io := bufio.NewReader(bytes.NewReader(d1))
		start := 0
		b, err := readBundle(io, &start, len(d1))

		assert.Nil(t, err)
		assert.Equal(t, 2, len(b.Messages))
	})

	t.Run("should append data(18 nulls) to bundle", func(t *testing.T) {

		d1 := append(d, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

		io := bufio.NewReader(bytes.NewReader(d1))
		start := 0
		b, err := readBundle(io, &start, len(d1))

		assert.Nil(t, err)
		assert.Equal(t, 2, len(b.Messages))
	})

	t.Run("append data(0,0,0,1) to bundle(error expected)", func(t *testing.T) {

		d1 := append(d, 0, 0, 0, 1)

		io := bufio.NewReader(bytes.NewReader(d1))
		start := 0
		_, err := readBundle(io, &start, len(d1))

		assert.NotNil(t, err)
	})

}
