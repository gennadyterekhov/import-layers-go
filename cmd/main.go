package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gennadyterekhov/levelslib/internal/collector"
	"github.com/gennadyterekhov/levelslib/internal/data"
	"github.com/gennadyterekhov/levelslib/internal/finalizer"
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

	commonData := data.New()
	dataCollector := collector.New(commonData)
	finalizerAnalyzer := finalizer.New(commonData, dataCollector)

	multichecker.Main(
		finalizerAnalyzer.Analyzer,
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
