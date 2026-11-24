package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gennadyterekhov/import-layers-go/internal/analyzer"
	"github.com/gennadyterekhov/import-layers-go/internal/cli"
	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	if usesNativeDriver(os.Args[1:]) {
		cfg, err := config.FromFile()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		singlechecker.Main(analyzer.New(cfg).Analyzer)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "import-layers-go: %v\n", err)
		os.Exit(2)
	}
	os.Exit(cli.Run(cwd, os.Args[1:], os.Stdout, os.Stderr))
}

// Preserve the Go vet protocol and the native singlechecker diagnostic mode.
func usesNativeDriver(args []string) bool {
	for index := 0; index < len(args); index++ {
		arg := args[index]
		flag, _, hasValue := strings.Cut(arg, "=")
		switch flag {
		case "-flags", "-V", "-json", "-test", "-fix", "-debug", "-tags", "-cpuprofile", "-memprofile", "-trace":
			return true
		case "-c", "--config":
			if !hasValue && index+1 < len(args) {
				index++
			}
			continue
		case "--":
			return false
		}
		if strings.HasSuffix(arg, ".cfg") {
			return true
		}
	}
	return false
}
