package metric

import "github.com/kirillrogovoy/pullk/github"

// DiffSize contains the calculated data
type DiffSize struct {
	average averageByDev
}

// Description of the metric
func (a *DiffSize) Description() string {
	return "What is the average diff size of the PR?"
}

// Calculate the average age of a PR in total and by developer
func (a *DiffSize) Calculate(pullRequests github.PullRequests) error {
	a.average.reset()

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		size := *(pr.DiffSize)
		a.average.add(float64(size), pr.User.Login)
	}

	return nil
}

// Converts the calculated data to a string
func (a *DiffSize) String() string {
	return a.average.string("bytes")
}
