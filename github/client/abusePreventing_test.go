package client

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAbusePreventing(t *testing.T) {
	t.Run("Works when there is a successful response after waiting for the timeout", func(t *testing.T) {
		timesCalled := 0

		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			switch timesCalled {
			case 0:
				return &http.Response{
					StatusCode: 403,
					Header: http.Header{
						"Retry-After": []string{"0.1"}, // means 100ms
						"Success":     []string{"no"},
					},
				}, nil
			case 1:
				return successfulResponse()
			default:
				panic("Should not be called")
			}
		}

		client := abusePreventing{
			HTTPClient: httpClientMock{response},
		}

		t1 := time.Now()
		res, err := client.Do(dummyRequest())
		t2 := time.Now()
		secondsPassed := t2.Sub(t1).Seconds()

		require.Nil(t, err)
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Equal(t, 2, timesCalled)
		require.True(
			t,
			secondsPassed >= 0.1 && secondsPassed < 0.2,
			fmt.Sprintf("100ms should pass because of Retry-Later, passed %f", secondsPassed),
		)
	})

	t.Run("Works when Retry-After repeats more than once", func(t *testing.T) {
		timesCalled := 0

		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			switch timesCalled {
			case 0, 1:
				return &http.Response{
					StatusCode: 403,
					Header: http.Header{
						"Retry-After": []string{"0.1"}, // means 100ms
						"Success":     []string{"no"},
					},
				}, nil
			case 2:
				return successfulResponse()
			default:
				panic("Should not be called")
			}
		}

		client := abusePreventing{
			HTTPClient: httpClientMock{response},
		}

		t1 := time.Now()
		res, err := client.Do(dummyRequest())
		t2 := time.Now()
		secondsPassed := t2.Sub(t1).Seconds()

		require.Nil(t, err)
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Equal(t, 3, timesCalled)
		require.True(
			t,
			secondsPassed >= 0.2 && secondsPassed < 0.21,
			fmt.Sprintf("200ms should pass because of Retry-Later, passed %f", secondsPassed),
		)
	})

	t.Run("Works when the response is successful", func(t *testing.T) {
		client := abusePreventing{
			HTTPClient: httpClientMock{successfulResponse},
		}

		res, err := client.Do(dummyRequest())

		require.Nil(t, err)
		require.Equal(t, "yes", res.Header.Get("Success"))
	})

	t.Run("Fails right away if there was an error", func(t *testing.T) {
		client := abusePreventing{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return nil, fmt.Errorf("Some weird network error")
			}},
		}

		res, err := client.Do(dummyRequest())

		require.Nil(t, res)
		require.Contains(t, err.Error(), "Some weird network error")
	})

	t.Run("Works when there is an error after waiting for the timeout", func(t *testing.T) {
		timesCalled := 0

		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			switch timesCalled {
			case 0:
				return &http.Response{
					StatusCode: 403,
					Header: http.Header{
						"Retry-After": []string{"0.1"}, // means 100ms
						"Success":     []string{"no"},
					},
				}, nil
			case 1:
				return nil, fmt.Errorf("The weirdest network error ever")
			default:
				panic("Should not be called")
			}
		}

		client := abusePreventing{
			HTTPClient: httpClientMock{response},
		}

		res, err := client.Do(dummyRequest())

		require.Nil(t, res)
		require.Contains(t, err.Error(), "The weirdest network error ever")
	})

	t.Run("Ignores Retry-Later if it can't be parsed as float", func(t *testing.T) {
		timesCalled := 0

		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			switch timesCalled {
			case 0:
				return &http.Response{
					StatusCode: 403,
					Header: http.Header{
						"Retry-After": []string{"pitty"}, // broken float
						"Success":     []string{"no"},
					},
				}, nil
			case 1:
				return successfulResponse() // actually, shouldn't be called in this test
			default:
				panic("Should not be called")
			}
		}

		client := abusePreventing{
			HTTPClient: httpClientMock{response},
		}

		res, err := client.Do(dummyRequest())

		require.Nil(t, err)
		require.Equal(t, "no", res.Header.Get("Success"))
		require.Equal(t, 1, timesCalled)
	})
}
