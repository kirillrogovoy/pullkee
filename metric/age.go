package metric

import "github.com/kirillrogovoy/pullkee/github"

// Age contains the calculated data
type Age struct {
	average averageList
}

// Description of the metric
func (m *Age) Description() string {
	return "What is the average age of the PR for the particular author?"
}

// Calculate the average age of a PR in total and by developer
func (m *Age) Calculate(pullRequests []github.PullRequest) error {
	a := averageMap{}
	a.reset()

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		delta := pr.MergedAt.Sub(pr.CreatedAt).Hours() / 24
		a.add(delta, pr.User.Login)
	}

	m.average = a.toList()
	return nil
}

// Converts the calculated data to a string
func (m *Age) String() string {
	return m.average.string("days")
}
