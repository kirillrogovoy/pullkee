package page

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	client "github.com/kirillrogovoy/pullk/github/client"
)

// Rest does stuff
func Rest(httpClient client.HTTPClient, firstPageResponse http.Response) ([]http.Response, error) {
	responses := []http.Response{}

	cur := firstPageResponse
	for {
		next, err := nextPage(httpClient, cur)
		if err != nil {
			return nil, err
		}
		if next == nil {
			return responses, nil
		}

		responses = append(responses, *next)
		cur = *next
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
