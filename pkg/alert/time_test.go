package alert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSystemTimerIsNotOver(t *testing.T) {
	timer := NewSystemTimer(time.Hour)
	timer.Start()

	over := timer.Over()
	time.Sleep(100 * time.Millisecond)
	over2 := timer.Over()

	assert.False(t, over)
	assert.False(t, over2)
}

func TestSystemTimerIsOver(t *testing.T) {
	timer := NewSystemTimer(100 * time.Millisecond)
	timer.Start()

	time.Sleep(250 * time.Millisecond)
	over := timer.Over()

	assert.True(t, over)
}

// Logical because it is always after Epoch
func TestSystemTimerAlwaysReturnsTrueIfNotStarted(t *testing.T) {
	timer := NewSystemTimer(time.Millisecond)

	res1 := timer.Over()
	time.Sleep(2 * time.Millisecond)
	res2 := timer.Over()
	res3 := timer.Over()

	assert.True(t, res1)
	assert.True(t, res2)
	assert.True(t, res3)
}
