package rules

import (
	"context"
	"strings"

	"config-audit/internal/config"
	"config-audit/internal/finding"
)

type DebugLoggingRule struct{}

func (DebugLoggingRule) ID() string {
	return "debug-logging"
}

func (rule DebugLoggingRule) Check(_ context.Context, doc config.Document) []finding.Finding {
	var findings []finding.Finding

	doc.Walk(func(node config.Node) {
		if value, ok := config.ScalarBool(node.Value); ok && value && keyHasAny(node.Key, "debug", "debugmode") {
			findings = append(findings, newFinding(
				rule.ID(),
				finding.SeverityLow,
				node,
				"debug mode is enabled",
				"disable debug mode in production",
			))
			return
		}

		value, ok := config.ScalarString(node.Value)
		if !ok {
			return
		}

		level := strings.ToLower(strings.TrimSpace(value))
		if level != "debug" && level != "trace" {
			return
		}

		if keyHasAny(node.Key, "level", "loglevel") || pathHasAny(node.Path, "log", "logging", "logger") {
			findings = append(findings, newFinding(
				rule.ID(),
				finding.SeverityLow,
				node,
				"logging is configured in debug or trace mode",
				"use info or a stricter log level in production",
			))
		}
	})

	return findings
}
