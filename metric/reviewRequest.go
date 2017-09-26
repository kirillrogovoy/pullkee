package metric

import "github.com/kirillrogovoy/pullkee/github"

// ReviewRequest contains the calculated data
type ReviewRequest struct {
	counter counterMap
}

// Description of the metric
func (m *ReviewRequest) Description() string {
	return "How often one is requested for a review?"
}

// Calculate the data
func (m *ReviewRequest) Calculate(pullRequests []github.PullRequest) error {
	m.counter = counterMap{}
	for _, pr := range pullRequests {
		for _, requestee := range *pr.ReviewRequests {
			author := requestee.Login
			if _, ok := m.counter[author]; !ok {
				m.counter[author] = &counter{author, 0}
			}
			m.counter[author].Count++
		}
	}
	return nil
}

// Converts the calculated data to a string
func (m *ReviewRequest) String() string {
	return m.counter.string()
}
