package checker

import (
	"context"

	"config-audit/internal/config"
	"config-audit/internal/finding"
	"config-audit/internal/rules"
)

type Checker struct {
	rules []rules.Rule
}

func New(ruleSet []rules.Rule) Checker {
	return Checker{rules: append([]rules.Rule(nil), ruleSet...)}
}

func NewDefault() Checker {
	return New(rules.DefaultRules())
}

func (checker Checker) Check(ctx context.Context, doc config.Document) []finding.Finding {
	var findings []finding.Finding

	for _, rule := range checker.rules {
		if err := ctx.Err(); err != nil {
			return findings
		}
		findings = append(findings, rule.Check(ctx, doc)...)
	}

	for i := range findings {
		if findings[i].Source == "" {
			findings[i].Source = doc.Source
		}
	}

	finding.Sort(findings)
	return findings
}
