package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/kirillrogovoy/pullk/github"
)

// Repository is a main point of retrieving any information which is aware of all the caches
type Repository struct {
	API      *github.API
	Settings Settings
	pulls    *github.PullRequests
}

// Settings is a struct to hold all the specifics regarding the fetching
type Settings struct {
	Repo        string
	Limit       int
	ResetCaches bool
}

// Pulls fetches all Pull Requests going through the file system cache
func (r *Repository) Pulls() (github.PullRequests, error) {
	defer r.enableCache()

	if err := r.fillPulls(); err != nil {
		return nil, err
	}

	r.applyLimit()

	if err := r.attachDetails(); err != nil {
		// silently try to write at least what we have to the local cache
		r.writeCache()
		return nil, err
	}

	if err := r.writeCache(); err != nil {
		return nil, err
	}

	return *r.pulls, nil
}

func (r *Repository) fillPulls() error {
	if !r.Settings.ResetCaches && r.pulls == nil {
		if err := r.readCache(); err != nil {
			return err
		}
	}

	if r.pulls == nil || len(*r.pulls) < r.Settings.Limit {
		pulls, err := r.API.PullsClosed(r.Settings.Repo)
		if err != nil {
			return err
		}
		r.pulls = &pulls
	}

	return nil
}

func (r *Repository) applyLimit() {
	pulls := *r.pulls
	limit := r.Settings.Limit
	len := len(pulls)

	if limit > len || limit == 0 {
		limit = len
	}

	pulls = pulls[:limit]
	r.pulls = &pulls
}

func (r *Repository) attachDetails() error {
	return r.pulls.ParallelLoop(func(pr *github.PullRequest) error {
		if pr.DiffSize == nil {
			size, err := r.API.DiffSize(r.Settings.Repo, pr.Number)
			if err != nil {
				return err
			}
			pr.DiffSize = &size
		}

		if pr.ReviewRequests == nil {
			users, err := r.API.ReviewRequests(r.Settings.Repo, pr.Number)
			if err != nil {
				return err
			}
			pr.ReviewRequests = &users
		}

		if pr.ReviewComments == nil {
			comments, err := r.API.Comments(r.Settings.Repo, pr.Number, "review")
			if err != nil {
				return err
			}
			pr.ReviewComments = &comments
		}

		if pr.IssueComments == nil {
			comments, err := r.API.Comments(r.Settings.Repo, pr.Number, "issue")
			if err != nil {
				return err
			}
			pr.IssueComments = &comments
		}

		return nil
	})
}

func (r *Repository) readCache() error {
	data, err := ioutil.ReadFile(r.cacheFilepath())
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			reportFsError(err)
		}
		return nil
	}

	r.pulls = &github.PullRequests{}
	return json.Unmarshal(data, r.pulls)
}

func (r *Repository) writeCache() error {
	filepath := r.cacheFilepath()
	dir := path.Dir(filepath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	jsoned, err := json.Marshal(r.pulls)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, jsoned, os.ModePerm)
}

func (r *Repository) cacheFilepath() string {
	return path.Join(os.TempDir(), "pullk", r.Settings.Repo, "cache.json")
}

func (r *Repository) enableCache() {
	r.Settings.ResetCaches = false
}

func reportFsError(err error) {
	log.Printf("File system error occurred while accessing the cache: %s\n", err)
}
