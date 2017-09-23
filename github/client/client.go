package github

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPClient is in interface for a HTTP client which "transform" a http.Request into http.Response (like http.Client)
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Log is a Printf kind of function which is used to report for logging
type Log func(message string)

// Client is a HTTPClient aware of Github API rules (like rate limiting, http codes, etc.)
type Client struct {
	HTTPClient
	LastResponse *http.Response
	Log
}

// Options is a set of configurable options for Client
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

// New creates a new instance of Client chaining all the different clients together
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

/*
 * func composeHTTPError(req *http.Request, res *http.Response) error {
 *     dump, err := httputil.DumpResponse(res, true)
 *
 *     if err != nil {
 *         panic(err)
 *     }
 *
 *     return fmt.Errorf(
 *         "HTTP Request failed.\nURL: %s\n\n%s",
 *         req.URL.String(),
 *         dump,
 *     )
 * }
 */
