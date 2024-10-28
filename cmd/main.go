package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gennadyterekhov/levelslib/internal/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

const ConfigFileName = `config.json`
const LevelslibLevel = 1

type Level struct {
	Name        string
	Order       int
	Description string // optional
}
type Config struct {
	ReportUnleveled bool
	ModeNumbers     bool
	Levels          []Level
}

func main() {
	// cfg := getConfig()

	multichecker.Main(
		getAnalyzers(nil)...,
	)
}

func getConfig() *Config {
	currentFile, err := os.Executable()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(currentFile), ConfigFileName))
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return &cfg
}

func getAnalyzers(cfg *Config) []*analysis.Analyzer {
	checks := []*analysis.Analyzer{
		analyzer.New(),
	}

	return checks
}
