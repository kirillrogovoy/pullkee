package page

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRest(t *testing.T) {
	t.Run("Returns an error when a subsequent request returns fails", func(t *testing.T) {
		goodLink := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next"`
		input, _ := successfulResponse()
		input.Header.Add("Link", goodLink)

		res, err := Rest(httpClientMock{func() (*http.Response, error) {
			return nil, fmt.Errorf("Some weird network error")
		}}, *input)

		require.Nil(t, res)
		require.Contains(t, err.Error(), "Some weird network error")
	})

	t.Run("Returns an empty array when can't parse the first page's Link header", func(t *testing.T) {
		input, _ := successfulResponse()
		rest, err := Rest(httpClientMock{successfulResponse}, *input)

		require.Nil(t, err)
		require.Empty(t, rest)
	})

	t.Run("Returns the rest correctly", func(t *testing.T) {
		linkPage := func(n int) string {
			return fmt.Sprintf("<https://api.github.com/user/repos?page=%d&per_page=100>; rel=\"next\"", n)
		}
		input, _ := successfulResponse()
		input.Header.Add("Link", linkPage(2))

		timesCalled := 0
		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			res := &http.Response{
				Header: http.Header{},
			}
			switch timesCalled {
			case 0:
				res.Header.Add("Link", linkPage(3))
			case 1:
				res.Header.Add("Link", linkPage(4))
			case 2:
			default:
				panic("Should not be called")
			}
			return res, nil
		}
		rest, err := Rest(httpClientMock{response}, *input)

		require.Nil(t, err)
		require.Equal(t, 3, timesCalled)
		require.Len(t, rest, 3)
	})
}

func TestNextURL(t *testing.T) {
	goodLink := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next", <https://api.github.com/user/repos?page=50&per_page=100>; rel="last"`

	t.Run("Returns nil if the header doesn't exist", func(t *testing.T) {
		input, _ := successfulResponse()
		res, err := nextPage(httpClientMock{successfulResponse}, *input)

		require.Nil(t, res)
		require.Nil(t, err)
	})

	t.Run("Returns nil if the header is of a wrong format", func(t *testing.T) {
		input, _ := successfulResponse()
		input.Header.Add("Link", "Total rubbish")
		res, err := nextPage(httpClientMock{successfulResponse}, *input)

		require.Nil(t, res)
		require.Nil(t, err)
	})

	t.Run("Return err if the request has failed", func(t *testing.T) {
		input, _ := successfulResponse()
		input.Header.Add("Link", goodLink)

		res, err := nextPage(httpClientMock{func() (*http.Response, error) {
			return nil, fmt.Errorf("Some weird network error")
		}}, *input)

		require.Nil(t, res)
		require.Contains(t, err.Error(), "Some weird network error")
	})

	t.Run("Return response if it's successful", func(t *testing.T) {
		input, _ := successfulResponse()
		input.Header.Add("Link", goodLink)

		res, err := nextPage(httpClientMock{successfulResponse}, *input)

		require.Nil(t, err)
		require.Equal(t, "yes", res.Header.Get("Success"))
	})
}

func TestLinkURL(t *testing.T) {
	t.Run("Found when the header has multiple rels", func(t *testing.T) {
		header := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next", <https://api.github.com/user/repos?page=50&per_page=100>; rel="last"`
		expected := "https://api.github.com/user/repos?page=3&per_page=100"
		actual := extractLinkURL(header, "next")

		require.Equal(t, expected, actual)
	})

	t.Run("Found the header only has one rel", func(t *testing.T) {
		header := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next"`
		expected := "https://api.github.com/user/repos?page=3&per_page=100"
		actual := extractLinkURL(header, "next")

		require.Equal(t, expected, actual)
	})
}

type httpClientMock struct {
	response func() (*http.Response, error)
}

func (h httpClientMock) Do(request *http.Request) (*http.Response, error) {
	return h.response()
}

func dummyRequest() *http.Request {
	req, _ := http.NewRequest("GET", "http://example.com/url1", nil)
	return req
}

func successfulResponse() (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Success": []string{"yes"},
		},
	}, nil
}
