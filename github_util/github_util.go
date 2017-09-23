package repository

import (
	"fmt"
	"log"

	"github.com/kirillrogovoy/pullk/cache"
	"github.com/kirillrogovoy/pullk/github"
)

// Pulls fetches the list of pull requests directly from the API
func Pulls(a github.API, limit int) ([]github.PullRequest, error) {
	return a.ClosedPullRequests(limit)
}

// FillDetails calls .FillDetails for each PR in prs in parallel
func FillDetails(a github.API, c cache.Cache, prs []github.PullRequest) error {
	ch := make(chan error)

	for i, p := range prs {
		go (func(i int, p github.PullRequest) {
			var err error

			cacheKey := fmt.Sprintf("pr%d", p.Number)
			found, err := c.Get(cacheKey, &p)
			if err != nil {
				reportFsError(err)
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
					reportFsError(err)
				}
				ch <- nil
			}
		})(i, p)
	}

	var err error

	for range prs {
		err = <-ch
	}

	return err
}

func reportFsError(err error) {
	log.Printf("File system error occurred while accessing the cache: %s\n", err)
}
