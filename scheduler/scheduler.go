package scheduler

import (
	"time"
)

type Job func()

// Schedule runs a job at a given interval (ms) for an infinite amount of time.
// It returns a ticker that can be stopped by calling ticker.Stop()
func Schedule(job Job, interval int) *time.Ticker {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				job()
			}
		}
	}()

	return ticker
}
