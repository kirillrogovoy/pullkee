package metric

import "github.com/kirillrogovoy/pullkee/github"

// Author contains the calculated data
type Author struct {
	counter counterMap
}

// Description of the metric
func (m *Author) Description() string {
	return "How many PRs did one create?"
}

// Calculate how often each developer is an author
func (m *Author) Calculate(pullRequests []github.PullRequest) error {
	m.counter = counterMap{}
	for _, pr := range pullRequests {
		author := pr.User.Login
		if _, ok := m.counter[author]; !ok {
			m.counter[author] = &counter{author, 0}
		}
		m.counter[author].Count++
	}
	return nil
}

// Converts the calculated data to a string
func (m *Author) String() string {
	return m.counter.string()
}
