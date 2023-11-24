package osc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimetag(t *testing.T) {
	t.Run("should create a TimeTag", func(t *testing.T) {
		ti := time.Now()
		tt := NewTimetag(ti)

		assert.Equal(t, ti, tt.Time())
	})

	t.Run("should create a timetag from timetag", func(t *testing.T) {
		ti := time.Now()
		tt := NewTimetagFromTimetag(timeToTimetag(ti))

		assert.True(t, ti.Equal(tt.time))
	})

	t.Run("should expires in about a minute", func(t *testing.T) {
		ti := time.Now().Add(time.Minute)
		tt := NewTimetag(ti)

		assert.True(t, tt.ExpiresIn() > 59*time.Second)
	})
}
