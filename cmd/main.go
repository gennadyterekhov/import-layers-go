package main

import (
	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"github.com/gennadyterekhov/import-layers-go/internal/finalizer"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	cfg := config.FromFile()
	finalizerAnalyzer := finalizer.New(cfg)

	multichecker.Main(
		finalizerAnalyzer.Analyzer,
	)
}
