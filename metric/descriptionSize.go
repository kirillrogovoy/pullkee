package metric

import "github.com/kirillrogovoy/pullkee/github"

// DescriptionSize contains the calculated data
type DescriptionSize struct {
	average averageList
}

// Description of the metric
func (m *DescriptionSize) Description() string {
	return "What is the average size of the Pull Request description?"
}

// Calculate the data
func (m *DescriptionSize) Calculate(pullRequests []github.PullRequest) error {
	a := averageMap{}
	a.reset()

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		a.add(float64(len(pr.Body)), pr.User.Login)
	}

	m.average = a.toList()
	return nil
}

// Converts the calculated data to a string
func (m *DescriptionSize) String() string {
	return m.average.string("chars")
}
