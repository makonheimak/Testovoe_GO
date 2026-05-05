package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"config-audit/internal/app"
	"config-audit/internal/finding"
	"config-audit/internal/grpcapi"
	"config-audit/internal/httpapi"
	"config-audit/internal/output"
)

const (
	ExitOK       = 0
	ExitFindings = 1
	ExitError    = 2
)

func Run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	options, err := ParseOptions(args, stderr)
	if err != nil {
		if errors.Is(err, ErrUsage) {
			_, _ = fmt.Fprintln(stderr, err)
			return ExitError
		}
		_, _ = fmt.Fprintln(stderr, "error:", err)
		return ExitError
	}

	analyzer := app.NewDefaultAnalyzer()
	ctx := context.Background()

	if options.HTTP {
		_, _ = fmt.Fprintf(stdout, "HTTP API listening on %s\n", options.HTTPAddr)
		if err := httpapi.Run(ctx, options.HTTPAddr, analyzer); err != nil {
			_, _ = fmt.Fprintln(stderr, "error:", err)
			return ExitError
		}
		return ExitOK
	}

	if options.GRPC {
		_, _ = fmt.Fprintf(stdout, "gRPC API listening on %s\n", options.GRPCAddr)
		if err := grpcapi.Run(ctx, options.GRPCAddr, analyzer); err != nil {
			_, _ = fmt.Fprintln(stderr, "error:", err)
			return ExitError
		}
		return ExitOK
	}

	var findingsOutputErr error
	findings, err := analyze(ctx, analyzer, options, stdin)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "error:", err)
		return ExitError
	}

	if options.JSON {
		findingsOutputErr = output.WriteJSON(stdout, findings)
	} else {
		findingsOutputErr = output.WriteText(stdout, findings)
	}
	if findingsOutputErr != nil {
		_, _ = fmt.Fprintln(stderr, "error:", findingsOutputErr)
		return ExitError
	}

	if len(findings) > 0 && !options.Silent {
		return ExitFindings
	}
	return ExitOK
}

func analyze(ctx context.Context, analyzer app.Analyzer, options Options, stdin io.Reader) ([]finding.Finding, error) {
	if options.UseStdin {
		return analyzer.AnalyzeReader(ctx, stdin, "stdin")
	}
	return analyzer.AnalyzePath(ctx, options.InputPath, options.Recursive, options.CheckPermission)
}
