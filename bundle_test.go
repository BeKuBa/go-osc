package osc

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBundle(t *testing.T) {
	t.Run("should append data to bundle", func(t *testing.T) {
		bundle := NewBundle(time.Now())

		bundle.Append(NewMessage("/a", "test"))
		bundle.Append(NewMessage("/b", "test2"))

		d, err := bundle.MarshalBinary()
		assert.Nil(t, err)

		d = append(d, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

		io := bufio.NewReader(bytes.NewReader(d))
		start := 0
		b, err := readBundle(io, &start, len(d))

		assert.Nil(t, err)
		assert.Equal(t, 2, len(b.Messages))
	})
}
