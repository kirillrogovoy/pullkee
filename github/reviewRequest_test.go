package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReviewRequest(t *testing.T) {
	t.Run("Works on good response", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(`[{"login": "User1"}]`)),
				}, nil
			}},
			RepoName: "someuser/somerepo",
		}

		users, err := a.ReviewRequests(1)
		require.Nil(t, err)
		require.Equal(t, []User{{"User1"}}, users)
	})

	t.Run("Fails when there is an error fetching the response", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return nil, fmt.Errorf("Some weird network error")
			}},
			RepoName: "someuser/somerepo",
		}

		_, err := a.ReviewRequests(1)
		require.EqualError(t, err, "Some weird network error")
	})
}
