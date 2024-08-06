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

// NewThrottleClient throttle time in ms
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
