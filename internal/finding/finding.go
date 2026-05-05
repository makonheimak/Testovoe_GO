package finding

import "sort"

type Severity string

const (
	SeverityLow    Severity = "LOW"
	SeverityMedium Severity = "MEDIUM"
	SeverityHigh   Severity = "HIGH"
)

type Finding struct {
	Source         string   `json:"source,omitempty"`
	RuleID         string   `json:"rule_id"`
	Severity       Severity `json:"severity"`
	Path           string   `json:"path,omitempty"`
	Message        string   `json:"message"`
	Recommendation string   `json:"recommendation"`
}

func (s Severity) Valid() bool {
	switch s {
	case SeverityLow, SeverityMedium, SeverityHigh:
		return true
	default:
		return false
	}
}

func (s Severity) Rank() int {
	switch s {
	case SeverityHigh:
		return 3
	case SeverityMedium:
		return 2
	case SeverityLow:
		return 1
	default:
		return 0
	}
}

func Sort(findings []Finding) {
	sort.SliceStable(findings, func(i, j int) bool {
		left := findings[i]
		right := findings[j]

		if left.Severity.Rank() != right.Severity.Rank() {
			return left.Severity.Rank() > right.Severity.Rank()
		}
		if left.Source != right.Source {
			return left.Source < right.Source
		}
		if left.Path != right.Path {
			return left.Path < right.Path
		}
		return left.RuleID < right.RuleID
	})
}
