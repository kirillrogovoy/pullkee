package github

import "net/http"

// retrying is a HTTPClient which re-tries the request for a number of times in case of network errors
type retrying struct {
	HTTPClient // "back-end" HTTPClient to use for actual HTTP queries
	MaxRetries int
}

// Do is HTTPClient.Do
func (c retrying) Do(req *http.Request) (*http.Response, error) {
	res, err := c.HTTPClient.Do(req)

	retriesLeft := c.MaxRetries
	for err != nil && retriesLeft > 0 {
		retriesLeft--
		res, err = c.HTTPClient.Do(req)
	}

	return res, err
}
