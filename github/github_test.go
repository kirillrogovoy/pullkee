package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Run("Works on correct response", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(`{"full_name": "someuser/somerepo"}`)),
				}, nil
			}},
			RepoName: "someuser/somerepo",
		}

		repo := &Repository{}
		err := a.Get("/some-url", repo)
		require.Nil(t, err)
		require.Equal(t, "someuser/somerepo", repo.FullName)
	})

	t.Run("Fails when the request failed", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return nil, fmt.Errorf("Some weird network error")
			}},
			RepoName: "someuser/somerepo",
		}

		err := a.Get("/some-url", nil)
		require.Equal(t, "Some weird network error", err.Error())
	})

	t.Run("Fails when the body is nil", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       nil,
				}, nil
			}},
			RepoName: "someuser/somerepo",
		}

		err := a.Get("/some-url", nil)
		require.NotNil(t, err)
	})

	t.Run("Fails when couldn't read the body", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(errorReader{}),
				}, nil
			}},
			RepoName: "someuser/somerepo",
		}

		err := a.Get("/some-url", nil)
		require.Equal(t, "Some weird reader error", err.Error())
	})

	t.Run("Fails when the body is not correct JSON", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(`I'm not JSON`)),
				}, nil
			}},
			RepoName: "someuser/somerepo",
		}

		err := a.Get("/some-url", nil)
		require.Contains(t, err.Error(), "invalid character 'I'")
	})
}

type httpClientMock struct {
	response func() (*http.Response, error)
}

func (h httpClientMock) Do(request *http.Request) (*http.Response, error) {
	return h.response()
}

type errorReader struct{}

func (e errorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("Some weird reader error")
}
