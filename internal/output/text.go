package output

import (
	"fmt"
	"io"

	"config-audit/internal/finding"
)

func WriteText(writer io.Writer, findings []finding.Finding) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(writer, "No issues found.")
		return err
	}

	for _, item := range findings {
		if item.Source != "" {
			if _, err := fmt.Fprintf(writer, "%s [%s] %s %s\n", item.Severity, item.RuleID, item.Source, formatPath(item.Path)); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(writer, "%s [%s] %s\n", item.Severity, item.RuleID, formatPath(item.Path)); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintf(writer, "  Problem: %s\n", item.Message); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(writer, "  Recommendation: %s\n", item.Recommendation); err != nil {
			return err
		}
	}

	return nil
}

func formatPath(path string) string {
	if path == "" {
		return "(root)"
	}
	return path
}
