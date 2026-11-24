package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gennadyterekhov/import-layers-go/internal/analyzer"
	"github.com/gennadyterekhov/import-layers-go/internal/layers"
)

type options struct {
	configPath string
	graph      bool
	help       bool
	paths      []string
}

func Run(cwd string, args []string, stdout, stderr io.Writer) int {
	opts, err := parseArguments(args)
	if err != nil {
		return reportError(stderr, err)
	}
	if opts.help {
		fmt.Fprintln(stdout, "Usage: import-layers-go [-c|--config PATH] [--graph] [path ...]")
		fmt.Fprintln(stdout, "Paths may be project directories, Go files, or Go package patterns.")
		return 0
	}

	debugOutput := stdout
	if opts.graph {
		debugOutput = stderr
	}
	result, err := analyzer.CheckProject(analyzer.Request{
		Directory: cwd, ConfigPath: opts.configPath, Paths: opts.paths,
		OnFile: func(path string) { fmt.Fprintln(debugOutput, "processing file "+path) },
	})
	if err != nil {
		return reportError(stderr, err)
	}
	if opts.graph {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result.Graph); err != nil {
			return reportError(stderr, err)
		}
	}

	for _, message := range result.Errors {
		fmt.Fprintln(stderr, message)
	}
	for _, issue := range result.Violations {
		fmt.Fprintf(stderr, "%s:%d -> %s: %s (%s -> %s)\n",
			issue.File, issue.Line, issue.Import, layers.ViolationMessage, issue.FromLayer, issue.ToLayer)
	}
	switch {
	case len(result.Errors) > 0:
		return 2
	case len(result.Violations) > 0:
		return 1
	default:
		if !opts.graph {
			fmt.Fprintln(stdout, "OK")
		}
		return 0
	}
}

func reportError(stderr io.Writer, err error) int {
	fmt.Fprintf(stderr, "import-layers-go: %v\n", err)
	return 2
}

func parseArguments(args []string) (options, error) {
	var opts options
	for i := 0; i < len(args); i++ {
		arg := args[i]
		key, value, hasValue := strings.Cut(arg, "=")
		switch key {
		case "-h", "--help":
			opts.help = true
		case "--graph":
			if hasValue {
				return opts, errors.New("--graph does not take a value")
			}
			opts.graph = true
		case "-c", "--config":
			if !hasValue && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				value = args[i]
			}
			if value == "" {
				return opts, fmt.Errorf("%s requires a path to a config file", key)
			}
			opts.configPath = value
		case "--":
			opts.paths = append(opts.paths, args[i+1:]...)
			return opts, nil
		default:
			if strings.HasPrefix(arg, "-") {
				return opts, fmt.Errorf("unknown option %s", arg)
			}
			opts.paths = append(opts.paths, arg)
		}
	}
	return opts, nil
}
