package app

import (
	"context"
	"fmt"
	"io"
	"os"

	"config-audit/internal/checker"
	"config-audit/internal/config"
	"config-audit/internal/dirscan"
	"config-audit/internal/filemode"
	"config-audit/internal/finding"
)

type Analyzer struct {
	checker checker.Checker
}

func NewAnalyzer(checker checker.Checker) Analyzer {
	return Analyzer{checker: checker}
}

func NewDefaultAnalyzer() Analyzer {
	return NewAnalyzer(checker.NewDefault())
}

func (analyzer Analyzer) AnalyzeBytes(ctx context.Context, data []byte, source string) ([]finding.Finding, error) {
	doc, err := config.Decode(data, source)
	if err != nil {
		return nil, err
	}

	findings := analyzer.checker.Check(ctx, doc)
	finding.Sort(findings)
	return findings, nil
}

func (analyzer Analyzer) AnalyzeReader(ctx context.Context, reader io.Reader, source string) ([]finding.Finding, error) {
	data, source, err := config.ReadAll(reader, source)
	if err != nil {
		return nil, err
	}
	return analyzer.AnalyzeBytes(ctx, data, source)
}

func (analyzer Analyzer) AnalyzeFile(ctx context.Context, path string, checkPermissions bool) ([]finding.Finding, error) {
	data, source, err := config.ReadFile(path)
	if err != nil {
		return nil, err
	}

	findings, err := analyzer.AnalyzeBytes(ctx, data, source)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", source, err)
	}

	if checkPermissions {
		permissionFindings, err := filemode.Check(source)
		if err != nil {
			return nil, err
		}
		findings = append(findings, permissionFindings...)
	}

	finding.Sort(findings)
	return findings, nil
}

func (analyzer Analyzer) AnalyzePath(ctx context.Context, path string, recursive bool, checkPermissions bool) ([]finding.Finding, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat input path %q: %w", path, err)
	}

	if !info.IsDir() {
		return analyzer.AnalyzeFile(ctx, path, checkPermissions)
	}

	if !recursive {
		return nil, fmt.Errorf("%q is a directory; use --recursive to scan config files", path)
	}

	paths, err := dirscan.FindConfigs(path)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no supported config files found in %q", path)
	}

	var allFindings []finding.Finding
	for _, configPath := range paths {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		findings, err := analyzer.AnalyzeFile(ctx, configPath, checkPermissions)
		if err != nil {
			return nil, err
		}
		allFindings = append(allFindings, findings...)
	}

	finding.Sort(allFindings)
	return allFindings, nil
}
