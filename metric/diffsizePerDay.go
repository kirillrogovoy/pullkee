package metric

import (
	"time"

	"github.com/kirillrogovoy/pullk/github"
)

// DiffSizePerDay contains the calculated data
type DiffSizePerDay struct {
	average averageByDev
}

// Description of the metric
func (a *DiffSizePerDay) Description() string {
	return "How many bytes/day of diffs does one generate?"
}

// Calculate the average age of a PR in total and by developer
func (a *DiffSizePerDay) Calculate(pullRequests github.PullRequests) error {
	a.average.reset()
	now := time.Now()
	oldestCreated := now

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		if pr.CreatedAt.Before(oldestCreated) {
			oldestCreated = pr.CreatedAt
		}

		size := *(pr.DiffSize)
		a.average.add(float64(size), pr.User.Login)
	}

	daysPassed := int(now.Sub(oldestCreated).Hours() / 24)
	a.average.setCount(daysPassed)

	return nil
}

// Converts the calculated data to a string
func (a *DiffSizePerDay) String() string {
	return a.average.string("bytes/day")
}
