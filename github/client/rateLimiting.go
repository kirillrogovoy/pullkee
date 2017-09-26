package client

import (
	"net/http"
	"time"
)

// rateLimiting is a HTTPClient which limits the rate of HTTP queries
// towards the server to avoid abuse mechanism triggering
type rateLimiting struct {
	HTTPClient                    // "back-end" HTTPClient to use for actual HTTP queries
	RateLimiter *<-chan time.Time // Example: time.Tick(time.Millisecond * 75)
}

// Do is HTTPClient.Do
func (c rateLimiting) Do(req *http.Request) (*http.Response, error) {
	if c.RateLimiter != nil {
		<-*c.RateLimiter
	}
	return c.HTTPClient.Do(req)

}
