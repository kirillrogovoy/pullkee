package github

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Run("Works fine, saves the last response, writes logs", func(t *testing.T) {
		var lastLog string
		logMock := Log(func(message string) {
			lastLog = message
		})

		client := Client{
			HTTPClient: httpClientMock{successfulResponse},
			Log:        &logMock,
		}

		res, err := client.Do(dummyRequest())

		require.Nil(t, err)
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Exactly(t, res, client.LastResponse)
	})
}

func TestNew(t *testing.T) {
	opts := Options{
		Credentials: nil,
		RateLimiter: nil,
		MaxRetries:  1,
		Log:         nil,
	}

	client := New(httpClientMock{successfulResponse}, opts)
	res, err := client.Do(dummyRequest())

	require.Nil(t, err)
	require.Equal(t, "yes", res.Header.Get("Success"))
}

type httpClientMock struct {
	response func() (*http.Response, error)
}

func (h httpClientMock) Do(request *http.Request) (*http.Response, error) {
	return h.response()
}

func dummyRequest() *http.Request {
	req, _ := http.NewRequest("GET", "http://example.com/url1", nil)
	return req
}

func successfulResponse() (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Success": []string{"yes"},
		},
	}, nil
}
