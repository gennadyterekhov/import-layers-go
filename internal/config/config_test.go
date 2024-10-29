package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanGetConfigFromFile(t *testing.T) {
	cfg := FromFile()

	assert.Equal(t, 4, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/project"))
	assert.Equal(t, 3, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/data"))
	assert.Equal(t, 2, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/config"))
	assert.Equal(t, 1, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/finalizer"))
}
