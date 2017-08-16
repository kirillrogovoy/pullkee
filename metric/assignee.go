package metric

import "github.com/kirillrogovoy/pullk/github"

// Assignee contains the calculated data
type Assignee struct {
	counter counterByDev
}

// Description of the metric
func (a *Assignee) Description() string {
	return "How often one is the assignee?"
}

// Calculate how often each developer is an assignee
func (a *Assignee) Calculate(pullRequests github.PullRequests) error {
	a.counter = counterByDev{}
	for _, pr := range pullRequests {
		for _, assignee := range pr.Assignees {
			author := assignee.Login
			if _, ok := a.counter[author]; !ok {
				a.counter[author] = &counter{author, 0}
			}
			a.counter[author].Count++
		}
	}
	return nil
}

// Converts the calculated data to a string
func (a *Assignee) String() string {
	return a.counter.string()
}
