package github

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/kirillrogovoy/pullkee/github/page"
	"github.com/pkg/errors"
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
	Comments       *[]Comment
}

// IsMerged tells if PullRequest was really merged, not just closed
func (p PullRequest) IsMerged() bool {
	return p.State == "closed" && !p.MergedAt.IsZero()
}

// FillDetails makes additional requests to fill details about the Pull Request (such as diff size)
func (p *PullRequest) FillDetails(a API) error {
	if p.DiffSize == nil {
		size, err := a.DiffSize(p.Number)
		if err != nil {
			return errors.Wrap(err, "diff size")
		}
		p.DiffSize = &size
	}

	if p.ReviewRequests == nil {
		users, err := a.ReviewRequests(p.Number)
		if err != nil {
			return errors.Wrap(err, "review requests")
		}
		p.ReviewRequests = &users
	}

	if p.Comments == nil {
		comments, err := a.Comments(p.Number)
		if err != nil {
			return errors.Wrap(err, "comments")
		}
		p.Comments = &comments
	}

	return nil
}

// ClosedPullRequests fetches a list of closed Pull Requests with a `limit`
func (a APIv3) ClosedPullRequests(limit int) ([]PullRequest, error) {
	perPage := 100
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/pulls?state=closed&per_page=%d&page=1",
		a.RepoName,
		perPage,
	)
	req, _ := http.NewRequest("GET", url, nil)

	prs := []PullRequest{}
	pageLimit := int(math.Ceil(float64(limit) / float64(perPage)))
	if limit <= 0 {
		pageLimit = int(math.Inf(1))
	}
	if err := page.All(a.HTTPClient, *req, &prs, pageLimit); err != nil {
		return nil, err
	}

	len := len(prs)
	if limit > len || limit == 0 {
		limit = len
	}
	return prs[:limit], nil
}
