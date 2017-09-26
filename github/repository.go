package github

import "fmt"

// Repository is a representation of a Github repository which is accessible via API
type Repository struct {
	FullName string `json:"full_name"`
}

// Repository fetches the remote repository data
func (a APIv3) Repository() (*Repository, error) {
	repo := &Repository{}

	if err := a.Get(fmt.Sprintf("https://api.github.com/repos/%s", a.RepoName), repo); err != nil {
		return nil, err
	}

	return repo, nil
}
