package repository

import (
	"fmt"
	"testing"

	"github.com/kirillrogovoy/pullk/github"
	"github.com/stretchr/testify/require"
)

var pullsFromAPI = []github.PullRequest{github.PullRequest{
	Body: "PR from API",
}}

var pullsFromCache = []github.PullRequest{github.PullRequest{
	Body: "PR from cache",
}}

func TestPulls(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		a := apiMock{}

		pulls, err := Pulls(a, 0)

		require.Equal(t, pullsFromAPI, pulls)
		require.Nil(t, err)
	})

	t.Run("Fails when couldn't fetch something from the API", func(t *testing.T) {
		a := apiMock{
			err: fmt.Errorf("Network failed"),
		}

		_, err := Pulls(a, 0)

		require.EqualError(t, err, "Network failed")
	})
}

func TestFillDetails(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		prs := []github.PullRequest{
			github.PullRequest{Number: 1},
			github.PullRequest{Number: 2},
			github.PullRequest{Number: 3},
			github.PullRequest{Number: 4},
		}

		a := apiMock{}
		c := newCacheMock()

		err := FillDetails(a, c, prs)

		require.Nil(t, err)
		require.Len(t, prs, 4)

		require.Equal(t, 100, *prs[0].DiffSize)
		require.Equal(t, 100, *prs[1].DiffSize)
		require.Equal(t, 100, *prs[2].DiffSize)
		require.Equal(t, 100, *prs[3].DiffSize)

		require.Equal(t, prs[0], c.store["pr1"])
		require.Equal(t, prs[3], c.store["pr4"])
	})

	t.Run("Works even when had a cache read error", func(t *testing.T) {
		prs := []github.PullRequest{
			github.PullRequest{Number: 1},
		}

		a := apiMock{}
		c := newCacheMock()
		c.getErr = fmt.Errorf("Nasty cache error")

		err := FillDetails(a, c, prs)

		require.Nil(t, err)
		require.Len(t, prs, 1)

		require.Equal(t, 100, *prs[0].DiffSize)
	})

	t.Run("Works even when had a cache write error", func(t *testing.T) {
		prs := []github.PullRequest{
			github.PullRequest{Number: 1},
		}

		a := apiMock{}
		c := newCacheMock()
		c.setErr = fmt.Errorf("Nasty cache error")

		err := FillDetails(a, c, prs)

		require.Nil(t, err)
		require.Len(t, prs, 1)

		require.Equal(t, 100, *prs[0].DiffSize)
	})

	t.Run("Fails when couldn't fetch details", func(t *testing.T) {
		prs := []github.PullRequest{
			github.PullRequest{Number: 1},
		}

		a := apiMock{
			err: fmt.Errorf("Fetching error"),
		}

		c := newCacheMock()
		err := FillDetails(a, c, prs)

		require.EqualError(t, err, "Fetching error")
	})
}

type cacheMock struct {
	store  map[string]interface{}
	found  bool
	getErr error
	setErr error
}

func newCacheMock() cacheMock {
	c := cacheMock{}
	c.store = map[string]interface{}{}
	return c
}

func (c cacheMock) Set(key string, target interface{}) error {
	if c.setErr != nil {
		return c.setErr
	}

	c.store[key] = target
	return nil
}

func (c cacheMock) Get(key string, target interface{}) (bool, error) {
	if c.getErr != nil {
		return false, c.getErr
	}

	if c.found {
		t := target.(*[]github.PullRequest)
		*t = pullsFromCache
	}
	return c.found, c.getErr
}

type apiMock struct {
	err error
}

func (a apiMock) ClosedPullRequests(limit int) ([]github.PullRequest, error) {
	if a.err != nil {
		return nil, a.err
	}
	return pullsFromAPI, nil
}

func (a apiMock) Get(url string, target interface{}) error {
	panic("not implemented")
}

func (a apiMock) Repository() (*github.Repository, error) {
	panic("not implemented")
}

func (a apiMock) DiffSize(number int) (int, error) {
	if a.err != nil {
		return 0, a.err
	}
	return 100, nil
}

func (a apiMock) Comments(number int) ([]github.Comment, error) {
	return []github.Comment{github.Comment{Body: "Neat!"}}, nil
}

func (a apiMock) ReviewRequests(number int) ([]github.User, error) {
	return []github.User{github.User{
		Login: "User1",
	}}, nil
}
