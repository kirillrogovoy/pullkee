package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffsize(t *testing.T) {
	t.Run("Works on good response", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Header: http.Header{
						"Content-Length": []string{"42"},
					},
				}, nil
			}},
			RepoName: "someuser/somerepo",
		}

		size, err := a.DiffSize(1)
		require.Nil(t, err)
		require.Equal(t, 42, size)
	})

	t.Run("Fails when Content-Length header is missing", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Header:     http.Header{},
				}, nil
			}},
			RepoName: "someuser/somerepo",
		}

		_, err := a.DiffSize(1)
		require.EqualError(t, err, "Expected Content-Length in response")
	})

	t.Run("Fails when there is an error fetching the response", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return nil, fmt.Errorf("Some weird network error")
			}},
			RepoName: "someuser/somerepo",
		}

		_, err := a.DiffSize(1)
		require.EqualError(t, err, "Some weird network error")
	})
}
