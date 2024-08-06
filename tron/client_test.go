package tron

import (
	"sync"
	"testing"
	"time"
)

func TestThrottleClient(t *testing.T) {
	throttleTime := 500
	throttleClient := NewThrottleClient(throttleTime)
	var wg sync.WaitGroup

	then := time.Now()
	numRuns := 3

	for i := 0; i < numRuns; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			// we don't care about the response or error
			throttleClient.Get("https://www.google.com")
		}()
	}

	wg.Wait()

	diff := time.Now().Sub(then)
	expectedDifference := int64(throttleTime * numRuns)

	if diff.Milliseconds() < expectedDifference {
		t.Fatal("Expected 3 requests to take at least 1.5 seconds")
	}
}
