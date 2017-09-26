package page_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	. "github.com/kirillrogovoy/pullkee/github/page"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	t.Run("Works on successful responses", func(t *testing.T) {
		linkPage := func(n int) string {
			return fmt.Sprintf("<https://api.github.com/user/repos?page=%d&per_page=100>; rel=\"next\"", n)
		}

		timesCalled := 0
		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			json := fmt.Sprintf("[{\"keyX\": \"val%d\"}]", timesCalled)
			res := &http.Response{
				Header: http.Header{},
				Body:   ioutil.NopCloser(strings.NewReader(json)),
			}
			switch timesCalled {
			case 0:
				res.Header.Add("Link", linkPage(3))
			case 1:
				res.Header.Add("Link", linkPage(4))
			case 2:
				// nothing
			default:
				panic("Should not be called")
			}
			return res, nil
		}

		actual := &[]SomeStruct{}
		err := All(httpClientMock{response}, http.Request{}, actual, 99)

		expected := &[]SomeStruct{
			{"val0"},
			{"val1"},
			{"val2"},
		}

		require.Nil(t, err)
		require.Equal(t, 3, timesCalled)
		require.Equal(t, expected, actual)
	})

	t.Run("Limit works", func(t *testing.T) {
		linkPage := func(n int) string {
			return fmt.Sprintf("<https://api.github.com/user/repos?page=%d&per_page=100>; rel=\"next\"", n)
		}

		timesCalled := 0
		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			json := fmt.Sprintf("[{\"keyX\": \"val%d\"}]", timesCalled)
			res := &http.Response{
				Header: http.Header{},
				Body:   ioutil.NopCloser(strings.NewReader(json)),
			}
			switch timesCalled {
			case 0:
				res.Header.Add("Link", linkPage(3))
			case 1:
				res.Header.Add("Link", linkPage(4))
			case 2:
				// nothing
			default:
				panic("Should not be called")
			}
			return res, nil
		}

		actual := &[]SomeStruct{}
		err := All(httpClientMock{response}, http.Request{}, actual, 2)

		expected := &[]SomeStruct{
			{"val0"},
			{"val1"},
		}

		require.Nil(t, err)
		require.Equal(t, 2, timesCalled)
		require.Equal(t, expected, actual)
	})

	t.Run("Fails when couldn't fetch the first page", func(t *testing.T) {
		result := &[]SomeStruct{}

		err := All(httpClientMock{func() (*http.Response, error) {
			return nil, fmt.Errorf("Some weird network error")
		}}, *dummyRequest(), result, 99)

		require.Equal(t, []SomeStruct{}, *result)
		require.Contains(t, err.Error(), "Some weird network error")
	})

	t.Run("Fails when response body is absent", func(t *testing.T) {
		result := &[]SomeStruct{}
		req := dummyRequest()

		err := All(httpClientMock{func() (*http.Response, error) {
			return &http.Response{
				Request: req,
			}, nil
		}}, *req, result, 99)

		require.Equal(t, []SomeStruct{}, *result)
		require.Contains(t, err.Error(), "Expected res.Body not to be nil")
	})

	t.Run("Fails when there was an error reading response body", func(t *testing.T) {
		result := &[]SomeStruct{}
		req := dummyRequest()

		err := All(httpClientMock{func() (*http.Response, error) {
			return &http.Response{
				Request: req,
				Body:    ioutil.NopCloser(errorReader{}),
			}, nil
		}}, *req, result, 99)

		require.Equal(t, []SomeStruct{}, *result)
		require.Equal(t, "Some weird reader error", err.Error())
	})

	t.Run("Fails when couldn't fetch the rest of pages", func(t *testing.T) {
		goodLink := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next"`

		timesCalled := 0
		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			switch timesCalled {
			case 0:
				return &http.Response{
					Header: http.Header{
						"Link": []string{goodLink},
					},
					Body: ioutil.NopCloser(strings.NewReader(`[{"keyX": "val1"}]`)),
				}, nil
			case 1:
				return nil, fmt.Errorf("Some weird network error")
			default:
				panic("Should not be called")
			}
		}

		result := &[]SomeStruct{}
		err := All(httpClientMock{response}, *dummyRequest(), result, 99)

		require.Equal(t, []SomeStruct{}, *result)
		require.Contains(t, err.Error(), "Some weird network error")
	})

	t.Run("Stops fetching when there is no Link header in the response", func(t *testing.T) {
		result := &[]SomeStruct{}
		err := All(httpClientMock{func() (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`[{"keyX": "val1"}]`)),
			}, nil
		}}, *dummyRequest(), result, 99)

		require.Nil(t, err)
		require.Equal(t, []SomeStruct{{KeyX: "val1"}}, *result)
	})

	t.Run("Stops fetching when couldn't parse the Link header", func(t *testing.T) {
		result := &[]SomeStruct{}
		err := All(httpClientMock{func() (*http.Response, error) {
			return &http.Response{
				Header: http.Header{
					"Link": []string{"Total rubbish"},
				},
				Body: ioutil.NopCloser(strings.NewReader(`[{"keyX": "val1"}]`)),
			}, nil
		}}, *dummyRequest(), result, 99)

		require.Nil(t, err)
		require.Equal(t, []SomeStruct{{KeyX: "val1"}}, *result)
	})

	t.Run("Fails when target is not a pointer", func(t *testing.T) {
		result := []SomeStruct{}
		err := All(httpClientMock{successfulResponse}, *dummyRequest(), result, 99)

		require.NotNil(t, err)
	})

	t.Run("Fails when target is not a pointer to a slice", func(t *testing.T) {
		result := &SomeStruct{}
		err := All(httpClientMock{successfulResponse}, *dummyRequest(), result, 99)

		require.NotNil(t, err)
	})

	t.Run("Fails when a subsequent request gets wrong kind of JSON", func(t *testing.T) {
		goodLink := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next"`

		timesCalled := 0
		response := func() (*http.Response, error) {
			defer func() { timesCalled++ }()
			switch timesCalled {
			case 0:
				return &http.Response{
					Header: http.Header{
						"Link": []string{goodLink},
					},
					Body: ioutil.NopCloser(strings.NewReader(`[{"keyX": "val1"}]`)),
				}, nil
			case 1:
				return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(`"JSON, but not an array"`)),
				}, nil
			default:
				panic("Should not be called")
			}
		}

		result := &[]SomeStruct{}
		err := All(httpClientMock{response}, *dummyRequest(), result, 99)

		require.Equal(t, []SomeStruct{}, *result)
		require.Contains(t, err.Error(), "cannot unmarshal string into Go value of type []page_test.SomeStruct")
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

type SomeStruct struct {
	KeyX string `json:"keyX"`
}

type errorReader struct{}

func (e errorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("Some weird reader error")
}
