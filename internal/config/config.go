package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gennadyterekhov/import-layers-go/internal/layers"
	"github.com/gennadyterekhov/import-layers-go/internal/project"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug       bool
	IgnoreTests bool
	Root        project.Root
	RootDir     string
	Layers      layers.Policy
	LayerNames  []string
}

func FromFile() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return Read(cwd, "")
}

func Read(cwd, configPath string) (*Config, error) {
	defaultPath, err := findDefault(cwd)
	if err != nil {
		return nil, err
	}
	root := cwd
	if defaultPath != "" {
		root = filepath.Dir(defaultPath)
	} else if moduleRoot, err := project.FindRoot(cwd); err == nil {
		root = moduleRoot
	}

	if configPath == "" {
		if defaultPath == "" {
			return nil, errors.New("import_layers.json or import_layers.yaml not found")
		}
		configPath = defaultPath
	} else if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(cwd, configPath)
	}

	cfg, err := decode(configPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", configPath, err)
	}
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	cfg.Root = project.Root(canonicalRoot)
	return cfg, nil
}

func findDefault(cwd string) (string, error) {
	for dir := cwd; ; dir = filepath.Dir(dir) {
		for _, name := range []string{"import_layers.json", "import_layers.yaml", "import_layers.yml"} {
			candidate := filepath.Join(dir, name)
			info, err := os.Stat(candidate)
			switch {
			case err == nil && info.Mode().IsRegular():
				return candidate, nil
			case err != nil && !errors.Is(err, os.ErrNotExist):
				return "", err
			}
		}
		if filepath.Dir(dir) == dir {
			return "", nil
		}
	}
}

func decode(filename string) (*Config, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var data map[string]any
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(raw, &data)
	default:
		err = json.Unmarshal(raw, &data)
	}
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New("expected a configuration mapping")
	}

	names, err := layerNames(data["layers"])
	if err != nil {
		return nil, err
	}
	cfg := &Config{Layers: layers.New(names), LayerNames: names}
	if value, present := data["rootDir"]; present {
		rootDir, ok := value.(string)
		if !ok {
			return nil, errors.New("rootDir must be a directory path")
		}
		cfg.RootDir = rootDir
	}
	cfg.Debug, err = booleanField(data, "debug")
	if err != nil {
		return nil, err
	}
	ignoreTests := "ignoreTests"
	if _, present := data[ignoreTests]; !present {
		ignoreTests = "ignore_tests"
	}
	cfg.IgnoreTests, err = booleanField(data, ignoreTests)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func layerNames(value any) ([]string, error) {
	invalid := errors.New("layers must be a non-empty list of unique strings")
	values, ok := value.([]any)
	if !ok || len(values) == 0 {
		return nil, invalid
	}
	seen := make(map[string]bool, len(values))
	names := make([]string, 0, len(values))
	for _, value := range values {
		name, ok := value.(string)
		if !ok || strings.TrimSpace(name) == "" || seen[name] {
			return nil, invalid
		}
		seen[name] = true
		names = append(names, name)
	}
	return names, nil
}

func booleanField(data map[string]any, key string) (bool, error) {
	value, present := data[key]
	if !present {
		return false, nil
	}
	result, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("%s must be a boolean", key)
	}
	return result, nil
}
