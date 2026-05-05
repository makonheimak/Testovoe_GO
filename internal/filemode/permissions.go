package filemode

import (
	"fmt"
	"os"

	"config-audit/internal/finding"
)

const RuleID = "file-permissions"

func Check(path string) ([]finding.Finding, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat config file %q: %w", path, err)
	}

	mode := info.Mode().Perm()
	if mode == 0 {
		return nil, nil
	}

	if mode&0o002 != 0 {
		return []finding.Finding{newPermissionFinding(path, finding.SeverityHigh, mode, "config file is writable by other users")}, nil
	}

	if mode&0o020 != 0 {
		return []finding.Finding{newPermissionFinding(path, finding.SeverityMedium, mode, "config file is writable by group users")}, nil
	}

	if mode&0o004 != 0 {
		return []finding.Finding{newPermissionFinding(path, finding.SeverityMedium, mode, "config file is readable by other users")}, nil
	}

	return nil, nil
}

func newPermissionFinding(source string, severity finding.Severity, mode os.FileMode, message string) finding.Finding {
	return finding.Finding{
		Source:         source,
		RuleID:         RuleID,
		Severity:       severity,
		Path:           "file.mode",
		Message:        fmt.Sprintf("%s: %04o", message, mode),
		Recommendation: "restrict config file permissions to the application owner",
	}
}
