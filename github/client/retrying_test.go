package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetrying(t *testing.T) {
	t.Run("It retries again if there is a network error", func(t *testing.T) {
		timesCalled := 0

		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			switch timesCalled {
			case 0:
				return nil, fmt.Errorf("Some weird network error")
			case 1:
				return nil, fmt.Errorf("Another weird network error")
			case 2:
				return successfulResponse()
			default:
				panic("Should not be called")
			}
		}

		client := retrying{
			HTTPClient: httpClientMock{response},
			MaxRetries: 2,
		}

		res, err := client.Do(dummyRequest())

		require.Nil(t, err)
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Equal(t, 3, timesCalled)
	})

	t.Run("It gives up retrying after some number of tries", func(t *testing.T) {
		timesCalled := 0

		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			return nil, fmt.Errorf("Some weird network error which doesn't go away, try #%d", timesCalled+1)
		}

		client := retrying{
			HTTPClient: httpClientMock{response},
			MaxRetries: 2,
		}

		res, err := client.Do(dummyRequest())

		require.Nil(t, res)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "Some weird network error which doesn't go away, try #3")
		require.Equal(t, 3, timesCalled, fmt.Sprintf("Should try for 3 times before giving up. Tried %d times", timesCalled))
	})
}
