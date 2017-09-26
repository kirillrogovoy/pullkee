// Package metric contains all the metrics that analyze Pull Requests in
// different ways.
package metric

import "github.com/kirillrogovoy/pullkee/github"

// Metric is a common interface for every metric in metric/*.go
type Metric interface {
	Description() string
	Calculate(pullRequests []github.PullRequest) error
	String() string
}

// Metrics returns the list of all available metrics
func Metrics() []Metric {
	return []Metric{
		&Age{},
		&AgeAssignee{},
		&Assignee{},
		&ReviewRequest{},
		&DiffSize{},
		&DiffSizePerDay{},
		&Author{},
		&AssigneeMatrix{},
		&CommentCharsPerDay{},
		&AuthorComments{},
		&DescriptionSize{},
	}
}
