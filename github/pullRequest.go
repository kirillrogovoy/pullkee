package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kirillrogovoy/pullk/github/page"
)

// PullRequest is a representation of the Pull Request the Github API returns
type PullRequest struct {
	Number         int       `json:"number"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	MergedAt       time.Time `json:"merged_at,omitempty"`
	User           User      `json:"user"`
	State          string    `json:"state"`
	Assignees      []User    `json:"assignees"`
	DiffURL        string    `json:"diff_url"`
	DiffSize       *int
	ReviewRequests *[]User
	ReviewComments *[]Comment
	IssueComments  *[]Comment
}

// IsMerged tells if PullRequest was really merged, not closed
func (p PullRequest) IsMerged() bool {
	return p.State == "closed" && !p.MergedAt.IsZero()
}

// Comments returns both ReviewComments and IssueComments
func (p PullRequest) Comments() []Comment {
	return append(*p.ReviewComments, *p.IssueComments...)
}

// PullRequests is a container for a list of PullRequest instances
type PullRequests []*PullRequest

// ParallelLoop executes a loop in parallel over the list
func (p *PullRequests) ParallelLoop(f func(pr *PullRequest) error) error {
	c := make(chan error)

	for _, pr := range *p {
		go (func(pr *PullRequest) {
			c <- f(pr)
		})(pr)
	}

	for _ = range *p {
		err := <-c
		if err != nil {
			return err
		}
	}

	return nil
}

// PullsClosed fetches a list of closed Pull Requests
func (a *API) PullsClosed(repo string) (PullRequests, error) {
	fetcher := pullsPageFetcher{a, repo}
	pages, err := page.FetchAll(fetcher)
	if err != nil {
		return nil, err
	}

	var result PullRequests
	for _, page := range pages {
		var pulls PullRequests
		if err := json.Unmarshal(HTTPBody(page.Response), &pulls); err != nil {
			return nil, err
		}
		result = append(result, pulls...)
	}

	return result, nil
}

type pullsPageFetcher struct {
	a    *API
	repo string
}

func (p pullsPageFetcher) GetPage(page int) (*http.Response, error) {
	return p.a.send(
		request(
			fmt.Sprintf(
				"https://api.github.com/repos/%s/pulls?state=closed&per_page=100&page=%d",
				p.repo,
				page,
			),
		),
	)
}
