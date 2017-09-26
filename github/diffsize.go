package github

import (
	"fmt"
	"net/http"
	"strconv"
)

// DiffSize fetches the size of the diff of the particular Pull Request given `number`
func (a APIv3) DiffSize(number int) (int, error) {
	req, _ := http.NewRequest("HEAD", fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d", a.RepoName, number), nil)
	req.Header.Add("Accept", "application/vnd.github.diff")
	res, err := a.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}

	length := res.Header.Get("Content-Length")
	if length == "" {
		return 0, fmt.Errorf("Expected Content-Length in response")
	}

	lengthInt, _ := strconv.Atoi(length)
	return lengthInt, nil
}
