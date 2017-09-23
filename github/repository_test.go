package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	t.Run("Works on good response", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(`{"full_name": "someuser/somerepo"}`)),
				}, nil
			}},
			RepoName: "someuser/somerepo",
		}

		repo, err := a.Repository()
		require.Nil(t, err)
		require.Equal(t, "someuser/somerepo", repo.FullName)
	})

	t.Run("Fails when there is an error fetching the response", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return nil, fmt.Errorf("Dogs have chewed the wires")
			}},
			RepoName: "someuser/somerepo",
		}

		repo, err := a.Repository()
		require.EqualError(t, err, "Dogs have chewed the wires")
		require.Nil(t, repo)
	})
}
