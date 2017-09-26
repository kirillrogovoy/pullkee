package metric

import (
	"time"

	"github.com/kirillrogovoy/pullk/github"
)

// DiffSizePerDay contains the calculated data
type DiffSizePerDay struct {
	average averageList
}

// Description of the metric
func (m *DiffSizePerDay) Description() string {
	return "How many bytes/day of diffs does one generate?"
}

// Calculate the average age of a PR in total and by developer
func (m *DiffSizePerDay) Calculate(pullRequests []github.PullRequest) error {
	a := averageMap{}
	a.reset()
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
		a.add(float64(size), pr.User.Login)
	}

	daysPassed := int(now.Sub(oldestCreated).Hours() / 24)
	a.setCount(daysPassed)

	m.average = a.toList()
	return nil
}

// Converts the calculated data to a string
func (m *DiffSizePerDay) String() string {
	return m.average.string("bytes/day")
}
