package metric

import "github.com/kirillrogovoy/pullk/github"

// ReviewRequest contains the calculated data
type ReviewRequest struct {
	counter counterByDev
}

// Description of the metric
func (a *ReviewRequest) Description() string {
	return "How often one is requested for a review?"
}

// Calculate the data
func (a *ReviewRequest) Calculate(pullRequests []github.PullRequest) error {
	a.counter = counterByDev{}
	for _, pr := range pullRequests {
		for _, requestee := range *pr.ReviewRequests {
			author := requestee.Login
			if _, ok := a.counter[author]; !ok {
				a.counter[author] = &counter{author, 0}
			}
			a.counter[author].Count++
		}
	}
	return nil
}

// Converts the calculated data to a string
func (a *ReviewRequest) String() string {
	return a.counter.string()
}
