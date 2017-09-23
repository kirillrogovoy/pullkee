package metric

import "github.com/kirillrogovoy/pullk/github"

// Age contains the calculated data
type Age struct {
	average averageByDev
}

// Description of the metric
func (a *Age) Description() string {
	return "What is the average age of the PR for the particular author?"
}

// Calculate the average age of a PR in total and by developer
func (a *Age) Calculate(pullRequests []github.PullRequest) error {
	a.average.reset()

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		delta := pr.MergedAt.Sub(pr.CreatedAt).Hours() / 24
		a.average.add(delta, pr.User.Login)
	}

	return nil
}

// Converts the calculated data to a string
func (a *Age) String() string {
	return a.average.string("days")
}
