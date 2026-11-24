package config

import "testing"

func TestCanGetConfigFromFile(t *testing.T) {
	cfg, err := FromFile()
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"project", "layers", "config", "filechecker", "analyzer", "cli"} {
		directory := "internal/" + name
		packagePath := "github.com/gennadyterekhov/import-layers-go/" + directory
		if layer := cfg.Layers.Match(packagePath, directory); layer.Name != directory {
			t.Errorf("matching %q: got %q, want %q", packagePath, layer.Name, directory)
		}
	}
}
