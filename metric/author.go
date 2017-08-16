package metric

import "github.com/kirillrogovoy/pullk/github"

// Author contains the calculated data
type Author struct {
	counter counterByDev
}

// Description of the metric
func (a *Author) Description() string {
	return "How many PRs did one create?"
}

// Calculate how often each developer is an author
func (a *Author) Calculate(pullRequests github.PullRequests) error {
	a.counter = counterByDev{}
	for _, pr := range pullRequests {
		author := pr.User.Login
		if _, ok := a.counter[author]; !ok {
			a.counter[author] = &counter{author, 0}
		}
		a.counter[author].Count++
	}
	return nil
}

// Converts the calculated data to a string
func (a *Author) String() string {
	return a.counter.string()
}
