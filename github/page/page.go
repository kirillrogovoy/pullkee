// Package page provides an utility to fetch paginated resources
package page

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/kirillrogovoy/pullkee/github/client"
)

// All fetches multiple pages given only a request for the first one and unmarshals them into `target`.
// JSON of each response must be an array and `target` must be a pointer to a slice of the according type.
func All(
	httpClient client.HTTPClient,
	firstPageRequest http.Request,
	target interface{},
	pageLimit int,
) error {
	first, err := httpClient.Do(&firstPageRequest)
	if err != nil {
		return err
	}

	rest, err := getRest(httpClient, *first, pageLimit-1)
	if err != nil {
		return err
	}

	all := append([]http.Response{*first}, rest...)

	return unmarshalResponses(all, target)
}

func getRest(
	httpClient client.HTTPClient,
	firstPageResponse http.Response,
	limit int,
) ([]http.Response, error) {
	responses := []http.Response{}

	cur := firstPageResponse

	for i := 0; ; {
		if i+1 > limit {
			return responses, nil
		}
		next, err := nextPage(httpClient, cur)
		if err != nil {
			return nil, err
		}
		if next == nil {
			return responses, nil
		}

		responses = append(responses, *next)
		cur = *next
		i++
	}
}

func nextPage(httpClient client.HTTPClient, prevPageResponse http.Response) (*http.Response, error) {
	link := prevPageResponse.Header.Get("Link")
	if link == "" {
		return nil, nil
	}

	nextURL := extractLinkURL(link, "next")
	if nextURL == "" {
		return nil, nil
	}

	req, _ := http.NewRequest("GET", nextURL, nil)
	return httpClient.Do(req)
}

func extractLinkURL(header string, rel string) string {
	for _, link := range strings.Split(header, ",") {
		if strings.Contains(link, fmt.Sprintf("rel=\"%s\"", rel)) {
			matches := regexp.MustCompile(`<(.*)>`).FindStringSubmatch(link)
			if len(matches) > 1 {
				return matches[1]
			}
		}
	}

	return ""
}

func unmarshalResponses(responses []http.Response, target interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(runtime.Error); ok {
				panic(e)
			}
			err = e.(error)
		}
	}()

	targetRefl := reflect.ValueOf(target).Elem()
	tmp, err := createSliceOfSameType(target)
	if err != nil {
		return err
	}

	for _, res := range responses {
		if err := unmarshalResponse(res, tmp); err != nil {
			return err
		}

		targetRefl = reflect.AppendSlice(targetRefl, reflect.ValueOf(tmp).Elem())
	}

	reflect.ValueOf(target).Elem().Set(targetRefl)
	return nil
}

func unmarshalResponse(res http.Response, target interface{}) error {
	if res.Body == nil {
		url := "<unknown>"
		if res.Request != nil && res.Request.URL != nil {
			url = res.Request.URL.String()
		}
		return fmt.Errorf("Expected res.Body not to be nil. URL: %s", url)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

func createSliceOfSameType(original interface{}) (interface{}, error) {
	targetRefl := reflect.ValueOf(original)
	resultSlice := reflect.New(targetRefl.Elem().Type())

	if resultSlice.Elem().Kind() != reflect.Slice {
		return nil, fmt.Errorf("Expected target to be a pointer to a slice")
	}

	return resultSlice.Interface(), nil
}
