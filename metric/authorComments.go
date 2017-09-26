package metric

import "github.com/kirillrogovoy/pullk/github"

// AuthorComments contains the calculated data
type AuthorComments struct {
	average averageList
}

// Description of the metric
func (m *AuthorComments) Description() string {
	return "How many comments does one have on his own PR?"
}

// Calculate the data
func (m *AuthorComments) Calculate(pullRequests []github.PullRequest) error {
	a := averageMap{}
	a.reset()

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
		a.add(float64(count), pr.User.Login)
	}

	m.average = a.toList()
	return nil
}

// Converts the calculated data to a string
func (m *AuthorComments) String() string {
	return m.average.string("comments")
}
