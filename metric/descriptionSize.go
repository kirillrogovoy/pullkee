package metric

import "github.com/kirillrogovoy/pullk/github"

// DescriptionSize contains the calculated data
type DescriptionSize struct {
	average averageByDev
}

// Description of the metric
func (a *DescriptionSize) Description() string {
	return "What is the average size of the Pull Request description?"
}

// Calculate the data
func (a *DescriptionSize) Calculate(pullRequests github.PullRequests) error {
	a.average.reset()

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		a.average.add(float64(len(pr.Body)), pr.User.Login)
	}

	return nil
}

// Converts the calculated data to a string
func (a *DescriptionSize) String() string {
	return a.average.string("chars")
}
