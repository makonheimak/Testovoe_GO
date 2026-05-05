package rules

import (
	"strings"
	"unicode"

	"config-audit/internal/config"
	"config-audit/internal/finding"
)

func newFinding(ruleID string, severity finding.Severity, node config.Node, message string, recommendation string) finding.Finding {
	return finding.Finding{
		RuleID:         ruleID,
		Severity:       severity,
		Path:           node.Path,
		Message:        message,
		Recommendation: recommendation,
	}
}

func normalizeToken(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func pathHasAny(path string, words ...string) bool {
	normalized := normalizeToken(path)
	for _, word := range words {
		if strings.Contains(normalized, normalizeToken(word)) {
			return true
		}
	}
	return false
}

func keyHasAny(key string, words ...string) bool {
	return pathHasAny(key, words...)
}

func isExternalSecretReference(value string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	return strings.HasPrefix(trimmed, "$") ||
		strings.HasPrefix(trimmed, "${") ||
		strings.HasPrefix(trimmed, "env:") ||
		strings.HasPrefix(trimmed, "vault:") ||
		strings.HasPrefix(trimmed, "secret://") ||
		strings.HasPrefix(trimmed, "file:")
}
