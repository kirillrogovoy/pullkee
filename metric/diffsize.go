package metric

import "github.com/kirillrogovoy/pullk/github"

// DiffSize contains the calculated data
type DiffSize struct {
	average averageList
}

// Description of the metric
func (m *DiffSize) Description() string {
	return "What is the average diff size of the PR?"
}

// Calculate the average age of a PR in total and by developer
func (m *DiffSize) Calculate(pullRequests []github.PullRequest) error {
	a := averageMap{}
	a.reset()

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		size := *(pr.DiffSize)
		a.add(float64(size), pr.User.Login)
	}

	m.average = a.toList()
	return nil
}

// Converts the calculated data to a string
func (m *DiffSize) String() string {
	return m.average.string("bytes")
}
