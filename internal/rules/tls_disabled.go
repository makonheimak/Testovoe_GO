package rules

import (
	"context"

	"config-audit/internal/config"
	"config-audit/internal/finding"
)

type TLSDisabledRule struct{}

func (TLSDisabledRule) ID() string {
	return "tls-disabled"
}

func (rule TLSDisabledRule) Check(_ context.Context, doc config.Document) []finding.Finding {
	var findings []finding.Finding

	doc.Walk(func(node config.Node) {
		value, ok := config.ScalarBool(node.Value)
		if !ok {
			return
		}

		if value && (keyHasAny(node.Key, "insecure", "skipverify", "insecureskipverify") || keyHasAny(node.Key, "disabletls", "tlsdisabled")) {
			findings = append(findings, newFinding(
				rule.ID(),
				finding.SeverityHigh,
				node,
				"TLS or certificate verification is explicitly disabled",
				"enable TLS and certificate verification for external connections",
			))
			return
		}

		if !value && keyHasAny(node.Key, "enabled", "verify", "verifycert", "verifycertificate") && pathHasAny(node.Path, "tls", "ssl", "https", "certificate", "cert") {
			findings = append(findings, newFinding(
				rule.ID(),
				finding.SeverityHigh,
				node,
				"TLS or certificate verification is disabled",
				"enable TLS and certificate verification for external connections",
			))
		}
	})

	return findings
}
