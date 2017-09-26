package github

import "fmt"

// ReviewRequests fetches a list of users which were requested to do a review
func (a APIv3) ReviewRequests(number int) ([]User, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/requested_reviewers", a.RepoName, number)
	users := []User{}
	err := a.Get(url, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
