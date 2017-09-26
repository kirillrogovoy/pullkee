package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateLimiting(t *testing.T) {
	t.Run("Rate limit works", func(t *testing.T) {
		limiter := time.Tick(time.Millisecond * 50)

		client := rateLimiting{
			HTTPClient:  httpClientMock{successfulResponse},
			RateLimiter: &limiter,
		}

		t1 := time.Now()

		res, err := client.Do(dummyRequest())
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Nil(t, err)

		res, err = client.Do(dummyRequest())
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Nil(t, err)

		t2 := time.Now()
		secondsPassed := t2.Sub(t1).Seconds()

		require.True(
			t,
			secondsPassed >= 0.1 && secondsPassed < 0.15,
			fmt.Sprintf("100ms should pass because of RateLimit, passed %f", secondsPassed),
		)
	})

	t.Run("Still works if there is no RateLimiter", func(t *testing.T) {
		client := rateLimiting{
			HTTPClient: httpClientMock{successfulResponse},
		}

		res, err := client.Do(dummyRequest())
		require.Equal(t, "yes", res.Header.Get("Success"))
		require.Nil(t, err)
	})
}
