package client

import (
	"net/http"
	"strconv"
	"time"
)

// abusePreventing is a HTTPClient which retries a query in case of Retry-Later response header
// https://developer.github.com/v3/guides/best-practices-for-integrators/#dealing-with-abuse-rate-limits
type abusePreventing struct {
	HTTPClient // "back-end" HTTPClient to use for actual HTTP queries
}

// Do is HTTPClient.Do
func (c abusePreventing) Do(req *http.Request) (*http.Response, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	is403 := res.StatusCode == 403
	retryAfter, parseError := strconv.ParseFloat(res.Header.Get("Retry-After"), 64)

	if is403 && parseError == nil {
		c.waitFor(retryAfter)
		res, err = c.Do(req)
	}
	return res, err
}

func (c abusePreventing) waitFor(seconds float64) {
	duration := time.Duration(seconds * float64(time.Second))
	timer := time.NewTimer(duration)
	<-timer.C
}
