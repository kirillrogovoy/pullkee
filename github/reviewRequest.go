package github

import "fmt"

// ReviewRequestsResponse is a representation of the /requested_reviewers API response
type ReviewRequestsResponse struct {
	Users []User `json:"users"`
}

// ReviewRequests fetches a list of users which were requested to do a review
func (a APIv3) ReviewRequests(number int) ([]User, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/requested_reviewers", a.RepoName, number)
	response := ReviewRequestsResponse{}
	err := a.Get(url, &response)
	if err != nil {
		return nil, err
	}

	return response.Users, nil
}
