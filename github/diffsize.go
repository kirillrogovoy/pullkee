package github

import "fmt"

// DiffSize fetches the diff of the particular pull request
func (a *API) DiffSize(repo string, number int) (int, error) {
	req := request(fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d", repo, number))
	req.Header.Add("Accept", "application/vnd.github.diff")
	res, err := a.send(req)
	if err != nil {
		return 0, err
	}

	return len(HTTPBody(res)), nil
}
