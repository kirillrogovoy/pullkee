package github

import (
	"encoding/json"
	"fmt"
)

// ReviewRequests fetches the diff of the particular pull request
func (a *API) ReviewRequests(repo string, number int) ([]User, error) {
	req := request(fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/requested_reviewers", repo, number))
	req.Header.Add("Accept", "application/vnd.github.black-cat-preview+json")
	res, err := a.send(req)
	if err != nil {
		return nil, err
	}

	body := HTTPBody(res)
	users := []User{}
	err = json.Unmarshal(body, &users)
	return users, err
}
