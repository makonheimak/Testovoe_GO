package rules

import (
	"context"

	"config-audit/internal/config"
	"config-audit/internal/finding"
)

type Rule interface {
	ID() string
	Check(context.Context, config.Document) []finding.Finding
}

func DefaultRules() []Rule {
	return []Rule{
		DebugLoggingRule{},
		PlainPasswordRule{},
		BindAllRule{},
		TLSDisabledRule{},
		WeakAlgorithmRule{},
	}
}
