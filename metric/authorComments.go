package metric

import "github.com/kirillrogovoy/pullk/github"

// AuthorComments contains the calculated data
type AuthorComments struct {
	average averageByDev
}

// Description of the metric
func (a *AuthorComments) Description() string {
	return "How many comments does one have on his own PR?"
}

// Calculate the data
func (a *AuthorComments) Calculate(pullRequests []github.PullRequest) error {
	a.average.reset()

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		count := 0
		for _, comment := range *pr.Comments {
			if comment.User.Login != pr.User.Login {
				count++
			}
		}
		a.average.add(float64(count), pr.User.Login)
	}

	return nil
}

// Converts the calculated data to a string
func (a *AuthorComments) String() string {
	return a.average.string("comments")
}
