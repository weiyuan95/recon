package tron

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestThrottleClient(t *testing.T) {
	throttleClient := NewThrottleClient(1)
	var wg sync.WaitGroup

	then := time.Now()

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			// we don't care about the response or error
			throttleClient.Get("https://www.google.com")
		}()
	}

	wg.Wait()

	diff := time.Now().Sub(then)
	fmt.Println("diff ->", diff.Seconds())

	if diff.Seconds() < 5 {
		t.Fatal("Expected 5 requests to take at least 10 seconds")
	}
}
