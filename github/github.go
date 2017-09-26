// Package github provides low-level tools to retrieve information from the Github API
package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kirillrogovoy/pullk/github/client"
)

// API is an interface for a collection of methods to retrieve information from Github API
type API interface {
	Get(url string, target interface{}) error
	Repository() (*Repository, error)
	ClosedPullRequests(limit int) ([]PullRequest, error)
	DiffSize(number int) (int, error)
	Comments(number int) ([]Comment, error)
	ReviewRequests(number int) ([]User, error)
}

// APIv3 is an implementation of API which works with Github REST API (v3)
type APIv3 struct {
	HTTPClient client.HTTPClient
	RepoName   string
}

// User is a representation of a Github user (e.g. an author of a Pull Request)
type User struct {
	Login string `json:"login"`
}

// Get makes an HTTP request, checks the response, reads the body and unmarshals it to the `target`
func (a APIv3) Get(url string, target interface{}) error {
	// According to the tests of http.Request an error might only occur on an invalid method which is not the case
	req, _ := http.NewRequest("GET", url, nil)

	res, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if res.Body == nil {
		return fmt.Errorf("Expected res.Body not to be nil. URL: %s", req.URL)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}
