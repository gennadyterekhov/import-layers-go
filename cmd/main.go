package main

import (
	"github.com/gennadyterekhov/import-layers-go/internal/analyzer"
	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg := config.FromFile()

	singlechecker.Main(
		analyzer.New(cfg).Analyzer,
	)
}
