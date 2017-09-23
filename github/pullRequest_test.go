package github

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClosedPullRequests(t *testing.T) {
	t.Run("Works on good response", func(t *testing.T) {
		pulls, err := successfulClosedPullRequests(0)
		require.Nil(t, err)

		require.Equal(t, 1, pulls[0].Number)
		require.Equal(t, "Body1", pulls[0].Body)

		require.Equal(t, 2, pulls[1].Number)
		require.Equal(t, "Body2", pulls[1].Body)
	})

	t.Run("Works when limit is applied", func(t *testing.T) {
		pulls, err := successfulClosedPullRequests(1)
		require.Nil(t, err)

		require.Len(t, pulls, 1)
		require.Equal(t, 1, pulls[0].Number)
		require.Equal(t, "Body1", pulls[0].Body)
	})

	t.Run("Works when limit is greater than len", func(t *testing.T) {
		pulls, err := successfulClosedPullRequests(999)
		require.Nil(t, err)

		require.Len(t, pulls, 2)
		require.Equal(t, 1, pulls[0].Number)
		require.Equal(t, "Body1", pulls[0].Body)

		require.Equal(t, 2, pulls[1].Number)
		require.Equal(t, "Body2", pulls[1].Body)
	})

	t.Run("Fails when there is an error fetching the response", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return nil, fmt.Errorf("Dogs have chewed the wires")
			}},
			RepoName: "someuser/somerepo",
		}

		pulls, err := a.ClosedPullRequests(0)
		require.EqualError(t, err, "Dogs have chewed the wires")
		require.Nil(t, pulls)
	})
}

func TestIsMerged(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		pull := PullRequest{
			State:    "closed",
			MergedAt: time.Now(),
		}

		require.True(t, pull.IsMerged())
	})

	t.Run("Negative", func(t *testing.T) {
		pull := PullRequest{
			State:    "closed",
			MergedAt: time.Time{},
		}

		require.False(t, pull.IsMerged())
	})
}

func TestFillDetails(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		pr := PullRequest{
			Number: 11,
		}

		a := apiMock{}

		err := pr.FillDetails(a)

		require.Nil(t, err)
		require.Equal(t, 100, *pr.DiffSize)
		require.Equal(t, "Neat!", (*pr.Comments)[0].Body)
		require.Equal(t, "User1", (*pr.ReviewRequests)[0].Login)
	})

	t.Run("Fails when couldn't fetch the diff size", func(t *testing.T) {
		pr := PullRequest{
			Number: 11,
		}

		err := pr.FillDetails(apiMock{
			diffSizeErr: fmt.Errorf("Weird error"),
		})
		require.EqualError(t, err, "Weird error")
	})

	t.Run("Fails when couldn't fetch the comments", func(t *testing.T) {
		pr := PullRequest{
			Number: 11,
		}

		err := pr.FillDetails(apiMock{
			commentsErr: fmt.Errorf("Weird error"),
		})
		require.EqualError(t, err, "Weird error")
	})

	t.Run("Fails when couldn't fetch the review requests", func(t *testing.T) {
		pr := PullRequest{
			Number: 11,
		}

		err := pr.FillDetails(apiMock{
			reviewRequestsErr: fmt.Errorf("Weird error"),
		})
		require.EqualError(t, err, "Weird error")
	})
}

func successfulClosedPullRequests(limit int) ([]PullRequest, error) {
	goodLink := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next"`

	timesCalled := 0
	response := func() (*http.Response, error) {
		defer func() { timesCalled++ }()
		switch timesCalled {
		case 0:
			json := `[{"number": 1, "body": "Body1"}]`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(json)),
				Header: http.Header{
					"Link": []string{goodLink},
				},
			}, nil
		case 1:
			json := `[{"number": 2, "body": "Body2"}]`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(json)),
			}, nil
		default:
			panic("should never be called")
		}
	}
	a := APIv3{
		HTTPClient: httpClientMock{response},
		RepoName:   "someuser/somerepo",
	}

	return a.ClosedPullRequests(limit)
}

type apiMock struct {
	diffSizeErr       error
	commentsErr       error
	reviewRequestsErr error
}

func (a apiMock) Get(url string, target interface{}) error {
	panic("not implemented")
}

func (a apiMock) Repository() (*Repository, error) {
	panic("not implemented")
}

func (a apiMock) ClosedPullRequests(limit int) ([]PullRequest, error) {
	panic("not implemented")
}

func (a apiMock) DiffSize(number int) (int, error) {
	if a.diffSizeErr != nil {
		return 0, a.diffSizeErr
	}

	return 100, nil
}

func (a apiMock) Comments(number int) ([]Comment, error) {
	if a.commentsErr != nil {
		return nil, a.commentsErr
	}

	return []Comment{Comment{Body: "Neat!"}}, nil
}

func (a apiMock) ReviewRequests(number int) ([]User, error) {
	if a.reviewRequestsErr != nil {
		log.Println("not nil")
		return nil, a.reviewRequestsErr
	}

	return []User{User{"User1"}}, nil
}
