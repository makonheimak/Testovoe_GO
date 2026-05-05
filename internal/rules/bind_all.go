package rules

import (
	"context"
	"strings"

	"config-audit/internal/config"
	"config-audit/internal/finding"
)

type BindAllRule struct{}

func (BindAllRule) ID() string {
	return "bind-all-interfaces"
}

func (rule BindAllRule) Check(_ context.Context, doc config.Document) []finding.Finding {
	var findings []finding.Finding

	doc.Walk(func(node config.Node) {
		value, ok := config.ScalarString(node.Value)
		if !ok {
			return
		}

		if !strings.Contains(strings.TrimSpace(value), "0.0.0.0") {
			return
		}

		if keyHasAny(node.Key, "host", "address", "addr", "bind", "listen", "endpoint", "url") {
			findings = append(findings, newFinding(
				rule.ID(),
				finding.SeverityMedium,
				node,
				"service is bound to all network interfaces",
				"bind to a private interface or enforce firewall and access-list restrictions",
			))
		}
	})

	return findings
}
