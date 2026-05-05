package rules_test

import (
	"context"
	"testing"

	"config-audit/internal/checker"
	"config-audit/internal/config"
	"config-audit/internal/finding"
)

func TestDefaultRulesFindUnsafeConfig(t *testing.T) {
	input := []byte(`
{
  "log": {"level": "debug"},
  "database": {"password": "plain-text"},
  "server": {"host": "0.0.0.0"},
  "tls": {"enabled": false},
  "storage": {"digest-algorithm": "MD5"}
}`)

	doc, err := config.Decode(input, "config.json")
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	findings := checker.NewDefault().Check(context.Background(), doc)

	assertHasRule(t, findings, "debug-logging", finding.SeverityLow)
	assertHasRule(t, findings, "plain-password", finding.SeverityHigh)
	assertHasRule(t, findings, "bind-all-interfaces", finding.SeverityMedium)
	assertHasRule(t, findings, "tls-disabled", finding.SeverityHigh)
	assertHasRule(t, findings, "weak-algorithm", finding.SeverityHigh)
}

func TestDefaultRulesDetectYAMLUnsafeConfig(t *testing.T) {
	input := []byte(`
debug: true
server:
  bind: 0.0.0.0
tls:
  verify: false
storage:
  digest-algorithm: MD5
database:
  password: plain-text
`)

	doc, err := config.Decode(input, "config.yaml")
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	findings := checker.NewDefault().Check(context.Background(), doc)

	assertHasRule(t, findings, "debug-logging", finding.SeverityLow)
	assertHasRule(t, findings, "plain-password", finding.SeverityHigh)
	assertHasRule(t, findings, "bind-all-interfaces", finding.SeverityMedium)
	assertHasRule(t, findings, "tls-disabled", finding.SeverityHigh)
	assertHasRule(t, findings, "weak-algorithm", finding.SeverityHigh)
}

func TestDefaultRulesDetectIndividualCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		filename string
		ruleID   string
		severity finding.Severity
	}{
		{
			name:     "debug boolean",
			input:    `{"debug":true}`,
			filename: "config.json",
			ruleID:   "debug-logging",
			severity: finding.SeverityLow,
		},
		{
			name:     "trace logging",
			input:    `{"logging":{"level":"trace"}}`,
			filename: "config.json",
			ruleID:   "debug-logging",
			severity: finding.SeverityLow,
		},
		{
			name:     "plain token",
			input:    `{"auth":{"token":"abc123"}}`,
			filename: "config.json",
			ruleID:   "plain-password",
			severity: finding.SeverityHigh,
		},
		{
			name:     "bind all",
			input:    `{"server":{"listen":"0.0.0.0:8080"}}`,
			filename: "config.json",
			ruleID:   "bind-all-interfaces",
			severity: finding.SeverityMedium,
		},
		{
			name:     "insecure skip verify",
			input:    `{"client":{"insecureSkipVerify":true}}`,
			filename: "config.json",
			ruleID:   "tls-disabled",
			severity: finding.SeverityHigh,
		},
		{
			name:     "sha1 algorithm",
			input:    `{"crypto":{"hash":"SHA1"}}`,
			filename: "config.json",
			ruleID:   "weak-algorithm",
			severity: finding.SeverityHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := config.Decode([]byte(tt.input), tt.filename)
			if err != nil {
				t.Fatalf("Decode returned error: %v", err)
			}

			findings := checker.NewDefault().Check(context.Background(), doc)
			assertHasRule(t, findings, tt.ruleID, tt.severity)
		})
	}
}

func TestDefaultRulesIgnoreSafeConfig(t *testing.T) {
	input := []byte(`
{
  "log": {"level": "info"},
  "database": {"password": "${DB_PASSWORD}"},
  "server": {"host": "127.0.0.1"},
  "tls": {"enabled": true},
  "storage": {"digest-algorithm": "SHA-256"}
}`)

	doc, err := config.Decode(input, "config.json")
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	findings := checker.NewDefault().Check(context.Background(), doc)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %#v", findings)
	}
}

func assertHasRule(t *testing.T, findings []finding.Finding, ruleID string, severity finding.Severity) {
	t.Helper()

	for _, item := range findings {
		if item.RuleID == ruleID && item.Severity == severity {
			return
		}
	}

	t.Fatalf("expected finding %s with severity %s, got %#v", ruleID, severity, findings)
}
