package mouse_driver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	movementWaitTime := 10 * time.Millisecond
	timeBetweenActions := 20 * time.Millisecond
	conf := New(movementWaitTime, timeBetweenActions)

	assert.NotNil(t, conf)
	assert.Equal(t, movementWaitTime, conf.MovementWaitTime)
	assert.Equal(t, timeBetweenActions, conf.TimeBetweenActions)
	assert.Len(t, conf.Points, 4)
	assert.Equal(t, 100.0, conf.Points[0].X)
	assert.Equal(t, 100.0, conf.Points[0].Y)
}
