package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gennadyterekhov/import-layers-go/internal/project"
	"gopkg.in/yaml.v3"
)

const fileName = `import_layers.yaml`

type yamlConfig struct {
	Debug                          bool     `yaml:"Debug"`
	ReportDirsWithoutAssignedLayer bool     `yaml:"ReportDirsWithoutAssignedLayer"`
	OnlyAdjacent                   bool     `yaml:"OnlyAdjacent"`
	Layers                         []string `yaml:"Layers"`
}

type Config struct {
	debug                          bool
	reportDirsWithoutAssignedLayer bool
	onlyAdjacent                   bool
	layers                         map[string]int // keys are substrings of pkg paths
}

func FromFile() *Config {
	confPath, err := getFullConfigPath()
	if err != nil {
		panic(err)
	}

	rawFile, err := os.ReadFile(confPath)
	if err != nil {
		panic(err)
	}

	var yamlCfg yamlConfig

	if err = yaml.Unmarshal(rawFile, &yamlCfg); err != nil {
		panic(err)
	}

	if yamlCfg.Debug {
		fmt.Println("creating config from file")
		fmt.Println("raw Layers from yaml", yamlCfg.Layers)
	}

	var cfg Config
	cfg.layers = make(map[string]int, 0)
	cfg.reportDirsWithoutAssignedLayer = yamlCfg.ReportDirsWithoutAssignedLayer
	cfg.debug = yamlCfg.Debug

	ln := len(yamlCfg.Layers)
	for i, layer := range yamlCfg.Layers {
		cfg.layers[layer] = ln - i
	}

	if yamlCfg.Debug {
		fmt.Println("processed layers into map")
		fmt.Println("final Layers map", cfg.Layers())
	}

	return &cfg
}

func FromMap(layers map[string]int) *Config {
	var cfg Config
	cfg.reportDirsWithoutAssignedLayer = true
	cfg.layers = layers

	return &cfg
}

func getFullConfigPath() (string, error) {
	rootDir, err := project.GetProjectRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		rootDir, fileName,
	), nil
}

func (c *Config) ReportDirsWithoutAssignedLayer() bool {
	return c.reportDirsWithoutAssignedLayer
}

func (c *Config) OnlyAdjacent() bool {
	return c.onlyAdjacent
}

// GetLayer yes this is O(n). returns 0 only if fullPkgName is not in config (ignored)
func (c *Config) GetLayer(fullPkgName string) int {
	for substr, layer := range c.layers {
		if strings.Contains(fullPkgName, substr) {
			return layer
		}
	}
	return 0
}

func (c *Config) Layers() map[string]int {
	return c.layers
}

func (c *Config) Debug() bool {
	return c.debug
}
