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
func (m *AssigneeMatrix) Description() string {
	return "How often does one developer pick another as an assignee?"
}

// Calculate how often each developer is an assignee
func (m *AssigneeMatrix) Calculate(pullRequests []github.PullRequest) error {
	m.counter = map[string]map[string]*counter{}
	for _, pr := range pullRequests {
		author := pr.User.Login
		if _, ok := m.counter[author]; !ok {
			m.counter[author] = map[string]*counter{}
		}
		for _, assignee := range pr.Assignees {
			name := assignee.Login
			if _, ok := m.counter[author][name]; !ok {
				m.counter[author][name] = &counter{name, 0}
			}
			m.counter[author][name].Count++
		}
	}
	return nil
}

// Converts the calculated data to a string
func (m *AssigneeMatrix) String() string {
	keys := []string{}

	for key := range m.counter {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	result := ""
	for _, key := range keys {
		author := m.counter[key]
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
