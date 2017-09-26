package client

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// errorWrapping is a HTTPClient which translates HTTP Response codes >= 300 into errors
type errorWrapping struct {
	HTTPClient // "back-end" HTTPClient to use for actual HTTP queries
}

// Do is HTTPClient.Do
func (c errorWrapping) Do(req *http.Request) (*http.Response, error) {
	res, err := c.HTTPClient.Do(req)

	if res.StatusCode >= 300 {
		err = fmt.Errorf(
			"Wrong HTTP response code: %d, Details:\n%s",
			res.StatusCode,
			composeHTTPError(req, res),
		)
		res = nil
	}

	return res, err
}

func composeHTTPError(req *http.Request, res *http.Response) error {
	dump, err := httputil.DumpResponse(res, true)

	if err != nil {
		return fmt.Errorf("<Error printing details: %s>", err.Error())
	}

	return fmt.Errorf(
		"HTTP Request failed.\nURL: %s\n\n%s",
		req.URL.String(),
		dump,
	)
}
