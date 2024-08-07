package scheduler

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	foo := 0
	ticker := Schedule(func() {
		foo++
	}, 500)

	// on the 2000th ms, the job wouldn't have run,
	// so we wait a little bit longer
	time.Sleep(2100 * time.Millisecond)
	ticker.Stop()

	assert.Equal(t, 4, foo)
}
