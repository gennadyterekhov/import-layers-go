package main

import (
	"fmt"
	"os"

	"github.com/gennadyterekhov/import-layers-go/internal/analyzer"
	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg, err := config.FromFile()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	singlechecker.Main(
		analyzer.New(cfg).Analyzer,
	)
}
