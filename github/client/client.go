// Package client contains an implementation of a HTTP client aware of
// all the Github rules like rate limiting or authenticating
package client

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPClient is in interface for a HTTP client which "transforms"
// a http.Request into http.Response (like http.Client)
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Log is a which is used to report for logging
type Log func(message string)

// Client is a top-level HTTPClient that remembers the last
// response and can log information out during the process
type Client struct {
	HTTPClient
	LastResponse *http.Response
	Log
}

// Options is a set of configurable options to create a whole chain of clients
type Options struct {
	*Credentials
	RateLimiter *<-chan time.Time
	MaxRetries  int
	Log
}

// Do is HTTPClient.Do
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.HTTPClient.Do(req)

	c.LastResponse = res

	if c.Log != nil {
		c.Log(fmt.Sprintf(
			"DONE - %s: %s",
			req.Method,
			req.URL.String(),
		))
	}

	return res, err
}

// New creates a new instance of Client chaining all the clients together given Options
func New(httpClient HTTPClient, opts Options) Client {
	retrying := retrying{
		HTTPClient: httpClient,
		MaxRetries: opts.MaxRetries,
	}
	rateLimiting := rateLimiting{
		HTTPClient:  retrying,
		RateLimiter: opts.RateLimiter,
	}
	abusePreventing := abusePreventing{
		HTTPClient: rateLimiting,
	}
	auth := authenticating{
		HTTPClient:  abusePreventing,
		Credentials: opts.Credentials,
	}
	errorWrapping := errorWrapping{
		HTTPClient: auth,
	}
	client := Client{
		HTTPClient: errorWrapping,
		Log:        opts.Log,
	}

	return client
}
