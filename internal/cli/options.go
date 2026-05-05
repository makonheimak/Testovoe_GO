package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

type Options struct {
	InputPath       string
	UseStdin        bool
	Silent          bool
	Recursive       bool
	CheckPermission bool
	JSON            bool
	HTTP            bool
	HTTPAddr        string
	GRPC            bool
	GRPCAddr        string
}

var ErrUsage = errors.New("usage error")

func ParseOptions(args []string, usageWriter io.Writer) (Options, error) {
	var options Options
	var shortSilent bool

	flagSet := flag.NewFlagSet("config-audit", flag.ContinueOnError)
	flagSet.SetOutput(usageWriter)

	flagSet.BoolVar(&shortSilent, "s", false, "do not return exit code 1 when findings are found")
	flagSet.BoolVar(&options.Silent, "silent", false, "do not return exit code 1 when findings are found")
	flagSet.BoolVar(&options.UseStdin, "stdin", false, "read config from stdin instead of a file")
	flagSet.BoolVar(&options.Recursive, "recursive", false, "recursively scan a directory for JSON/YAML configs")
	flagSet.BoolVar(&options.CheckPermission, "check-permissions", false, "check config file permissions using os.Stat")
	flagSet.BoolVar(&options.JSON, "json", false, "print findings as JSON")
	flagSet.BoolVar(&options.HTTP, "http", false, "run REST API server instead of CLI analysis")
	flagSet.StringVar(&options.HTTPAddr, "addr", ":8080", "address for --http mode")
	flagSet.BoolVar(&options.GRPC, "grpc", false, "run gRPC API server instead of CLI analysis")
	flagSet.StringVar(&options.GRPCAddr, "grpc-addr", ":9090", "address for --grpc mode")

	flagSet.Usage = func() {
		_, _ = fmt.Fprintln(usageWriter, "Usage:")
		_, _ = fmt.Fprintln(usageWriter, "  config-audit [flags] <config-file>")
		_, _ = fmt.Fprintln(usageWriter, "  config-audit --stdin [flags]")
		_, _ = fmt.Fprintln(usageWriter, "  config-audit --http --addr :8080")
		_, _ = fmt.Fprintln(usageWriter, "  config-audit --grpc --grpc-addr :9090")
		_, _ = fmt.Fprintln(usageWriter)
		_, _ = fmt.Fprintln(usageWriter, "Flags:")
		flagSet.PrintDefaults()
	}

	if err := flagSet.Parse(args); err != nil {
		return Options{}, fmt.Errorf("%w: %v", ErrUsage, err)
	}

	options.Silent = options.Silent || shortSilent
	remaining := flagSet.Args()

	if options.HTTP && options.GRPC {
		return Options{}, fmt.Errorf("%w: choose either --http or --grpc", ErrUsage)
	}

	if options.HTTP {
		if len(remaining) != 0 {
			return Options{}, fmt.Errorf("%w: --http does not accept a config path", ErrUsage)
		}
		return options, nil
	}

	if options.GRPC {
		if len(remaining) != 0 {
			return Options{}, fmt.Errorf("%w: --grpc does not accept a config path", ErrUsage)
		}
		return options, nil
	}

	if options.UseStdin {
		if len(remaining) != 0 {
			return Options{}, fmt.Errorf("%w: --stdin does not accept a config path", ErrUsage)
		}
		return options, nil
	}

	if len(remaining) != 1 {
		return Options{}, fmt.Errorf("%w: expected exactly one config file or directory path", ErrUsage)
	}

	options.InputPath = remaining[0]
	return options, nil
}
