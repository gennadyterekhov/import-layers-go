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
	IgnoreTests bool     `yaml:"ignore_tests"`
	Debug       bool     `yaml:"debug"`
	Layers      []string `yaml:"layers"`
}

type Config struct {
	debug       bool
	ignoreTests bool
	layers      map[string]int // keys are substrings of pkg paths
}

func FromFile() (*Config, error) {
	absolutePathToConfig, err := getFullConfigPath()
	if err != nil {
		return nil, fmt.Errorf("could not get path to config")
	}

	rawFile, err := os.ReadFile(absolutePathToConfig)
	if err != nil {
		return nil, fmt.Errorf("config file %v could not be read", absolutePathToConfig)
	}

	var yamlCfg yamlConfig

	if err = yaml.Unmarshal(rawFile, &yamlCfg); err != nil {
		return nil, fmt.Errorf("config file %v could not be decoded from yaml", absolutePathToConfig)
	}

	var cfg Config
	cfg.layers = make(map[string]int)
	cfg.debug = yamlCfg.Debug
	// don't read it from config so it's always false
	//cfg.ignoreTests = yamlCfg.IgnoreTests

	ln := len(yamlCfg.Layers)
	for i, layer := range yamlCfg.Layers {
		cfg.layers[layer] = ln - i
	}

	return &cfg, nil
}

func FromMap(layers map[string]int) *Config {
	var cfg Config
	cfg.layers = layers
	cfg.debug = true
	return &cfg
}

func (c *Config) SetIgnoreTests(v bool) *Config {
	c.ignoreTests = v
	return c
}

func (c *Config) ReportDirsWithoutAssignedLayer() bool {
	return false
}

func (c *Config) OnlyAdjacent() bool {
	return false
}

func (c *Config) IgnoreTests() bool {
	return c.ignoreTests
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

func getFullConfigPath() (string, error) {
	rootDir, err := project.GetProjectRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		rootDir, fileName,
	), nil
}
