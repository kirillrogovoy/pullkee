package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComments(t *testing.T) {
	t.Run("Works on good response", func(t *testing.T) {
		goodLink := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next"`

		timesCalled := 0
		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			switch timesCalled {
			case 0:
				json := `[{"user": {"login": "User1"}, "body": "Body1"}]`
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(json)),
					Header: http.Header{
						"Link": []string{goodLink},
					},
				}, nil
			case 1:
				json := `[{"user": {"login": "User2"}, "body": "Body2"}]`
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(json)),
				}, nil
			case 2:
				json := `[{"user": {"login": "User3"}, "body": "Body3"}]`
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(json)),
					Header: http.Header{
						"Link": []string{goodLink},
					},
				}, nil
			case 3:
				json := `[{"user": {"login": "User4"}, "body": "Body4"}]`
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(json)),
				}, nil
			default:
				panic("Should not be called")
			}
		}
		a := APIv3{
			HTTPClient: httpClientMock{response},
			RepoName:   "someuser/somerepo",
		}

		expected := []Comment{
			{User{"User1"}, "Body1"},
			{User{"User2"}, "Body2"},
			{User{"User3"}, "Body3"},
			{User{"User4"}, "Body4"},
		}

		comments, err := a.Comments(1)
		require.Nil(t, err)
		require.Equal(t, expected, comments)
	})

	t.Run("Fails when there is an error fetching the response", func(t *testing.T) {
		a := APIv3{
			HTTPClient: httpClientMock{func() (*http.Response, error) {
				return nil, fmt.Errorf("Dogs have chewed the wires")
			}},
			RepoName: "someuser/somerepo",
		}

		comments, err := a.Comments(1)
		require.EqualError(t, err, "Dogs have chewed the wires")
		require.Nil(t, comments)
	})
}
