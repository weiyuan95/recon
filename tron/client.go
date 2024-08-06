package tron

import (
	"net/http"
	"sync"
	"time"
)

type ThrottleClient struct {
	mu       sync.Mutex
	throttle int
	client   http.Client
}

// NewThrottleClient A thin wrapper around the net/http client
// which throttles requests by a given amount of milliseconds. This process
// is thread safe and only allows one request to be made at any given time.
// This is necessary since the TronGrid API has a rate limit of 5 requests per second.
func NewThrottleClient(throttle int) *ThrottleClient {
	return &ThrottleClient{
		throttle: throttle,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ThrottleClient) Get(url string) (*http.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	resp, err := c.client.Get(url)

	time.Sleep(time.Duration(c.throttle) * time.Millisecond)

	return resp, err
}
