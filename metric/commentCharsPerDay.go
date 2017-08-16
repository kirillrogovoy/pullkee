package metric

import (
	"time"

	"github.com/kirillrogovoy/pullk/github"
)

// CommentCharsPerDay contains the calculated data
type CommentCharsPerDay struct {
	average averageByDev
}

// Description of the metric
func (a *CommentCharsPerDay) Description() string {
	return "How many chars/day of comments does one generate?"
}

// Calculate the average age of a PR in total and by developer
func (a *CommentCharsPerDay) Calculate(pullRequests github.PullRequests) error {
	a.average.reset()
	now := time.Now()
	oldestCreated := now

	for _, pr := range pullRequests {
		if pr.CreatedAt.Before(oldestCreated) {
			oldestCreated = pr.CreatedAt
		}

		for _, comment := range pr.Comments() {
			size := len(comment.Body)
			a.average.add(float64(size), comment.User.Login)
		}
	}

	daysPassed := int(now.Sub(oldestCreated).Hours() / 24)
	a.average.setCount(daysPassed)

	return nil
}

// Converts the calculated data to a string
func (a *CommentCharsPerDay) String() string {
	return a.average.string("chars/day")
}
