package metric

import "github.com/kirillrogovoy/pullkee/github"

// Assignee contains the calculated data
type Assignee struct {
	counter counterMap
}

// Description of the metric
func (m *Assignee) Description() string {
	return "How often one is the assignee?"
}

// Calculate how often each developer is an assignee
func (m *Assignee) Calculate(pullRequests []github.PullRequest) error {
	m.counter = counterMap{}
	for _, pr := range pullRequests {
		for _, assignee := range pr.Assignees {
			author := assignee.Login
			if _, ok := m.counter[author]; !ok {
				m.counter[author] = &counter{author, 0}
			}
			m.counter[author].Count++
		}
	}
	return nil
}

// Converts the calculated data to a string
func (m *Assignee) String() string {
	return m.counter.string()
}
