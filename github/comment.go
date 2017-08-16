package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kirillrogovoy/pullk/github/page"
)

// Comment is a representation of a regular or review comment
type Comment struct {
	User User   `json:"user"`
	Body string `json:"body"`
}

// Comments fetches the diff of the particular pull request
func (a *API) Comments(repo string, number int, commentType string) ([]Comment, error) {
	fetcher := commentPageFetcher{a, repo, number, commentType}
	pages, err := page.FetchAll(fetcher)
	if err != nil {
		return nil, err
	}
	result := []Comment{}

	for _, page := range pages {
		var comments []Comment
		if err := json.Unmarshal(HTTPBody(page.Response), &comments); err != nil {
			return nil, err
		}
		result = append(result, comments...)
	}

	return result, nil
}

type commentPageFetcher struct {
	a           *API
	repo        string
	number      int
	commentType string
}

func (c commentPageFetcher) GetPage(page int) (*http.Response, error) {
	var urlPart string

	switch c.commentType {
	case "review":
		urlPart = "pulls"
	case "issue":
		urlPart = "issues"
	default:
		return nil, fmt.Errorf("commentType should be either 'review' or 'issue'")
	}

	return c.a.send(
		request(
			fmt.Sprintf(
				"https://api.github.com/repos/%s/%s/%d/comments?per_page=100&page=%d",
				c.repo,
				urlPart,
				c.number,
				page,
			),
		),
	)
}
