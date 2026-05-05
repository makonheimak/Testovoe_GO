package rules

import (
	"context"
	"strings"

	"config-audit/internal/config"
	"config-audit/internal/finding"
)

type PlainPasswordRule struct{}

func (PlainPasswordRule) ID() string {
	return "plain-password"
}

func (rule PlainPasswordRule) Check(_ context.Context, doc config.Document) []finding.Finding {
	var findings []finding.Finding

	doc.Walk(func(node config.Node) {
		value, ok := config.ScalarString(node.Value)
		if !ok {
			return
		}

		value = strings.TrimSpace(value)
		if value == "" || isExternalSecretReference(value) {
			return
		}

		if keyHasAny(node.Key, "password", "passwd", "passphrase", "secret", "token", "apikey", "apiKey") {
			findings = append(findings, newFinding(
				rule.ID(),
				finding.SeverityHigh,
				node,
				"sensitive value appears to be stored directly in the config",
				"move secrets to environment variables or a secret manager",
			))
		}
	})

	return findings
}
