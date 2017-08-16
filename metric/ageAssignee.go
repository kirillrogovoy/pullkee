package metric

import "github.com/kirillrogovoy/pullk/github"

// AgeAssignee contains the calculated data
type AgeAssignee struct {
	average averageByDev
}

// Description of the metric
func (a *AgeAssignee) Description() string {
	return "What is the average age of the PR when one is the assignee?"
}

// Calculate the average age of a PR in total and by developer
func (a *AgeAssignee) Calculate(pullRequests github.PullRequests) error {
	a.average.reset()

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		delta := pr.MergedAt.Sub(pr.CreatedAt).Hours() / 24
		for _, assignee := range pr.Assignees {
			a.average.add(delta, assignee.Login)
		}
	}

	return nil
}

// Converts the calculated data to a string
func (a *AgeAssignee) String() string {
	return a.average.string("days")
}
