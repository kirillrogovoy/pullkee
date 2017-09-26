package metric

import "github.com/kirillrogovoy/pullk/github"

// AgeAssignee contains the calculated data
type AgeAssignee struct {
	average averageList
}

// Description of the metric
func (m *AgeAssignee) Description() string {
	return "What is the average age of the PR when one is the assignee?"
}

// Calculate the average age of a PR in total and by developer
func (m *AgeAssignee) Calculate(pullRequests []github.PullRequest) error {
	a := averageMap{}
	a.reset()

	for _, pr := range pullRequests {
		if !pr.IsMerged() {
			continue
		}

		delta := pr.MergedAt.Sub(pr.CreatedAt).Hours() / 24
		for _, assignee := range pr.Assignees {
			a.add(delta, assignee.Login)
		}
	}

	m.average = a.toList()
	return nil
}

// Converts the calculated data to a string
func (m *AgeAssignee) String() string {
	return m.average.string("days")
}
