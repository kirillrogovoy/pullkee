package metric

import (
	"time"

	"github.com/kirillrogovoy/pullkee/github"
)

// CommentCharsPerDay contains the calculated data
type CommentCharsPerDay struct {
	average averageList
}

// Description of the metric
func (m *CommentCharsPerDay) Description() string {
	return "How many chars/day of comments does one generate?"
}

// Calculate the average age of a PR in total and by developer
func (m *CommentCharsPerDay) Calculate(pullRequests []github.PullRequest) error {
	a := averageMap{}
	a.reset()
	now := time.Now()
	oldestCreated := now

	for _, pr := range pullRequests {
		if pr.CreatedAt.Before(oldestCreated) {
			oldestCreated = pr.CreatedAt
		}

		for _, comment := range *pr.Comments {
			size := len(comment.Body)
			a.add(float64(size), comment.User.Login)
		}
	}

	daysPassed := int(now.Sub(oldestCreated).Hours() / 24)
	a.setCount(daysPassed)

	m.average = a.toList()
	return nil
}

// Converts the calculated data to a string
func (m *CommentCharsPerDay) String() string {
	return m.average.string("chars/day")
}
