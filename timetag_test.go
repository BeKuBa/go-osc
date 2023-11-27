package osc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimetag(t *testing.T) {
	t.Run("should create an immediate timetag", func(t *testing.T) {
		tt := NewImmediateTimetag()

		assert.Equal(t, tt.Time(), time.Date(1900, time.January, 1, 1, 0, 0, 1, time.Local))
	})

	t.Run("should create a TimeTag", func(t *testing.T) {
		ti := time.Now()
		tt := NewTimetagFromTime(ti)

		assert.True(t, tt.Time().Equal(ti))
	})

	t.Run("should expires in about a minute", func(t *testing.T) {
		ti := time.Now().Add(time.Minute)
		tt := NewTimetagFromTime(ti)

		actual := tt.ExpiresIn().Round(time.Millisecond)

		assert.True(t, actual == 60*time.Second)
	})

	t.Run("should marshall binary an immediate tag", func(t *testing.T) {
		tt := NewImmediateTimetag()

		actual, err := tt.MarshalBinary()
		assert.Nil(t, err)

		assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 1}, actual)
	})
}
