package github

import (
	"fmt"
	"net/http"
)

// errorWrapping is a HTTPClient which translates HTTP Response codes >= 300 into errors
type errorWrapping struct {
	HTTPClient // "back-end" HTTPClient to use for actual HTTP queries
}

// Do is HTTPClient.Do
func (c errorWrapping) Do(req *http.Request) (*http.Response, error) {
	res, err := c.HTTPClient.Do(req)

	if res.StatusCode >= 300 {
		err = fmt.Errorf("URL: %s\nWrong HTTP response code: %d", req.URL, res.StatusCode)
		res = nil
	}

	return res, err
}
