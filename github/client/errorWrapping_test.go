package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorWrapping(t *testing.T) {
	t.Run("Successfully turns a 404 into an error", func(t *testing.T) {
		request := dummyRequest()

		client := errorWrapping{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 404,
				}, nil
			}},
		}

		res, err := client.Do(request)

		require.Nil(t, res)
		require.Contains(t, err.Error(), "Wrong HTTP response code: 404")
	})

	t.Run("Still works when the response code is OK", func(t *testing.T) {
		request := dummyRequest()

		client := errorWrapping{
			HTTPClient: httpClientMock{successfulResponse},
		}

		res, err := client.Do(dummyRequest())

		require.Nil(t, err)
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Equal(t, http.Header{}, request.Header)
	})
}
