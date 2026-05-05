package rules

import (
	"context"
	"strings"

	"config-audit/internal/config"
	"config-audit/internal/finding"
)

type WeakAlgorithmRule struct{}

var weakAlgorithms = map[string]string{
	"md5":      "MD5",
	"sha1":     "SHA-1",
	"des":      "DES",
	"3des":     "3DES",
	"rc4":      "RC4",
	"blowfish": "Blowfish",
}

func (WeakAlgorithmRule) ID() string {
	return "weak-algorithm"
}

func (rule WeakAlgorithmRule) Check(_ context.Context, doc config.Document) []finding.Finding {
	var findings []finding.Finding

	doc.Walk(func(node config.Node) {
		value, ok := config.ScalarString(node.Value)
		if !ok {
			return
		}

		if !keyHasAny(node.Key, "algorithm", "digest", "hash", "cipher", "signature") {
			return
		}

		normalized := normalizeToken(value)
		for token, displayName := range weakAlgorithms {
			if normalized == token || strings.Contains(normalized, token) {
				findings = append(findings, newFinding(
					rule.ID(),
					finding.SeverityHigh,
					node,
					"weak or deprecated algorithm is configured: "+displayName,
					"replace it with a modern algorithm appropriate for the use case",
				))
				return
			}
		}
	})

	return findings
}
