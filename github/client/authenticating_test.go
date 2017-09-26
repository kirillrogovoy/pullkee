package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	t.Run("Auth headers are set, when Credentials are present", func(t *testing.T) {
		request := dummyRequest()

		client := authenticating{
			HTTPClient: httpClientMock{successfulResponse},
			Credentials: &Credentials{
				Username:            "User1",
				PersonalAccessToken: "Token1",
			},
		}

		res, err := client.Do(request)

		require.Nil(t, err)
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Equal(t, http.Header{
			"Authorization": []string{"Basic VXNlcjE6VG9rZW4x"},
			"User-Agent":    []string{"User1"},
		}, request.Header)
	})

	t.Run("Still works when Credentials aren't present", func(t *testing.T) {
		request := dummyRequest()

		client := authenticating{
			HTTPClient: httpClientMock{successfulResponse},
		}

		res, err := client.Do(dummyRequest())

		require.Nil(t, err)
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Equal(t, http.Header{}, request.Header)
	})
}
