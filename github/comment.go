package github

import (
	"fmt"
	"math"
	"net/http"

	"github.com/kirillrogovoy/pullk/github/page"
)

// Comment is a representation of a Github issue comment
type Comment struct {
	User User   `json:"user"`
	Body string `json:"body"`
}

// Comments fetches all the comments of a Pull Request given its `number`
func (a APIv3) Comments(number int) ([]Comment, error) {
	allComments := []Comment{}

	types := []string{"pulls", "issues"}
	for _, commentType := range types {
		comments := []Comment{}
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/%d/comments?per_page=100",
			a.RepoName,
			commentType,
			number,
		)
		req, _ := http.NewRequest("GET", url, nil)

		pageLimit := int(math.Inf(1))
		if err := page.All(a.HTTPClient, *req, &comments, pageLimit); err != nil {
			return nil, err
		}
		allComments = append(allComments, comments...)
	}

	return allComments, nil
}
