// Package util is a high-level API for Github
package util

import (
	"fmt"
	"log"

	"github.com/kirillrogovoy/pullkee/cache"
	"github.com/kirillrogovoy/pullkee/github"
	"github.com/pkg/errors"
)

// Pulls fetches the list of pull requests directly from the API
func Pulls(a github.API, limit int) ([]github.PullRequest, error) {
	return a.ClosedPullRequests(limit)
}

// FillDetails calls .FillDetails for each PR in prs in parallel.
// It returns a channel which will never be closed, so the caller
// should expect len(prs) values from it
func FillDetails(
	a github.API,
	c cache.Cache,
	prs []github.PullRequest,
) chan error {
	ch := make(chan error, len(prs))

	for i, p := range prs {
		go (func(i int, p github.PullRequest) {
			var err error

			cacheKey := fmt.Sprintf("pr%d", p.Number)
			found, err := c.Get(cacheKey, &p)
			if err != nil {
				reportFsError(errors.Wrap(err, "getting cache"))
			}

			if !found {
				err = p.FillDetails(a)
			}

			if err != nil {
				ch <- err
			} else {
				// since p is a copy of i-th elem, we explicitly assign it to prs[i] to make the actual change
				prs[i] = p
				if err := c.Set(cacheKey, p); err != nil {
					reportFsError(errors.Wrap(err, "setting cache"))
				}
				ch <- nil
			}
		})(i, p)
	}

	return ch
}

func reportFsError(err error) {
	log.Printf("File system error occurred while accessing the cache: %s\n", err)
}
