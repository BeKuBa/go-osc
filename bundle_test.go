package osc_test

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/bekuba/go-osc"
	"github.com/stretchr/testify/assert"
)

func TestBundle(t *testing.T) {
	bundle := osc.NewBundle(time.Now())

	err := bundle.Append(osc.NewMessage("/a", "test"))
	assert.Nil(t, err)
	err = bundle.Append(osc.NewMessage("/b", "test2"))
	assert.Nil(t, err)

	d, err := bundle.MarshalBinary()
	assert.Nil(t, err)

	t.Run("should read bundle without padding", func(t *testing.T) {
		io := bufio.NewReader(bytes.NewReader(d))
		start := 0
		b, err := osc.ReadBundle(io, &start, len(d))

		assert.Nil(t, err)
		assert.Equal(t, 2, len(b.Messages))
	})

	t.Run("should read bundle with 4 bytes padded", func(t *testing.T) {

		d1 := append(d, 0, 0, 0, 0)

		io := bufio.NewReader(bytes.NewReader(d1))
		start := 0
		b, err := osc.ReadBundle(io, &start, len(d1))

		assert.Nil(t, err)
		assert.Equal(t, 2, len(b.Messages))
	})

	t.Run("should read bundle with 18 bytes padded", func(t *testing.T) {

		d1 := append(d, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

		io := bufio.NewReader(bytes.NewReader(d1))
		start := 0
		b, err := osc.ReadBundle(io, &start, len(d1))

		assert.Nil(t, err)
		assert.Equal(t, 2, len(b.Messages))
	})

	t.Run("should fail read bundle when padding is not well formatted", func(t *testing.T) {

		d1 := append(d, 0, 0, 0, 1)

		io := bufio.NewReader(bytes.NewReader(d1))
		start := 0
		_, err := osc.ReadBundle(io, &start, len(d1))

		assert.NotNil(t, err)
	})

}
