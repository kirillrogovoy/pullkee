package metric

import (
	"fmt"
	"sort"

	"github.com/kirillrogovoy/pullk/github"
)

// AssigneeMatrix contains the calculated data
type AssigneeMatrix struct {
	counter map[string]map[string]*counter
}

// Description of the metric
func (a *AssigneeMatrix) Description() string {
	return "How often does one developer pick another as an assignee?"
}

// Calculate how often each developer is an assignee
func (a *AssigneeMatrix) Calculate(pullRequests github.PullRequests) error {
	a.counter = map[string]map[string]*counter{}
	for _, pr := range pullRequests {
		author := pr.User.Login
		if _, ok := a.counter[author]; !ok {
			a.counter[author] = map[string]*counter{}
		}
		for _, assignee := range pr.Assignees {
			name := assignee.Login
			if _, ok := a.counter[author][name]; !ok {
				a.counter[author][name] = &counter{name, 0}
			}
			a.counter[author][name].Count++
		}
	}
	return nil
}

// Converts the calculated data to a string
func (a *AssigneeMatrix) String() string {
	keys := []string{}

	for key := range a.counter {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	result := ""
	for _, key := range keys {
		author := a.counter[key]
		result += fmt.Sprintf("Author %s:\n", key)

		var cs counters
		for _, c := range author {
			cs = append(cs, *c)
		}
		sort.Sort(sort.Reverse(cs))
		for _, c := range cs {
			result += fmt.Sprintf("\t%s: %d\n", c.Name, c.Count)
		}
	}

	return result
}
